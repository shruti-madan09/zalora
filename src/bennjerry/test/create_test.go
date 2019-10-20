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

func TestCreateDataUnAuthorized(t *testing.T) {
	/*
		Test Scenario: Calling create api without auth token in header
		Expectation: Response with status code 401
	*/
	mysqlc.Init()
	logger.Init()
	route := gin.Default()
	route.POST("/bennjerry/", authenticator.IsAuthorized, bennjerry.CreateData)

	// Creating mock request for create functionality
	postData := []byte(`{
			"name": "Name of Ice Cream",
			"productId": "test123"
	}`)
	data := url.Values{}
	data.Set("data", string(postData))
	req, reqErr := http.NewRequest(http.MethodPost, "/bennjerry/", bytes.NewBufferString(data.Encode()))
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

func TestCreateDataEmptyRequest(t *testing.T) {
	/*
		Testing Scenario: Calling create api with empty post form data in request
		Expectation: Appropriate error response
	*/
	mysqlc.Init()
	logger.Init()
	route := gin.Default()
	route.POST("/bennjerry/", authenticator.IsAuthorized, bennjerry.CreateData)

	// Creating mock request for create functionality
	data := url.Values{}
	req, reqErr := http.NewRequest(http.MethodPost, "/bennjerry/", bytes.NewBufferString(data.Encode()))
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
			t.Fatalf("Expected response {success: false, id:0, message: %s} but got {success: %v, id: %d,"+
				" message: %s}\n", constants.RequestInvalidErrorMessage, resp.Success, resp.Id, resp.Message)
		}
	}
	mysqlc.DBClosing()
}

func TestCreateDataUnReadableRequest(t *testing.T) {
	/*
		Testing Scenario: Calling create api with invalid json in post form key: `data`
		Expectation: Appropriate error response
	*/
	mysqlc.Init()
	logger.Init()
	route := gin.Default()
	route.POST("/bennjerry/", authenticator.IsAuthorized, bennjerry.CreateData)

	// Creating mock request for create functionality
	postData := []byte(`{
		"name": "Name of Ice Cream""productId": "test123"
	}`)
	data := url.Values{}
	data.Set("data", string(postData))
	req, reqErr := http.NewRequest(http.MethodPost, "/bennjerry/", bytes.NewBufferString(data.Encode()))
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
		if resp.Success || resp.Id != 0 {
			t.Fatalf("Expected response {success: false, id: 0, message: unMarshall error message} but got"+
				" {success: %v, id: %d, message: %s}\n", resp.Success, resp.Id, resp.Message)
		}
	}
	mysqlc.DBClosing()
}

func TestCreateData(t *testing.T) {
	/*
		Testing Scenario: Calling create api with correct request data and headers
		Expectation: Success response and entry in DB
		** this DB entry will be used to test read, update, delete apis and will eventually will be cleaned up from DB
	*/
	mysqlc.Init()
	logger.Init()
	route := gin.Default()
	route.POST("/bennjerry/", authenticator.IsAuthorized, bennjerry.CreateData)

	// Creating mock request for create functionality
	postData := []byte(`{
			"productId": "test123",
			"name": "Name of Ice Cream",
			"image_closed": "Link of closed image",
			"image_open": "Link of open image",
			"description": "Description of Ice Cream",
			"story": "Story of Ice Cream",
			"sourcing_values": ["List", "of", "sourcing", "values"],
			"ingredients": ["List", "of", "ingredients"],
			"allergy_info": "Allergy related information",
			"dietary_certifications": "Name of dietary certifications"
	}`)
	data := url.Values{}
	data.Set("data", string(postData))
	req, reqErr := http.NewRequest(http.MethodPost, "/bennjerry/", bytes.NewBufferString(data.Encode()))
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
		if !resp.Success || resp.Id == 0 || resp.Message != constants.CreateSuccessMessage {
			t.Fatalf("Expected response {success: true, id: non-zero, message: %s} but got"+
				" {sucess: %v, id: %d, message: %s}\n", constants.CreateSuccessMessage, resp.Success, resp.Id,
				resp.Message)
		}
	}
	mysqlc.DBClosing()
}

func TestCreateDataDuplicateRequest(t *testing.T) {
	/*
		Testing Scenario: Calling create api with duplicate request as previous scenario
		Expectation: Appropriate error response
		** product_id in table product has a unique constraint to avoid duplicate data
	*/
	mysqlc.Init()
	logger.Init()
	route := gin.Default()
	route.POST("/bennjerry/", authenticator.IsAuthorized, bennjerry.CreateData)

	// Creating mock request for create functionality
	postData := []byte(`{
			"productId": "test123",
			"name": "Name of Ice Cream",
			"image_closed": "Link of closed image",
			"image_open": "Link of open image",
			"description": "Description of Ice Cream",
			"story": "Story of Ice Cream",
			"sourcing_values": ["List", "of", "sourcing", "values"],
			"ingredients": ["List", "of", "ingredients"],
			"allergy_info": "Allergy related information",
			"dietary_certifications": "Name of dietary certifications"
	}`)
	data := url.Values{}
	data.Set("data", string(postData))
	req, reqErr := http.NewRequest(http.MethodPost, "/bennjerry/", bytes.NewBufferString(data.Encode()))
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
		t.Fatalf("Expected to get status %d but got %d\n", http.StatusOK, recorder.Code)
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
		if resp.Success || resp.Id != 0 || resp.Message != constants.GenericErrorMessage {
			t.Fatalf("Expected response {success: false, id: 0, message: %s} but got"+
				" {sucess: %v, id: %d, message: %s}\n", constants.GenericErrorMessage, resp.Success, resp.Id,
				resp.Message)
		}
	}
	mysqlc.DBClosing()
}
