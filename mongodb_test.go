package mongodb

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestGet(t *testing.T) {
	mdb := Must("test")
	defer mdb.Close()

	fmt.Println(mdb.ListDatabaseNames())
	fmt.Println(mdb.ListCollectionNames())
	var rec bson.M
	nul, err := mdb.FindId("ods_ad_biz", 1659565190724636, &rec)
	if err != nil {
		panic(err)
	}

	fmt.Println(nul, rec)
}

func TestClient_ListDatabaseNames(t *testing.T) {
	mdb := Must("test")

	fmt.Println(mdb.ListDatabaseNames())
}

func TestClient_ListCollectionNames(t *testing.T) {
	mdb := Must("test")

	fmt.Println(mdb.ListCollectionNames())
}

func TestClient_Count(t *testing.T) {
	mdb := Must("test")

	fmt.Println(mdb.Count("rtm_macro_mac"))
	fmt.Println(mdb.Count("rtm_macro_mac", bson.M{"_id": "rUZNJrakx5QCVX8l9vGr"}))
}

func TestClient_Find(t *testing.T) {
	mdb := Must("test")
	var result []bson.M
	err := mdb.Find("rtm_macro_mac", mdb.ALL, &result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	err = mdb.Find("rtm_macro_mac", bson.M{"_id": "rUZNJrakx5QCVX8l9vGr"}, &result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	var result2 []bson.M
	err = mdb.Find("rtm_macro_mac", bson.M{"_id": "rUZNJrakx5QCVX8l9vGr"}, &result2, options.Find().SetProjection(bson.M{"mtime": 1}).SetSkip(1).SetLimit(2).SetSort(bson.M{"code": -1}))
	if err != nil {
		panic(err)
	}
	fmt.Println(result2)
}

func TestClient_FindId(t *testing.T) {
	mdb := Must("test")
	var result bson.M
	nvl, err := mdb.FindId("rtm_macro_mac", "rUZNJrakx5QCVX8l9vGr", &result)
	if err != nil {
		panic(err)
	}
	fmt.Println(nvl, result)
}

func TestClient_FindOne(t *testing.T) {
	mdb := Must("test")
	var result bson.M
	nvl, err := mdb.FindOne("rtm_macro_mac", bson.M{"_id": "rUZNJrakx5QCVX8l9vGr"}, &result, options.FindOne().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		panic(err)
	}
	fmt.Println(nvl, result)
}

func TestClient_FindWith(t *testing.T) {
	mdb := Must("test")
	err := mdb.FindWith("rtm_macro_mac", bson.M{"_id": "rUZNJrakx5QCVX8l9vGr"}, func(cur *mongo.Cursor) error {
		var rec bson.M
		for cur.Next(nil) {
			cur.Decode(&rec)
			fmt.Println(rec)
		}
		return nil
	}, options.Find().SetProjection(bson.M{"mtime": 1}))
	if err != nil {
		panic(err)
	}

}

func TestClient_Distinct(t *testing.T) {
	mdb := Must("test")
	fmt.Println(mdb.Distinct("rtm_macro_mac", "code", mdb.ALL))
}

func TestClient_FindIdAndReplace(t *testing.T) {
	mdb := Must("test")
	var rec bson.M
	result, err := mdb.FindIdAndReplace("test", "123", bson.M{"code": "xxx123456"}, &rec, options.FindOneAndReplace().SetUpsert(true).SetReturnDocument(options.After))
	if err != nil {
		panic(err)
	}
	fmt.Println(result, rec)
}

func TestClient_FindIdAndUpdate(t *testing.T) {
	mdb := Must("test")
	var rec bson.M
	result, err := mdb.FindIdAndUpdate("test", "123", bson.M{"$set": bson.M{"code": "xxx123"}}, &rec, options.FindOneAndUpdate().SetUpsert(true))
	if err != nil {
		panic(err)
	}
	fmt.Println(result, rec)
}

func TestClient_FindIdAndDelete(t *testing.T) {
	mdb := Must("test")
	var rec bson.M
	result, err := mdb.FindIdAndDelete("test", "123", &rec)
	if err != nil {
		panic(err)
	}
	fmt.Println(result, rec)
}

func TestClient_InsertOne(t *testing.T) {
	mdb := Must("test")
	result, err := mdb.InsertOne("test", bson.M{"code": "another"})
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func TestClient_InsertMany(t *testing.T) {
	mdb := Must("test")
	result, err := mdb.InsertMany("test", []interface{}{bson.M{"code": "another"}})
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func TestClient_ReplaceId(t *testing.T) {
	mdb := Must("test")
	var rec bson.M
	fmt.Println(mdb.FindId("test", mdb.ObjectId("5eddbaee58b7a483c9226ae2"), &rec))
	fmt.Println(rec)
}
