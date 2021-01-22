package notifier

import (
	"github.com/agrim123/gatekeeper/pkg/logger"
)

type Notifier interface {
	Name() string
	Notify(message string) error
}

func AttachFallbackNotifier(notifer Notifier) func(string) {
	return func(message string) {
		err := notifer.Notify(message)
		if err != nil {
			logger.Error("Notifier: %s failed. Fallback to default notifier", logger.Underline(notifer.Name()))
			logger.Notifier(message)
		}
	}
}
