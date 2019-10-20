package model

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"bennjerry/structs"
	"constants"
	"logger"
)

func InsertIntoProduct(txn *sql.Tx, iceCreamData *structs.IceCreamDataStruct) (int, bool) {
	/*
		To take ice cream data and insert it into product table
	*/
	funcName := "InsertIntoProduct"
	if iceCreamData == nil {
		return 0, false
	}
	query := "INSERT INTO product (product_id, name, description, story, image_closed, image_opened, allergy"
	if iceCreamData.DietaryCertifications != "" {
		query += ", dietary_certification_id"
	}
	query += ")"
	query += " VALUES ('" + iceCreamData.ProductId + "'"
	query += ", '" + strings.Replace(iceCreamData.Name, "'", "''", -1) + "'"
	query += ", '" + strings.Replace(iceCreamData.Description, "'", "''", -1) + "'"
	query += ", '" + strings.Replace(iceCreamData.Story, "'", "''", -1) + "'"
	query += ", '" + strings.Replace(iceCreamData.ImageClosed, "'", "''", -1) + "'"
	query += ", '" + strings.Replace(iceCreamData.ImageOpened, "'", "''", -1) + "'"
	query += ", '" + strings.Replace(iceCreamData.AllergyInfo, "'", "''", -1) + "'"
	if iceCreamData.DietaryCertifications != "" {
		dietaryCertificationId := SelectFromDietaryCertification(txn, []string{iceCreamData.DietaryCertifications})
		query += ", '" + strconv.Itoa(dietaryCertificationId[0].Id) + "'"
	}
	query += ")"
	insert, err := txn.Exec(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
	} else {
		id, _ := insert.LastInsertId()
		return int(id), true
	}
	return 0, false
}

func InsertIntoSourcingValue(txn *sql.Tx, nameMap map[string]bool) bool {
	/*
		To take map {name: true}, and insert into sourcingvalue, if it doesn't exist already
	*/
	funcName := "InsertIntoSourcingValue"
	for name := range nameMap {
		query := "INSERT IGNORE INTO sourcingvalue (name)" +
			" VALUES ('" + strings.Replace(name, "'", "''", -1) + "')"
		_, err := txn.Exec(query)
		if err != nil {
			logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
				constants.MySQLQueryRunErrorMessage, err.Error())
			return false
		}
	}
	return true
}

func InsertIntoIngredient(txn *sql.Tx, nameMap map[string]bool) bool {
	/*
		To take map {name: true} and insert into ingredient, if it doesn't exist already
	*/
	funcName := "InsertIntoIngredient"
	for name := range nameMap {
		query := "INSERT IGNORE INTO ingredient (name)" +
			" VALUES ('" + strings.Replace(name, "'", "''", -1) + "')"
		_, err := txn.Exec(query)
		if err != nil {
			logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
				constants.MySQLQueryRunErrorMessage, err.Error())
			return false
		}
	}
	return true
}

func InsertIntoDietaryCertification(txn *sql.Tx, nameMap map[string]bool) bool {
	/*
		To take map {name: true} and insert into dietarycertification, if it doesn't exist already
	*/
	funcName := "InsertIntoDietaryCertification"
	for name := range nameMap {
		query := "INSERT IGNORE INTO dietarycertification (name)" +
			" VALUES ('" + strings.Replace(name, "'", "''", -1) + "')"
		_, err := txn.Exec(query)
		if err != nil {
			logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
				constants.MySQLQueryRunErrorMessage, err.Error())
			return false
		}
	}
	fmt.Println("insertign to InsertIntoDietaryCertification", nameMap)
	return true
}

func InsertToProductSourcingValueById(txn *sql.Tx, productIdPk int, sourcingValueId int) bool {
	/*
		To take product_id (primary key of product table) and sourcingvalue_id
		and insert into product_sourcingvalue table
	*/
	funcName := "InsertToProductSourcingValueById"
	query := "INSERT INTO product_sourcingvalue (product_id, sourcingvalue_id)" +
		" VALUES (" + strconv.Itoa(productIdPk) + ", " + strconv.Itoa(sourcingValueId) + ")"
	_, err := txn.Exec(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
		return false
	}
	return true
}

func InsertIntoProductSourcingValue(txn *sql.Tx, productIdPK int, sourcingValues []string) bool {
	/*
		To take a list of sourcing values and insert into product_sourcingvalue table
	*/
	sourcingValueDataFromDB := SelectFromSourcingValue(txn, sourcingValues)
	for _, data := range sourcingValueDataFromDB {
		success := InsertToProductSourcingValueById(txn, productIdPK, data.Id)
		if !success {
			return false
		}
	}
	return true
}

func InsertToProductIngredientById(txn *sql.Tx, productIdPK int, ingredientId int) bool {
	/*
		To take product_id (primary key of product table) and ingredient_id and insert into product_ingredient table
	*/
	funcName := "InsertToProductIngredientById"
	query := "INSERT INTO product_ingredient (product_id, ingredient_id)" +
		" VALUES (" + strconv.Itoa(productIdPK) + ", " + strconv.Itoa(ingredientId) + ")"
	_, err := txn.Exec(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
		return false
	}
	return true
}

func InsertIntoProductIngredient(txn *sql.Tx, productId int, ingredients []string) bool {
	/*
		To take a list of ingredients and insert into product_ingredient table
	*/
	IngredientDataFromDB := SelectFromIngredient(txn, ingredients)
	for _, data := range IngredientDataFromDB {
		success := InsertToProductIngredientById(txn, productId, data.Id)
		if !success {
			return false
		}
	}
	return true
}
