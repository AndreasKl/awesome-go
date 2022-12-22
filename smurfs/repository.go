package smurfs

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) []Smurf
	Get(ctx context.Context, id uuid.UUID) (Smurf, error)
}
