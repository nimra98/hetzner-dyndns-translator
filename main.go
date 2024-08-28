package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nimra98/hetzner-dyndns-translator/hetzner_dns"
)

const VERSION = "1.4.0"

func main() {
	// Set release mode for Gin
	gin.SetMode(gin.ReleaseMode)

	// Get port from ENV variable or set to 3000 if not provided
	port := "3000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
		log.Printf("Port supplied through ENV: %s", port)
	} else {
		log.Print("Port not supplied through ENV, using default port 3000")
	}

	// Get Auth Token to access this service from ENV variable
	auth_token := ""
	if os.Getenv("AUTH_TOKEN") == "" {
		log.Print("No AUTH_TOKEN provided in ENV.")
		log.Print("This service will be open to the public!")
	} else {
		auth_token = os.Getenv("AUTH_TOKEN")
		log.Print("AUTH_TOKEN provided in ENV")
	}

	// Create url structure depending on AUTH_TOKEN
	// If Auth Token is provided, the URL will be /dyndns/:authtoken/:recordName/:zoneName/:token/:value
	// If Auth Token is not provided, the URL will be /dyndns/:recordName/:zoneName/:token/:value
	// The Auth Token will be checked in the handler function

	r := gin.Default()

	// Route abhängig von AUTH_TOKEN erstellen
	if auth_token == "" {
		// Route ohne Auth Token
		r.GET("/dyndns/:recordName/:zoneName/:token/:value", handleRequest)
	} else {
		// Route mit Auth Token
		r.GET("/dyndns/:authToken/:recordName/:zoneName/:token/:value", func(c *gin.Context) {
			// Auth Token überprüfen
			if c.Param("authToken") != auth_token {
				c.String(http.StatusUnauthorized, "unauthorized")
				return
			}
			handleRequest(c)
		})
	}

	log.Printf("Starting server on Port " + port + " ...")
	log.Printf("Translator Version: %s", VERSION)

	r.Run(":" + port)
}

// Handler-Funktion für die Route
func handleRequest(c *gin.Context) {
	token := c.Param("token")
	zoneName := c.Param("zoneName")
	recordName := c.Param("recordName")
	value := c.Param("value")

	if token == "" || zoneName == "" || recordName == "" || value == "" {
		c.String(http.StatusBadRequest, "badreq")
		log.Printf("badreq - token: %s, zoneName: %s, recordName: %s, value: %s", token, zoneName, recordName, value)
		return
	}

	log.Printf("%s-%s-%s", zoneName, recordName, value)
	dns := hetzner_dns.NewHetznerDNS(token)
	err := dns.PatchRecord(zoneName, recordName, value)
	if err != nil {
		c.String(http.StatusInternalServerError, "error")
		log.Printf("Failed to patch DNS record: %v", err)
		log.Printf("error - error: %s, token: %s, zoneName: %s, recordName: %s, value: %s", err.Error(), token, zoneName, recordName, value)
		return
	}

	c.String(http.StatusOK, "success")
	log.Printf("ok - token: %s, zoneName: %s, recordName: %s, value: %s", token, zoneName, recordName, value)
}
