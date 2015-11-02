package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

// ServeIndex serves the front-end's index file
func ServeIndex(c *gin.Context) {
	auth, err := aws.EnvAuth()

	if err != nil {
		log.Fatal(err)
	}

	client := s3.New(auth, aws.USEast)
	bucket := client.Bucket("tbbr.me")
	data, err := bucket.Get("index.html")

	if err != nil {
		log.Fatal(err)
	}

	c.Data(http.StatusOK, "text/html", data)
}
