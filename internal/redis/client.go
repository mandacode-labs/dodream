package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client wraps the Redis client for Dodream-specific operations.
type Client struct {
	client *redis.Client
}

// NewClient creates a new Redis client from the given URL.
func NewClient(redisURL string) (*Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis url: %w", err)
	}

	return &Client{
		client: redis.NewClient(opt),
	}, nil
}

// Close closes the Redis client connection.
func (c *Client) Close() error {
	return c.client.Close()
}

// Ping checks the Redis connection.
func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// AddCollectionCard adds a card to a collection's card set.
func (c *Client) AddCollectionCard(ctx context.Context, collectionID string, cardID string) error {
	key := fmt.Sprintf("collection:%s:cards", collectionID)
	return c.client.SAdd(ctx, key, cardID).Err()
}

// RemoveCollectionCard removes a card from a collection's card set.
func (c *Client) RemoveCollectionCard(ctx context.Context, collectionID string, cardID string) error {
	key := fmt.Sprintf("collection:%s:cards", collectionID)
	return c.client.SRem(ctx, key, cardID).Err()
}

// GetCollectionCards returns all card IDs in a collection.
func (c *Client) GetCollectionCards(ctx context.Context, collectionID string) ([]string, error) {
	key := fmt.Sprintf("collection:%s:cards", collectionID)
	return c.client.SMembers(ctx, key).Result()
}

// SetCardState stores the minimal state info for a user-card pair.
func (c *Client) SetCardState(ctx context.Context, userID, cardID string, nextReviewAt time.Time, stateID string, priorityScore float64) error {
	key := fmt.Sprintf("state:user:%s:card:%s", userID, cardID)
	return c.client.HSet(ctx, key, map[string]any{
		"next_review_at": nextReviewAt.Format(time.RFC3339),
		"state_id":       stateID,
		"priority_score": priorityScore,
	}).Err()
}

// GetCardState retrieves the minimal state info for a user-card pair.
func (c *Client) GetCardState(ctx context.Context, userID, cardID string) (map[string]string, error) {
	key := fmt.Sprintf("state:user:%s:card:%s", userID, cardID)
	return c.client.HGetAll(ctx, key).Result()
}

// UpdatePriorityScore updates the priority score for a user-card in a collection's sorted set.
func (c *Client) UpdatePriorityScore(ctx context.Context, userID, collectionID, cardID string, score float64) error {
	key := fmt.Sprintf("state:scores:user:%s:collection:%s", userID, collectionID)
	return c.client.ZAdd(ctx, key, redis.Z{Score: score, Member: cardID}).Err()
}

// GetNextCards returns the top N cards by priority score for a user in a collection.
func (c *Client) GetNextCards(ctx context.Context, userID, collectionID string, limit int64) ([]redis.Z, error) {
	key := fmt.Sprintf("state:scores:user:%s:collection:%s", userID, collectionID)
	return c.client.ZRevRangeWithScores(ctx, key, 0, limit-1).Result()
}

// RemovePriorityScore removes a card from the priority sorted set.
func (c *Client) RemovePriorityScore(ctx context.Context, userID, collectionID, cardID string) error {
	key := fmt.Sprintf("state:scores:user:%s:collection:%s", userID, collectionID)
	return c.client.ZRem(ctx, key, cardID).Err()
}

// Pipeline returns a Redis pipeline for batch operations.
func (c *Client) Pipeline() redis.Pipeliner {
	return c.client.Pipeline()
}
