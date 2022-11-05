package main

/* ==============================================
Copyright (c) Eensymachines
Developed by 		: kneerunjun@gmail.com
Developed on 		: OCT'22
All the middleware tests here
============================================== */
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*-----------------------------
Helper functions, private utility functions
-----------------------------*/
// getTestMongoColl : will get a mongo collection using the test connection
// tests are run on host machines and hence the mongo access is using the localhost connection
// mongo runs from within the container and hence the port that the container gets shared is importatn
func getTestMongoColl(name, dbname string) *mongo.Collection {
	todo, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(todo, options.Client().ApplyURI("mongodb://localhost:37017"))
	if err != nil {
		log.Errorf("getTestMongoColl: failed to connect to database %s", err)
	}
	if client.Ping(todo, nil) != nil {
		log.Error("getTestMongoColl: not connected to database, ping failed")
	}
	return client.Database(dbname).Collection(name)
}

// newTestGinContext : since we are testing on the middleware level we need to make a new mock context
func newTestGinContext(body interface{}, method, path string) *gin.Context {
	byt, _ := json.Marshal(body)
	reader := bytes.NewReader(byt)
	return &gin.Context{
		Request: &http.Request{
			Method: method,
			URL:    &url.URL{Path: fmt.Sprintf("http://sample.com/mockrequest/%s", path)},
			Proto:  "HTTP/1.1",
			Header: map[string][]string{
				"Accept-Encoding": {"application/json"},
			},
			Body: io.NopCloser(reader),
		},
	}
}
func TestAccountPayload(t *testing.T) {
	// Wil test the unmarshalling of the account into the route payload
	// https://gosamples.dev/struct-to-io-reader/
	byt, _ := json.Marshal(&UserAccount{})
	reader := bytes.NewReader(byt)
	ctx := &gin.Context{
		Request: &http.Request{
			Method: "POST",
			URL:    &url.URL{Path: "http://sample.com/applications/accounts"},
			Proto:  "HTTP/1.1",
			Header: map[string][]string{
				"Accept-Encoding": {"application/json"},
			},
			Body: io.NopCloser(reader),
		},
	}
	AccountPayload(ctx)
	acc, _ := ctx.Get("account")
	assert.NotNil(t, acc, "unpected fail to set accout object in context")
	// Setting an account with values
	sampleAcc := &UserAccount{
		Email: "john.dore@gmail.com",
		Phone: "8980982093",
		Title: "John Dore",
	}
	byt, _ = json.Marshal(sampleAcc)
	reader = bytes.NewReader(byt)
	ctx.Request.Body = io.NopCloser(reader)
	AccountPayload(ctx)
	val, _ := ctx.Get("account")
	accVal, _ := val.(Account)
	assert.NotNil(t, acc, "Unexpected nil value on the account")
	// Then checking for the value to determine if the values we sent in are the same that come out
	assert.Equal(t, accVal.GetEmail(), sampleAcc.Email, "could not match verify the email")
	assert.Equal(t, accVal.GetTitle(), sampleAcc.Title, "could not match verify the title")
	assert.Equal(t, accVal.GetPhone(), sampleAcc.Phone, "could not match verify the phone")

	// Setting a nil account in context
	byt, _ = json.Marshal(nil)
	reader = bytes.NewReader(byt)
	ctx.Request.Body = io.NopCloser(reader)
	AccountPayload(ctx)
	acc, _ = ctx.Get("account")
	assert.NotNil(t, acc, "unexpected nil account has been set in context")
}

// TestInsertOne : wrote this test when the object id being inserted was having problems with _id
func TestInsertOne(t *testing.T) {
	// ======== setting up the test
	ua, err := JsonSampleRandomAccount()
	if err != nil {
		t.Error("TestInsertOne: failed to get random sample account")
	}
	t.Log(ua)
	// ======== database connection
	coll := getTestMongoColl("accounts", "eensydb")
	assert.NotNil(t, coll, "Unexpected null collection, cannot proceed with test")
	result, err := coll.InsertOne(context.TODO(), ua)
	// ========= Testing
	assert.Nil(t, err, "failed to insert simple account")
	assert.NotNil(t, result.InsertedID, "unexpected nil id from insertion")
	t.Logf("Inserted id %v", result.InsertedID)
	// ========= Getting the inserted document
	sr := coll.FindOne(context.TODO(), bson.M{
		"_id": result.InsertedID,
	})
	assert.Nil(t, sr.Err(), "unexpected error when getting the user account")
	usrAcc := UserAccount{}
	err = sr.Decode(&usrAcc)
	assert.Nil(t, err, "Unexpected error when decoding user account")
	assert.Equal(t, ua.Email, usrAcc.Email, "Email of the accounts do not match")
	assert.Equal(t, ua.Title, usrAcc.Title, "Email of the accounts do not match")
	// And then when you are done inserting you can delete the account
	// ========= Cleaning up
	coll.DeleteOne(context.TODO(), bson.M{
		"email": ua.Email,
	})
}

// TestAccPostMddlware : this will test posting a new account to the database via middleware
// Mocking up the http request and then pushing it in
func TestAccPostMddlware(t *testing.T) {
	ua, err := JsonSampleRandomAccount()
	if err != nil {
		t.Error("TestInsertOne: failed to get random sample account")
	}
	t.Log(ua)
	ctx := newTestGinContext(ua, "POST", "accounts")

	// We woudl have to make a database connection to have it injected in the context
	todo, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(todo, options.Client().ApplyURI("mongodb://localhost:37017"))
	if err != nil {
		log.Panicf("failed to connect to database %s", err)
	}
	collAccs := client.Database("eensydb").Collection("accounts")
	ctx.Set("coll", collAccs)
	AccountPayload(ctx)
	Accounts(ctx)
	assert.Equal(t, ctx.Request.Response.StatusCode, http.StatusCreated, "http status code not as expected")
}
