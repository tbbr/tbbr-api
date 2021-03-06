package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mattes/migrate/migrate"
	"github.com/tbbr/tbbr-api/auth"
	"github.com/tbbr/tbbr-api/controllers"
	"github.com/tbbr/tbbr-api/database"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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

func getDBUrl() string {
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
	return dbURL
}

func migrateDB() {
	fmt.Println("TBBR-API - Running Migrations located in: " + os.Getenv("TBBR_MIGRATIONS_DIR"))

	allErrors, ok := migrate.UpSync(getDBUrl(), os.Getenv("TBBR_MIGRATIONS_DIR"))
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

	database.DBCon, err = gorm.Open("postgres", getDBUrl())

	if err != nil {
		fmt.Printf("TBBR-API - Error occurred %s\n", err)
	} else {
		fmt.Printf("TBBR-API - Connection setup with database\n")
		fmt.Println("TBBR-API - Pinging open connections:", database.DBCon.DB().Stats().OpenConnections)
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
		transactions := authorized.Group("/transactions")
		{
			transactions.GET("", controllers.TransactionIndex)
			transactions.POST("", controllers.TransactionCreate)

			transactions.PATCH("/:id", controllers.TransactionUpdate)
			transactions.DELETE("/:id", controllers.TransactionDelete)
		}
		groups := authorized.Group("/groups")
		{
			groups.GET("", controllers.GroupIndex)
			groups.POST("", controllers.GroupCreate)

			groups.GET("/:id", controllers.GroupShow)
			groups.POST("/:id/join", controllers.GroupJoin)
			groups.PATCH("/:id", controllers.GroupUpdate)
			groups.DELETE("/:id", controllers.GroupDelete)
		}
		groupTransactions := authorized.Group("/group-transactions")
		{
			groupTransactions.GET("", controllers.GroupTransactionIndex)
			groupTransactions.POST("", controllers.GroupTransactionCreate)

			// groupTransactions.PATCH("/:id", controllers.GroupTransactionUpdate)
			groupTransactions.DELETE("/:id", controllers.GroupTransactionDelete)
		}
		groupMembers := authorized.Group("/group-members")
		{
			groupMembers.GET("", controllers.GroupMemberIndex)
			groupMembers.POST("", controllers.GroupMemberCreate)
			groupMembers.DELETE("/:id", controllers.GroupMemberDelete)
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
		deviceTokens := authorized.Group("/device-tokens")
		{
			deviceTokens.POST("", controllers.DeviceTokenCreate)
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
				switch err := e.Meta.(type) {
				case appError.Err:
					errors = append(errors, err)
				}
			}
			// Use Status of first error
			// TODO: Currently c.JSON doesn't set the content type properly
			// a fix exists in v1.2 follow this thread: https://github.com/gin-gonic/gin/issues/762
			if len(errors) > 0 {
				c.JSON(errors[0].Status, gin.H{"errors": errors})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"errors": c.Errors.Errors()})
			}
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
