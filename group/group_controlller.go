package group

import (
	"net/http"
	"payup/user"

	"github.com/gin-gonic/gin"
)

// Index When the group's index is routed to
// this handler will run. Generally, it will
// come with some query parameters like limit and offset
// @returns an array of group structs
func Index(c *gin.Context) {
	group := Group{
		Name:        "Family Group",
		Description: "Thingies about things",
		Users: []user.User{
			{
				Name:  "Maaz Ali",
				Email: "maazali40@gmail.com",
			},
			{
				Name:  "Test User",
				Email: "testuser@gmail.com",
			},
		},
	}

	c.JSON(http.StatusOK, group)
}

// Show is used to show one specific group, returns a group struct
// @returns a group struct
func Show(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"groupShow": c.Param("id")})
}

// Create is used to create one specific group, it'll come with some form data
// @returns a group struct
func Create(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"groupCreate": "someContent"})
}

// Update is used to update a specific group, it'll also come with some form data'
// @returns a group struct
func Update(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"groupUpdate": "someContent"})
}

// Delete is used to delete one specific group with a `id`
func Delete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"groupDelete": "someContent"})
}
