package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nimra98/hetzner-dyndns-translator/hetzner_dns"
)

// VERSION holds the current application version. 
// It can be overridden during build time using: -ldflags "-X main.VERSION=1.0.0"
var VERSION = "v1.0.0"

func main() {
	// Disable Gin's debug output for a cleaner production log
	gin.SetMode(gin.ReleaseMode)

	// Determine the listening port from environment or use 3000 as default
	port := "3000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
		log.Printf("Port supplied through ENV: %s", port)
	} else {
		log.Print("Port not supplied through ENV, using default port 3000")
	}

	// Service authentication: if SERVICE_AUTH_TOKEN is set, all requests must include it in the URL
	service_auth_token := ""
	if os.Getenv("SERVICE_AUTH_TOKEN") == "" {
		log.Print("No SERVICE_AUTH_TOKEN provided in ENV.")
		log.Print("This service will be open to the public!")
	} else {
		service_auth_token = os.Getenv("SERVICE_AUTH_TOKEN")
		log.Print("SERVICE_AUTH_TOKEN provided in ENV")
	}

	r := gin.New()
	
	// Optional: logging of full requests (caution: might leak API tokens if not used carefully)
	if os.Getenv("SHOW_HETZNER_API_TOKEN") == "true" {
		r.Use(gin.Logger())
	}

	// Define routes based on whether authentication is required
	if service_auth_token == "" {
		// Public route: /dyndns/:recordName/:zoneName/:token/:value
		r.GET("/dyndns/:recordName/:zoneName/:token/:value", handleRequest)
	} else {
		// Protected route: /dyndns/:authToken/:recordName/:zoneName/:token/:value
		r.GET("/dyndns/:authToken/:recordName/:zoneName/:token/:value", func(c *gin.Context) {
			// Validate service-level authentication
			if c.Param("authToken") != service_auth_token {
				c.String(http.StatusUnauthorized, "unauthorized")
				log.Print("Unauthorized request with invalid Auth Token: " + c.Param("authToken") + " - from IP: " + c.ClientIP())
				return
			}
			handleRequest(c)
		})
	}

	log.Printf("Starting server on Port " + port + " ...")
	log.Printf("Translator Version: %s", VERSION)
	if service_auth_token == "" {
		log.Print("Awaiting requests in format /dyndns/:recordName/:zoneName/:hetzner_api_token/:value")
	} else {
		log.Print("Awaiting requests in format /dyndns/:service_authToken/:recordName/:zoneName/:hetzner_api_token/:value")
	}

	// Start the HTTP server
	r.Run(":" + port)
}

// handleRequest processes the incoming DynDNS update request.
// It extracts parameters from the URL, initializes the correct Hetzner DNS provider,
// and performs the record update (Patch).
func handleRequest(c *gin.Context) {
	token := c.Param("token")       // Hetzner API Token
	zoneName := c.Param("zoneName")   // DNS Zone (e.g. example.com)
	recordName := c.Param("recordName") // Subdomain (e.g. fritzbox)
	value := c.Param("value")       // New IP address

	// Basic validation of required parameters
	if token == "" || zoneName == "" || recordName == "" || value == "" {
		c.String(http.StatusBadRequest, "badreq")
		log.Printf("badreq - token: %s, zoneName: %s, recordName: %s, value: %s", token, zoneName, recordName, value)
		return
	}

	log.Printf("Updating %s.%s with new IP %s", recordName, zoneName, value)
	
	// Create provider (Legacy or Cloud) based on HETZNER_API_VERSION env variable
	dns := hetzner_dns.NewHetznerDNS(token)
	
	// Execute the update
	err := dns.PatchRecord(zoneName, recordName, value)
	if err != nil {
		c.String(http.StatusInternalServerError, "error")
		log.Printf("Failed to patch DNS record: %v", err)
		log.Printf("error - error: %s, token: %s, zoneName: %s, recordName: %s, value: %s", err.Error(), token, zoneName, recordName, value)
		return
	}

	// Return success to the client (e.g. Fritz!Box)
	c.String(http.StatusOK, "success")
	
	if os.Getenv("SHOW_HETZNER_API_TOKEN") == "true" {
		log.Printf("Transaction ok, DNS updated - token: %s, zoneName: %s, recordName: %s, value: %s", token, zoneName, recordName, value)
	} else {
		log.Printf("Transaction ok, DNS updated - token: redacted, zoneName: %s, recordName: %s, value: %s", zoneName, recordName, value)
	}
}
