package main

import (
	"fmt"
	"payup/controllers"
	"payup/database"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func main() {

	configRuntime()
	bootstrap()
	startGin()
}

func configRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

func bootstrap() {

	var err error
	database.DBCon, err = gorm.Open("postgres", "user=maazali dbname=payup_backup sslmode=disable")

	if err != nil {
		fmt.Printf("Error occurred %s\n", err)
	} else {
		fmt.Printf("Connection setup with database\n")
		fmt.Printf("Pinging: %s \n", database.DBCon.DB().Ping())
	}
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
			groups.GET("/", controllers.GroupIndex)
			groups.POST("/", controllers.GroupCreate)

			groups.GET("/:id", controllers.GroupShow)
			groups.PUT("/:id", controllers.GroupUpdate)
			groups.DELETE("/:id", controllers.GroupDelete)
		}

		users := v1.Group("/users")
		{
			users.GET("/", controllers.UserIndex)

			users.GET("/:id", controllers.UserShow)
			users.POST("/", controllers.UserCreate)
			users.PUT("/:id", controllers.UserUpdate)
			users.DELETE("/:id", controllers.UserDelete)
		}
	}

	// Listen and server on 0.0.0.0:8080
	router.Run(":8080")
}
