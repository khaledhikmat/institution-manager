package campaign

import (
	"context"
	"fmt"
)

// In-memory campaigns database
var campaigns = []Campaign{
	{
		ID:            "100",
		Type:          ConfirmedCampaign,
		PaymentType:   PrePledgePaymentCampaign,
		Name:          "CPRE Gaza Rebuild",
		Description:   "This is an initial contribution from the good people of the earth. But we will make the zionists cough up the rest.",
		InstitutionID: "100",
		Virtual:       true,
		ImageURL:      "",
		Goal:          10000000,
		Currency:      "USD",
		Duration:      1000,
	},
	{
		ID:            "200",
		Type:          ConfirmedCampaign,
		PaymentType:   PostPledgePaymentCampaign,
		Name:          "CPOST Khan Yunis Rebuild",
		Description:   "This is an initial contribution from the good people of the earth. But we will make the zionists cough up the rest.",
		InstitutionID: "100",
		Virtual:       true,
		ImageURL:      "",
		Goal:          10000,
		Currency:      "USD",
		Duration:      1000,
	},
	{
		ID:            "300",
		Type:          UnconfirmedCampaign,
		Name:          "U Rafah Rebuild",
		Description:   "This is an initial contribution from the good people of the earth. But we will make the zionists cough up the rest.",
		InstitutionID: "100",
		Virtual:       true,
		ImageURL:      "",
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

func (c *campaignService) GetCampaigns(search string, opts ...Opt) ([]Campaign, error) {

	filter := Campaign{}
	for _, applyOpt := range opts {
		applyOpt(filter)
	}

	return campaigns, nil
}

func (c *campaignService) GetCampaign(id string) (Campaign, error) {
	for _, cmp := range campaigns {
		if cmp.ID == id {
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
