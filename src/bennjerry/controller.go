package bennjerry

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"

	"bennjerry/model"
	"bennjerry/structs"
	"constants"
	"logger"
	"utils"
)

func CreateData(ginContext *gin.Context) {
	/*
		To save information of a new ice cream product
		Sample Url: "http://host/bennjerry/"
		Request Method: POST
		Request Data:
		{
			"data": {
				"productId": "123",
				"name": "Name of Ice Cream",
				"description": "Description of Ice Cream",
				"story": "Story of Ice Cream",
				"image_closed": "Link of closed image",
				"image_open": "Link of open image",
				"sourcing_values": ["List", "of", "sourcing", "values"],
				"ingredients": ["List", "of", "ingredients"],
				"allergy_info": "Allergy related information",
				"dietary_certifications": "Name of dietary certifications"
			}
		}
		Response Data:
		{
			"message": "Success/Error message",
			"success": true/false
			"id": 12/0,
		}
	*/
	var (
		isAuthorized  bool
		iceCreamData  *structs.IceCreamDataStruct
		response      *structs.CreateUpdateDeleteResponse
		responseBytes []byte
		responseErr   error
		logIdentifier = "bennjerry.CreateData"
	)

	// defer block to recover and log if any error occurs
	defer func() {
		if r := recover(); r != nil {
			logger.ZaloraStatsLogger.Error(constants.BenNJerryLogBucketName, logIdentifier,
				constants.GenericErrorMessage, string(debug.Stack()))
			serializer := utils.GetSerializer(constants.JsonSerializerType)
			serializer.ReturnError(ginContext, http.StatusInternalServerError, constants.GenericErrorMessage)
			return
		}
	}()

	//If request is authorized, is_authorized = 1 would've been dropped in ginContext object by the auth middleware
	if isAuth, isAuthExists := ginContext.Get(constants.IsAuthorizedKeyName); isAuthExists {
		isAuthorized = isAuth.(int) == 1
	}
	if !isAuthorized {
		// Returning error response if request is not authorized
		serializer := utils.GetSerializer(constants.JsonSerializerType)
		serializer.ReturnUnAuthorized(ginContext, constants.UnAuthorizedErrorMessage)
		return
	}

	postData := ginContext.DefaultPostForm("data", "{}")
	// converting post form data to structure
	umMarshalErr := json.Unmarshal([]byte(postData), &iceCreamData)
	if umMarshalErr != nil {
		logger.ZaloraStatsLogger.Error(constants.BenNJerryLogBucketName, logIdentifier, constants.UnMarshalErrorString,
			umMarshalErr.Error())
		response = &structs.CreateUpdateDeleteResponse{
			Message: umMarshalErr.Error(),
		}
	} else if iceCreamData == nil || iceCreamData.ProductId == "" {
		response = &structs.CreateUpdateDeleteResponse{
			Message: constants.RequestInvalidErrorMessage,
		}
	} else {
		// Calling function to execute queries in an atomic transaction
		idList, success := model.InsertRecord([]*structs.IceCreamDataStruct{iceCreamData})
		if success && len(idList) > 0 {
			response = &structs.CreateUpdateDeleteResponse{
				Success: true,
				Message: constants.CreateSuccessMessage,
				Id:      idList[0],
			}
		} else {
			response = &structs.CreateUpdateDeleteResponse{
				Message: constants.GenericErrorMessage,
			}
		}
	}
	// converting response structure to []byte
	responseBytes, responseErr = json.Marshal(response)
	// calling common util function to send json response (utils/common.go)
	serializer := utils.GetSerializer(constants.JsonSerializerType)
	if responseErr != nil {
		logger.ZaloraStatsLogger.Error(constants.BenNJerryLogBucketName, logIdentifier,
			constants.JsonSerializationErrorMessage, responseErr.Error())
		serializer.ReturnError(ginContext, http.StatusInternalServerError,
			constants.JsonSerializationErrorMessage+" %v\"", responseErr)
	}
	serializer.ReturnOk(ginContext, responseBytes)
}

