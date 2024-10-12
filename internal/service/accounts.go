package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func CreateAccount() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fmt.Println("Create Account")
	}
}

func GetAccount() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fmt.Println("Get Account")
	}
}
