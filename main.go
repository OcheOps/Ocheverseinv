package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var (
	token  string
	orgID  int64
	client *github.Client
	router *gin.Engine
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token = os.Getenv("GITHUB_TOKEN")
	orgIDStr := os.Getenv("ORG_ID")
	orgID, err = strconv.ParseInt(orgIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Invalid ORG_ID environment variable: %v", err)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client = github.NewClient(tc)

	router = gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", indexHandler)
	router.POST("/add", addHandler)

	log.Fatal(router.Run(":8080"))
}

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func addHandler(c *gin.Context) {
	username := c.PostForm("github")
	_, _, err := client.Organizations.AddOrgMembership(orgID, username, nil)
	if err != nil {
		log.Printf("Error adding user %s to organization: %v", username, err)
		c.String(http.StatusInternalServerError, "Error adding user to organization")
		return
	}

	c.String(http.StatusOK, "OK, Check your EMAIL")
}