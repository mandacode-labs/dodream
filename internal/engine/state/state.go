// Package state provides the SRS state domain.
package state

import (
	"time"

	"github.com/google/uuid"
)

// ID is a unique identifier for a State.
type ID string

// String returns the string representation of the ID.
func (id ID) String() string {
	return string(id)
}

// SRSState holds FSRS-related learning metrics for a user-card pair.
type SRSState struct {
	nextReviewAt    time.Time
	interval        float64
	easeFactor      float64
	totalReviews    int
	totalSuccessful int
	streak          int
	reviewsLast1d   int
	reviewsLast3d   int
	reviewsLast7d   int
	avgResponseTime *time.Duration
	lastReviewAt    *time.Time
}

// NextReviewAt returns the next scheduled review time.
func (s *SRSState) NextReviewAt() time.Time { return s.nextReviewAt }

// Interval returns the current review interval in days.
func (s *SRSState) Interval() float64 { return s.interval }

// EaseFactor returns the ease factor for SRS calculations.
func (s *SRSState) EaseFactor() float64 { return s.easeFactor }

// TotalReviews returns the total number of reviews.
func (s *SRSState) TotalReviews() int { return s.totalReviews }

// TotalSuccessful returns the number of successful reviews.
func (s *SRSState) TotalSuccessful() int { return s.totalSuccessful }

// Streak returns the current streak of successful reviews.
func (s *SRSState) Streak() int { return s.streak }

// ReviewsLast1d returns the number of reviews in the last 1 day.
func (s *SRSState) ReviewsLast1d() int { return s.reviewsLast1d }

// ReviewsLast3d returns the number of reviews in the last 3 days.
func (s *SRSState) ReviewsLast3d() int { return s.reviewsLast3d }

// ReviewsLast7d returns the number of reviews in the last 7 days.
func (s *SRSState) ReviewsLast7d() int { return s.reviewsLast7d }

// AvgResponseTime returns the average response time, or nil if not recorded.
func (s *SRSState) AvgResponseTime() *time.Duration { return s.avgResponseTime }

// LastReviewAt returns the last review time, or nil if never reviewed.
func (s *SRSState) LastReviewAt() *time.Time { return s.lastReviewAt }

// State represents the SRS learning state for a user-card pair.
type State struct {
	id        ID
	userID    string
	cardID    string
	srs       SRSState
	createdAt time.Time
	updatedAt time.Time
}

// New creates a new State for a user-card pair with default SRS values.
func New(userID string, cardID string) *State {
	now := time.Now()
	return &State{
		id:     ID(uuid.New().String()),
		userID: userID,
		cardID: cardID,
		srs: SRSState{
			nextReviewAt: now,
			interval:     0,
			easeFactor:   2.5,
		},
		createdAt: now,
		updatedAt: now,
	}
}

// ID returns the unique identifier of the state.
func (s *State) ID() ID { return s.id }

// UserID returns the user identifier.
func (s *State) UserID() string { return s.userID }

// CardID returns the card identifier.
func (s *State) CardID() string { return s.cardID }

// SRS returns the SRS learning metrics.
func (s *State) SRS() SRSState { return s.srs }

// CreatedAt returns the timestamp when the state was created.
func (s *State) CreatedAt() time.Time { return s.createdAt }

// UpdatedAt returns the timestamp when the state was last modified.
func (s *State) UpdatedAt() time.Time { return s.updatedAt }

// NextReviewAt returns the next scheduled review time.
func (s *State) NextReviewAt() time.Time { return s.srs.nextReviewAt }

// Interval returns the current review interval in days.
func (s *State) Interval() float64 { return s.srs.interval }

// EaseFactor returns the ease factor for SRS calculations.
func (s *State) EaseFactor() float64 { return s.srs.easeFactor }

// TotalReviews returns the total number of reviews.
func (s *State) TotalReviews() int { return s.srs.totalReviews }

// TotalSuccessful returns the number of successful reviews.
func (s *State) TotalSuccessful() int { return s.srs.totalSuccessful }

// Streak returns the current streak of successful reviews.
func (s *State) Streak() int { return s.srs.streak }

// ReviewsLast1d returns the number of reviews in the last 1 day.
func (s *State) ReviewsLast1d() int { return s.srs.reviewsLast1d }

// ReviewsLast3d returns the number of reviews in the last 3 days.
func (s *State) ReviewsLast3d() int { return s.srs.reviewsLast3d }

// ReviewsLast7d returns the number of reviews in the last 7 days.
func (s *State) ReviewsLast7d() int { return s.srs.reviewsLast7d }

// AvgResponseTime returns the average response time, or nil if not recorded.
func (s *State) AvgResponseTime() *time.Duration { return s.srs.avgResponseTime }

// LastReviewAt returns the last review time, or nil if never reviewed.
func (s *State) LastReviewAt() *time.Time { return s.srs.lastReviewAt }

// SetNextReviewAt updates the next review time.
func (s *State) SetNextReviewAt(t time.Time) {
	s.srs.nextReviewAt = t
	s.updatedAt = time.Now()
}

// SetInterval updates the review interval.
func (s *State) SetInterval(interval float64) {
	s.srs.interval = interval
	s.updatedAt = time.Now()
}

// SetEaseFactor updates the ease factor.
func (s *State) SetEaseFactor(ease float64) {
	s.srs.easeFactor = ease
	s.updatedAt = time.Now()
}

// IncrementTotalReviews increments the total review count.
func (s *State) IncrementTotalReviews() {
	s.srs.totalReviews++
	s.updatedAt = time.Now()
}

// IncrementTotalSuccessful increments the successful review count.
func (s *State) IncrementTotalSuccessful() {
	s.srs.totalSuccessful++
	s.updatedAt = time.Now()
}

// SetStreak updates the streak count.
func (s *State) SetStreak(streak int) {
	s.srs.streak = streak
	s.updatedAt = time.Now()
}

// SetReviewsLast1d updates the 1-day review count.
func (s *State) SetReviewsLast1d(count int) {
	s.srs.reviewsLast1d = count
	s.updatedAt = time.Now()
}

// SetReviewsLast3d updates the 3-day review count.
func (s *State) SetReviewsLast3d(count int) {
	s.srs.reviewsLast3d = count
	s.updatedAt = time.Now()
}

// SetReviewsLast7d updates the 7-day review count.
func (s *State) SetReviewsLast7d(count int) {
	s.srs.reviewsLast7d = count
	s.updatedAt = time.Now()
}

// SetAvgResponseTime updates the average response time.
func (s *State) SetAvgResponseTime(d time.Duration) {
	s.srs.avgResponseTime = &d
	s.updatedAt = time.Now()
}

// SetLastReviewAt updates the last review time.
func (s *State) SetLastReviewAt(t time.Time) {
	s.srs.lastReviewAt = &t
	s.updatedAt = time.Now()
}

// SetID sets the state ID.
func (s *State) SetID(id ID) {
	s.id = id
}

// SetCreatedAt sets the created timestamp.
func (s *State) SetCreatedAt(t time.Time) {
	s.createdAt = t
}

// SetSRS sets the entire SRS state.
func (s *State) SetSRS(srs SRSState) {
	s.srs = srs
	s.updatedAt = time.Now()
}
