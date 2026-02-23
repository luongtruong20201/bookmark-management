package worker_test

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	bookmarkHandler "github.com/luongtruong20201/bookmark-management/internal/handlers/bookmark"
	model "github.com/luongtruong20201/bookmark-management/internal/models"
	bookmarkRepo "github.com/luongtruong20201/bookmark-management/internal/repositories/bookmark"
	queueRepo "github.com/luongtruong20201/bookmark-management/internal/repositories/queue"
	bookmarkService "github.com/luongtruong20201/bookmark-management/internal/services/bookmark"
	queueService "github.com/luongtruong20201/bookmark-management/internal/services/queue"
	"github.com/luongtruong20201/bookmark-management/internal/test/fixture"
	"github.com/luongtruong20201/bookmark-management/internal/worker"
	"github.com/luongtruong20201/bookmark-management/pkg/redis"
	"github.com/luongtruong20201/bookmark-management/pkg/stringutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const (
	testQueueName = "test-bookmark-import"
	numWorkers    = 2
	importUserID  = "550e8400-e29b-41d4-a716-446655440000"
)

func makeTestWorkerEngine(t *testing.T, ctx context.Context, importMessages []*queueService.ImportMessage) (worker.Engine, *gorm.DB) {
	t.Helper()

	db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})

	redisClient := redis.InitMockRedis(t)
	repo := queueRepo.NewRedisQueue(redisClient, testQueueName)

	for _, msg := range importMessages {
		msgBytes, err := json.Marshal(msg)
		require.NoError(t, err)
		err = repo.PushMessage(ctx, msgBytes)
		require.NoError(t, err)
	}

	bookmarkRepository := bookmarkRepo.NewBookmark(db)
	keyGen := stringutils.NewKeyGen()
	bookmarkSvc := bookmarkService.NewBookmarkSvc(bookmarkRepository, keyGen)
	handler := bookmarkHandler.NewWorkerHandler(bookmarkSvc)

	engine := worker.NewEngine(repo, handler, numWorkers)
	return engine, db
}

func TestWorker_ProcessesImportMessagesAndWritesToDB(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	importMessages := []*queueService.ImportMessage{
		{
			UID: importUserID,
			Bookmarks: []*queueService.ImportBookmarkInput{
				{Description: "Test Import One", URL: "https://example.com/one"},
				{Description: "Test Import Two", URL: "https://example.com/two"},
			},
		},
	}

	engine, db := makeTestWorkerEngine(t, ctx, importMessages)

	var countBefore int64
	err := db.Model(&model.Bookmark{}).Where("user_id = ?", importUserID).Count(&countBefore).Error
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		engine.Start(ctx)
	}()

	time.Sleep(2 * time.Second)
	cancel()
	wg.Wait()

	var countAfter int64
	err = db.Model(&model.Bookmark{}).Where("user_id = ?", importUserID).Count(&countAfter).Error
	require.NoError(t, err)

	added := int(countAfter - countBefore)
	assert.Equal(t, 2, added, "expected 2 new bookmarks from import message")
}
