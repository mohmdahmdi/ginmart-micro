package controllers

import (
	"go-micro/config"
	"go-micro/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddReview(c *gin.Context){
	var review models.Review

	userId, exist := c.Get("user_id")
	if !exist || userId == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. login is required"})
		return
	}

	if err := c.ShouldBindJSON(&review); err != nil {
		log.Println("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	query := `INSERT INTO reviews (product_id, user_id, rating, comment) 
	          VALUES ($1, $2, $3, $4) RETURNING id`
	err := config.DB.QueryRow(query, review.ProductId, userId, review.Rating, review.Comment).Scan(&review.ID)

	if err != nil {
		log.Println("Error inserting review:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Review added successfully",
		"review": review,
	})}

func GetReviews(c *gin.Context){
	var reviews [] models.Review
	rows , err := config.DB.Query("SELECT id, product_id, user_id, rating, comment FROM products")
	if err != nil {
		log.Println("Error fetching reviews: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Review
		err := rows.Scan(&r.ID, &r.ProductId, &r.UserId, &r.Rating, &r.Comment)
		if err != nil {
			log.Println("Error scanning row: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read review data"})
			return
		}
		reviews = append(reviews, r)
	}

	c.JSON(http.StatusOK, gin.H{
		"reviews": reviews,
	})
}

func DeleteReview(c *gin.Context){
	productId := c.Param("id")
	reviewId := c.Param("id")
	query := "DELETE FROM reviews where product_id = $1 AND id = $2"
	result, err := config.DB.Exec(query, productId, reviewId)
	if err != nil {
		log.Println("Error deleting review:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting affected rows:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}