package realtimer

import (
	"context"
	"fmt"

	"github.com/ably/ably-go/ably"
)

type ablyService struct {
	CanxCtx  context.Context
	AblyRest *ably.REST
}

func NewAblyService(canxCtx context.Context, apiKey, client string) (IService, error) {
	// Setup Ably REST gateway
	// REST client is all we need. Realtimeclient is not needed at the server.
	// Also....since we are authenticating using an API KEY, no need for token refresh
	ablyr, err := ably.NewREST(ably.WithKey(apiKey),
		ably.WithClientID(client),
	)
	if err != nil {
		fmt.Println("Failed to start ably client", err)
		return nil, err
	}

	return &ablyService{
		CanxCtx:  canxCtx,
		AblyRest: ablyr,
	}, nil
}

func (c *ablyService) Token() (string, error) {
	// WARNING: We are assuming that the user is authenticated because they are accessing HTMX
	tokenDetail, err := c.AblyRest.Auth.RequestToken(c.CanxCtx, nil)
	if err != nil {
		return "", err
	}

	return tokenDetail.Token, nil
}

func (c *ablyService) Externalize(channelName string, messageName string, data any) error {
	channel := c.AblyRest.Channels.Get(channelName)
	err := channel.Publish(c.CanxCtx, messageName, data)
	if err != nil {
		return err
	}

	return nil
}

func (c *ablyService) Finalize() {
	// TODO:
}
