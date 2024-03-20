package campaign

import (
	"time"
)

type CampaignOpt func(c Campaign)

func WithInstitution(i string) CampaignOpt {
	return func(c Campaign) {
		c.InstitutionID = i
	}
}

func WithVirtual(v bool) CampaignOpt {
	return func(c Campaign) {
		c.Virtual = v
	}
}

func WithGoal(g int64) CampaignOpt {
	return func(c Campaign) {
		c.Goal = g
	}
}

func WithCurrency(curr string) CampaignOpt {
	return func(c Campaign) {
		c.Currency = curr
	}
}

type Campaign struct {
	Id                 string    `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	InstitutionID      string    `json:"institutionId"`
	Virtual            bool      `json:"virtual"`
	ImageUrl           string    `json:"imageUrl"`
	Goal               int64     `json:"goal"`
	Currency           string    `json:"currency"`
	Duration           int64     `json:"duration"`
	GoalUSD            int64     `json:"goalUsd"`
	AwayFromGoal       int64     `json:"awayFromGoal"`
	MaxPledgeAmount    int64     `json:"maxPledgeAmount"`
	MinPledgeAmount    int64     `json:"minPledgeAmount"`
	AvgPledgeAmount    int64     `json:"avgPledgeAmount"`
	TotPledgeAmount    int64     `json:"totPledgeAmount"`
	TotPledgeAmountUSD int64     `json:"totPledgeAmountUSD"`
	Pledges            int64     `json:"pledges"`
	Donors             int64     `json:"donors"`
	MostGenerousDonor  string    `json:"mostGenerous"`
	StartTime          time.Time `json:"startTime"`
	EndTime            time.Time `json:"endTime"`
}
