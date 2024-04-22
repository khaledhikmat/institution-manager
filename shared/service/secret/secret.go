package secret

import (
	"context"

	dapr "github.com/dapr/go-sdk/client"

	"github.com/khaledhikmat/institution-manager/shared/equates"
)

type secretService struct {
	CanxCtx    context.Context
	Daprclient dapr.Client
}

func NewSecretService(canxCtx context.Context, client dapr.Client) (IService, error) {
	return &secretService{
		CanxCtx:    canxCtx,
		Daprclient: client,
	}, nil
}

func (c *secretService) Creds(alias string) (map[string]string, error) {
	opt := map[string]string{
		"version": "2",
	}

	secretValue, err := c.Daprclient.GetSecret(c.CanxCtx, equates.InstitutionManagerSecrets, alias, opt)
	if err != nil {
		return nil, err
	}

	return secretValue, nil
}

func (c *secretService) Finalize() {
	// TODO:
}
