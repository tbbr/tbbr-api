package main

import (
	"fmt"
	"os"
	"runtime"

	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mattes/migrate/migrate"
	"github.com/tbbr/tbbr-api/auth"
	"github.com/tbbr/tbbr-api/controllers"
	"github.com/tbbr/tbbr-api/database"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/tbbr/tbbr-api/app-error"
)

func main() {
	configRuntime()
	migrateDB()
	bootstrap()
	startGin()
}

func configRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("TBBR-API - Running with %d CPUs\n", nuCPU)
}

func migrateDB() {
	fmt.Println("TBBR-API - Running Migrations")

	var dbURL string
	if os.Getenv("TBBR_DB_PASSWORD") == "" {
		dbURL = fmt.Sprintf("postgres://%s@localhost:5432/%s?sslmode=disable",
			os.Getenv("TBBR_DB_USER"),
			os.Getenv("TBBR_DB_NAME"),
		)
	} else {
		dbURL = fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable",
			os.Getenv("TBBR_DB_USER"),
			os.Getenv("TBBR_DB_PASSWORD"),
			os.Getenv("TBBR_DB_NAME"),
		)
	}

	allErrors, ok := migrate.UpSync(dbURL, "./migrations")
	if !ok {
		fmt.Println("TBBR-API Migrations failed!")
		fmt.Println(allErrors)
		fmt.Println("TBBR-API exiting...")
		os.Exit(1)
	}

	fmt.Println("TBBR-API - Migrations Finished!")
}

func bootstrap() {

	var err error

	database.DBCon, err = gorm.Open("postgres",
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("TBBR_DB_USER"),
			os.Getenv("TBBR_DB_PASSWORD"),
			os.Getenv("TBBR_DB_NAME"),
		),
	)

	if err != nil {
		fmt.Printf("TBBR-API - Error occurred %s\n", err)
	} else {
		fmt.Printf("TBBR-API - Connection setup with database\n")
		fmt.Printf("TBBR-API - Pinging: %s \n", database.DBCon.DB().Ping())
	}
}

func startGin() {
	// Creates a gin router with default middlewares:
	// logger and recovery (crash-free) middlewares
	router := gin.Default()
	router.RedirectTrailingSlash = true

	router.Use(Cors())
	router.Use(handleErrors())

	router.GET("", controllers.ServeIndex)

	router.NoRoute(controllers.ServeIndex)

	authorized := router.Group("api", OAuthMiddleware())
	{
		friendships := authorized.Group("/friendships")
		{
			friendships.GET("", controllers.FriendshipIndex)

			friendships.GET("/:id", controllers.FriendshipShow)
		}
		groups := authorized.Group("/groups")
		{
			groups.GET("", controllers.GroupIndex)
			groups.POST("", controllers.GroupCreate)

			groups.GET("/:id", controllers.GroupShow)
			groups.PATCH("/:id", controllers.GroupUpdate)
			groups.DELETE("/:id", controllers.GroupDelete)
		}
		transactions := authorized.Group("/transactions")
		{
			transactions.GET("", controllers.TransactionIndex)
			transactions.POST("", controllers.TransactionCreate)

			transactions.PATCH("/:id", controllers.TransactionUpdate)
			transactions.DELETE("/:id", controllers.TransactionDelete)
		}
		tokens := router.Group("api/tokens")
		{
			tokens.POST("/oauth/grant", controllers.TokenOAuthGrant)
		}
		users := authorized.Group("/users")
		{
			users.GET("", controllers.UserIndex)
			users.GET("/:id", controllers.UserShow)
			users.PATCH("/:id", controllers.UserUpdate)
			users.DELETE("/:id", controllers.UserDelete)
		}
	}

	// Listen and serve on 0.0.0.0:8090
	router.Run(":8090")
}

// Cors - Enables cors for the api
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Writer.Header().Add("Access-Control-Allow-Methods", "HEAD, GET, PATCH, POST, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func handleErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			errors := []appError.Err{}

			for _, e := range c.Errors {
				err := e.Meta.(appError.Err)
				errors = append(errors, err)
			}
			// Use Status of first error
			c.JSON(errors[0].Status, gin.H{"errors": errors})
		}
	}
}

// OAuthMiddleware handles validation of the authoriation code
func OAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			c.Abort()
			errors := []appError.Err{appError.AuthorizationMissing}
			c.JSON(errors[0].Status, gin.H{"errors": errors})
			return
		}

		token, err := auth.GetToken(authorization)

		if err != nil {
			c.Abort()
			errors := []appError.Err{err.(appError.Err)}
			c.JSON(errors[0].Status, gin.H{"errors": errors})
			return
		}

		// Check that token hasn't expired here
		if token.Expired() {
			c.Abort()
			errors := []appError.Err{appError.AccessTokenExpired}
			c.JSON(errors[0].Status, gin.H{"errors": errors})
		}

		// Attach the current user's id onto the context
		c.Keys = make(map[string]interface{})
		c.Keys["CurrentUserID"] = token.UserID

		c.Next()

	}
}
