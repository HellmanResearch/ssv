package controller

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/bloxapp/ssv/protocol/v2/types"
)

// OnTimeout is trigger upon timeout for the given height
func (c *Controller) OnTimeout(logger *zap.Logger, msg types.EventMsg) error {
	// TODO add validation

	timeoutData, err := msg.GetTimeoutData()
	if err != nil {
		return errors.Wrap(err, "failed to get timeout data")
	}
	instance := c.StoredInstances.FindInstance(timeoutData.Height)
	if instance == nil {
		return errors.New("instance is nil")
	}
	decided, _ := instance.IsDecided()
	if decided {
		return nil
	}
	return instance.UponRoundTimeout(logger)
}
