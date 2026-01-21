package interfaces

import (
	"context"
)

type IStorage interface {
	Save(ctx context.Context, userID int64, username, text string) (string, error)
	Count() int
}
