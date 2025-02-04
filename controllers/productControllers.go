package controllers

import (
	"database/sql"
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

func GetProduct(c *gin.Context) {
	var id = c.Param("id")
	var query = "SELECT id, name, description, price, category, stock_count FROM products WHERE id = $1"
	row := config.DB.QueryRow(query, id)
	var product models.Product
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Category, &product.StockCount)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			log.Println("Error fetching product: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		}
		return
	}

	// Return the product data as a JSON response
	c.JSON(http.StatusOK, gin.H{"product": product})
}


func UpdateProduct(c *gin.Context){}

func DeleteProduct(c *gin.Context){}