func ReadData(ginContext *gin.Context) {
	/*
		To fetch information of an ice cream by providing product_id
		Sample Url: "http://host/bennjerry/2190/"
		Request Method: GET
		Request Data: product_id to be provided in the url, e.g. 2190 in sample url
		Response Data:
		{
			"message": "Success/Error message",
			"success": true / false,
			"data": {
				"productId": "123",
				"name": "Name of Ice Cream",
				"description": "Description of Ice Cream",
				"story": "Story of Ice Cream",
				"image_closed": "Link of closed image",
				"image_open": "Link of open image",
				"sourcing_values": ["List", "of", "sourcing", "values"],
				"ingredients": ["List", "of", "ingredients"],
				"allergy_info": "Allergy related information",
				"dietary_certifications": "Name of dietary certifications"
			}
		}
	*/
	var (
		isAuthorized  bool
		response      *structs.ReadResponse
		responseBytes []byte
		responseErr   error
		logIdentifier = "bennjerry.ReadData"
	)

	// defer block to recover and log if any error occurs
	defer func() {
		if r := recover(); r != nil {
			logger.ZaloraStatsLogger.Error(constants.BenNJerryLogBucketName, logIdentifier,
				constants.GenericErrorMessage, string(debug.Stack()))
			serializer := utils.GetSerializer(constants.JsonSerializerType)
			serializer.ReturnError(ginContext, http.StatusInternalServerError, constants.GenericErrorMessage)
			return
		}
	}()

	// If request is authorized, is_authorized = 1 would've been dropped in ginContext object by the auth middleware
	if isAuth, isAuthExists := ginContext.Get(constants.IsAuthorizedKeyName); isAuthExists {
		isAuthorized = isAuth.(int) == 1
	}
	if !isAuthorized {
		// Returning error response if request is not authorized
		serializer := utils.GetSerializer(constants.JsonSerializerType)
		serializer.ReturnUnAuthorized(ginContext, constants.UnAuthorizedErrorMessage)
		return
	}

	productId := ginContext.Params.ByName("product_id")
	// fetching data from product table using product id
	// success: false, if some error occurs while running the query
	// success: true, productData: {}, if requested product_id is not found or is inactive
	productData, success := model.SelectFromProductByProductId(productId)
	if !success {
		response = &structs.ReadResponse{
			Message: constants.GenericErrorMessage,
		}
	} else if productData.ProductId == "" {
		response = &structs.ReadResponse{
			Message: constants.NoRecordsFoundMessage,
		}
	} else {
		response = &structs.ReadResponse{
			Success: true,
			Message: constants.ReadSuccessMessage,
		}
		response.Data = &structs.IceCreamDataStruct{
			Id:          productData.Id,
			ProductId:   productData.ProductId,
			Name:        productData.Name,
			Description: productData.Description,
			Story:       productData.Story,
			ImageClosed: productData.ImageClosed,
			ImageOpened: productData.ImageOpened,
			AllergyInfo: productData.Allergy,
		}
		// Id of a Dietary Certification is in product table as a foreign key
		// Using the same to fetch it's name from dietarycertification table
		if productData.DietaryCertificationId != 0 {
			dietaryCertification := model.SelectFromDietaryCertificationById(productData.DietaryCertificationId)
			if dietaryCertification != nil {
				response.Data.DietaryCertifications = dietaryCertification.Name
			}
		}
		// Fetching list of sourcing values from relation table of product and sourcing value
		response.Data.SourcingValues = model.SelectSourcingValueNameByProductIdPK(productData.Id)
		// Fetching list of ingredients from relation table of product and ingredient
		response.Data.Ingredients = model.SelectIngredientNameFromProductIngredientByProductIdPK(productData.Id)
	}
	// Converting response structure to []byte
	responseBytes, responseErr = json.Marshal(response)
	// Calling common util function to send json response (utils/common.go)
	serializer := utils.GetSerializer(constants.JsonSerializerType)
	if responseErr != nil {
		logger.ZaloraStatsLogger.Error(constants.BenNJerryLogBucketName, logIdentifier,
			constants.JsonSerializationErrorMessage, responseErr.Error())
		serializer.ReturnError(ginContext, http.StatusInternalServerError,
			constants.JsonSerializationErrorMessage+" %v\"", responseErr)
	}
	serializer.ReturnOk(ginContext, responseBytes)
}

