package membershipactor

import (
	"context"
	"fmt"
)

const (
	miscStateKey         string = "misc"
	mainStateKey         string = "main"
	transactionsStateKey string = "transactions"
)

type miscState struct {
	ExchangeRate float64 `json:"exchangeRate"`
}

func (t *MembershipActor) saveState(ctx context.Context, stateKey string) {
	if stateKey == miscStateKey {
		t.GetStateManager().Set(ctx, stateKey, t.miscState)
	}

	if stateKey == mainStateKey {
		t.GetStateManager().Set(ctx, stateKey, t.mainState)
	}

	if stateKey == transactionsStateKey {
		t.GetStateManager().Set(ctx, stateKey, t.transactionsState)
	}

	t.GetStateManager().Save(ctx)
}

func (t *MembershipActor) getState(ctx context.Context, stateKey string) {
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
		c, err := MembershipService.GetMembership(t.ID())
		if err != nil {
			fmt.Println("getState error", err)
			return
		}

		t.mainState = c
		t.GetStateManager().Set(ctx, stateKey, t.mainState)
		return
	}

	if !exist && stateKey == transactionsStateKey {
		// Load from domain database
		ps, err := MembershipService.GetTransactions()
		if err != nil {
			fmt.Println("getState error", err)
			return
		}

		t.transactionsState = ps
		t.GetStateManager().Set(ctx, stateKey, t.transactionsState)
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

	if stateKey == transactionsStateKey {
		t.GetStateManager().Get(ctx, stateKey, &t.transactionsState)
	}
}
