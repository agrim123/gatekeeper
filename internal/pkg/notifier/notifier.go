package notifier

import "github.com/spf13/viper"

type Notifier interface {
	Notify(message string) error
	FallbackNotify(message string)
}

func GetNotifier() Notifier {
	switch viper.GetString("notifier.type") {
	case "slack":
		return Slack{
			Hook: viper.GetString("notifier.slack.hook"),
		}
	default:
		return NewDefaultNotifier()
	}
}
