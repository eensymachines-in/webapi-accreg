package main

/* ==============================================
Copyright (c) Eensymachines
Developed by 		: kneerunjun@gmail.com
Developed on 		: OCT'22
This was developed to read and write test data and make it available as in memory slice
test data resides in json alongside
============================================== */
import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
)

// JsonSampleUserAccounts : shall read the entire file and then send back the block result
// limit 		: if you want only a sample of the entire seed and not the entire dump
//
/*
	seed, err := JsonSampleUserAccounts(-1)
	if err !=nil{
		return err
	}
	fmt.Printf("length of the sntire seed data is %d", len(seed))
	seed, err = JsonSampleUserAccounts(10)
	fmt.Println("We are expecting to have only 10 items in the sample seed")

*/
func JsonSampleUserAccounts(limit int) ([]*UserAccount, error) {
	byt, err := ioutil.ReadFile("useraccs.json")
	if err != nil {
		return nil, err
	}
	result := []*UserAccount{} // the result we send back
	if err := json.Unmarshal(byt, &result); err != nil {
		return nil, err
	}
	if limit > 0 {
		return result[:limit], nil
	}
	return result, nil // case where the entire dump is requested

}

// JsonSampleRandomAccount : from the entire sample of user accounts this will get the random account
func JsonSampleRandomAccount() (*UserAccount, error) {
	sample, err := JsonSampleUserAccounts()
	if err != nil {
		return nil, err
	}
	rndIndx := rand.Intn(len(sample) - 1)
	return sample[rndIndx], nil
}
