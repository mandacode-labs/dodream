package api

import (
	"github.com/mandacode-labs/dodream/internal/card"
	"github.com/mandacode-labs/dodream/internal/collection"
	"github.com/mandacode-labs/dodream/internal/deck"
	"github.com/mandacode-labs/dodream/internal/studyevent"
	"github.com/mandacode-labs/dodream/internal/user"
	"github.com/mandacode-labs/dodream/pkg/oas"
)

func mapUser(u *user.User) *oas.User {
	return &oas.User{
		ID:         oas.NewOptString(u.ID().String()),
		Nickname:   oas.NewOptString(u.Nickname()),
		ProviderID: oas.NewOptString(u.ProviderID()),
		CreatedAt:  oas.NewOptDateTime(u.CreatedAt()),
		UpdatedAt:  oas.NewOptDateTime(u.UpdatedAt()),
	}
}

func mapCard(c *card.Card) *oas.Card {
	return &oas.Card{
		ID:        oas.NewOptString(c.ID().String()),
		Question:  oas.NewOptString(c.Question()),
		Hint:      oas.NewOptString(c.Hint()),
		Content:   oas.NewOptString(c.Content()),
		Creator:   oas.NewOptString(c.CreatorID()),
		CreatedAt: oas.NewOptDateTime(c.CreatedAt()),
		UpdatedAt: oas.NewOptDateTime(c.UpdatedAt()),
	}
}

func mapDeck(d *deck.Deck) *oas.Deck {
	return &oas.Deck{
		ID:        oas.NewOptString(d.ID().String()),
		Name:      oas.NewOptString(d.Name()),
		Creator:   oas.NewOptString(d.CreatorID()),
		CreatedAt: oas.NewOptDateTime(d.CreatedAt()),
		UpdatedAt: oas.NewOptDateTime(d.UpdatedAt()),
	}
}

func mapCollection(c *collection.Collection) *oas.Collection {
	return &oas.Collection{
		ID:        oas.NewOptString(c.ID().String()),
		Name:      oas.NewOptString(c.Name()),
		Creator:   oas.NewOptString(c.CreatorID()),
		CreatedAt: oas.NewOptDateTime(c.CreatedAt()),
		UpdatedAt: oas.NewOptDateTime(c.UpdatedAt()),
	}
}

func mapStudyEvent(e *studyevent.Event) *oas.StudyEvent {
	se := &oas.StudyEvent{
		ID:           oas.NewOptString(e.ID().String()),
		CollectionID: oas.NewOptString(e.CollectionID()),
		CardID:       oas.NewOptString(e.CardID()),
		UserID:       oas.NewOptString(e.UserID()),
		EventType:    oas.NewOptString(string(e.EventType())),
		CreatedAt:    oas.NewOptDateTime(e.CreatedAt()),
	}
	if e.Quality() != nil {
		se.Quality = oas.NewOptInt(*e.Quality())
	}
	if e.ResponseTime() != nil {
		se.ResponseTimeNs = oas.NewOptInt64(e.ResponseTime().Nanoseconds())
	}
	if e.PreviousDeckID() != nil {
		se.PreviousDeckID = oas.NewOptString(*e.PreviousDeckID())
	}
	return se
}
