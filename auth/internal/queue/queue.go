package queue

type Producer interface {
	SendMessage(key string, value string) (partition int32, offset int64, err error)
	GetTopic() string
}
