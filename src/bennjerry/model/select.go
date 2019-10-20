package model

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"constants"
	"logger"
	"mysqlc"
)

func SelectFromProductByProductId(productId string) (*Product, bool) {
	/*
		To take product_id and select columns from product table
	*/
	funcName := "SelectFromProductByProductId"
	query := "SELECT id, product_id, name, description, story, image_closed, image_opened, allergy," +
		" dietary_certification_id FROM product WHERE is_inactive = 0 and product_id = '" + productId + "'"
	selectQ, err := mysqlc.MySqlDB.Query(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
	} else {
		product := &Product{}
		for selectQ.Next() {
			err := selectQ.Scan(&product.Id, &product.ProductId, &product.Name, &product.Description, &product.Story,
				&product.ImageClosed, &product.ImageOpened, &product.Allergy, &product.DietaryCertificationId)
			if err != nil {
				logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
					constants.MySQLSelectScanErrorMessage, err.Error())
			}
		}
		return product, true
	}
	return nil, false
}

func SelectIdFromProductByProductId(productId string) (int, bool) {
	/*
		To take product_id and select id from product table
	*/
	funcName := "SelectIdFromProductByProductId"
	query := "SELECT id FROM product WHERE product_id = '" + productId + "'"
	selectQ, err := mysqlc.MySqlDB.Query(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
	} else {
		var id int
		for selectQ.Next() {
			err := selectQ.Scan(&id)
			if err != nil {
				logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
					constants.MySQLSelectScanErrorMessage, err.Error())
			}
			fmt.Println("id", id)
		}
		fmt.Println("returning ", id)
		return id, true
	}
	return 0, false
}

func SelectFromSourcingValue(txn *sql.Tx, nameList []string) []*Property {
	/*
		To take list of names and select id, name from sourcingvalue table
	*/
	funcName := "SelectFromSourcingValue"
	result := make([]*Property, 0)
	lenNameList := len(nameList)
	query := "SELECT id, name FROM sourcingvalue WHERE name In ("
	for index := 0; index < lenNameList-1; index++ {
		query += "'" + nameList[index] + "', "
	}
	query += "'" + nameList[lenNameList-1] + "')"
	selectQ, err := txn.Query(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
	} else {
		for selectQ.Next() {
			sourcingValue := &Property{}
			err := selectQ.Scan(&sourcingValue.Id, &sourcingValue.Name)
			if err != nil {
				logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
					constants.MySQLSelectScanErrorMessage, err.Error())
			} else {
				result = append(result, sourcingValue)
			}
		}
	}
	return result
}

func SelectFromIngredient(txn *sql.Tx, nameList []string) []*Property {
	/*
		To take list of names and select id, name from ingredient table
	*/
	funcName := "SelectFromIngredient"
	result := make([]*Property, 0)
	lenNameList := len(nameList)
	query := "SELECT id, name FROM ingredient WHERE name In ("
	for index := 0; index < lenNameList-1; index++ {
		query += "'" + nameList[index] + "', "
	}
	query += "'" + nameList[lenNameList-1] + "')"
	selectQ, err := txn.Query(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
	} else {
		for selectQ.Next() {
			ingredient := &Property{}
			err := selectQ.Scan(&ingredient.Id, &ingredient.Name)
			if err != nil {
				logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
					constants.MySQLSelectScanErrorMessage, err.Error())
			} else {
				result = append(result, ingredient)
			}
		}
	}
	return result
}

func SelectFromDietaryCertification(txn *sql.Tx, nameList []string) []*Property {
	/*
		To take list of names and select id, name from dietarycertification table
	*/
	funcName := "SelectFromDietaryCertification"
	result := make([]*Property, 0)
	lenNameList := len(nameList)
	query := "SELECT id, name FROM dietarycertification WHERE name In ("
	for index := 0; index < lenNameList-1; index++ {
		query += "'" + nameList[index] + "', "
	}
	query += "'" + nameList[lenNameList-1] + "')"
	selectQ, err := txn.Query(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
	} else {
		for selectQ.Next() {
			dietaryCertification := &Property{}
			err := selectQ.Scan(&dietaryCertification.Id, &dietaryCertification.Name)
			if err != nil {
				logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
					constants.MySQLSelectScanErrorMessage, err.Error())
			} else {
				result = append(result, dietaryCertification)
			}
		}
	}
	return result
}

