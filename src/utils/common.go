package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"constants"
)

type Serializer struct {
	contentType string
}

func GetSerializer(serializerType string) *Serializer {
	/*
	To send content type based on serializer type (e.g. json, proto)
	Arguments: string
	Returns: *Serializer
	*/
	switch serializerType {
	default:
		return &Serializer{
			contentType: constants.JsonContentType,
		}
	}
}

func (serializer *Serializer) ReturnOk(ginContext *gin.Context, result []byte) {
	/*
	To send http success response with response as []byte
	*/
	if ginContext.IsAborted() {
		return
	}
	ginContext.Abort()
	ginContext.Data(http.StatusOK, serializer.contentType, result)
}

func (serializer *Serializer) ReturnError(ginContext *gin.Context, code int, sFmt string, v ...interface{}) {
	/*
		To send http error response with relevant error code and error response as []byte
	*/
	if ginContext.IsAborted() {
		return
	}
	ginContext.Abort()
	ginContext.Data(code, serializer.contentType, []byte(fmt.Sprintf("{\"error\":\"%s\"}",
		fmt.Sprintf(sFmt, v...))))
}

func (serializer *Serializer) ReturnUnAuthorized(ginContext *gin.Context, sFmt string, v ...interface{}) {
	/*
		To send http error response with relevant error code and error response as []byte
	*/
	if ginContext.IsAborted() {
		return
	}
	ginContext.Abort()
	ginContext.Data(http.StatusUnauthorized, serializer.contentType,
		[]byte(fmt.Sprintf("{\"error\":\"%s\"}", fmt.Sprintf(sFmt, v...))))
}

func ListToMap(list []string) map[string]bool {
	/*
	To iterate over a list of string and convert it to a map of string and boolean
	Arguments: []string
	Returns: map[string]bool
	 */
	result := make(map[string]bool)
	for _, each := range list {
		result[each] = true
	}
	return result
}

func ListOfStringCompare(sl1 []string, sl2 []string) bool {
	if len(sl1) != len(sl2) {
		return false
	}
	for index, value := range sl1 {
		if sl2[index] != value {
			return false
		}
	}
	return true
}