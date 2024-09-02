package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nimra98/hetzner-dyndns-translator/hetzner_dns"
)

// Override Version with build flag -ldflags "-X main.VERSION=1.0.0"
var VERSION = "v1.0.0"

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
	service_auth_token := ""
	if os.Getenv("SERVICE_AUTH_TOKEN") == "" {
		log.Print("No SERVICE_AUTH_TOKEN provided in ENV.")
		log.Print("This service will be open to the public!")
	} else {
		service_auth_token = os.Getenv("SERVICE_AUTH_TOKEN")
		log.Print("SERVICE_AUTH_TOKEN provided in ENV")
	}

	// Create url structure depending on SERVICE_AUTH_TOKEN
	// If Auth Token is provided, the URL will be /dyndns/:authtoken/:recordName/:zoneName/:token/:value
	// If Auth Token is not provided, the URL will be /dyndns/:recordName/:zoneName/:token/:value
	// The Auth Token will be checked in the handler function

	r := gin.New()
	if os.Getenv("SHOW_HETZNER_API_TOKEN") == "true" {
		// add logger to gin to show full request
		r.Use(gin.Logger())
	}

	// Route abhängig von SERVICE_AUTH_TOKEN erstellen
	if service_auth_token == "" {
		// Route ohne Auth Token
		r.GET("/dyndns/:recordName/:zoneName/:token/:value", handleRequest)
	} else {
		// Route mit Auth Token
		r.GET("/dyndns/:authToken/:recordName/:zoneName/:token/:value", func(c *gin.Context) {
			// Auth Token überprüfen
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

	log.Printf("Updating %s.%s with new IP %s", recordName, zoneName, value)
	dns := hetzner_dns.NewHetznerDNS(token)
	err := dns.PatchRecord(zoneName, recordName, value)
	if err != nil {
		c.String(http.StatusInternalServerError, "error")
		log.Printf("Failed to patch DNS record: %v", err)
		log.Printf("error - error: %s, token: %s, zoneName: %s, recordName: %s, value: %s", err.Error(), token, zoneName, recordName, value)
		return
	}

	c.String(http.StatusOK, "success")
	if os.Getenv("SHOW_HETZNER_API_TOKEN") == "true" {
		log.Printf("Transaktion ok, DNS updated - token: %s, zoneName: %s, recordName: %s, value: %s", token, zoneName, recordName, value)
	} else {
		log.Printf("Transaktion ok, DNS updated - token: redacted, zoneName: %s, recordName: %s, value: %s", zoneName, recordName, value)
	}
}
