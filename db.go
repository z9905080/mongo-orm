package mongo_orm

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MgoDB struct {
	isClone       bool
	lastInsertID  interface{}
	lastInsertIDs []interface{}
	mError        error
	client        *mongo.Client
	ctx           context.Context
	uri           string
	database      string
	collection    string
}

func New(opts ...DBOption) *MgoDB {

	// init setting option
	opt := DBOptions{}
	for _, o := range opts {
		o(&opt)
	}

	return &MgoDB{
		ctx:        opt.ctx,
		uri:        opt.uri,
		database:   opt.database,
		collection: opt.collection,
	}
}

func (m *MgoDB) ConnectDB(opt ...*options.ClientOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var err error

	var optList []*options.ClientOptions
	optList = append(optList, options.Client().ApplyURI(m.uri))
	optList = append(optList, opt...)

	m.client, err = mongo.Connect(ctx, optList...)
	if err != nil {
		return err
	}

	return nil
}

func (m *MgoDB) SetCollection(collection string) *MgoDB {
	m.collection = collection
	return m.clone()
}

func (m *MgoDB) GetClient() *mongo.Client {
	return m.clone().client
}

func (m *MgoDB) DB() *mongo.Database {
	return m.GetClient().Database(m.database)
}

func (m *MgoDB) Error() error {
	return m.mError
}

func (m *MgoDB) check() bool {
	if m.mError != nil {
		return false
	}

	if m.database == "" {
		return false
	}

	if m.collection == "" {
		return false
	}
	return true
}

// clone Clone database instance
func (m *MgoDB) clone() *MgoDB {
	if m.isClone {
		return m
	}
	db := &MgoDB{
		mError:     nil,
		isClone:    true,
		client:     m.client,
		ctx:        m.ctx,
		uri:        m.uri,
		database:   m.database,
		collection: m.collection,
	}
	return db
}
