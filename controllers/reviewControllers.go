package controllers

import (
	"go-micro/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddReview(c *gin.Context){}

func GetReviews(c *gin.Context){}

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