package main

import (
	"context"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestDBConnection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://srvmongo:27017"))
	if err != nil {
		log.Panicf("failed to connect to database %s", err)
	}
	assert.NotNil(t, client, "uexpected nil valiue for the client")
}
