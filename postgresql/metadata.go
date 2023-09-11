package postgresql

import (
	"github.com/google/uuid"
	"time"
)

// Metadata is a common struct that is embedded in most data models.
type Metadata struct {
	// ID is a unique identifier for the current record.
	ID uuid.UUID `json:"id" bun:"id,pk,type:uuid"`
	// CreatedAt is the date at which the current record was created.
	CreatedAt time.Time `json:"createdAt,omitempty" bun:"created_at,notnull"`
	// UpdatedAt is the date at which the current record was last updated. It can be empty on creation.
	UpdatedAt *time.Time `json:"updatedAt,omitempty" bun:"updated_at"`
}

func NewMetadata(id uuid.UUID, createdAt time.Time, updatedAt *time.Time) Metadata {
	return Metadata{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
