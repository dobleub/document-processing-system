package interfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IMongoCollection interface {
	SetCollectionName(name string)
	SetDBContext(ctx context.Context)
	Clone(opts ...*options.CollectionOptions) (*mongo.Collection, error)
	Name(string) string
	Database(*mongo.Database) *mongo.Database
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	UpdateByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	EstimatedDocumentCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error)
	Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult
	FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
	Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error)
	GetIndexes(ctx context.Context) ([]mongo.IndexModel, error)
	CreateIndex(ctx context.Context, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error)
	DropIndex(ctx context.Context, name string, opts ...*options.DropIndexesOptions) error
	CreateManyIndexes(ctx context.Context, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error)
	DropManyIndexes(ctx context.Context, names []string, opts ...*options.DropIndexesOptions) ([]string, error)
	Indexes() mongo.IndexView
	Drop(ctx context.Context) error
}

type MongoCollection struct {
	CollectionName string
	DB             *mongo.Database
}

func (c *MongoCollection) SetCollectionName(name string) {
	c.CollectionName = name
}

func (c *MongoCollection) SetDBContext(ctx context.Context) {
	c.DB = ctx.Value(MongodbKey).(MongoDBContext).DB
}

func (c *MongoCollection) Clone(opts ...*options.CollectionOptions) (*mongo.Collection, error) {
	return c.DB.Collection(c.CollectionName).Clone(opts...)
}

func (c *MongoCollection) Name() string {
	return c.DB.Collection(c.CollectionName).Name()
}

func (c *MongoCollection) Database() *mongo.Database {
	return c.DB.Collection(c.CollectionName).Database()
}

func (c *MongoCollection) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	return c.DB.Collection(c.CollectionName).BulkWrite(ctx, models, opts...)
}

func (c *MongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return c.DB.Collection(c.CollectionName).InsertOne(ctx, document, opts...)
}

func (c *MongoCollection) InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return c.DB.Collection(c.CollectionName).InsertMany(ctx, documents, opts...)
}

func (c *MongoCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return c.DB.Collection(c.CollectionName).DeleteOne(ctx, filter, opts...)
}

func (c *MongoCollection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return c.DB.Collection(c.CollectionName).DeleteMany(ctx, filter, opts...)
}

func (c *MongoCollection) UpdateByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.DB.Collection(c.CollectionName).UpdateByID(ctx, id, update, opts...)
}

func (c *MongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.DB.Collection(c.CollectionName).UpdateOne(ctx, filter, update, opts...)
}

func (c *MongoCollection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.DB.Collection(c.CollectionName).UpdateMany(ctx, filter, update, opts...)
}

func (c *MongoCollection) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	return c.DB.Collection(c.CollectionName).ReplaceOne(ctx, filter, replacement, opts...)
}

func (c *MongoCollection) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	return c.DB.Collection(c.CollectionName).Aggregate(ctx, pipeline, opts...)
}

func (c *MongoCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return c.DB.Collection(c.CollectionName).CountDocuments(ctx, filter, opts...)
}

func (c *MongoCollection) EstimatedDocumentCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	return c.DB.Collection(c.CollectionName).EstimatedDocumentCount(ctx, opts...)
}

func (c *MongoCollection) Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	return c.DB.Collection(c.CollectionName).Distinct(ctx, fieldName, filter, opts...)
}

func (c *MongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, int32, error) {
	cursor, err := c.DB.Collection(c.CollectionName).Find(ctx, filter, opts...)
	if err != nil {
		return nil, 0, err
	}
	count, err := c.DB.Collection(c.CollectionName).CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return cursor, int32(count), nil
}

func (c *MongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return c.DB.Collection(c.CollectionName).FindOne(ctx, filter, opts...)
}

func (c *MongoCollection) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	return c.DB.Collection(c.CollectionName).FindOneAndDelete(ctx, filter, opts...)
}

func (c *MongoCollection) FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	return c.DB.Collection(c.CollectionName).FindOneAndReplace(ctx, filter, replacement, opts...)
}

func (c *MongoCollection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	return c.DB.Collection(c.CollectionName).FindOneAndUpdate(ctx, filter, update, opts...)
}

func (c *MongoCollection) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	return c.DB.Collection(c.CollectionName).Watch(ctx, pipeline, opts...)
}

func (c *MongoCollection) GetIndexes(ctx context.Context) ([]mongo.IndexModel, error) {
	cursor, err := c.DB.Collection(c.CollectionName).Indexes().List(ctx)
	if err != nil {
		return nil, err
	}

	indexList := []mongo.IndexModel{}
	cursor.All(ctx, &indexList)

	return indexList, nil
}

func (c *MongoCollection) CreateIndex(ctx context.Context, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error) {
	res, err := c.DB.Collection(c.CollectionName).Indexes().CreateOne(ctx, model, opts...)
	return res, err
}

func (c *MongoCollection) DropIndex(ctx context.Context, name string, opts ...*options.DropIndexesOptions) error {
	_, err := c.DB.Collection(c.CollectionName).Indexes().DropOne(ctx, name, opts...)
	return err
}

func (c *MongoCollection) CreateManyIndexes(ctx context.Context, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	res, err := c.DB.Collection(c.CollectionName).Indexes().CreateMany(ctx, models, opts...)
	return res, err
}

func (c *MongoCollection) DropManyIndexes(ctx context.Context, names []string, opts ...*options.DropIndexesOptions) ([]string, error) {
	res := []string{}

	for _, name := range names {
		tmpRes, tmpErr := c.DB.Collection(c.CollectionName).Indexes().DropOne(ctx, name, opts...)
		if tmpErr != nil {
			return res, tmpErr
		}
		res = append(res, tmpRes.String())
	}

	return res, nil
}

func (c *MongoCollection) Indexes() mongo.IndexView {
	return c.DB.Collection(c.CollectionName).Indexes()
}

func (c *MongoCollection) Drop(ctx context.Context) error {
	return c.DB.Collection(c.CollectionName).Drop(ctx)
}
