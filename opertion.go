package mongo_orm

import (
	"context"
	MgoError "github.com/z9905080/mongo-orm/error"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

func (m *MgoDB) First(ctx context.Context, filter interface{}, result interface{}) *MgoDB {

	if result == nil {
		m.mError = MgoError.ErrReflectNil
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
		m.mError = MgoError.ErrReflectNil
		return m
	}

	if !m.check() {
		return m
	}

	if reflect.TypeOf(result).Elem().Kind() != reflect.Slice {
		m.mError = MgoError.ErrReflectNonSlice
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

func (m *MgoDB) Insert(ctx context.Context, result interface{}, opts ...*options.InsertOneOptions) *MgoDB {

	if result == nil {
		m.mError = MgoError.ErrReflectNil
		return m
	}

	if !m.check() {
		return m
	}

	rCollection := m.GetClient().Database(m.database).Collection(m.collection)

	insertData, err := rCollection.InsertOne(ctx, result, opts...)
	if err != nil {
		m.mError = err
		return m
	}

	m.lastInsertID = insertData.InsertedID

	return m
}

func (m *MgoDB) Inserts(ctx context.Context, result []interface{}, opts ...*options.InsertManyOptions) *MgoDB {

	if result == nil {
		m.mError = MgoError.ErrReflectNil
		return m
	}

	if !m.check() {
		return m
	}

	rCollection := m.GetClient().Database(m.database).Collection(m.collection)

	insertData, err := rCollection.InsertMany(ctx, result, opts...)
	if err != nil {
		m.mError = err
		return m
	}

	m.lastInsertIDs = insertData.InsertedIDs

	return m
}
