package main

import (
	"context"
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// ValidateForCreate: validates the account details, checkls for errenous values
// will check to see if the email, phone and the title fall in a pattern
/*
	// your sample code here
*/
func ValidateForCreate(acc Account) error {
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

// ValidateForUpdate: validates the account for its details
// will check to see if the update`able fields are valid
// Incase the fields arent populated, or have zero value the check is missed
// Email is not an update`able` field
/*
	// your sample code here
*/
func ValidateForUpdate(acc Account) error {
	if acc.GetPhone() != "" {
		if matched, _ := regexp.MatchString(patternPhone, acc.GetPhone()); !matched {
			return fmt.Errorf("invalid email for the account")
		}
	}
	if acc.GetTitle() != "" {
		if matched, _ := regexp.MatchString(patternTitle, acc.GetTitle()); !matched {
			return fmt.Errorf("invalid email for the account")
		}
	}
	return nil
}

type AccountFilter func(Account) bson.M // this can filter accounts on various criteria

// CheckExists: Checks to see if there is exactly one account with the same email
// Will error in case the count of the documents is not equal to 1
// can customize the filter on which existence of the document is based
//
/*
	if CheckExists(acc, coll, func(acc Account) bson.M {
		return bson.M{"email": acc.GetEmail()}
	}) != nil {
		// error handling code here
	}
*/
func CheckExists(acc Account, coll *mongo.Collection, af AccountFilter) error {
	count, err := coll.CountDocuments(context.TODO(), af(acc))
	if err != nil {
		return fmt.Errorf("CheckExists: failed to get account")
	}
	if count != int64(1) {
		return fmt.Errorf("CheckExists: No account found")
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

// UpdateAccount: updates the title and phone of a single account
//
/*
	if err := UpdateAccount(acc, coll); err != nil {
		ThrowErr(fmt.Errorf("Accounts: Failed query to update accounts %s", err), log.WithFields(log.Fields{
			"email": acc.GetEmail(),
		}), http.StatusInternalServerError, c)
		return
	}
*/
func UpdateAccount(acc Account, coll *mongo.Collection) error {
	flt := bson.M{
		"email": acc.GetEmail(),
	}
	patch := bson.M{"title": acc.GetTitle(), "phone": acc.GetPhone()}
	result, err := coll.UpdateOne(context.Background(), flt, patch)
	if err != nil {
		return fmt.Errorf("UpdateOne: failed query, check database connection")
	}
	log.Infof("account updated : %v", result.UpsertedID)
	return nil
}

// ArchiveAccount : removes the account details from the collection only to be archived in some other collection
//
/*
	if err := ArchiveAccount(oid, coll); err != nil {
		ThrowErr(fmt.Errorf("Accounts: Failed query to delete account %s", err), log.WithFields(log.Fields{
			"_id": id,
		}), http.StatusInternalServerError, c)
		return
	}
*/
func ArchiveAccount(oid primitive.ObjectID, coll *mongo.Collection) error {
	flt := bson.M{"_id": oid}
	sr := coll.FindOne(context.TODO(), flt)
	archived := &UserAccount{}
	if err := sr.Decode(archived); err != nil {
		return err
	}
	// Archived replica, if this succeeds then we can remove from main collection
	// on this collection there isnt any unique key constraint on email, phone
	_, err := coll.Database().Collection("archvaccounts").InsertOne(context.TODO(), archived)
	if err != nil {
		return err
	}
	// when the account information is backed up its ready to be deleted
	// we move the sameto archived collecitons
	coll.DeleteOne(context.TODO(), flt)
	return nil
}
