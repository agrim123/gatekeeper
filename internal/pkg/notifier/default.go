package notifier

import "fmt"

type Default struct{}

func NewDefaultNotifier() Default {
	return Default{}
}

func (d Default) Notify(message string) error {
	fmt.Println("[Notifier] " + message)
	return nil
}
