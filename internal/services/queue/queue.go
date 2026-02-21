// Package queue provides abstractions for pushing background jobs to a message
// queue, such as bookmark import tasks that are processed asynchronously.
package queue

import (
	"context"

	"github.com/luongtruong20201/bookmark-management/internal/repositories/queue"
)

// Service defines the interface for enqueuing background jobs.
// Currently it supports sending bookmark import jobs to the queue.
//
//go:generate mockery --name Service --filename queue.go
type Service interface {
	// SendBookmarkJob splits the provided bookmark slice into smaller batches and
	// enqueues one or more bookmark import jobs for the given user ID.
	// The provided context is used to propagate deadlines and cancellation to the
	// underlying queue repository.
	SendBookmarkJob(ctx context.Context, uid string, bookmarks []*ImportBookmarkInput) error
}

// queueService is a concrete implementation of Service that delegates
// queue operations to a Repository.
type queueService struct {
	repo queue.Repository
}

// NewService creates a new queue service using the provided repository
// implementation. The repository is responsible for persisting messages
// to the underlying queue (e.g., Redis list, message broker, etc.).
func NewService(repo queue.Repository) Service {
	return &queueService{
		repo: repo,
	}
}
