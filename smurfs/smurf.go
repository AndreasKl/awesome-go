package smurfs

import "github.com/google/uuid"

type Name string

type Smurf struct {
	ID     uuid.UUID
	Name   Name
	Height uint
}
