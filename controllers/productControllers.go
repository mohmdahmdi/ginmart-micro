package controllers

import (
	"go-micro/config"
	"go-micro/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddProduct(c *gin.Context){}

func GetProducts(c *gin.Context){
	var products []models.Product
	rows , err := config.DB.Query("SELECT id, name, description, price, category, stock_count FROM products")
	if err != nil {
		log.Println("Error fetching products: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Category, &p.StockCount)
		if err != nil {
			log.Println("Error scanning row: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read product data"})
			return
		}
		products = append(products, p)
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

func GetProduct(c *gin.Context){}

func UpdateProduct(c *gin.Context){}

func DeleteProduct(c *gin.Context){}