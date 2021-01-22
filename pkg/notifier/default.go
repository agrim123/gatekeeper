package notifier

import (
	"github.com/agrim123/gatekeeper/pkg/logger"
)

type Default struct{}

func NewDefaultNotifier() Default {
	return Default{}
}

func (d Default) Notify(message string) error {
	logger.Notifier(message)
	return nil
}

func (d Default) Name() string {
	return "default"
}
