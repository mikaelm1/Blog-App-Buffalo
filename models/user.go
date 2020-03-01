package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              uuid.UUID `json:"id" db:"id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	Username        string    `json:"username" db:"username"`
	Email           string    `json:"email" db:"email"`
	Admin           bool      `json:"admin" db:"admin"`
	PasswordHash    string    `json:"-" db:"password_hash"`
	Password        string    `json:"-" db:"-"`
	PasswordConfirm string    `json:"-" db:"-"`
}

// Create validates and creates a new User.
func (u *User) Create(tx *pop.Connection) (*validate.Errors, error) {
	u.Email = strings.ToLower(u.Email)
	u.Admin = false
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return validate.NewErrors(), errors.WithStack(err)
	}
	u.PasswordHash = string(pwdHash)
	return tx.ValidateAndCreate(u)
}

// Authorize checks user's password for logging in
func (u *User) Authorize(tx *pop.Connection) error {
	err := tx.Where("email = ?", strings.ToLower(u.Email)).First(u)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			// couldn't find an user with that email address
			return errors.New("User not found.")
		}
		return errors.WithStack(err)
	}
	// confirm that the given password matches the hashed password from the db
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return errors.New("Invalid password.")
	}
	return nil
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.EmailIsPresent{Name: "Email", Field: u.Email},
		&validators.StringIsPresent{Field: u.Username, Name: "Username"},
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&validators.StringsMatch{Name: "Password", Field: u.Password, Field2: u.PasswordConfirm, Message: "Passwords do not match."},
		&UsernameNotTaken{Name: "Username", Field: u.Username, tx: tx},
		&EmailNotTaken{Name: "Email", Field: u.Email, tx: tx},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

type UsernameNotTaken struct {
	Name  string
	Field string
	tx    *pop.Connection
}

func (v *UsernameNotTaken) IsValid(errors *validate.Errors) {
	query := v.tx.Where("username = ?", v.Field)
	queryUser := User{}
	err := query.First(&queryUser)
	if err == nil {
		// found a user with same username
		errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("The username %s is not available.", v.Field))
	}
}

type EmailNotTaken struct {
	Name  string
	Field string
	tx    *pop.Connection
}

// IsValid performs the validation check for unique emails
func (v *EmailNotTaken) IsValid(errors *validate.Errors) {
	query := v.tx.Where("email = ?", v.Field)
	queryUser := User{}
	err := query.First(&queryUser)
	if err == nil {
		// found a user with the same email
		errors.Add(validators.GenerateKey(v.Name), "An account with that email already exists.")
	}
}
