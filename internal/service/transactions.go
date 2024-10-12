package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func CreateTransaction() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fmt.Println("Create Transaction")
	}
}
