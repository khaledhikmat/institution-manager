package notification

type IService interface {
	Notifiy(from, subject, body string, to []string) error
}
