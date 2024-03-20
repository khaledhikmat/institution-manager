package memberactor

import (
	"context"
	"fmt"
)

const (
	miscStateKey         string = "misc"
	mainStateKey         string = "main"
	institutionsStateKey string = "institutions"
	membershipsStateKey  string = "memberships"
	transactionsStateKey string = "transactions"
	pledgesStateKey      string = "pledges"
)

type miscState struct {
	ExchangeRate float64 `json:"exchangeRate"`
}

func (t *MemberActor) saveState(ctx context.Context, stateKey string) {
	if stateKey == miscStateKey {
		t.GetStateManager().Set(ctx, stateKey, t.miscState)
	}

	if stateKey == mainStateKey {
		t.GetStateManager().Set(ctx, stateKey, t.mainState)
	}

	if stateKey == institutionsStateKey {
		t.GetStateManager().Set(ctx, stateKey, t.institutionsState)
	}

	if stateKey == membershipsStateKey {
		t.GetStateManager().Set(ctx, stateKey, t.membershipsState)
	}

	if stateKey == transactionsStateKey {
		t.GetStateManager().Set(ctx, stateKey, t.transactionsState)
	}

	if stateKey == pledgesStateKey {
		t.GetStateManager().Set(ctx, stateKey, t.pledgesState)
	}

	t.GetStateManager().Save(ctx)
}

func (t *MemberActor) getState(ctx context.Context, stateKey string) {
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
		c, err := MemberService.GetMember(t.ID())
		if err != nil {
			fmt.Println("getState error", err)
			return
		}

		t.mainState = c
		t.GetStateManager().Set(ctx, stateKey, t.mainState)
		return
	}

	if !exist && stateKey == institutionsStateKey {
		// Load from domain database
		ps, err := MemberService.GetInstitutions(t.ID())
		if err != nil {
			fmt.Println("getState error", err)
			return
		}

		t.institutionsState = ps
		t.GetStateManager().Set(ctx, stateKey, t.institutionsState)
		return
	}

	if !exist && stateKey == membershipsStateKey {
		// Load from domain database
		ps, err := MemberService.GetMemberships(t.ID())
		if err != nil {
			fmt.Println("getState error", err)
			return
		}

		t.membershipsState = ps
		t.GetStateManager().Set(ctx, stateKey, t.membershipsState)
		return
	}

	if !exist && stateKey == transactionsStateKey {
		// Load from domain database
		ps, err := MemberService.GetTransactions(t.ID())
		if err != nil {
			fmt.Println("getState error", err)
			return
		}

		t.transactionsState = ps
		t.GetStateManager().Set(ctx, stateKey, t.transactionsState)
		return
	}

	if !exist && stateKey == pledgesStateKey {
		// Load from domain database
		ps, err := MemberService.GetPledges(t.ID())
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

	if stateKey == institutionsStateKey {
		t.GetStateManager().Get(ctx, stateKey, &t.institutionsState)
	}

	if stateKey == membershipsStateKey {
		t.GetStateManager().Get(ctx, stateKey, &t.membershipsState)
	}

	if stateKey == transactionsStateKey {
		t.GetStateManager().Get(ctx, stateKey, &t.transactionsState)
	}

	if stateKey == pledgesStateKey {
		t.GetStateManager().Get(ctx, stateKey, &t.pledgesState)
	}
}
