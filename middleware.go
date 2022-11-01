package main

/* ==============================================
Copyright (c) Eensymachines
Developed by 		: kneerunjun@gmail.com
Developed on 		: OCT'22
Has all the handlers for the routes. Typical  func(c * gin.Context) implementations
============================================== */
import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
==================
- CORS enabling all cross origin requests for all verbs except OPTIONS
- this will be applied to all api across the board during the develpment stages
- do not apply this middleware though for routes that deliver web static content
====================
*/
func CORS(c *gin.Context) {
	// First, we add the headers with need to enable CORS
	// Make sure to adjust these headers to your needs
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")
	// Second, we handle the OPTIONS problem
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		// Everytime we receive an OPTIONS request,
		// we just return an HTTP 200 Status Code
		// Like this, Angular can now do the real
		// request using any other method than OPTIONS
		c.AbortWithStatus(http.StatusOK)
	}
}

func DBCollection(cl *mongo.Client, dbName, collName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("coll", cl.Database(dbName).Collection(collName))
	}
}

// ThrowErr :is able to digest error
// will also log the appropriate error debug fields and message
// will pack the context with response code
//
/*
	ThrowErr(fmt.Errorf("Accounts: invalid account in payload, cannot be nil"), log.WithFields(log.Fields{
			"payload": val,
	}), http.StatusBadRequest, c)
*/
func ThrowErr(e error, le *log.Entry, code int, c *gin.Context) {
	le.Error(e) // error gets logged
	switch code {
	case http.StatusBadRequest:
		c.AbortWithStatusJSON(code, gin.H{
			"err": "One or more inputs is nil/invalid/duplicate, check and send again",
		})
	case http.StatusInternalServerError:
		c.AbortWithStatusJSON(code, gin.H{
			"err": "One or more operations on the server has failed, please try after sometime",
		})
	case http.StatusNotFound:
		c.AbortWithStatusJSON(code, gin.H{
			"err": "One or ore resources you were looking for was not found",
		})
	}
}

// AccountPayload : desrialixing the account payload from the request
//
/*
	// from within the downstream handlers
	// this is howyou get the acccount as payload
	val, ok := c.Get("account")
	if !ok || val == nil {
		ThrowErr(fmt.Errorf("Accounts: invalid account in payload, cannot be nil"), log.WithFields(log.Fields{
			"payload": val,
		}), http.StatusBadRequest, c)
		return
	}
*/
func AccountPayload(c *gin.Context) {
	if c.Request.Method == "POST" || c.Request.Method == "PATCH" || c.Request.Method == "PUT" {
		// the verb tells me if the incoming request has the payload
		payload := &UserAccount{}
		if err := c.BindJSON(payload); err != nil {
			ThrowErr(fmt.Errorf("AccountPayload %s", err), log.WithFields(log.Fields{}), http.StatusBadRequest, c)
			return
		}
		c.Set("account", payload)
	}
}

// Accounts : when CRUD on collection of accounts
// Handles posting of new accounts
// handles getting index of accounts
func Accounts(c *gin.Context) {
	val, _ := c.Get("coll")
	coll := val.(*mongo.Collection)
	if c.Request.Method == "POST" {
		// posting a new account
		val, ok := c.Get("account")
		if !ok || val == nil {
			ThrowErr(fmt.Errorf("Accounts: invalid account in payload, cannot be nil"), log.WithFields(log.Fields{
				"payload": val,
			}), http.StatusBadRequest, c)
			return
		}
		acc, _ := val.(Account)
		// Now that we have an interface to the account
		if err := ValidateForCreate(acc); err != nil {
			// Invalid account details
			ThrowErr(fmt.Errorf("Accounts: Invalid account details"), log.WithFields(log.Fields{
				"email": acc.GetEmail(),
				"phone": acc.GetPhone(),
				"title": acc.GetTitle(),
			}), http.StatusBadRequest, c)
			return
		}
		if CheckDuplicate(acc, coll) != nil {
			// Duplicate account exists
			// email, phone are unique
			ThrowErr(fmt.Errorf("Accounts: Duplicate accounts detected"), log.WithFields(log.Fields{
				"email": acc.GetEmail(),
				"phone": acc.GetPhone(),
			}), http.StatusBadRequest, c)
			return
		}
		if err := CreateNewAccount(acc, coll); err != nil {
			// Error creating a new account
			ThrowErr(fmt.Errorf("Accounts: Failed query to create accounts %s", err), log.WithFields(log.Fields{}), http.StatusInternalServerError, c)
			return
		}
		// The account has been created
		c.AbortWithStatus(http.StatusCreated)
		return
	}
}

// AccountDetails : getting, modifying single account details
//
/*
 */
func AccountDetails(c *gin.Context) {
	val, _ := c.Get("coll")
	coll := val.(*mongo.Collection)
	id, _ := c.Params.Get("accid")          // unique id of the account as string
	oid, _ := primitive.ObjectIDFromHex(id) // object id of the account

	if c.Request.Method == "PUT" {
		// except email other details can be changed
		val, ok := c.Get("account")
		if !ok || val == nil {
			// Could not read payload
			ThrowErr(fmt.Errorf("Accounts: invalid account in payload, cannot be nil"), log.WithFields(log.Fields{
				"payload": val,
			}), http.StatusBadRequest, c)
			return
		}
		acc, _ := val.(Account)
		// email is an unique identifier for any account
		if ValidateForUpdate(acc) != nil {
			// Invalid account details, email identifies the account being modified
			ThrowErr(fmt.Errorf("Accounts: Invalid account details to validate"), log.WithFields(log.Fields{
				"email": acc.GetEmail(),
				"title": acc.GetTitle(),
				"phone": acc.GetPhone(),
			}), http.StatusBadRequest, c)
			return
		}
		if CheckExists(acc, coll, func(acc Account) bson.M {
			return bson.M{"email": acc.GetEmail()}
		}) != nil {
			// Account sought to be updated cannot be found
			// email, phone are unique
			ThrowErr(fmt.Errorf("Accounts: No account found"), log.WithFields(log.Fields{
				"email": acc.GetEmail(),
			}), http.StatusNotFound, c)
			return
		}
		if err := UpdateAccount(acc, coll); err != nil {
			ThrowErr(fmt.Errorf("Accounts: Failed query to update accounts %s", err), log.WithFields(log.Fields{
				"email": acc.GetEmail(),
			}), http.StatusInternalServerError, c)
			return
		}
		c.AbortWithStatus(http.StatusOK)
		return
	} else if c.Request.Method == "DELETE" {
		if CheckExists(&UserAccount{ID: oid}, coll, func(acc Account) bson.M {
			return bson.M{"_id": oid}
		}) != nil {
			ThrowErr(fmt.Errorf("Accounts: No account found"), log.WithFields(log.Fields{
				"id": id,
			}), http.StatusNotFound, c)
			return
		}
		// Data cannot be ever deleted permanently
		// Keeping up with the rules of the modern internet we need to still stove away the account details
		if err := ArchiveAccount(oid, coll); err != nil {
			ThrowErr(fmt.Errorf("Accounts: Failed query to delete account %s", err), log.WithFields(log.Fields{
				"_id": id,
			}), http.StatusInternalServerError, c)
			return
		}
		c.AbortWithStatus(http.StatusOK)
		return
	}
}
