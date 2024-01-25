package preference

import (
	"context"

	b "github.com/forumGamers/octo-cats/pkg/base"
)

type PreferenceRepo interface {
	Create(ctx context.Context, userId string) (UserPreference, error)
	FindByUserId(ctx context.Context, userId string) (UserPreference, error)
	UpdateTags(ctx context.Context, userId string, tags []TagPreference) error
}

type PreferenceRepoImpl struct{ b.BaseRepo }

type PreferenceService interface {
	CreateUserNewTags(ctx context.Context, data UserPreference, newData []string) []TagPreference
}

type PreferenceServiceImpl struct{ Repo PreferenceRepo }
