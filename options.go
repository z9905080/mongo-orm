package mongo_orm

import "context"

type DBOptions struct {
	ctx        context.Context
	uri        string
	database   string
	collection string
}

type DBOption func(*DBOptions)

func SetDatabase(database string) DBOption {
	return func(o *DBOptions) {
		o.database = database
	}
}

func SetCollection(collection string) DBOption {
	return func(o *DBOptions) {
		o.collection = collection
	}
}

func SetUri(uri string) DBOption {
	return func(o *DBOptions) {
		o.uri = uri
	}
}

func SetContext(ctx context.Context) DBOption {
	return func(o *DBOptions) {
		o.ctx = ctx
	}
}

