package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (cc *Client) DBListCollectionNames(db string, filter interface{}) ([]string, error) {
	return cc.Client.Database(db).ListCollectionNames(nil, filter)
}

func (cc *Client) DBCollection(db string, cl string, opts ...*options.CollectionOptions) *mongo.Collection {
	if cc.collectionOptions != nil {
		opts = append(opts, cc.collectionOptions)
	}
	return cc.Database(db).Collection(cl, opts...)
}

func (cc *Client) DBCount(db string, cl string, filters ...interface{}) (ret int64, err error) {
	if len(filters) == 0 {
		return cc.Database(db).Collection(cl, cc.collectionOptions).EstimatedDocumentCount(nil)
	} else {
		return cc.Database(db).Collection(cl, cc.collectionOptions).CountDocuments(nil, filters[0])
	}

}

func (cc *Client) DBFindId(db string, cl string, id interface{}, ret interface{}, opts ...*options.FindOneOptions) (not bool, err error) {
	err = cc.Database(db).Collection(cl, cc.collectionOptions).FindOne(nil, bson.M{"_id": id}, opts...).Decode(ret)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) DBFindOne(db string, cl string, filter interface{}, ret interface{}, opts ...*options.FindOneOptions) (not bool, err error) {
	err = cc.Database(db).Collection(cl, cc.collectionOptions).FindOne(nil, filter, opts...).Decode(ret)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) DBFind(db string, cl string, filter interface{}, ret interface{}, opts ...*options.FindOptions) (err error) {
	cur, err := cc.Database(db).Collection(cl, cc.collectionOptions).Find(nil, filter, opts...)
	if err == nil {
		err = cur.All(nil, ret)
	}
	return
}

func (cc *Client) DBFindWith(db string, cl string, filter interface{}, with func(cur *mongo.Cursor) error, opts ...*options.FindOptions) (err error) {
	cur, err := cc.Database(db).Collection(cl, cc.collectionOptions).Find(nil, filter, opts...)
	if err == nil {
		defer cur.Close(nil)
		err = with(cur)
	}
	return
}

func (cc *Client) DBDistinct(db string, cl string, fieldName string, filter interface{}, opts ...*options.DistinctOptions) (ret []interface{}, err error) {
	ret, err = cc.Database(db).Collection(cl, cc.collectionOptions).Distinct(nil, fieldName, filter, opts...)
	return
}

