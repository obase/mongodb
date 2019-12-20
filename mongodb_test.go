package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestDemo(t *testing.T) {
	opts := options.Client().SetAuth()
}
