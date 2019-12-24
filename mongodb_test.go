package mongodb

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestGet(t *testing.T) {
	 mdb := Must("test")
	 defer mdb.Close()

	fmt.Println(mdb.CountWith("zj_event_changed", bson.M{"_id":"11112"}))

}