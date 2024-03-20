package institutionactor

import (
	"context"
	"fmt"
)

const (
	mainStateKey string = "main"
	miscStateKey string = "misc"
)

type miscState struct {
	ExchangeRate float64 `json:"exchangeRate"`
}

func (t *InstitutionActor) saveState(ctx context.Context, stateKey string) {
	if stateKey == miscStateKey {
		t.GetStateManager().Set(ctx, stateKey, t.miscState)
	}

	if stateKey == mainStateKey {
		t.GetStateManager().Set(ctx, stateKey, t.mainState)
	}

	t.GetStateManager().Save(ctx)
}

func (t *InstitutionActor) getState(ctx context.Context, stateKey string) {
	exist, err := t.GetStateManager().Contains(ctx, stateKey)
	if err != nil {
		fmt.Println("getState error", err)
		return
	}

	// If the state does not exist, initialize a new state struct and return it
	if !exist && stateKey == mainStateKey {
		// Load from domain database
		c, err := InstitutionService.GetInstitution(t.ID())
		if err != nil {
			fmt.Println("getState error", err)
			return
		}

		t.mainState = c
		t.GetStateManager().Set(ctx, stateKey, t.mainState)
		return
	}

	if !exist && stateKey == miscStateKey {
		// TODO: Load from exchange service
		t.miscState = miscState{}
		t.miscState.ExchangeRate = 1
		t.GetStateManager().Set(ctx, stateKey, t.miscState)
		return
	}

	// If the state exists, retrieve from store into the actor struct
	if stateKey == mainStateKey {
		t.GetStateManager().Get(ctx, stateKey, &t.mainState)
		return
	}

	if stateKey == miscStateKey {
		t.GetStateManager().Get(ctx, stateKey, &t.miscState)
		return
	}
}
