package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (cc *Client) DBListCollectionNames(db string, filter interface{}) ([]string, error) {
	if filter == nil {
		filter = emptyFilter
	}
	return cc.Client.Database(db, cc.dbopts...).ListCollectionNames(nil, filter)
}

func (cc *Client) DBCollection(db string, cl string) *mongo.Collection {
	return cc.Database(db, cc.dbopts...).Collection(cl,cc.clopts...)
}

func (cc *Client) DBCount(db string, cl string, filter interface{}) (ret int64, err error) {

	if filter == nil {
		return cc.Database(db, cc.dbopts...).Collection(cl,cc.clopts...).EstimatedDocumentCount(nil)
	} else {
		return cc.Database(db, cc.dbopts...).Collection(cl,cc.clopts...).CountDocuments(nil, filter)
	}

}

func (cc *Client) DBFindId(db string, cl string, id interface{}, ret interface{}) (not bool, err error) {
	err = cc.Database(db, cc.dbopts...).Collection(cl,cc.clopts...).FindOne(nil, bson.M{"_id": id}).Decode(ret)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) DBFindOne(db string, cl string, filter interface{}, ret interface{}) (not bool, err error) {
	if filter == nil {
		filter = emptyFilter
	}
	err = cc.Database(db, cc.dbopts...).Collection(cl,cc.clopts...).FindOne(nil, filter).Decode(ret)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) DBFindAll(db string, cl string, filter interface{}, ret interface{}, opts ...*options.FindOptions) (err error) {
	if filter == nil {
		filter = emptyFilter
	}
	cur, err := cc.Database(db, cc.dbopts...).Collection(cl,cc.clopts...).Find(nil, filter, opts...)
	if err == nil {
		err = cur.All(nil, ret)
	}
	return
}

func (cc *Client) DBFindWith(db string, cl string, filter interface{}, fn func(cur *mongo.Cursor) error, opts ...*options.FindOptions) (err error) {
	if filter == nil {
		filter = emptyFilter
	}
	cur, err := cc.Database(db, cc.dbopts...).Collection(cl,cc.clopts...).Find(nil, filter, opts...)
	if err == nil {
		defer cur.Close(nil)
		err = fn(cur)
	}
	return
}

func (cc *Client) DBDistinct(db string, cl string, fieldName string, filter interface{}, opts ...*options.DistinctOptions) (ret []interface{}, err error) {
	if filter == nil {
		filter = emptyFilter
	}
	ret, err = cc.Database(db, cc.dbopts...).Collection(cl,cc.clopts...).Distinct(nil, fieldName, filter, opts...)
	return
}
