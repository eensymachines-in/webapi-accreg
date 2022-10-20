package main

/* ==============================================
Copyright (c) Eensymachines
Developed by 		: kneerunjun@gmail.com
Developed on 		: OCT'22
Eensymachines accounts need to be maintained over an api endpoint
containerized application can help do that
============================================== */
import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	FVerbose, FLogF bool
	logFile         string
)

const (
	// getting connected to a container that runs mongo
	CONNSTR = "mongodb://srvmongo:27017"
	DBNAME  = "eensydb"
)

func init() {
	flag.BoolVar(&FVerbose, "verbose", false, "Level of log messages")
	flag.BoolVar(&FLogF, "flog", false, "Log output direction")
	// Setting up log configuration for the api
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	log.SetReportCaller(false)
	// By default the log output is stdout and the level is info
	log.SetOutput(os.Stdout)     // FLogF will set it main, but dfault is stdout
	log.SetLevel(log.DebugLevel) // default level info debug but FVerbose will set it main
	logFile = os.Getenv("LOGF")
}

func main() {
	/*Log setup : level of logging and direction of logging*/
	flag.Parse() // command line flags are parsed
	log.WithFields(log.Fields{
		"verbose": FVerbose,
		"flog":    FLogF,
	}).Info("Log configuration..")
	if FVerbose {
		log.SetLevel(log.DebugLevel)
	}
	if FLogF {
		lf, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Failed to connect to log file, kindly check the privileges")
		} else {
			log.Infof("Check log file for entries @ %s", logFile)
			log.SetOutput(lf)
		}
	}
	log.Info("Now starting account services..")
	defer log.Warn("Now shutting down account services..")
	// ------------Setting up the database connections
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(CONNSTR))
	if err != nil {
		log.Panicf("failed to connect to database %s", err)
	}
	collAccs := client.Database(DBNAME).Collection("accounts")
	indexes, _ := collAccs.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{primitive.E{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
	})
	fmt.Println(len(indexes))
	// ------------- Database connection setup
	defer func() {
		log.Warn("Smart bill: Now diconnecting database connection")
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"app":       "acccountservices",
			"logs":      logFile,
			"verblog":   FVerbose,
			"logtofile": FLogF,
		})
	})
	r.Use(CORS)
	accounts := r.Group("/accounts")
	accounts.Use(DBCollection(client, "eensydb", "accounts"))
	accounts.POST("", AccountPayload, Accounts) // creates new account
	// incase of get and delete the AccountPayload middleware will be in action
	// incase of put, patch the account details need to be read back
	account := accounts.Group("/:accid", AccountPayload)
	// CRUD on single account
	account.GET("", AccountDetails)
	account.PATCH("", AccountDetails)
	account.DELETE("", AccountDetails)
	log.Fatal(r.Run(":8080"))
}
