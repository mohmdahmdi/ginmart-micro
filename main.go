package main

import (
	"go-micro/controllers"
	"go-micro/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    router.Use(middlewares.LoggerMiddleWare())

    router.POST("/auth/register", controllers.Register)
    router.POST("/auth/login", controllers.Login)

    router.POST("/products", middlewares.AuthMiddleware(), controllers.AddProduct) 
    router.GET("/products", controllers.GetProducts)
    router.GET("/products/:id", controllers.GetProduct)
    router.PUT("/products/:id", middlewares.AuthMiddleware(), controllers.UpdateProduct)
    router.DELETE("/products/:id", middlewares.AuthMiddleware(), controllers.DeleteProduct)

    router.POST("/products/:id/reviews", middlewares.AuthMiddleware(), controllers.AddReview)
    router.GET("/products/:id/reviews", controllers.GetReviews)
		router.DELETE("/products/:id/reviews/:id", controllers.DeleteReview)

    router.Run()
}
