package events

type PublisherInterface interface {
	Publish(eventType string, payload any, metadata map[string]string) error
	Close() error
}
