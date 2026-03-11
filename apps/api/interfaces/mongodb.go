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
	BulkWrite(models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)
	InsertOne(document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	UpdateByID(id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	ReplaceOne(filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
	Aggregate(pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
	CountDocuments(filter interface{}, opts ...*options.CountOptions) (int64, error)
	EstimatedDocumentCount(opts ...*options.EstimatedDocumentCountOptions) (int64, error)
	Distinct(fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error)
	Find(filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	FindOneAndDelete(filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult
	FindOneAndReplace(filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult
	FindOneAndUpdate(filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
	Watch(pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error)
	GetIndexes() ([]mongo.IndexModel, error)
	CreateIndex(model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error)
	DropIndex(name string, opts ...*options.DropIndexesOptions) error
	CreateManyIndexes(models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error)
	DropManyIndexes(names []string, opts ...*options.DropIndexesOptions) ([]string, error)
	Indexes() mongo.IndexView
	Drop() error
}

type MongoCollection struct {
	CollectionName string
	DB             *mongo.Database
	ctx            context.Context
}

func (c *MongoCollection) SetCollectionName(name string) {
	c.CollectionName = name
}

func (c *MongoCollection) SetDBContext(ctx context.Context) {
	c.ctx = ctx
	c.DB = c.ctx.Value(MongodbKey).(*MongoDBContext).DB
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

func (c *MongoCollection) BulkWrite(models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	return c.DB.Collection(c.CollectionName).BulkWrite(c.ctx, models, opts...)
}

func (c *MongoCollection) InsertOne(document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return c.DB.Collection(c.CollectionName).InsertOne(c.ctx, document, opts...)
}

func (c *MongoCollection) InsertMany(documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return c.DB.Collection(c.CollectionName).InsertMany(c.ctx, documents, opts...)
}

func (c *MongoCollection) DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return c.DB.Collection(c.CollectionName).DeleteOne(c.ctx, filter, opts...)
}

func (c *MongoCollection) DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return c.DB.Collection(c.CollectionName).DeleteMany(c.ctx, filter, opts...)
}

func (c *MongoCollection) UpdateByID(id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.DB.Collection(c.CollectionName).UpdateByID(c.ctx, id, update, opts...)
}

func (c *MongoCollection) UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.DB.Collection(c.CollectionName).UpdateOne(c.ctx, filter, update, opts...)
}

func (c *MongoCollection) UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.DB.Collection(c.CollectionName).UpdateMany(c.ctx, filter, update, opts...)
}

func (c *MongoCollection) ReplaceOne(filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	return c.DB.Collection(c.CollectionName).ReplaceOne(c.ctx, filter, replacement, opts...)
}

func (c *MongoCollection) Aggregate(pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	return c.DB.Collection(c.CollectionName).Aggregate(c.ctx, pipeline, opts...)
}

func (c *MongoCollection) CountDocuments(filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return c.DB.Collection(c.CollectionName).CountDocuments(c.ctx, filter, opts...)
}

func (c *MongoCollection) EstimatedDocumentCount(opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	return c.DB.Collection(c.CollectionName).EstimatedDocumentCount(c.ctx, opts...)
}

func (c *MongoCollection) Distinct(fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	return c.DB.Collection(c.CollectionName).Distinct(c.ctx, fieldName, filter, opts...)
}

func (c *MongoCollection) Find(filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, int32, error) {
	cursor, err := c.DB.Collection(c.CollectionName).Find(c.ctx, filter, opts...)
	if err != nil {
		return nil, 0, err
	}
	count, err := c.DB.Collection(c.CollectionName).CountDocuments(c.ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return cursor, int32(count), nil
}

func (c *MongoCollection) FindOne(filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return c.DB.Collection(c.CollectionName).FindOne(c.ctx, filter, opts...)
}

func (c *MongoCollection) FindOneAndDelete(filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	return c.DB.Collection(c.CollectionName).FindOneAndDelete(c.ctx, filter, opts...)
}

func (c *MongoCollection) FindOneAndReplace(filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	return c.DB.Collection(c.CollectionName).FindOneAndReplace(c.ctx, filter, replacement, opts...)
}

func (c *MongoCollection) FindOneAndUpdate(filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	return c.DB.Collection(c.CollectionName).FindOneAndUpdate(c.ctx, filter, update, opts...)
}

func (c *MongoCollection) Watch(pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	return c.DB.Collection(c.CollectionName).Watch(c.ctx, pipeline, opts...)
}

func (c *MongoCollection) GetIndexes() ([]mongo.IndexModel, error) {
	cursor, err := c.DB.Collection(c.CollectionName).Indexes().List(c.ctx)
	if err != nil {
		return nil, err
	}

	indexList := []mongo.IndexModel{}
	cursor.All(c.ctx, &indexList)

	return indexList, nil
}

func (c *MongoCollection) CreateIndex(model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error) {
	res, err := c.DB.Collection(c.CollectionName).Indexes().CreateOne(c.ctx, model, opts...)
	return res, err
}

func (c *MongoCollection) DropIndex(name string, opts ...*options.DropIndexesOptions) error {
	_, err := c.DB.Collection(c.CollectionName).Indexes().DropOne(c.ctx, name, opts...)
	return err
}

func (c *MongoCollection) CreateManyIndexes(models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	res, err := c.DB.Collection(c.CollectionName).Indexes().CreateMany(c.ctx, models, opts...)
	return res, err
}

func (c *MongoCollection) DropManyIndexes(names []string, opts ...*options.DropIndexesOptions) ([]string, error) {
	res := []string{}

	for _, name := range names {
		tmpRes, tmpErr := c.DB.Collection(c.CollectionName).Indexes().DropOne(c.ctx, name, opts...)
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

func (c *MongoCollection) Drop() error {
	return c.DB.Collection(c.CollectionName).Drop(c.ctx)
}
