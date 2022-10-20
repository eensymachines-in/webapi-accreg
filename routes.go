package main

/* ==============================================
Copyright (c) Eensymachines
Developed by 		: kneerunjun@gmail.com
Developed on 		: OCT'22
Eensymachines accounts need to be maintained over an api endpoint
containerized application can help do that
============================================== */
import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

/*==================
- CORS enabling all cross origin requests for all verbs except OPTIONS
- this will be applied to all api across the board during the develpment stages
- do not apply this middleware though for routes that deliver web static content
====================*/
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

// AccountPayload : desrialixing the account payload from the request
func AccountPayload(c *gin.Context) {
	if c.Request.Method == "POST" || c.Request.Method == "PATCH" || c.Request.Method == "PUT" {
		// the verb tells me if the incoming request has the payload
		payload := &UserAccount{}
		if err := c.BindJSON(payload); err != nil {
			log.WithFields(log.Fields{
				"payload": "debug your payload here",
			}).Error("AccountPayload:failed to bind account payload")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
			return
		}
		c.Set("account", payload)
	}
}

// Accounts : when CRUD on collection of accounts
func Accounts(c *gin.Context) {
	val, _ := c.Get("coll")
	coll := val.(*mongo.Collection)
	if c.Request.Method == "POST" {
		// posting a new account
		val, ok := c.Get("account")
		if !ok || val == nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		acc, _ := val.(Account)
		// Now that we have an interface to the account
		if Validate(acc) != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if CheckDuplicate(acc, coll) != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if CreateNewAccount(acc, coll) != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
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
	if c.Request.Method == "PUT" {
		// except email other details can be changed
		// email is an unique identifier for any account

	}
}
