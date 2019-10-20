package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"authenticator"
	"bennjerry"
	"bennjerry/structs"
	"constants"
	"logger"
	"mysqlc"
	"utils"
)

func TestReadDataUnAuthorized(t *testing.T) {
	/*
		Test Scenario: Calling read api without auth token in header
		Expectation: Response with status code 401
	*/
	mysqlc.Init()
	logger.Init()
	route := gin.Default()
	route.GET("/bennjerry/:product_id/", authenticator.IsAuthorized, bennjerry.ReadData)

	// Creating mock request for read functionality
	req, reqErr := http.NewRequest(http.MethodGet, "/bennjerry/test123/", nil)
	if reqErr != nil {
		t.Fatalf("Couldn't create request: %v\n", reqErr)
	}

	// Creating a response recorder to inspect the response
	recorder := httptest.NewRecorder()

	// Performing the request
	route.ServeHTTP(recorder, req)

	// Checking to see if the response was what you expected
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status code %d but got %d\n", http.StatusUnauthorized, recorder.Code)
	}
	mysqlc.DBClosing()
}

func TestReadDataNoRecordFound(t *testing.T) {
	/*
		Testing Scenario: Calling read api with product_id that doesn't exist in DB
		Expectation: Appropriate error response
	*/
	mysqlc.Init()
	logger.Init()
	route := gin.Default()
	route.GET("/bennjerry/:product_id/", authenticator.IsAuthorized, bennjerry.ReadData)

	// Creating mock request for read functionality
	req, reqErr := http.NewRequest(http.MethodGet, "/bennjerry/test456/", nil)
	if reqErr != nil {
		t.Fatalf("Couldn't create request: %v\n", reqErr)
	}
	// Generating token for authorization
	jwtToken, tokenErr := authenticator.GenerateJWT()
	if tokenErr != nil {
		t.Fatalf("Couldn't generate token %s\n", tokenErr.Error())
	}
	req.Header.Add(constants.JWTTokenKeyNameInHeader, jwtToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Creating a response recorder to inspect the response
	recorder := httptest.NewRecorder()

	// Performing the request
	route.ServeHTTP(recorder, req)

	// Checking to see if the response was what you expected
	if recorder.Code != http.StatusOK {
		t.Fatalf("Expected status code %d but got %d\n", http.StatusOK, recorder.Code)
	} else {
		respBytes, respErr := ioutil.ReadAll(recorder.Body)
		if respErr != nil {
			t.Fatalf("Error while reading response %s\n", respErr.Error())
		}
		resp := &structs.ReadResponse{}
		unMarshallErr := json.Unmarshal(respBytes, resp)
		if unMarshallErr != nil {
			t.Fatalf("Error while parsing response %s\n", unMarshallErr.Error())
		}
		if resp.Success || resp.Data != nil || resp.Message != constants.NoRecordsFoundMessage {
			t.Fatalf("Expected response {success: false, data: nil, message: %s} but got"+
				" {sucess: %v, data: %v, message: %s}\n", constants.NoRecordsFoundMessage, resp.Success, resp.Data,
				resp.Message)
		}
	}
	mysqlc.DBClosing()
}

func TestReadData(t *testing.T) {
	/*
		Testing Scenario: Calling read api with correct url params and request headers
		Expectation: Success response with complete information of requested product_id
	*/
	mysqlc.Init()
	logger.Init()
	route := gin.Default()
	route.GET("/bennjerry/:product_id/", authenticator.IsAuthorized, bennjerry.ReadData)

	// Creating mock request for read functionality
	req, reqErr := http.NewRequest(http.MethodGet, "/bennjerry/test123/", nil)
	if reqErr != nil {
		t.Fatalf("Couldn't create request: %v\n", reqErr)
	}
	// Generating token for authorization
	jwtToken, tokenErr := authenticator.GenerateJWT()
	if tokenErr != nil {
		t.Fatalf("Couldn't generate token %s\n", tokenErr.Error())
	}
	req.Header.Add(constants.JWTTokenKeyNameInHeader, jwtToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Creating a response recorder to inspect the response
	recorder := httptest.NewRecorder()

	// Performing the request
	route.ServeHTTP(recorder, req)

	// Checking to see if the response was what you expected
	if recorder.Code != http.StatusOK {
		t.Fatalf("Expected status code %d but got %d\n", http.StatusOK, recorder.Code)
	} else {
		respBytes, respErr := ioutil.ReadAll(recorder.Body)
		if respErr != nil {
			t.Fatalf("Error while reading response %s\n", respErr.Error())
		}
		resp := &structs.ReadResponse{}
		unMarshallErr := json.Unmarshal(respBytes, resp)
		if unMarshallErr != nil {
			t.Fatalf("Error while parsing response %s\n", unMarshallErr.Error())
		}

		// Comparing received response with expected response
		expectedIceCreamData := &structs.IceCreamDataStruct{
			ProductId:             "test123",
			Name:                  "Name of Ice Cream",
			ImageClosed:           "Link of closed image",
			ImageOpened:           "Link of open image",
			Description:           "Description of Ice Cream",
			Story:                 "Story of Ice Cream",
			SourcingValues:        []string{"List", "of", "sourcing", "values"},
			Ingredients:           []string{"List", "of", "ingredients"},
			AllergyInfo:           "Allergy related information",
			DietaryCertifications: "Name of dietary certifications",
		}
		isDataMatching := resp.Data != nil && resp.Data.ProductId == expectedIceCreamData.ProductId &&
			resp.Data.Name == expectedIceCreamData.Name && resp.Data.ImageClosed == expectedIceCreamData.ImageClosed &&
			resp.Data.ImageOpened == expectedIceCreamData.ImageOpened &&
			resp.Data.Description == expectedIceCreamData.Description &&
			resp.Data.Story == expectedIceCreamData.Story &&
			utils.ListOfStringCompare(resp.Data.SourcingValues, expectedIceCreamData.SourcingValues) &&
			utils.ListOfStringCompare(resp.Data.Ingredients, expectedIceCreamData.Ingredients) &&
			resp.Data.AllergyInfo == expectedIceCreamData.AllergyInfo &&
			resp.Data.DietaryCertifications == expectedIceCreamData.DietaryCertifications

		if !resp.Success || !isDataMatching || resp.Message != constants.ReadSuccessMessage {
			t.Fatalf("Expected response {success: true, data: %v, message: %s} but got"+
				" {success: %v, data: %v, message: %s}\n", expectedIceCreamData, constants.ReadSuccessMessage,
				resp.Success, resp.Data, resp.Message)
		}
	}
	mysqlc.DBClosing()
}
