package controllers

import (
	"database/sql"
	"go-micro/config"
	"go-micro/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddProduct(c *gin.Context){
	userRole, exists := c.Get("role")
	if !exists || userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Admins only"})
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		log.Println("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	query := `INSERT INTO products (name, description, price, category, stock_count) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := config.DB.QueryRow(query, product.Name, product.Description, product.Price, product.Category, product.StockCount).Scan(&product.ID)

	if err != nil {
		log.Println("Error inserting product:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product added successfully",
		"product": product,
	})
}

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

	c.JSON(http.StatusOK, gin.H{"product": product})
}


func UpdateProduct(c *gin.Context){}

func DeleteProduct(c *gin.Context){
	productId := c.Param("id")
	query := "DELETE FROM products WHERE id = $1" 
	result, err := config.DB.Exec(query, productId)
	
	if err != nil {
		log.Println("error deleting product: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting affected rows:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}