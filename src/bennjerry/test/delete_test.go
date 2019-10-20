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
)

func TestDeleteDataUnAuthorized(t *testing.T) {
	/*
		Test Scenario: Calling delete api without auth token in header
		Expectation: Response with status code 401
	*/
	mysqlc.Init()
	logger.Init()
	route := gin.Default()
	route.DELETE("/bennjerry/:product_id/", authenticator.IsAuthorized, bennjerry.DeleteData)

	// Creating mock request for read functionality
	req, reqErr := http.NewRequest(http.MethodDelete, "/bennjerry/test123/", nil)
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

func TestDeleteDataNoRecordFound(t *testing.T) {
	/*
		Testing Scenario: Calling delete api with product_id that doesn't exist in DB
		Expectation: Appropriate error response
	*/
	mysqlc.Init()
	logger.Init()
	route := gin.Default()
	route.DELETE("/bennjerry/:product_id/", authenticator.IsAuthorized, bennjerry.DeleteData)

	// Creating mock request for delete functionality
	req, reqErr := http.NewRequest(http.MethodDelete, "/bennjerry/test456/", nil)
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
		resp := &structs.CreateUpdateDeleteResponse{}
		unMarshallErr := json.Unmarshal(respBytes, resp)
		if unMarshallErr != nil {
			t.Fatalf("Error while parsing response %s\n", unMarshallErr.Error())
		}
		if resp.Success || resp.Id != 0 || resp.Message != constants.NoRecordsFoundMessage {
			t.Fatalf("Expected response {success: false, id: 0, message: %s} but got"+
				" {success: %v, id: %d, message: %s}\n", constants.NoRecordsFoundMessage, resp.Success, resp.Id,
				resp.Message)
		}
	}
	mysqlc.DBClosing()
}

func TestSoftDeleteData(t *testing.T) {
	/*
		Testing Scenario: Calling soft delete api with correct url params, query params and request headers
		Expectation: Success response and record marked as inactive in DB
	*/
	mysqlc.Init()
	logger.Init()
	route := gin.Default()
	route.DELETE("/bennjerry/:product_id/", authenticator.IsAuthorized, bennjerry.DeleteData)

	// Creating mock request for delete functionality
	req, reqErr := http.NewRequest(http.MethodDelete, "/bennjerry/test123/", nil)
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
		resp := &structs.CreateUpdateDeleteResponse{}
		unMarshallErr := json.Unmarshal(respBytes, resp)
		if unMarshallErr != nil {
			t.Fatalf("Error while parsing response %s\n", unMarshallErr.Error())
		}
		if !resp.Success || resp.Id == 0 || resp.Message != constants.SoftDeleteSuccessMessage {
			t.Fatalf("Expected response {success: true, id: non-zero, message: %s} but got"+
				" {sucess: %v, id: %d, message: %s}\n", constants.SoftDeleteSuccessMessage, resp.Success, resp.Id,
				resp.Message)
		}
	}
	mysqlc.DBClosing()
}

func TestDeleteData(t *testing.T) {
	/*
		Testing Scenario: Calling delete api with correct url params, query params and request headers
		Expectation: Success response and record permanently deleted from DB
	*/
	mysqlc.Init()
	logger.Init()
	route := gin.Default()
	route.DELETE("/bennjerry/:product_id/", authenticator.IsAuthorized, bennjerry.DeleteData)

	// Creating mock request for delete functionality
	req, reqErr := http.NewRequest(http.MethodDelete, "/bennjerry/test123/?permanent=1", nil)
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
		resp := &structs.CreateUpdateDeleteResponse{}
		unMarshallErr := json.Unmarshal(respBytes, resp)
		if unMarshallErr != nil {
			t.Fatalf("Error while parsing response %s\n", unMarshallErr.Error())
		}
		if !resp.Success || resp.Id == 0 || resp.Message != constants.PermanentDeleteSuccessMessage {
			t.Fatalf("Expected response {success: true, id: non-zero, message: %s} but got"+
				" {sucess: %v, id: %d, message: %s}\n", constants.PermanentDeleteSuccessMessage, resp.Success, resp.Id,
				resp.Message)
		}
	}
	mysqlc.DBClosing()
}
