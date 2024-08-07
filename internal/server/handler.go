package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/worldline-go/wkafka"

	"github.com/worldline-go/molen/internal/decoder"
	"github.com/worldline-go/molen/internal/kafka"
)

type Handler struct {
	Ctx            context.Context //nolint:containedctx // no need
	Client         *wkafka.Client
	ClientAdmin    *kadm.Client
	Group          kafka.Group
	ProduceMessage func(ctx context.Context, records []*kgo.Record) error
}

// @Summary Create group
// @Description Create a specific group and not consume any message.
// @Router /v1/group [post]
// @Param payload body GroupRequest{} false "topic and group_id"
// @Accept application/json
// @Success 200 {object} APIRespond{}
// @failure 400 {object} APIRespond{}
// @failure 500 {object} APIRespond{}
func (h Handler) CreateGroup(c echo.Context) error {
	group := GroupRequest{}

	if err := c.Bind(&group); err != nil {
		return c.JSON(http.StatusBadRequest, APIRespond{Message: fmt.Sprintf("unable to bind request err: %s", err)})
	}

	if group.GroupID == "" {
		return c.JSON(http.StatusBadRequest, APIRespond{Message: "group_id is empty"})
	}

	if err := h.Group.CreateGroup(h.Ctx, group.GroupID, group.Topic, group.From); err != nil {
		return c.JSON(http.StatusBadRequest, APIRespond{Message: fmt.Sprintf("unable to create group err: %s", err)})
	}

	log.Ctx(c.Request().Context()).Debug().Str("group_id", group.GroupID).Str("topic", group.Topic).Msg("group created")

	return c.JSON(http.StatusOK, APIRespond{Message: fmt.Sprintf("successfully created group [%s] for topic [%s]", group.GroupID, group.Topic)})
}

// @Summary Publish message
// @Description Publish message(s) to kafka with topic and partition(optional)
// @Router /v1/publish [post]
// @Param topic query string true "topic name"
// @Param partition query int32 false "specific partition number"
// @Param key query string false "key"
// @Param raw query bool false "raw body"
// @Param payload body interface{} false "send key values" SchemaExample()
// @Accept application/json
// @Success 200 {object} APIRespond{}
// @failure 400 {object} APIRespond{}
// @failure 500 {object} APIRespond{}
func (h Handler) Publish(c echo.Context) error {
	publish := PublishRequest{}
	binder := echo.DefaultBinder{}
	if err := binder.BindQueryParams(c, &publish); err != nil {
		return c.JSON(http.StatusBadRequest, APIRespond{Message: fmt.Sprintf("unable to bind request err: %s", err)})
	}

	if publish.Topic == "" {
		return c.JSON(http.StatusBadRequest, APIRespond{Message: "topic is empty"})
	}

	body := c.Request().Body
	if body == nil {
		return c.JSON(http.StatusBadRequest, APIRespond{Message: "body is empty"})
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIRespond{Message: fmt.Sprintf("unable to read body err: %s", err)})
	}

	// key
	var key []byte
	if publish.Key != "" {
		key = []byte(publish.Key)
	}

	msg := []*kgo.Record{}
	if !publish.Raw {
		if !json.Valid(data) {
			return c.JSON(http.StatusBadRequest, APIRespond{Message: "body is not valid json"})
		}

		var msgs decoder.Messages
		if err := json.Unmarshal(data, &msgs); err != nil {
			return c.JSON(http.StatusBadRequest,
				APIRespond{Message: fmt.Sprintf("unable to unmarshal json err: %s", err)},
			)
		}

		for _, m := range msgs {
			msg = append(msg, &kgo.Record{
				Value:     m,
				Topic:     publish.Topic,
				Partition: publish.Partition,
				Key:       key,
			})
		}
	} else {
		msg = append(msg, &kgo.Record{
			Value:     data,
			Topic:     publish.Topic,
			Partition: publish.Partition,
			Key:       key,
		})
	}

	if err := h.ProduceMessage(h.Ctx, msg); err != nil {
		if errors.Is(err, kerr.UnknownTopicOrPartition) {
			return c.JSON(http.StatusBadRequest,
				APIRespond{Message: fmt.Sprintf("topic [%s] or specific partition does not exist", publish.Topic)},
			)
		}

		return c.JSON(http.StatusInternalServerError,
			APIRespond{Message: fmt.Sprintf("unable to write message err: %s", err)},
		)
	}

	log.Ctx(c.Request().Context()).Debug().Str("topic", publish.Topic).Str("data", string(data)).Msgf("published write")

	return c.JSON(http.StatusOK,
		APIRespond{Message: fmt.Sprintf("successfully published to [%s]", publish.Topic)},
	)
}
