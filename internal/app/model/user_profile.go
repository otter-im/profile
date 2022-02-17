package model

import (
	"github.com/google/uuid"
	"time"
)

type UserProfile struct {
	tableName struct{} `json:"-" pg:"user_profiles"`

	UserID      uuid.UUID `json:"user_id" pg:"user_id"`
	Username    string    `json:"username" pg:"username"`
	DisplayName string    `json:"display_name" pg:"display_name"`
	Bio         string    `json:"bio" pg:"bio"`
	Location    string    `json:"location" pg:"location"`
	Website     string    `json:"website" pg:"website"`
	Pronouns    string    `json:"pronouns" pg:"pronouns"`
	Birthday    time.Time `json:"date_of_birth" pg:"date_of_birth"`
	Avatar      []byte    `json:"avatar" pg:"avatar"`
}
