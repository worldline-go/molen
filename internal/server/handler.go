package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/worldline-go/wkafka"
)

type Handler struct {
	Ctx            context.Context //nolint:containedctx // no need
	Client         *wkafka.Client
	ProduceMessage func(ctx context.Context, data ...Message) error
}

type Message struct {
	Data      []byte
	Topic     string
	Partition int32
	Key       []byte
}

func (m Message) MarshalJSON() ([]byte, error) {
	return m.Data, nil
}

func (m Message) ProducerHook(r *wkafka.Record) *wkafka.Record {
	r.Topic = m.Topic
	r.Partition = m.Partition
	r.Value = m.Data
	r.Key = m.Key

	return r
}

// @Summary Publish message
// @Description Publish message(s) to kafka with topic and partition(optional)
// @Security ApiKeyAuth || OAuth2AccessCode
// @Router /publish [post]
// @Param topic query string true "topic name"
// @Param partition query int32 false "specific partition number"
// @Param key query string false "key"
// @Param payload body interface{} false "send key values" SchemaExample()
// @Accept application/json
// @Success 200 {object} APIRespond{}
// @failure 400 {object} APIRespond{}
// @failure 500 {object} APIRespond{}
func (h Handler) Publish(c echo.Context) error {
	topic := c.QueryParam("topic")
	if topic == "" {
		return c.JSON(http.StatusBadRequest, APIRespond{Error: "topic is empty"}) //nolint:wrapcheck // no need
	}
	partitionRaw := c.QueryParam("partition")

	// string to int32
	partition := int32(0)
	if partitionRaw != "" {
		partition64, err := strconv.ParseInt(partitionRaw, 10, 32)
		if err != nil {
			return c.JSON(http.StatusBadRequest, //nolint:wrapcheck // no need
				APIRespond{Error: fmt.Sprintf("unable to parse partition err: %s", err)},
			)
		}

		partition = int32(partition64)
	}

	body := c.Request().Body
	if body == nil {
		return c.JSON(http.StatusBadRequest, APIRespond{Error: "body is empty"}) //nolint:wrapcheck // no need
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIRespond{Error: fmt.Sprintf("unable to read body err: %s", err)}) //nolint:wrapcheck // no need
	}

	if !json.Valid(data) {
		return c.JSON(http.StatusBadRequest, APIRespond{Error: "body is not valid json"}) //nolint:wrapcheck // no need
	}

	var key []byte
	keyRaw := c.QueryParam("key")
	if keyRaw != "" {
		key = []byte(keyRaw)
	}

	msg := Message{
		Data:      data,
		Topic:     topic,
		Partition: partition,
		Key:       key,
	}

	if err := h.ProduceMessage(h.Ctx, msg); err != nil {
		if errors.Is(err, kerr.UnknownTopicOrPartition) {
			return c.JSON(http.StatusBadRequest, //nolint:wrapcheck // no need
				APIRespond{Error: fmt.Sprintf("topic [%s] or specific partition does not exist", topic)},
			)
		}

		return c.JSON(http.StatusInternalServerError, //nolint:wrapcheck // no need
			APIRespond{Error: fmt.Sprintf("unable to write message err: %s", err)},
		)
	}

	log.Ctx(c.Request().Context()).Debug().Str("topic", topic).Str("data", string(data)).Msgf("published write")

	return c.JSON(http.StatusOK, //nolint:wrapcheck // no need
		APIRespond{Message: fmt.Sprintf("successfully published to [%s]", topic)},
	)
}
