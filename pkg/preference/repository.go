package preference

import (
	"context"
	"time"

	b "github.com/forumGamers/octo-cats/pkg/base"
	"go.mongodb.org/mongo-driver/bson"
)

func NewPreferenceRepo() PreferenceRepo {
	return &PreferenceRepoImpl{b.NewBaseRepo(b.GetCollection(b.Like))}
}

func (r *PreferenceRepoImpl) Create(ctx context.Context, userId string) (UserPreference, error) {
	data := UserPreference{
		UserId:    userId,
		Tags:      []TagPreference{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if result, err := r.BaseRepo.Create(ctx, data); err != nil {
		return data, err
	} else {
		data.Id = result
	}
	return data, nil
}

func (r *PreferenceRepoImpl) FindByUserId(ctx context.Context, userId string) (data UserPreference, err error) {
	err = r.FindOneByQuery(ctx, bson.M{"userId": userId}, &data)
	return
}

func (r *PreferenceRepoImpl) UpdateTags(ctx context.Context, userId string, tags []TagPreference) error {
	_, err := r.UpdateOne(ctx, bson.M{"userId": userId}, bson.M{
		"$set": bson.M{
			"tags": tags,
		},
	})
	return err
}
