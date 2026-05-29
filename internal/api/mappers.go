package api

import (
	"github.com/mandacode-labs/dodream/internal/card"
	"github.com/mandacode-labs/dodream/internal/collection"
	"github.com/mandacode-labs/dodream/internal/deck"
	"github.com/mandacode-labs/dodream/internal/studyevent"
	"github.com/mandacode-labs/dodream/internal/user"
	gen "github.com/mandacode-labs/dodream/pkg/api/v1"
)

func mapUser(u *user.User) *gen.User {
	return &gen.User{
		ID:         gen.NewOptString(u.ID().String()),
		Nickname:   gen.NewOptString(u.Nickname()),
		ProviderID: gen.NewOptString(u.ProviderID()),
		CreatedAt:  gen.NewOptDateTime(u.CreatedAt()),
		UpdatedAt:  gen.NewOptDateTime(u.UpdatedAt()),
	}
}

func mapCard(c *card.Card) *gen.Card {
	return &gen.Card{
		ID:        gen.NewOptString(c.ID().String()),
		Question:  gen.NewOptString(c.Question()),
		Hint:      gen.NewOptString(c.Hint()),
		Content:   gen.NewOptString(c.Content()),
		Creator:   gen.NewOptString(c.Creator()),
		CreatedAt: gen.NewOptDateTime(c.CreatedAt()),
		UpdatedAt: gen.NewOptDateTime(c.UpdatedAt()),
	}
}

func mapDeck(d *deck.Deck) *gen.Deck {
	return &gen.Deck{
		ID:        gen.NewOptString(d.ID().String()),
		Name:      gen.NewOptString(d.Name()),
		Creator:   gen.NewOptString(d.Creator()),
		CreatedAt: gen.NewOptDateTime(d.CreatedAt()),
		UpdatedAt: gen.NewOptDateTime(d.UpdatedAt()),
	}
}

func mapCollection(c *collection.Collection) *gen.Collection {
	return &gen.Collection{
		ID:        gen.NewOptString(c.ID().String()),
		Name:      gen.NewOptString(c.Name()),
		Creator:   gen.NewOptString(c.Creator()),
		CreatedAt: gen.NewOptDateTime(c.CreatedAt()),
		UpdatedAt: gen.NewOptDateTime(c.UpdatedAt()),
	}
}

func mapStudyEvent(e *studyevent.Event) *gen.StudyEvent {
	se := &gen.StudyEvent{
		ID:           gen.NewOptString(e.ID().String()),
		CollectionID: gen.NewOptString(e.CollectionID()),
		CardID:       gen.NewOptString(e.CardID()),
		UserID:       gen.NewOptString(e.UserID()),
		EventType:    gen.NewOptString(string(e.EventType())),
		CreatedAt:    gen.NewOptDateTime(e.CreatedAt()),
	}
	if e.Quality() != nil {
		se.Quality = gen.NewOptInt(*e.Quality())
	}
	if e.ResponseTime() != nil {
		se.ResponseTimeNs = gen.NewOptInt64(e.ResponseTime().Nanoseconds())
	}
	if e.PreviousDeckID() != nil {
		se.PreviousDeckID = gen.NewOptString(*e.PreviousDeckID())
	}
	return se
}
