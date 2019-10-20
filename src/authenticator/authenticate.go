package authenticator

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"constants"
	"logger"
)

var MySigningKey = []byte("zaloracaptainjacksparrowsayshiassignment")

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["client"] = "Zalora Client"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	tokenString, err := token.SignedString(MySigningKey)
	if err != nil {
		fmt.Println("Error while generating token: ", err.Error())
		return "", err
	}
	return tokenString, nil
}

func IsAuthorized(ginContext *gin.Context) {
	logIdentifier := "authenticate.isAuthorized"
	requestedToken := ginContext.Request.Header.Get(constants.JWTTokenKeyNameInHeader)
	if requestedToken != "" {
		token, err := jwt.Parse(requestedToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("there was an error")
			}
			return MySigningKey, nil
		})
		if err != nil {
			logger.ZaloraStatsLogger.Error(constants.AuthLogBucketName, logIdentifier,
				constants.JWTTokenParseErrorMessage, err.Error())
		}
		if token.Valid {
			ginContext.Set("is_authorized", 1)
		}
	}
	ginContext.Next()
}
