package main

import (
	"bennjerry/structs"
	"bennjerry/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"constants"
	"logger"
	"mysqlc"
)

func main() {
	logIdentifier := "uploader"

	// connecting to mysql
	mysqlc.DBConnecting()
	logger.Init()

	// reading json file with list of IceCreamDataStruct
	jsonFile, err := os.Open("icecream.json")
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.BenNJerryLogBucketName, logIdentifier,
			"Error while reading file", err.Error())
		fmt.Println(err.Error())
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Converting byte data to structure
	var iceCreamData []*structs.IceCreamDataStruct
	umMarshalErr := json.Unmarshal(byteValue, &iceCreamData)
	if umMarshalErr != nil {
		logger.ZaloraStatsLogger.Error(constants.BenNJerryLogBucketName, logIdentifier, constants.UnMarshalErrorString,
			umMarshalErr.Error())
		fmt.Println(umMarshalErr.Error())
	} else {
		// Calling function to execute queries in an atomic transaction
		model.InsertRecord(iceCreamData)
	}

	// closing connection with mysql
	mysqlc.DBClosing()
}
