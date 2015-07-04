package main

import (
	"fmt"
	"payup/group"
	"payup/user"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func main() {

	configRuntime()
	startGin()
}

func configRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

func bootsrap() {

	db, err := gorm.Open("postgres", "user=test dbname=payup sslmode=disable")

	if err != nil {
		fmt.Printf("Error occurred", err)
	}

	db.DB().Ping()
}

func startGin() {
	// Creates a gin router with default middlewares:
	// logger and recovery (crash-free) middlewares
	router := gin.Default()

	// Handle assets and index.html file
	// router.Static("/", "index.html")
	router.Static("/assets", "./assets")

	v1 := router.Group("/api/v1")
	{
		groups := v1.Group("/groups")
		{
			groups.GET("/", group.Index)
			groups.POST("/", group.Create)

			groups.GET("/:id", group.Show)
			groups.PUT("/:id", group.Update)
			groups.DELETE("/:id", group.Delete)
		}

		users := v1.Group("/users")
		{
			users.GET("/", user.Index)

			users.GET("/:id", user.Show)
			users.POST("/", user.Create)
			users.PUT("/:id", user.Update)
			users.DELETE("/:id", user.Delete)
		}
	}

	// Listen and server on 0.0.0.0:8080
	router.Run(":8080")
}
