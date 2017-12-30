package models

import (
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	uuid "github.com/satori/go.uuid"
)

type Comment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Content   string    `json:"content" db:"content"`
	AuthorID  uuid.UUID `json:"author_id" db:"author_id"`
	PostID    uuid.UUID `json:"post_id" db:"post_id"`
	Author    User      `json:"-" db:"-"`
}

type Comments []Comment

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (c *Comment) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: c.Content, Name: "Content"},
	), nil
}
