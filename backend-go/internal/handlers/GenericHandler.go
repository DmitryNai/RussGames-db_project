package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListHandler(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet("db").(*gorm.DB)
		results := []interface{}{}
		if err := db.Model(model).Find(&results).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, results)
	}
}

func GetHandler(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet("db").(*gorm.DB)
		id := c.Param("id")
		obj := model
		if err := db.First(obj, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, obj)
	}
}

func CreateHandler(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet("db").(*gorm.DB)
		obj := model
		if err := c.BindJSON(obj); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(obj).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, obj)
	}
}

func DeleteHandler(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet("db").(*gorm.DB)
		id := c.Param("id")
		if err := db.Delete(model, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"deleted": id})
	}
}
