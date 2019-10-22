package model

import (
	"bennjerry/structs"
	"mysqlc"
	"utils"
)

var logIdentifier = "bennjerry.model."

func InsertRecord(iceCreamData []*structs.IceCreamDataStruct) ([]int, bool) {
	/*
		To take a list of ice cream data and insert data using an atomic transaction
		Arguments: List of ice cream data
		Return: List of ids of records inserted in table
	*/
	success := true
	sourcingValuesMap := make(map[string]bool)
	ingredientsMap := make(map[string]bool)
	dietaryCertificationsMap := make(map[string]bool)
	for _, iceCream := range iceCreamData {
		// Creating maps out of sourcing values, ingredients and dietary certifications
		// Using Maps to fetch and keep unique values of each property
		// These unique values will later be inserted into corresponding tables
		sourcingValuesMap = utils.ListToMap(iceCream.SourcingValues)
		ingredientsMap = utils.ListToMap(iceCream.Ingredients)
		if iceCream.DietaryCertifications != "" {
			dietaryCertificationsMap[iceCream.DietaryCertifications] = true
		}
	}
	// Creating mysql transaction
	// If an query operation returns success=false, transaction will be rolled back, else committed at last
	mySqlTxn, mySqlTxnErr := mysqlc.MySqlDB.Begin()
	if mySqlTxnErr != nil {
		panic(mySqlTxnErr.Error())
	}
	// Inserting data into sourcingvalue, ingredient, dietarycertification tables
	// All three tables have unique constraint on name column to avoid duplicate entries
	// Therefore, insert ignore query is being used to avoid error if name already exists in table
	success = InsertIntoSourcingValue(mySqlTxn, sourcingValuesMap)
	if !success {
		mySqlTxn.Rollback()
		return nil, false
	}
	success = InsertIntoIngredient(mySqlTxn, ingredientsMap)
	if !success {
		mySqlTxn.Rollback()
		return nil, false
	}
	success = InsertIntoDietaryCertification(mySqlTxn, dietaryCertificationsMap)
	if !success {
		mySqlTxn.Rollback()
		return nil, false
	}

	idList := make([]int, 0)
	for _, iceCream := range iceCreamData {
		id, success := InsertIntoProduct(mySqlTxn, iceCream)
		// success: false, if an error occurs while running the query
		// success: true, if record is inserted successfully
		// id: the primary key of the inserted record and will be 0 in case of error
		if !success {
			mySqlTxn.Rollback()
			return nil, false
		} else {
			idList = append(idList, id)
			// Data will be inserted to relation table of product and sourcing value
			if len(iceCream.SourcingValues) > 0 {
				success = InsertIntoProductSourcingValue(mySqlTxn, id, iceCream.SourcingValues)
				if !success {
					mySqlTxn.Rollback()
					return nil, false
				}
			}
			// Data will be inserted to relation table of product and ingredient
			if len(iceCream.Ingredients) > 0 {
				success = InsertIntoProductIngredient(mySqlTxn, id, iceCream.Ingredients)
				if !success {
					mySqlTxn.Rollback()
					return nil, false
				}
			}
		}
	}
	mySqlTxn.Commit()
	return idList, true
}

func UpdateRecord(id int, iceCreamData *structs.IceCreamDataStruct, fieldMap map[string]bool) bool {
	/*
		To take an id and ice cream data and update data for that id using an atomic transaction
		Arguments: List of ice cream data
		Return: Boolean to indicate success or failure
	*/
	success := true
	// Creating mysql transaction
	// If an query operation returns success=false, transaction will be rolled back, else committed at last
	mySqlTxn, mySqlTxnErr := mysqlc.MySqlDB.Begin()
	if mySqlTxnErr != nil {
		panic(mySqlTxnErr.Error())
	}
	// Updating data in product table
	success = UpdateProductById(mySqlTxn, id, iceCreamData, fieldMap)
	if !success {
		mySqlTxn.Rollback()
		return false
	}
	if _, exists := fieldMap["sourcing_values"]; exists {
		sourcingValuesMap := utils.ListToMap(iceCreamData.SourcingValues)
		// Inserting any sourcing value name that is not already in table
		success = InsertIntoSourcingValue(mySqlTxn, sourcingValuesMap)
		if !success {
			mySqlTxn.Rollback()
			return false
		}
		// Updating relation table of product and sourcingvalue
		success = UpdateProductSourcingValueByProductIdPK(mySqlTxn, id, sourcingValuesMap)
		if !success {
			mySqlTxn.Rollback()
			return false
		}
	}
	if _, exists := fieldMap["ingredients"]; exists {
		ingredientsMap := utils.ListToMap(iceCreamData.Ingredients)
		// Inserting any ingredient name that is not already in table
		success = InsertIntoIngredient(mySqlTxn, ingredientsMap)
		if !success {
			mySqlTxn.Rollback()
			return false
		}
		// Updating relation table of product and ingredient
		success = UpdateProductIngredientByProductIdPK(mySqlTxn, id, ingredientsMap)
		if !success {
			mySqlTxn.Rollback()
			return false
		}
	}
	mySqlTxn.Commit()
	return true
}

func DropRecord(id int) bool {
	/*
		To take an id and delete data for that id using an atomic transaction
		Arguments: id
		Return: Boolean to indicate success or failure
	*/
	success := true
	mySqlTxn, mySqlTxnErr := mysqlc.MySqlDB.Begin()
	if mySqlTxnErr != nil {
		panic(mySqlTxnErr.Error())
	}
	// Deleting references of record in relation table of product and sourcingvalue
	success = DeleteFromProductSourcingValueByProductIdPK(mySqlTxn, id)
	if !success {
		mySqlTxn.Rollback()
		return false
	}
	// Deleting references of record in relation table of product and ingredient
	success = DeleteFromProductIngredientByProductIdPK(mySqlTxn, id)
	if !success {
		mySqlTxn.Rollback()
		return false
	}
	// Deleting actual record from product table
	success = DeleteFromProductById(mySqlTxn, id)
	if !success {
		mySqlTxn.Rollback()
		return false
	}
	mySqlTxn.Commit()
	// Post deletion, there might be unused sourcing values, ingredients and dietary certifications
	// Deleting such unused entries from the tables to avoid stale data
	DeleteUnUsedSourcingValue()
	DeleteUnUsedIngredient()
	DeleteUnUsedDietaryCertification()
	return true
}
