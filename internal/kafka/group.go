package kafka

import (
	"context"
	"errors"

	"github.com/twmb/franz-go/pkg/kadm"
)

var errExit = errors.New("exit")

type Group struct {
	ClientAdmin *kadm.Client
}

func (g *Group) CreateGroup(ctx context.Context, groupID string, topic string, from string) error {
	v := kadm.Offsets{}
	if from == "" {
		v.AddOffset(topic, 0, 0, -1)
	} else {
		offsets, err := g.ClientAdmin.FetchOffsetsForTopics(ctx, from, topic)
		if err != nil {
			return err
		}

		v = offsets.Offsets()
	}

	_, err := g.ClientAdmin.CommitOffsets(ctx, groupID, v)

	return err
}