func UpdateData(ginContext *gin.Context) {
	/*
		To update information of an existing ice cream product by providing product_id
		Sample Url: "http://host/bennjerry/2190/"
		Request Method: PUT
		Request Data: product_id to be provided in the url, e.g. 2190 in sample url
		{
			"data": {
				"name": "Name of Ice Cream",
				"story": "Story of Ice Cream",
				"image_closed": "Link of closed image",
				"sourcing_values": ["List", "of", "sourcing", "values"],
				"allergy_info": "Allergy related information",
				"dietary_certifications": "Name of dietary certifications"
			},
			"fields": "name,story,image_closed,sourcing_values,allergy_info,dietary_certifications"
		}
		Response Data:
		{
			"message": "Success/Error message",
			"success": true/false
			"id": 12/0,
		}
	*/
	var (
		isAuthorized  bool
		iceCreamData  *structs.IceCreamDataStruct
		response      *structs.CreateUpdateDeleteResponse
		responseBytes []byte
		responseErr   error
		logIdentifier = "bennjerry.UpdateData"
	)

	// defer block to recover and log if any error occurs
	defer func() {
		if r := recover(); r != nil {
			logger.ZaloraStatsLogger.Error(constants.BenNJerryLogBucketName, logIdentifier,
				constants.GenericErrorMessage, string(debug.Stack()))
			serializer := utils.GetSerializer(constants.JsonSerializerType)
			serializer.ReturnError(ginContext, http.StatusInternalServerError, constants.GenericErrorMessage)
			return
		}
	}()

	//If request is authorized, is_authorized = 1 would've been dropped in ginContext object by the auth middleware
	if isAuth, isAuthExists := ginContext.Get(constants.IsAuthorizedKeyName); isAuthExists {
		isAuthorized = isAuth.(int) == 1
	}
	if !isAuthorized {
		// Returning error response if request is not authorized
		serializer := utils.GetSerializer(constants.JsonSerializerType)
		serializer.ReturnUnAuthorized(ginContext, constants.UnAuthorizedErrorMessage)
		return
	}

	productId := ginContext.Params.ByName("product_id")
	// fetching id (primary key) of ice cream product using product_id
	// success: false, if some error occurs while running the query
	// success: true, id: 0, if requested product_id is not found
	id, success := model.SelectIdFromProductByProductId(productId)
	if !success {
		response = &structs.CreateUpdateDeleteResponse{
			Message: constants.GenericErrorMessage,
		}
	} else if id == 0 {
		response = &structs.CreateUpdateDeleteResponse{
			Message: constants.NoRecordsFoundMessage,
		}
	} else {
		postData := ginContext.DefaultPostForm("data", "{}")
		postFields := ginContext.DefaultPostForm("fields", "")
		// converting post form data to structure
		umMarshalErr := json.Unmarshal([]byte(postData), &iceCreamData)
		if umMarshalErr != nil {
			logger.ZaloraStatsLogger.Error(constants.BenNJerryLogBucketName, logIdentifier,
				constants.UnMarshalErrorString, umMarshalErr.Error())
			response = &structs.CreateUpdateDeleteResponse{
				Message: umMarshalErr.Error(),
			}
		} else if iceCreamData == nil || iceCreamData.ProductId == "" || postFields == "" {
			response = &structs.CreateUpdateDeleteResponse{
				Message: constants.RequestInvalidErrorMessage,
			}
		} else {
			// The user may want to update only specific properties of an ice cream product
			// UnMarshalling of postData sets default values for the missing properties in request
			// These missing properties of the product can be overwritten in DB with default values of their datatype
			// Therefore information regarding which fields to be updated is mandatory in the request
			// 'fields' should be a comma separated string of names of all fields to be updated
			// Removing all spaces from postFields
			postFields = strings.Replace(postFields, " ", "", -1)
			// Splitting postFields on comma and converting it to map {"string":bool}
			fieldMap := utils.ListToMap(strings.Split(postFields, ","))
			// Calling function to execute queries in an atomic transaction
			success := model.UpdateRecord(id, iceCreamData, fieldMap)
			if success {
				response = &structs.CreateUpdateDeleteResponse{
					Success: true,
					Message: constants.UpdateSuccessMessage,
					Id:      id,
				}
			} else {
				response = &structs.CreateUpdateDeleteResponse{
					Message: constants.GenericErrorMessage,
				}
			}
		}
	}
	// converting response structure to []byte
	responseBytes, responseErr = json.Marshal(response)
	// Calling common util function to send json response (utils/common.go)
	serializer := utils.GetSerializer(constants.JsonSerializerType)
	if responseErr != nil {
		logger.ZaloraStatsLogger.Error(constants.BenNJerryLogBucketName, logIdentifier,
			constants.JsonSerializationErrorMessage, responseErr.Error())
		serializer.ReturnError(ginContext, http.StatusInternalServerError,
			constants.JsonSerializationErrorMessage+" %v\"", responseErr)
	}
	serializer.ReturnOk(ginContext, responseBytes)
}

