package controllers

import (
	"net/http"
	"gopost-be/config"
	"gopost-be/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateArticle(c *gin.Context) {
	var article models.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Validasi
	if len(article.Title) < 20 || len(article.Content) < 200 || len(article.Category) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed"})
		return
	}
	if article.Status != "publish" && article.Status != "draft" && article.Status != "thrash" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	query := `INSERT INTO articles (title, content, category, status, created_date, updated_date) 
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := config.SupabaseDB.QueryRow(query, article.Title, article.Content, article.Category, article.Status, time.Now(), time.Now()).Scan(&article.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": 201, "article": article})
}


func GetArticles(c *gin.Context) {
	limitParam := c.Param("limit")
	offsetParam := c.Param("offset")
	category := c.Query("category")

	var articles []models.Article
	var totalData int
	var query string
	var args []interface{}

	countQuery := `SELECT COUNT(*) FROM articles`
	countArgs := []interface{}{}

	if category != "" {
		countQuery += ` WHERE status = $1`
		countArgs = append(countArgs, category)
	}

	err := config.SupabaseDB.Get(&totalData, countQuery, countArgs...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total count", "details": err.Error()})
		return
	}

	query = `SELECT * FROM articles`
	args = []interface{}{}

	if category != "" {
		query += ` WHERE status = $1`
		args = append(args, category)
	}

	query += ` ORDER BY created_date DESC`

	if limitParam != "-" {
		limit, _ := strconv.Atoi(limitParam)
		offset, _ := strconv.Atoi(offsetParam)

		if category != "" {
			query += ` LIMIT $2 OFFSET $3`
			args = append(args, limit, offset)
		} else {
			query += ` LIMIT $1 OFFSET $2`
			args = append(args, limit, offset)
		}
	}

	err = config.SupabaseDB.Select(&articles, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":       200,
		"articles":   articles,
		"total_data": totalData,
	})
}



func GetArticleByID(c *gin.Context) {
	id := c.Param("id")

	var article models.Article
	query := `SELECT * FROM articles WHERE id = $1`
	err := config.SupabaseDB.Get(&article, query, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "article": article})
}

func UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	var updateData map[string]interface{}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	validStatus := map[string]bool{"publish": true, "draft": true, "thrash": true}
	for key, value := range updateData {
		valStr, ok := value.(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": key + " must be a string"})
			return
		}

		switch key {
		case "title":
			if len(valStr) < 20 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Title must be at least 20 characters"})
				return
			}
		case "content":
			if len(valStr) < 200 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Content must be at least 200 characters"})
				return
			}
		case "category":
			if len(valStr) < 3 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Category must be at least 3 characters"})
				return
			}
		case "status":
			if !validStatus[valStr] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be one of: publish, draft, thrash"})
				return
			}
		}
	}

	query := "UPDATE articles SET "
	args := []interface{}{}
	i := 1

	for key, value := range updateData {
		query += key + " = $" + strconv.Itoa(i) + ", "
		args = append(args, value)
		i++
	}

	query += "updated_date = $" + strconv.Itoa(i)
	args = append(args, time.Now())
	i++

	query += " WHERE id = $" + strconv.Itoa(i)
	args = append(args, id)

	_, err := config.SupabaseDB.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Article updated successfully"})
}


func DeleteArticle(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM articles WHERE id = $1`
	_, err := config.SupabaseDB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Article deleted"})
}
