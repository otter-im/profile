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
	tableName struct{} `json:"-" pg:"profile"`

	UserID      uuid.UUID `json:"user_id" pg:"id"`
	DisplayName string    `json:"display_name" pg:"display_name"`
	Description string    `json:"description" pg:"description"`
	Location    string    `json:"location" pg:"location"`
	Website     string    `json:"website" pg:"website"`
	Pronouns    string    `json:"pronouns" pg:"pronouns"`
	DateOfBirth time.Time `json:"date_of_birth" pg:"date_of_birth"`
	AvatarURL   string    `json:"avatar_url" pg:"avatar_url"`
}

type ProfileService interface {
	SelectProfile(ctx context.Context, id uuid.UUID) (*Profile, error)
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
	user := new(Profile)
	if err := datasource.Postgres().
		ModelContext(ctx, user).
		Where("id = ?", id).
		Select(); err != nil {
		return nil, err
	}
	return user, nil
}
