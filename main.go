package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var (
	token  string
	orgID  int64
	port   int
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

	portStr := os.Getenv("PORT")
	port, err = strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid PORT environment variable: %v", err)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client = github.NewClient(tc)

	router = gin.Default()

	router.GET("/", indexHandler)
	router.POST("/add", addHandler)

	log.Printf("Starting server on port %d", port)
	log.Fatal(router.Run(fmt.Sprintf(":%d", port)))
}

func indexHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello, World!")
}

func addHandler(c *gin.Context) {
	username := c.PostForm("github")
	invitationOpt := &github.OrgInvitationOpt{
		InviteeID: &username,
		Role:      github.String("direct_member"),
	}
	_, _, err := client.Organizations.CreateOrgInvitation(orgID, invitationOpt)
	if err != nil {
		log.Printf("Error adding user %s to organization: %v", username, err)
		c.String(http.StatusInternalServerError, "Error adding user to organization")
		return
	}

	c.String(http.StatusOK, "OK, Check your EMAIL")
}