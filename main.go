package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	err := app.Run(":3000")
	if err != nil {
		return
	}
}
