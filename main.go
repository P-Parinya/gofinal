package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/P-Parinya/gofinal/transactions"
	"github.com/gin-gonic/gin"
)

func authMiddleware(c *gin.Context) {
	log.Println("START authorization ...")
	token := c.GetHeader("Authorization")
	if token != "November 10, 2009" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "PERMISSION Denied."})
		c.Abort()
		return
	}
	c.Next()
	log.Println("END authorization ...")
}

func setRouting() *gin.Engine {
	log.Println("SETTING api route ...")

	r := gin.Default()
	customers := r.Group("/customers")
	customers.Use(authMiddleware)
	customers.POST("/", transactions.CreateCustomerHandler)
	customers.GET("/:id", transactions.GetCustomerByIDHandler)
	customers.GET("/", transactions.GetCustomerHandler)
	customers.PUT("/:id", transactions.UpdateCustomerHandler)
	customers.DELETE("/:id", transactions.DeleteCustomerByIDHandler)
	return r
}

func main() {
	fmt.Println("customer service")

	transactions.CreateTable()

	r := setRouting()

	r.Run(":2009")
	//run port ":2009"
}
