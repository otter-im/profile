package model

import (
	"context"
	"fmt"
	"github.com/go-redis/cache/v8"
	"github.com/google/uuid"
	"github.com/otter-im/profile/internal/app/datasource"
	"sync"
	"time"
)

var (
	profileDtoOnce     sync.Once
	profileDtoInstance *profileDto
)

type Profile struct {
	UserID      uuid.UUID `json:"user_id"`
	DisplayName *string   `json:"display_name,omitempty"`
	Description *string   `json:"description,omitempty"`
	Location    *string   `json:"location,omitempty"`
	Website     *string   `json:"website,omitempty"`
	Pronouns    *string   `json:"pronouns,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	AvatarURL   string    `json:"avatar_url"`
}

type ProfileService interface {
	SelectProfile(ctx context.Context, id string) (*Profile, error)
}

type profileDto struct{}

func ProfileDTO() *profileDto {
	profileDtoOnce.Do(func() {
		profileDtoInstance = &profileDto{}
	})
	return profileDtoInstance
}

func (p *profileDto) SelectProfile(ctx context.Context, id uuid.UUID) (*Profile, error) {
	user := new(Profile)
	if err := datasource.RedisCache().Once(&cache.Item{
		Ctx:   ctx,
		Key:   fmt.Sprintf("profile:%s", id.String()),
		Value: user,
		TTL:   15 * time.Minute,
		Do: func(item *cache.Item) (interface{}, error) {
			return selectProfile(ctx, id)
		},
	}); err != nil {
		return nil, err
	}
	return user, nil
}

func selectProfile(ctx context.Context, id uuid.UUID) (*Profile, error) {
	profile := new(Profile)
	err := datasource.DB().
		QueryRowContext(ctx, "SELECT * FROM profile WHERE id = $1", id).
		Scan(&profile.UserID, &profile.DisplayName, &profile.Description, &profile.Location, &profile.Website, &profile.Pronouns, &profile.DateOfBirth, &profile.AvatarURL)
	if err != nil {
		return nil, err
	}
	return profile, nil
}