func SelectFromDietaryCertificationById(id int) *Property {
	/*
		To take id (primary key) and select name from dietarycertification table
	*/
	funcName := "SelectFromDietaryCertificationById"
	query := "SELECT name FROM dietarycertification WHERE id = " + strconv.Itoa(id)
	selectQ, err := mysqlc.MySqlDB.Query(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
	} else {
		data := &Property{Id: id}
		for selectQ.Next() {
			data.Id = id
			err := selectQ.Scan(&data.Name)
			if err != nil {
				logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
					constants.MySQLSelectScanErrorMessage, err.Error())
			}
			return data
		}
	}
	return nil
}

func SelectSourcingValueNameByProductIdPK(productIdPK int) []string {
	/*
		To take product_id (primary key of product table) and select sourcingvalue name
	*/
	funcName := "SelectSourcingValueNameByProductIdPK"
	result := make([]string, 0)
	query := "SELECT sourcingvalue.name FROM product_sourcingvalue INNER JOIN sourcingvalue" +
		" ON product_sourcingvalue.sourcingvalue_id = sourcingvalue.id" +
		" WHERE product_sourcingvalue.product_id = " + strconv.Itoa(productIdPK)
	selectQ, err := mysqlc.MySqlDB.Query(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
	} else {
		var name string
		for selectQ.Next() {
			err := selectQ.Scan(&name)
			if err != nil {
				logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
					constants.MySQLSelectScanErrorMessage, err.Error())
			} else {
				result = append(result, name)
			}
		}
	}
	return result
}

func SelectFromProductSourcingValueByProductIdPK(txn *sql.Tx, productIdPK int) []*ProductProperty {
	/*
		To take product_id (primary key of product table) and select product id, sourcingvalue id, sourcingvalue name
	*/
	funcName := "SelectFromProductSourcingValueByProductIdPK"
	result := make([]*ProductProperty, 0)
	query := "SELECT product_sourcingvalue.product_id, product_sourcingvalue.sourcingvalue_id, sourcingvalue.name" +
		" FROM product_sourcingvalue INNER JOIN sourcingvalue" +
		" ON product_sourcingvalue.sourcingvalue_id = sourcingvalue.id" +
		" WHERE product_sourcingvalue.product_id = " + strconv.Itoa(productIdPK)
	selectQ, err := txn.Query(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
	} else {
		for selectQ.Next() {
			productProperty := &ProductProperty{}
			err := selectQ.Scan(&productProperty.ProductId, &productProperty.PropertyId, &productProperty.PropertyName)
			if err != nil {
				logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
					constants.MySQLSelectScanErrorMessage, err.Error())
			} else {
				result = append(result, productProperty)
			}
		}
	}
	return result
}

func SelectIngredientNameFromProductIngredientByProductIdPK(productIdPK int) []string {
	/*
		To take product_id (primary key of product table) and select ingredient name
	*/
	funcName := "SelectIngredientNameFromProductIngredientByProductIdPK"
	result := make([]string, 0)
	query := "SELECT ingredient.name FROM product_ingredient INNER JOIN ingredient ON" +
		" product_ingredient.ingredient_id = ingredient.id" +
		" WHERE product_ingredient.product_id = " + strconv.Itoa(productIdPK)
	selectQ, err := mysqlc.MySqlDB.Query(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
	} else {
		var name string
		for selectQ.Next() {
			err := selectQ.Scan(&name)
			if err != nil {
				logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
					constants.MySQLSelectScanErrorMessage, err.Error())
			}
			result = append(result, name)
		}
	}
	return result
}

func SelectFromProductIngredientByProductIdPK(txn *sql.Tx, productIdPK int) []*ProductProperty {
	/*
		To take product_id (primary key of product table) and select product id, ingredient id, ingredient name
	*/
	funcName := "SelectFromProductIngredientByProductIdPK"
	result := make([]*ProductProperty, 0)
	query := "SELECT product_ingredient.product_id, product_ingredient.ingredient_id, ingredient.name" +
		" FROM product_ingredient INNER JOIN ingredient" +
		" ON product_ingredient.ingredient_id = ingredient.id" +
		" WHERE product_ingredient.product_id = " + strconv.Itoa(productIdPK)
	selectQ, err := txn.Query(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
	} else {
		productProperty := &ProductProperty{}
		for selectQ.Next() {
			err := selectQ.Scan(&productProperty.ProductId, &productProperty.PropertyId, &productProperty.PropertyName)
			if err != nil {
				logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
					constants.MySQLSelectScanErrorMessage, err.Error())
			} else {
				result = append(result, productProperty)
			}
		}
	}
	return result
}
