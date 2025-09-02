package nosql_document

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log/slog"
)

type MongoDocument struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoDocument(url, databaseName string, dcOptions ...Option) (*MongoDocument, error) {
	if url == "" {
		return nil, errors.New("[MongoDocument] - url should not be empty")
	}

	if databaseName == "" {
		return nil, errors.New("[MongoDocument] - database should not be empty")
	}

	var opts documentOptions
	for _, option := range dcOptions {
		option(&opts)
	}

	client, err := mongo.Connect(
		options.Client().
			ApplyURI(url).
			SetTimeout(opts.timeout).
			SetAppName(opts.appName),
	)
	if err != nil {
		return nil, fmt.Errorf("[MongoDocument] - failed to connect to mongo: %w", err)
	}

	return &MongoDocument{
		client:   client,
		database: client.Database(databaseName),
	}, nil
}

func (document *MongoDocument) InsertDocument(collection string, data any) (*WriteResponse, error) {
	coll := document.database.Collection(collection)

	result, err := coll.InsertOne(context.TODO(), data)
	if err != nil {
		return nil, fmt.Errorf("[MongoDocument] - failed to insert document: %w", err)
	}

	return newWriteResponse(result.InsertedID, result.Acknowledged), nil
}

func (document *MongoDocument) FindDocument(collection, field string, matchValue any, decode any) error {
	coll := document.database.Collection(collection)
	filter := bson.D{{field, matchValue}}
	ctx := context.TODO()

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return fmt.Errorf("[MongoDocument] - failed on find document: %w", err)
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			slog.Error("[MongoDocument] - failed to close cursor: %s", err.Error())
		}
	}(cursor, ctx)

	if err = cursor.All(ctx, decode); err != nil {
		return fmt.Errorf("[MongoDocument] - failed on iterates the cursor and decodes each document into results: %w", err)
	}
	return nil
}
