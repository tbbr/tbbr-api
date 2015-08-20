package controllers

import (
	"net/http"
	"payup/database"
	"payup/models"

	"github.com/gin-gonic/gin"
)

// GroupIndex When the group's index is routed to
// this handler will run. Generally, it will
// come with some query parameters like limit and offset
// @returns an array of group structs
func GroupIndex(c *gin.Context) {
	var groups []models.Group
	database.DBCon.Limit(c.Param("limit")).Find(&groups)

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

// GroupShow is used to show one specific group, returns a group struct
// @returns a group struct
func GroupShow(c *gin.Context) {
	var group models.Group
	database.DBCon.First(&group, c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"group": group})
}

// // LoginJSON stuff
// type LoginJSON struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

// GroupCreate is used to create one specific group, it'll come with some form data
// @returns a group struct
func GroupCreate(c *gin.Context) {
	// group := models.Group{
	// 	Name:        c.PostForm("name"),
	// 	Description: c.PostForm("description"),
	// }
	var group models.Group
	c.Bind(&group)
	c.JSON(200, group)
	// database.DBCon.Create(&group)

	// c.JSON(http.StatusOK, gin.H{"group": c.PostForm("group")})
}

// GroupUpdate is used to update a specific group, it'll also come with some form data'
// @returns a group struct
func GroupUpdate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"groupUpdate": "someContent"})
}

// GroupDelete is used to delete one specific group with a `id`
func GroupDelete(c *gin.Context) {
	var group models.Group
	database.DBCon.First(&group, c.Param("id"))
	database.DBCon.Delete(&group)

	c.JSON(http.StatusOK, gin.H{"success": true})
}
