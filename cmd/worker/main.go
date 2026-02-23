// Package main runs the bookmark import worker: it connects to Redis and
// the database, starts an engine that pops messages from the "bookmark_import"
// queue, and processes each message with the bookmark worker handler (batch
// insert). Run alongside the API so that import jobs enqueued by the API
// are consumed here.
package main

import (
	"context"

	bookmarkHandler "github.com/luongtruong20201/bookmark-management/internal/handlers/bookmark"
	"github.com/luongtruong20201/bookmark-management/internal/infrastructure"
	"github.com/luongtruong20201/bookmark-management/internal/repositories/bookmark"
	"github.com/luongtruong20201/bookmark-management/internal/repositories/queue"
	queueRepo "github.com/luongtruong20201/bookmark-management/internal/repositories/queue"
	bookmarkService "github.com/luongtruong20201/bookmark-management/internal/services/bookmark"
	"github.com/luongtruong20201/bookmark-management/internal/worker"
	"github.com/luongtruong20201/bookmark-management/pkg/stringutils"
)

const (
	// numWorkers represent for number of worker in worker pool
	numWorkers = 5
)

// main wires Redis, DB, bookmark repository/service, and the worker handler,
// then runs the worker engine until the process exits.
func main() {
	redis := infrastructure.CreateRedis()
	db := infrastructure.CreateSqlDB()

	bookmarkRepo := bookmark.NewBookmark(db)
	keyGen := stringutils.NewKeyGen()
	bookmarkSvc := bookmarkService.NewBookmarkSvc(bookmarkRepo, keyGen)
	handler := bookmarkHandler.NewWorkerHandler(bookmarkSvc)

	repo := queueRepo.NewRedisQueue(redis, queue.QueueNameBookmarkImport)
	engine := worker.NewEngine(repo, handler, numWorkers)

	engine.Start(context.Background())
}
