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
func JsonSampleUserAccounts() ([]*UserAccount, error) {
	byt, err := ioutil.ReadFile("useraccs.json")
	if err != nil {
		return nil, err
	}
	result := []*UserAccount{} // the result we send back
	if err := json.Unmarshal(byt, &result); err != nil {
		return nil, err
	}
	return result, nil
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
