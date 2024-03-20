package campaign

import (
	"context"
	"fmt"
)

// In-memory campaigns database
var campaigns []Campaign = []Campaign{
	{
		Id:            "100",
		Name:          "Gaza Rebuild",
		Description:   "This is an initial contribution from the good people of the earth. But we will make the zionists cough up the rest.",
		InstitutionID: "100",
		Virtual:       true,
		ImageUrl:      "",
		Goal:          10000000,
		Currency:      "USD",
		Duration:      1000,
	},
	{
		Id:            "200",
		Name:          "Khan Yunis Rebuild",
		Description:   "This is an initial contribution from the good people of the earth. But we will make the zionists cough up the rest.",
		InstitutionID: "100",
		Virtual:       true,
		ImageUrl:      "",
		Goal:          10000,
		Currency:      "USD",
		Duration:      1000,
	},
	{
		Id:            "300",
		Name:          "Rafah Rebuild",
		Description:   "This is an initial contribution from the good people of the earth. But we will make the zionists cough up the rest.",
		InstitutionID: "100",
		Virtual:       true,
		ImageUrl:      "",
		Goal:          50000,
		Currency:      "USD",
		Duration:      1000,
	},
}

type campaignService struct {
	CanxCtx context.Context
}

func NewService(canxCtx context.Context) IService {
	return &campaignService{
		CanxCtx: canxCtx,
	}
}

func (c *campaignService) GetCampaigns(search string, opts ...CampaignOpt) ([]Campaign, error) {

	filter := Campaign{}
	for _, applyOpt := range opts {
		applyOpt(filter)
	}

	return campaigns, nil
}

func (c *campaignService) GetCampaign(id string) (Campaign, error) {
	for _, cmp := range campaigns {
		if cmp.Id == id {
			return cmp, nil
		}
	}

	return Campaign{}, fmt.Errorf("unable to find campaign [%s]", id)
}

func (c *campaignService) NewCampaign(cmp Campaign) (Campaign, error) {
	return Campaign{}, nil
}

func (c *campaignService) UpdateCampaign(cmp *Campaign) error {
	return nil
}

func (c *campaignService) Finalize() {
	// TODO:
}
