package dataloaders

import (
	"context"
	"fmt"
	"sync"
	"time"

	"learning/internal/user"
)

// UserLoader provides batched loading of users
type UserLoader struct {
	repo  user.Repository
	cache map[string]*user.User
	mutex sync.RWMutex

	// Batching fields
	batch      map[string][]chan userResult
	batchMutex sync.Mutex
	maxBatch   int
	wait       time.Duration
}

type userResult struct {
	user *user.User
	err  error
}

// NewUserLoader creates a new UserLoader
func NewUserLoader(repo user.Repository) *UserLoader {
	return &UserLoader{
		repo:     repo,
		cache:    make(map[string]*user.User),
		batch:    make(map[string][]chan userResult),
		maxBatch: 100,
		wait:     1 * time.Millisecond,
	}
}

// Load loads a single user by ID with batching
func (ul *UserLoader) Load(ctx context.Context, id string) (*user.User, error) {
	// Check cache first
	ul.mutex.RLock()
	if cached, exists := ul.cache[id]; exists {
		ul.mutex.RUnlock()
		return cached, nil
	}
	ul.mutex.RUnlock()

	// Create result channel
	result := make(chan userResult, 1)

	// Add to batch
	ul.batchMutex.Lock()
	ul.batch[id] = append(ul.batch[id], result)
	batch := ul.batch[id]

	// If this is the first request for this ID, start the batch timer
	if len(batch) == 1 {
		go ul.startBatch(id)
	}
	ul.batchMutex.Unlock()

	// Wait for result
	res := <-result
	return res.user, res.err
}

// LoadMany loads multiple users by IDs
func (ul *UserLoader) LoadMany(ctx context.Context, ids []string) ([]*user.User, []error) {
	users := make([]*user.User, len(ids))
	errors := make([]error, len(ids))

	for i, id := range ids {
		users[i], errors[i] = ul.Load(ctx, id)
	}

	return users, errors
}

// Prime adds a user to the cache
func (ul *UserLoader) Prime(ctx context.Context, id string, user *user.User) {
	ul.mutex.Lock()
	defer ul.mutex.Unlock()
	ul.cache[id] = user
}

// Clear removes a user from the cache
func (ul *UserLoader) Clear(ctx context.Context, id string) {
	ul.mutex.Lock()
	defer ul.mutex.Unlock()
	delete(ul.cache, id)
}

// startBatch starts a batch operation for the given ID
func (ul *UserLoader) startBatch(id string) {
	time.Sleep(ul.wait)

	ul.batchMutex.Lock()
	batch := ul.batch[id]
	delete(ul.batch, id)
	ul.batchMutex.Unlock()

	// Fetch the user
	user, err := ul.repo.GetByID(context.Background(), id)

	// Cache successful results
	if err == nil {
		ul.Prime(context.Background(), id, user)
	}

	// Send results to all waiters
	result := userResult{user: user, err: err}
	for _, ch := range batch {
		ch <- result
		close(ch)
	}
}

// batchGetUsers fetches multiple users in a single operation
func batchGetUsers(ctx context.Context, repo user.Repository, ids []string) ([]*user.User, []error) {
	// Map to store results
	userMap := make(map[string]*user.User)
	errors := make([]error, len(ids))

	// Fetch each user (in a real implementation, you might have a BatchGetByIDs method)
	for _, id := range ids {
		user, err := repo.GetByID(ctx, id)
		if err != nil {
			// Store error but continue processing other IDs
			userMap[id] = nil
		} else {
			userMap[id] = user
		}
	}

	// Build result slice in the same order as input IDs
	results := make([]*user.User, len(ids))
	for i, id := range ids {
		results[i] = userMap[id]
		if results[i] == nil && errors[i] == nil {
			errors[i] = fmt.Errorf("user not found: %s", id)
		}
	}

	return results, errors
}
