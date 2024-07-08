package kafka

import (
	"context"
	"errors"

	"github.com/twmb/franz-go/pkg/kadm"
)

var errAlreadyExists = errors.New("group already exists")

type Group struct {
	ClientAdmin *kadm.Client
}

func (g *Group) IsExists(ctx context.Context, groupID, topic string) (bool, error) {
	offsets, err := g.ClientAdmin.FetchOffsets(ctx, groupID)
	if err != nil {
		return false, err
	}

	for k := range offsets.KOffsets() {
		if k == topic {
			return true, nil
		}
	}

	return false, nil
}

func (g *Group) CreateGroup(ctx context.Context, groupID string, topic string, from string) error {
	if exists, err := g.IsExists(ctx, groupID, topic); err != nil {
		return err
	} else if exists {
		return errAlreadyExists
	}

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
