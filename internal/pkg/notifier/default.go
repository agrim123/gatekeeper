package notifier

import "fmt"

type Default struct{}

func NewDefaultNotifier() Default {
	return Default{}
}

func (d Default) Notify(message string) error {
	d.FallbackNotify(message)
	return nil
}

func (d Default) FallbackNotify(message string) {
	fmt.Println("[Notifier] " + message)
}
