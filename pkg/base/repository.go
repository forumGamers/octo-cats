package base

import (
	"context"

	cfg "github.com/forumGamers/octo-cats/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewBaseRepo(db *mongo.Collection) BaseRepo {
	return &BaseRepoImpl{db}
}

func (r *BaseRepoImpl) DeleteManyByQuery(ctx context.Context, filter any) error {
	_, err := r.DB.DeleteMany(ctx, filter)
	return err
}

func (r *BaseRepoImpl) DeleteOneById(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.DB.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *BaseRepoImpl) DeleteOneByQuery(ctx context.Context, query any) error {
	_, err := r.DB.DeleteOne(ctx, query)
	return err
}

func (r *BaseRepoImpl) FindOneById(ctx context.Context, id primitive.ObjectID, data any) error {
	return r.DB.FindOne(ctx, bson.M{"_id": id}).Decode(data)
}

func (r *BaseRepoImpl) InsertMany(ctx context.Context, data []any) (*mongo.InsertManyResult, error) {
	return r.DB.InsertMany(ctx, data)
}

func (r *BaseRepoImpl) Create(ctx context.Context, data any) (primitive.ObjectID, error) {
	result, err := r.DB.InsertOne(ctx, data)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func GetCollection(name CollectionName) *mongo.Collection {
	return cfg.DB.Collection(string(name))
}

func (r *BaseRepoImpl) FindOneByQuery(ctx context.Context, query any, result any) error {
	return r.DB.FindOne(ctx, query).Decode(result)
}

func (r *BaseRepoImpl) UpdateOneByQuery(ctx context.Context, id primitive.ObjectID, query any) (*mongo.UpdateResult, error) {
	return r.DB.UpdateByID(ctx, id, query)
}

func (r *BaseRepoImpl) UpdateOne(ctx context.Context, filter, update any) (*mongo.UpdateResult, error) {
	return r.DB.UpdateOne(ctx, filter, update)
}

func (r *BaseRepoImpl) FindByQuery(ctx context.Context, query any) (*mongo.Cursor, error) {
	return r.DB.Find(ctx, query)
}

func (r *BaseRepoImpl) GetSession() (mongo.Session, error) {
	return r.DB.Database().Client().StartSession()
}
