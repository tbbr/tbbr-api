package main

import (
	"fmt"
	"net/http"
	"os"
	"payup/app-error"
	"payup/auth"
	"payup/controllers"
	"payup/database"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
)

func main() {
	configRuntime()
	bootstrap()
	setupAuthProviders()
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

func setupAuthProviders() {
	goth.UseProviders(
		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), "http://localhost:8080/auth/facebook/callback"),
	)

	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return "facebook", nil
	}

}

func startGin() {
	// Creates a gin router with default middlewares:
	// logger and recovery (crash-free) middlewares
	router := gin.Default()
	router.RedirectTrailingSlash = true

	router.Use(Cors())
	router.Use(handleErrors())

	router.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to PayUp's API")
	})

	authorized := router.Group("", OAuthMiddleware())
	{
		groups := authorized.Group("/groups")
		{
			groups.GET("", controllers.GroupIndex)
			groups.POST("", controllers.GroupCreate)

			groups.GET("/:id", controllers.GroupShow)
			groups.PATCH("/:id", controllers.GroupUpdate)
			groups.DELETE("/:id", controllers.GroupDelete)
		}

		users := authorized.Group("/users")
		{
			users.GET("", controllers.UserIndex)
			users.POST("", controllers.UserCreate)
			users.GET("/:id", controllers.UserShow)
			users.PUT("/:id", controllers.UserUpdate)
			users.DELETE("/:id", controllers.UserDelete)
		}

		transactions := authorized.Group("/transactions")
		{
			transactions.POST("", controllers.TransactionCreate)
			transactions.GET("", controllers.TransactionIndex)

		}
		tokens := router.Group("/tokens")
		{
			tokens.POST("/oauth/grant", controllers.TokenOAuthGrant)
		}
	}

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
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
