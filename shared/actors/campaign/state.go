package campaignactor

import (
	"context"
	"fmt"
)

const (
	miscStateKey    string = "misc"
	mainStateKey    string = "main"
	pledgesStateKey string = "pledges"
)

type miscState struct {
	ExchangeRate float64 `json:"exchangeRate"`
}

func (t *CampaignActor) saveState(ctx context.Context, stateKey string) {
	if stateKey == miscStateKey {
		t.GetStateManager().Set(ctx, stateKey, t.miscState)
	}

	if stateKey == mainStateKey {
		t.GetStateManager().Set(ctx, stateKey, t.mainState)
	}

	if stateKey == pledgesStateKey {
		t.GetStateManager().Set(ctx, stateKey, t.pledgesState)
	}

	t.GetStateManager().Save(ctx)
}

func (t *CampaignActor) getState(ctx context.Context, stateKey string) {
	exist, err := t.GetStateManager().Contains(ctx, stateKey)
	if err != nil {
		fmt.Println("getState error", err)
		return
	}

	// If the state does not exist, initialize a new state struct and return it
	if !exist && stateKey == miscStateKey {
		// TODO: Load from exchange service
		t.miscState = miscState{}
		t.miscState.ExchangeRate = 1
		t.GetStateManager().Set(ctx, stateKey, t.miscState)
		return
	}

	if !exist && stateKey == mainStateKey {
		// Load from domain database
		c, err := CampaignService.GetCampaign(t.ID())
		if err != nil {
			fmt.Println("getState error", err)
			return
		}

		t.mainState = c
		t.GetStateManager().Set(ctx, stateKey, t.mainState)
		return
	}

	if !exist && stateKey == pledgesStateKey {
		// Load from domain database
		ps, err := MemberService.GetPledgesByCampaign(t.ID())
		if err != nil {
			fmt.Println("getState error", err)
			return
		}

		t.pledgesState = ps
		t.GetStateManager().Set(ctx, stateKey, t.pledgesState)
		return
	}

	// If the state exists, retrieve from store into the actor struct
	if stateKey == miscStateKey {
		t.GetStateManager().Get(ctx, stateKey, &t.miscState)
		return
	}

	if stateKey == mainStateKey {
		t.GetStateManager().Get(ctx, stateKey, &t.mainState)
		return
	}

	if stateKey == pledgesStateKey {
		t.GetStateManager().Get(ctx, stateKey, &t.pledgesState)
	}
}