func (cc *Client) DBFindIdAndUpdate(db string, cl string, id interface{}, update interface{}, ret interface{}, opts ...*options.FindOneAndUpdateOptions) (not bool, err error) {
	result := cc.Database(db).Collection(cl, cc.collectionOptions).FindOneAndUpdate(nil, bson.M{"_id": id}, update, opts...)
	if ret != nil {
		err = result.Decode(&ret)
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) DBFindIdAndReplace(db string, cl string, id interface{}, replace interface{}, ret interface{}, opts ...*options.FindOneAndReplaceOptions) (not bool, err error) {
	result := cc.Database(db).Collection(cl, cc.collectionOptions).FindOneAndReplace(nil, bson.M{"_id": id}, replace, opts...)
	if ret != nil {
		err = result.Decode(&ret)
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) DBFindIdAndDelete(db string, cl string, id interface{}, ret interface{}, opts ...*options.FindOneAndDeleteOptions) (not bool, err error) {
	result := cc.Database(db).Collection(cl, cc.collectionOptions).FindOneAndDelete(nil, bson.M{"_id": id}, opts...)
	if ret != nil {
		err = result.Decode(&ret)
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) DBFindOneAndUpdate(db string, cl string, filter interface{}, update interface{}, ret interface{}, opts ...*options.FindOneAndUpdateOptions) (not bool, err error) {
	result := cc.Database(db).Collection(cl, cc.collectionOptions).FindOneAndUpdate(nil, filter, update, opts...)
	if ret != nil {
		err = result.Decode(&ret)
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) DBFindOneAndReplace(db string, cl string, filter interface{}, replace interface{}, ret interface{}, opts ...*options.FindOneAndReplaceOptions) (not bool, err error) {
	result := cc.Database(db).Collection(cl, cc.collectionOptions).FindOneAndReplace(nil, filter, replace, opts...)
	if ret != nil {
		err = result.Decode(&ret)
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) DBFindOneAndDelete(db string, cl string, filter interface{}, ret interface{}, opts ...*options.FindOneAndDeleteOptions) (not bool, err error) {
	result := cc.Database(db).Collection(cl, cc.collectionOptions).FindOneAndDelete(nil, filter, opts...)
	if ret != nil {
		err = result.Decode(&ret)
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) DBInsertOne(db string, cl string, doc interface{}, opts ...*options.InsertOneOptions) (result *mongo.InsertOneResult, err error) {
	result, err = cc.Database(db).Collection(cl, cc.collectionOptions).InsertOne(nil, doc, opts...)
	if err != nil {
		return
	}
	return
}

func (cc *Client) DBInsertMany(db string, cl string, docs []interface{}, opts ...*options.InsertManyOptions) (result *mongo.InsertManyResult, err error) {
	result, err = cc.Database(db).Collection(cl, cc.collectionOptions).InsertMany(nil, docs, opts...)
	if err != nil {
		return
	}
	return
}

func (cc *Client) DBReplaceId(db string, cl string, id interface{}, replace interface{}, opts ...*options.ReplaceOptions) (result *mongo.UpdateResult, err error) {
	result, err = cc.Database(db).Collection(cl, cc.collectionOptions).ReplaceOne(nil, bson.M{"_id": id}, replace, opts...)
	return
}

func (cc *Client) DBReplaceOne(db string, cl string, filter interface{}, replace interface{}, opts ...*options.ReplaceOptions) (result *mongo.UpdateResult, err error) {
	result, err = cc.Database(db).Collection(cl, cc.collectionOptions).ReplaceOne(nil, filter, replace, opts...)
	return
}

func (cc *Client) DBUpdateId(db string, cl string, id interface{}, update interface{}, opts ...*options.UpdateOptions) (result *mongo.UpdateResult, err error) {
	result, err = cc.Database(db).Collection(cl, cc.collectionOptions).UpdateOne(nil, bson.M{"_id": id}, update, opts...)
	return
}

func (cc *Client) DBUpdateOne(db string, cl string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (result *mongo.UpdateResult, err error) {
	result, err = cc.Database(db).Collection(cl, cc.collectionOptions).UpdateOne(nil, filter, update, opts...)
	return
}

func (cc *Client) DBUpdateMany(db string, cl string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (result *mongo.UpdateResult, err error) {
	result, err = cc.Database(db).Collection(cl, cc.collectionOptions).UpdateMany(nil, filter, update, opts...)
	return
}

func (cc *Client) DBDeleteId(db string, cl string, id interface{}, update interface{}, opts ...*options.DeleteOptions) (result *mongo.DeleteResult, err error) {
	result, err = cc.Database(db).Collection(cl, cc.collectionOptions).DeleteOne(nil, bson.M{"_id": id}, opts...)
	return
}

func (cc *Client) DBDeleteOne(db string, cl string, filter interface{}, update interface{}, opts ...*options.DeleteOptions) (result *mongo.DeleteResult, err error) {
	result, err = cc.Database(db).Collection(cl, cc.collectionOptions).DeleteOne(nil, filter, opts...)
	return
}

// 必须注意: empty filter会删除整个集合数据
func (cc *Client) DBDeleteMany(db string, cl string, filter interface{}, update interface{}, opts ...*options.DeleteOptions) (result *mongo.DeleteResult, err error) {
	result, err = cc.Database(db).Collection(cl, cc.collectionOptions).DeleteMany(nil, filter, opts...)
	return
}

func (cc *Client) DBAggregate(db string, cl string, pipeline interface{}, ret interface{}, opts ...*options.AggregateOptions) (err error) {
	cur, err := cc.Database(db).Collection(cl, cc.collectionOptions).Aggregate(nil, pipeline, opts...)
	if err == nil {
		err = cur.All(nil, ret)
	}
	return
}

func (cc *Client) DBAggregateWith(db string, cl string, pipeline interface{}, with func(cur *mongo.Cursor), opts ...*options.AggregateOptions) (err error) {
	cur, err := cc.Database(db).Collection(cl, cc.collectionOptions).Aggregate(nil, pipeline, opts...)
	if err == nil {
		defer cur.Close(nil)
		with(cur)
	}
	return
}

func (cc *Client) DBBulkWrite(db string, cl string, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (result *mongo.BulkWriteResult, err error) {
	if len(models) == 0 {
		return
	}
	result, err = cc.Database(db).Collection(cl, cc.collectionOptions).BulkWrite(nil, models, opts...)
	return
}
