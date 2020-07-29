package mongo_orm

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"time"
)

type MgoDB struct {
	isClone    bool
	mError     error
	client     *mongo.Client
	ctx        context.Context
	uri        string
	database   string
	collection string
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

func (m *MgoDB) First(ctx context.Context, filter interface{}, result interface{}) *MgoDB {

	if result == nil {
		m.mError = errors.New("can't reflect nil pointer")
		return m
	}

	if !m.check() {
		return m
	}

	rCollection := m.GetClient().Database(m.database).Collection(m.collection)

	if err := rCollection.FindOne(ctx, filter).Decode(result); err != nil {
		m.mError = err
		return m
	}

	return m
}

func (m *MgoDB) Find(ctx context.Context, filter interface{}, result interface{}) *MgoDB {

	if result == nil {
		m.mError = errors.New("can't reflect nil pointer")
		return m
	}

	if !m.check() {
		return m
	}

	if reflect.TypeOf(result).Elem().Kind() != reflect.Slice {
		m.mError = errors.New("can't reflect not slice object")
		return m
	}

	rCollection := m.GetClient().Database(m.database).Collection(m.collection)

	cursor, err := rCollection.Find(ctx, filter)
	if err != nil {
		m.mError = err
		return m
	}

	getErr := cursor.All(ctx, result)
	if getErr != nil {
		m.mError = getErr
		return m
	}
	return m
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
