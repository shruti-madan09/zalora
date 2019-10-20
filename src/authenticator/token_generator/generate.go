package main

import (
	"fmt"

	"authenticator"
)

func main() {
	token, err := authenticator.GenerateJWT()
	if err == nil {
		fmt.Println("Token: ", token)
	}
}
