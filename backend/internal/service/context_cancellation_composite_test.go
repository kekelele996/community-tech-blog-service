package service

import (
	"context"
	"errors"
	"testing"
)

func TestContextCancellationPropagation(t *testing.T) {
	e := newTestEnv(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := e.articleRepo.FindByID(ctx, 9999); !errors.Is(err, context.Canceled) {
		t.Fatalf("article FindByID with canceled ctx err = %v, want context.Canceled", err)
	}
	if _, err := e.followRepo.AllFollowingIDs(ctx, 9999); !errors.Is(err, context.Canceled) {
		t.Fatalf("AllFollowingIDs with canceled ctx err = %v, want context.Canceled", err)
	}
	if _, _, err := e.notificationRepo.ListByUser(ctx, 9999, 1, 10, false); !errors.Is(err, context.Canceled) {
		t.Fatalf("notification ListByUser with canceled ctx err = %v, want context.Canceled", err)
	}
	if _, err := e.userRepo.FindByID(ctx, 9999); !errors.Is(err, context.Canceled) {
		t.Fatalf("user FindByID with canceled ctx err = %v, want context.Canceled", err)
	}
}
