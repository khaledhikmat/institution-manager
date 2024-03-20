package notification

import "context"

type notificationService struct {
	CanxCtx context.Context
}

func NewService(canxCtx context.Context) IService {
	return &notificationService{
		CanxCtx: canxCtx,
	}
}

func (c *notificationService) Notifiy(from, subject, body string, to []string) error {
	return nil
}
