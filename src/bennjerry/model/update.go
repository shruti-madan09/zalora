package model

import (
	"bennjerry/structs"
	"database/sql"
	"logger"
	"mysqlc"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"constants"
)

func UpdateProductById(txn *sql.Tx, id int, iceCreamData *structs.IceCreamDataStruct, fieldsMap map[string]bool) bool {
	/*
	To take ice cream data and update fields in product table based on map {fieldName: true}
	 */
	funcName := "UpdateProductById"
	query := "UPDATE product SET"
	if _, exists := fieldsMap["name"]; exists {
		query += " name = '"+iceCreamData.Name+"',"
	}
	if _, exists := fieldsMap["description"]; exists {
		query += " description = '"+iceCreamData.Description+"',"
	}
	if _, exists := fieldsMap["story"]; exists {
		query += " story = '"+iceCreamData.Story+"',"
	}
	if _, exists := fieldsMap["image_closed"]; exists {
		query += " image_closed = '"+iceCreamData.ImageClosed+"',"
	}
	if _, exists := fieldsMap["image_open"]; exists {
		query += " image_open = '"+iceCreamData.ImageOpened+"',"
	}
	if _, exists := fieldsMap["allergy_info"]; exists {
		query += " allergy_info = '"+iceCreamData.AllergyInfo+"',"
	}
	if _, exists := fieldsMap["dietary_certifications"]; exists {
		if iceCreamData.DietaryCertifications != "" {
			nameMap := map[string]bool {
				iceCreamData.DietaryCertifications: true,
			}
			success := InsertIntoDietaryCertification(txn, nameMap)
			if !success {
				return false
			}
			dietaryCertificationId := SelectFromDietaryCertification(txn, []string{iceCreamData.DietaryCertifications})
			query += " dietary_certification_id = '"+strconv.Itoa(dietaryCertificationId[0].Id)+"',"
		} else {
			query += " dietary_certification_id = NULL,"
		}
	}
	query = strings.TrimSuffix(query, ",")
	query += "WHERE id = " + strconv.Itoa(id)
	_, err := txn.Exec(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
    }
    return true
}

func UpdateProductSourcingValueByProductIdPK(txn *sql.Tx, productIdPK int, nameMap map[string]bool) bool {
	/*
	Take product_id (primary key of product table) and map {name: true}
	and update data in product_sourcingvalue table
	 */
	productPropertyData := SelectFromProductSourcingValueByProductIdPK(txn, productIdPK)
	existingNameMap := make(map[string][]int)
	for _, productProperty := range productPropertyData {
		if exists := nameMap[productProperty.PropertyName]; !exists {
			success := DeleteFromProductSourcingValueById(txn, productProperty.ProductId, productProperty.PropertyId)
			if !success {
				return false
			}
		} else {
			existingNameMap[productProperty.PropertyName] = []int{productProperty.ProductId, productProperty.PropertyId}
		}
	}
    for name := range nameMap {
    	if idList, exists := existingNameMap[name]; !exists {
    		success := InsertToProductSourcingValueById(txn, idList[0], idList[1])
    		if !success {
    			return false
			}
    	}
    }
    return true
}

func UpdateProductIngredientByProductIdPK(txn *sql.Tx, productIdPK int, nameMap map[string]bool) bool {
	/*
		Take product_id (primary key of product table) and map {name: true}
		and update data in product_ingredient table
	*/
	productPropertyData := SelectFromProductIngredientByProductIdPK(txn, productIdPK)
	existingNameMap := make(map[string][]int)
	for _, productProperty := range productPropertyData {
		if _, exists := nameMap[productProperty.PropertyName]; !exists {
			success := DeleteFromProductIngredientById(txn, productProperty.ProductId, productProperty.PropertyId)
			if !success {
				return false
			}
		} else {
			existingNameMap[productProperty.PropertyName] = []int{productProperty.ProductId, productProperty.PropertyId}
		}
	}
    for name := range nameMap {
    	if idList, exists := existingNameMap[name]; !exists {
			success := InsertToProductIngredientById(txn, idList[0], idList[1])
			if !success {
				return false
			}
    	}
    }
    return true
}

func UpdateProductIsInActiveByProductId(productId string) (int, bool) {
	/*
		Take product_id and update is_inactive = 1 in product table
	*/
	funcName := "UpdateProductIsInActiveByProductId"
	query := "UPDATE product SET is_inactive = 1 WHERE product_id = '" + productId + "'"
	update, err := mysqlc.MySqlDB.Exec(query)
	if err != nil {
		logger.ZaloraStatsLogger.Error(constants.MySQLLogBucketName, logIdentifier+funcName,
			constants.MySQLQueryRunErrorMessage, err.Error())
	} else {
		id, _ := update.RowsAffected()
		return int(id), true
	}
	return 0, false
}
