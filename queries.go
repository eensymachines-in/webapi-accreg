package main

import (
	"context"
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/* ==============================================
Copyright (c) Eensymachines
Developed by 		: kneerunjun@gmail.com
Developed on 		: OCT'22
Middleware queries for accounts in eensymachines
============================================== */

const (
	// link to go tools to test this email pattern
	patternEmail = `^[[:alnum:]]+[.\-_]{0,1}[[:alnum:]]*[@]{1}[[:alpha:]]+[.]{1}[[:alnum:]]{2,}[.]{0,1}[[:alnum:]]{0,}$`
	patternPhone = `^[0-9]{10}$` // we are assuming that we have phone numbers from India only
	// title of the account allowed is about 1 to 16 characters including numbers and a handful of special characters
	patternTitle = `^[a-zA-Z0-9_\-.]{1,16}$`
)

// Validate: validates the account details, checkls for errenous values
// will check to see if the email, phone and the title fall in a pattern
/*
	// your sample code here
*/
func Validate(acc Account) error {
	if matched, _ := regexp.MatchString(patternEmail, acc.GetEmail()); !matched {
		return fmt.Errorf("invalid email for the account")
	}
	if matched, _ := regexp.MatchString(patternPhone, acc.GetPhone()); !matched {
		return fmt.Errorf("invalid email for the account")
	}
	if matched, _ := regexp.MatchString(patternTitle, acc.GetTitle()); !matched {
		return fmt.Errorf("invalid email for the account")
	}
	return nil
}

// CheckDuplicate: Checks for duplicates on account unique fields
// Will check to see if email is duplicate or the phone is
// will error in case duplicate is found
//
/*
	// your sample code here
*/
func CheckDuplicate(acc Account, coll *mongo.Collection) error {
	emailFlt := bson.M{
		"email": acc.GetEmail(),
	}
	count, err := coll.CountDocuments(context.TODO(), emailFlt)
	if err != nil {
		return fmt.Errorf("failed to get account duplicates")
	}
	if count > int64(0) {
		return fmt.Errorf("account with duplicate email found")
	}
	phoneFlt := bson.M{
		"phone": acc.GetPhone(),
	}
	count, err = coll.CountDocuments(context.TODO(), phoneFlt)
	if err != nil {
		return fmt.Errorf("failed to get account duplicates")
	}
	if count > int64(0) {
		return fmt.Errorf("account with duplicate phone found")
	}
	return nil
}

// CreateNewAccount: registers new account to the database
//
/*
	// your sample code here
*/
func CreateNewAccount(acc Account, coll *mongo.Collection) error {
	result, err := coll.InsertOne(context.Background(), acc)
	if err != nil {
		return fmt.Errorf("InsertOne: failed query, check database connection")
	}
	log.Infof("new account inserted : %v", result.InsertedID)
	return nil
}
