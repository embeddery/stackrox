package types

import (
	"bitbucket.org/stack-rox/apollo/generated/api/v1"
)

// Notifier interface defines the contract that all plugins must satisfy
type Notifier interface {
	// Notify triggers the plugins to send a notification
	Notify(alert *v1.Alert) error
	// ProtoNotifier gets the proto version of the notifier
	ProtoNotifier() *v1.Notifier
	// Test sends a test message
	Test() error
}
