package institution

import (
	"context"
	"fmt"
)

// In-memory institutions database
var institutions []Institution = []Institution{
	{
		Id:      "100",
		Name:    "Unrwa",
		City:    "Gaza",
		State:   "Gaza Strip",
		Country: "Palestine",
	},
}

type institutionService struct {
	CanxCtx context.Context
}

func NewService(canxCtx context.Context) IService {
	return &institutionService{
		CanxCtx: canxCtx,
	}
}

func (c *institutionService) GetInstitutions(search string, opts ...InstitutionOpt) ([]Institution, error) {

	filter := Institution{}
	for _, applyOpt := range opts {
		applyOpt(filter)
	}

	return institutions, nil
}

func (c *institutionService) GetInstitution(id string) (Institution, error) {
	for _, ins := range institutions {
		if ins.Id == id {
			return ins, nil
		}
	}

	return Institution{}, fmt.Errorf("unable to find institution [%s]", id)
}

func (c *institutionService) NewInstitution(ins Institution) (Institution, error) {
	return Institution{}, nil
}

func (c *institutionService) UpdateInstitution(ins *Institution) error {
	return nil
}

func (c *institutionService) Finalize() {
	// TODO:
}
