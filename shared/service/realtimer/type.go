package realtimer

type IService interface {
	Token() (string, error)
	Externalize(channelName string, messageName string, data any) error

	Finalize()
}
