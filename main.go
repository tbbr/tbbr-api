package main

import (
	"fmt"
	"payup/database"
	"payup/group"
	"payup/user"
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
