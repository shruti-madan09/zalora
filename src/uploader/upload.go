package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"bennjerry/model"
	"bennjerry/structs"
	"mysqlc"
)

func main() {
	// connecting to mysql
	mysqlc.DBConnecting()

	// reading json file with list of IceCreamDataStruct
	jsonFile, err := os.Open("icecream.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Converting byte data to structure
	var iceCreamData []*structs.IceCreamDataStruct
	umMarshalErr := json.Unmarshal(byteValue, &iceCreamData)
	if umMarshalErr != nil {
		fmt.Println(umMarshalErr.Error())
	} else {
		// Calling function to execute queries in an atomic transaction
		model.InsertRecord(iceCreamData)
	}

	// closing connection with mysql
	mysqlc.DBClosing()
}
