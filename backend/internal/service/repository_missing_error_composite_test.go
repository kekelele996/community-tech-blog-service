package service

import (
	"context"
	"errors"
	"testing"

	"github.com/techblog/community/internal/repository"
)

func TestRepositoryMissingRecordErrorsComposite(t *testing.T) {
	e := newTestEnv(t)
	ctx := context.Background()

	if err := e.articleRepo.UpdateStatus(ctx, 9999, 1); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("update missing article status err = %v, want ErrNotFound", err)
	}
	if err := e.followRepo.Delete(ctx, 9999, 9998); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("delete missing follow err = %v, want ErrNotFound", err)
	}
	if err := e.notificationRepo.MarkRead(ctx, 9999, 9998); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("mark missing notification read err = %v, want ErrNotFound", err)
	}
	if err := e.collectionRepo.RemoveArticle(ctx, 9999, 9998); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("remove missing collection article err = %v, want ErrNotFound", err)
	}
}
