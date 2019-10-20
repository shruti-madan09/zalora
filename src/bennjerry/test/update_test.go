package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"

	"authenticator"
	"bennjerry"
	"bennjerry/structs"
	"constants"
	"logger"
	"mysqlc"
)

func TestUpdateDataUnAuthorized(t *testing.T) {
	/*
		Test Scenario: Calling update api without auth token in header
		Expectation: Response with status code 401
	*/
	mysqlc.Init()
	logger.Init()
	route := gin.Default()
	route.PUT("/bennjerry/", authenticator.IsAuthorized, bennjerry.UpdateData)

	// Creating mock request for create functionality
	postData := []byte(`{
			"name": "New Name of Ice Cream",
			"productId": "test123",
			"description": "New Description of Ice Cream"
	}`)
	data := url.Values{}
	data.Set("data", string(postData))
	data.Set("fields", "name,description")
	req, reqErr := http.NewRequest(http.MethodPut, "/bennjerry/", bytes.NewBufferString(data.Encode()))
	if reqErr != nil {
		t.Fatalf("Couldn't create request: %v\n", reqErr)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

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

func TestUpdateDataEmptyRequest(t *testing.T) {
	/*
		Testing Scenario: Calling update api with empty post form data in request
		Expectation: Appropriate error response
	*/
	mysqlc.Init()
	logger.Init()
	route := gin.Default()
	route.PUT("/bennjerry/:product_id/", authenticator.IsAuthorized, bennjerry.UpdateData)

	// Creating mock request for create functionality
	data := url.Values{}
	req, reqErr := http.NewRequest(http.MethodPut, "/bennjerry/test123/", bytes.NewBufferString(data.Encode()))
	if reqErr != nil {
		t.Fatalf("Couldn't create request: %v\n", reqErr)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
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
		if resp.Success || resp.Id != 0 || resp.Message != constants.RequestInvalidErrorMessage {
			t.Fatalf("Expected response {success: false, id: 0, message: %s} but got {success: %v, id: %d," +
				" message: %s}\n", constants.RequestInvalidErrorMessage, resp.Success, resp.Id, resp.Message)
		}
	}
	mysqlc.DBClosing()
}

func TestUpdateDataUnReadableRequest(t *testing.T) {
	/*
		Testing Scenario: Calling update api with invalid json in post form key: `data`
		Expectation: Appropriate error response
	*/
	mysqlc.Init()
	logger.Init()
	route := gin.Default()
	route.PUT("/bennjerry/:product_id/", authenticator.IsAuthorized, bennjerry.UpdateData)

	// Creating mock request for create functionality
	postData := []byte(`{
		"name": "New Name of Ice Cream""productId": "test123"
	}`)
	data := url.Values{}
	data.Set("data", string(postData))
	data.Set("fields", "name")
	req, reqErr := http.NewRequest(http.MethodPut, "/bennjerry/test123/", bytes.NewBufferString(data.Encode()))
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
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	// Creating a response recorder to inspect the response
	recorder := httptest.NewRecorder()

	// Performing the request
	route.ServeHTTP(recorder, req)

	// Checking to see if the response was what you expected
	if recorder.Code != http.StatusOK {
		t.Fatalf("Expected status %d but got %d\n", http.StatusOK, recorder.Code)
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
		if resp.Success || resp.Id != 0 {
			t.Fatalf("Expected response {success: false, id: 0, message: unMarshall error message} but got" +
				" {success: %v, id: %d, message: %s}\n", resp.Success, resp.Id, resp.Message)
		}
	}
	mysqlc.DBClosing()
}

func TestUpdateDataNoRecordFound(t *testing.T) {
	/*
		Testing Scenario: Calling update api with product_id that doesn't exist in DB
		Expectation: Appropriate error response
	 */
	mysqlc.Init()
	logger.Init()
	route := gin.Default()
	route.PUT("/bennjerry/:product_id/", authenticator.IsAuthorized, bennjerry.UpdateData)

	// Creating mock request for read functionality
	postData := []byte(`{
			"productId": "test456"
			"name": " New Name of Ice Cream",
			"description": "New Description of Ice Cream",
	}`)
	data := url.Values{}
	data.Set("data", string(postData))
	data.Set("fields", "name,description")
	req, reqErr := http.NewRequest(http.MethodPut, "/bennjerry/test456/", bytes.NewBufferString(data.Encode()))
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
			t.Fatalf("Expected response {success: false, id: 0, message: %s} but got" +
				" {sucess: %v, data: %v, message: %s}\n", constants.NoRecordsFoundMessage, resp.Success, resp.Id,
				resp.Message)
		}
	}
	mysqlc.DBClosing()
}

func TestUpdateData(t *testing.T) {
	/*
		Testing Scenario: Calling create api with correct request data and headers
		Expectation: Success response and new data in DB
	*/
	mysqlc.Init()
	logger.Init()
	route := gin.Default()
	route.PUT("/bennjerry/:product_id/", authenticator.IsAuthorized, bennjerry.UpdateData)

	// Creating mock request for create functionality
	postData := []byte(`{
			"productId": "test123",
			"name": " New Name of Ice Cream",
			"description": "New Description of Ice Cream"
	}`)
	data := url.Values{}
	data.Set("data", string(postData))
	data.Set("fields", "name,description")
	req, reqErr := http.NewRequest(http.MethodPut, "/bennjerry/test123/", bytes.NewBufferString(data.Encode()))
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
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

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
		if !resp.Success || resp.Id == 0 || resp.Message != constants.UpdateSuccessMessage {
			t.Fatalf("Expected response {success: true, id: non-zero, message: %s} but got" +
				" {sucess: %v, id: %d, message: %s}\n", constants.UpdateSuccessMessage, resp.Success, resp.Id,
				resp.Message)
		}
	}
	mysqlc.DBClosing()
}
