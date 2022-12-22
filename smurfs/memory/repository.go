package memory

import (
	"context"

	"github.com/google/uuid"

	"awesome/apperrors"
	"awesome/smurfs"
)

var _ smurfs.Repository = (*Repository)(nil)

var smurfsStore = []smurfs.Smurf{
	{ID: uuid.New(), Name: "Schlaubi Schlumpf", Height: 42},
	{ID: uuid.New(), Name: "Papa Schlumpf", Height: 23},
}

type Repository struct{}

func (r *Repository) List(_ context.Context) []smurfs.Smurf {
	return smurfsStore
}

func (r *Repository) Get(_ context.Context, id uuid.UUID) (smurfs.Smurf, error) {
	for _, smurf := range smurfsStore {
		if smurf.ID == id {
			return smurf, nil
		}
	}
	return smurfs.Smurf{}, apperrors.ErrEntityNotFound
}
