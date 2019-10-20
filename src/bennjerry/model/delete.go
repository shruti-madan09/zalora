package model

import (
	"database/sql"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"constants"
	"logger"
	"mysqlc"
)

func SoftDeleteFromProductByProductId(productId string) (int, bool) {
	/*
		To take product_id and mark record in product table as inactive
	*/
	return UpdateProductIsInActiveByProductId(productId)
}

func DeleteFromProductById(txn *sql.Tx, id int) bool {
	/*
		To take id as input and delete record from product table
	*/
	funcName := "DeleteFromProductById"
	query := "DELETE FROM product WHERE id = " + strconv.Itoa(id)
	_, err := txn.Exec(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
		return false
	}
	return true
}

func DeleteFromProductSourcingValueByProductIdPK(txn *sql.Tx, ProductIdPK int) bool {
	/*
		To take product_id(primary key of product table) and delete record from product_sourcingvalue table
	*/
	funcName := "DeleteFromProductSourcingValueByProductIdPK"
	query := "DELETE FROM product_sourcingvalue WHERE product_id = " + strconv.Itoa(ProductIdPK)
	_, err := txn.Exec(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
		return false
	}
	return true
}

func DeleteFromProductSourcingValueById(txn *sql.Tx, productIdPK int, sourcingValueId int) bool {
	/*
		To take product_id(primary key of product table) and sourcingvalue_id
		and delete record from product_sourcingvalue table
	*/
	funcName := "DeleteFromProductSourcingValueById"
	query := "DELETE FROM product_sourcingvalue" +
		" WHERE product_id = " + strconv.Itoa(productIdPK) + " AND sourcingvalue_id = " + strconv.Itoa(sourcingValueId)
	_, err := txn.Exec(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
		return false
	}
	return true
}

func DeleteFromProductIngredientByProductIdPK(txn *sql.Tx, productIdPK int) bool {
	/*
		To take product_id(primary key of product table) and delete record from product_ingredient table
	*/
	funcName := "DeleteFromProductIngredientByProductIdPK"
	query := "DELETE FROM product_ingredient WHERE product_id = " + strconv.Itoa(productIdPK)
	_, err := txn.Exec(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
		return false
	}
	return true
}

func DeleteFromProductIngredientById(txn *sql.Tx, productIdPK int, ingredientId int) bool {
	/*
		To take product_id(primary key of product table) and ingredient_id as input
		and delete record from product_ingredient table
	*/
	funcName := "DeleteFromProductIngredientByProductIdPK"
	query := "DELETE FROM product_ingredient" + " WHERE product_id = " + strconv.Itoa(productIdPK) +
		" AND ingredient_id = " + strconv.Itoa(ingredientId)
	_, err := txn.Exec(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
		return false
	}
	return true
}

func DeleteUnUsedSourcingValue() bool {
	/*
		Deleting unused data from sourcingvalue table
	*/
	funcName := "DeleteUnUsedSourcingValue"
	query := "DELETE sourcingvalue FROM sourcingvalue LEFT JOIN product_sourcingvalue" +
		" ON sourcingvalue.id = product_sourcingvalue.sourcingvalue_id" +
		" WHERE product_sourcingvalue.sourcingvalue_id is NULL"
	_, err := mysqlc.MySqlDB.Exec(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
		return false
	}
	return true
}

func DeleteUnUsedIngredient() bool {
	/*
		Deleting unused data from ingredient table
	*/
	funcName := "DeleteUnUsedSourcingValue"
	query := "DELETE ingredient FROM ingredient LEFT JOIN product_ingredient" +
		" ON ingredient.id = product_ingredient.ingredient_id WHERE product_ingredient.ingredient_id is NULL"
	_, err := mysqlc.MySqlDB.Exec(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
		return false
	}
	return true
}

func DeleteUnUsedDietaryCertification() bool {
	/*
		Deleting unused data from dietarycertification table
	*/
	funcName := "DeleteUnUsedDietaryCertification"
	query := "DELETE dietarycertification FROM dietarycertification LEFT JOIN product" +
		" ON dietarycertification.id = product.dietary_certification_id WHERE product.id is NULL"
	_, err := mysqlc.MySqlDB.Exec(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
		return false
	}
	return true
}