func DeleteData(ginContext *gin.Context) {
	/*
		To delete information of an existing ice cream product by providing product_id
		Sample Url: "http://host/bennjerry/2190/" or "http://host/bennjerry/2190/?permanent=1"
		Request Method: DELETE
		Request Data: product_id to be provided in the url, e.g. 2190 in sample url
		URL Param: permanent=1, if user wants data to be permanently removed from DB
		Response Data:
		{
			"message": "Success/Error message",
			"success": true/false
			"id": 12/0,
		}
	*/
	var (
		isAuthorized  bool
		response      *structs.CreateUpdateDeleteResponse
		responseBytes []byte
		responseErr   error
		logIdentifier = "bennjerry.DeleteData"
	)

	// defer block to recover and log if any error occurs
	defer func() {
		if r := recover(); r != nil {
			logger.ZaloraStatsLogger.Error(constants.BenNJerryLogBucketName, logIdentifier,
				constants.GenericErrorMessage, string(debug.Stack()))
			serializer := utils.GetSerializer(constants.JsonSerializerType)
			serializer.ReturnError(ginContext, http.StatusInternalServerError, constants.GenericErrorMessage)
			return
		}
	}()

	// If request is authorized, is_authorized = 1 would've been dropped in ginContext object by the auth middleware
	if isAuth, isAuthExists := ginContext.Get(constants.IsAuthorizedKeyName); isAuthExists {
		isAuthorized = isAuth.(int) == 1
	}
	if !isAuthorized {
		// Returning error response if request is not authorized
		serializer := utils.GetSerializer(constants.JsonSerializerType)
		serializer.ReturnUnAuthorized(ginContext, constants.UnAuthorizedErrorMessage)
		return
	}

	productId := ginContext.Params.ByName("product_id")
	// If URL param permanent=1 is not present then record will only be soft deleted i.e. marked as inactive
	isPermanentDelete := ginContext.DefaultQuery("permanent", "0") == "1"
	if isPermanentDelete {
		// Primary key, id needs to be fetched because references for record in other tables need to be deleted first
		// success: false, if an error occurs while running the query
		// success: true, id: 0, if record is not found for the product_id
		id, success := model.SelectIdFromProductByProductId(productId)
		if !success {
			response = &structs.CreateUpdateDeleteResponse{
				Message: constants.GenericErrorMessage,
			}
		} else if id == 0 {
			response = &structs.CreateUpdateDeleteResponse{
				Message: constants.NoRecordsFoundMessage,
			}
		} else {
			// Calling function to execute queries in an atomic transaction
			success = model.DropRecord(id)
			if success {
				response = &structs.CreateUpdateDeleteResponse{
					Success: true,
					Message: constants.PermanentDeleteSuccessMessage,
					Id:      id,
				}
			} else {
				response = &structs.CreateUpdateDeleteResponse{
					Message: constants.GenericErrorMessage,
				}
			}
		}
	} else {
		// For soft deletion record will exist in DB but will be marked as inactive by setting column is_inactive = 1
		// In read operation, an ice cream product will be fetched only if it's not inactive
		// success: false, if some error occurs while running query
		// success: true, id: 0, if requested product_id is not found
		id, success := model.SoftDeleteFromProductByProductId(productId)
		if !success {
			response = &structs.CreateUpdateDeleteResponse{
				Message: constants.GenericErrorMessage,
			}
		} else if id == 0 {
			response = &structs.CreateUpdateDeleteResponse{
				Message: constants.NoRecordsFoundMessage,
			}
		} else {
			response = &structs.CreateUpdateDeleteResponse{
				Success: true,
				Message: constants.SoftDeleteSuccessMessage,
				Id:      id,
			}
		}
	}
	// Converting response structure to []byte
	responseBytes, responseErr = json.Marshal(response)
	// Calling common util function to send json response (utils/common.go)
	serializer := utils.GetSerializer(constants.JsonSerializerType)
	if responseErr != nil {
		logger.ZaloraStatsLogger.Error(constants.BenNJerryLogBucketName, logIdentifier,
			constants.JsonSerializationErrorMessage, responseErr.Error())
		serializer.ReturnError(ginContext, http.StatusInternalServerError,
			constants.JsonSerializationErrorMessage+" %v\"", responseErr)
	}
	serializer.ReturnOk(ginContext, responseBytes)
}
