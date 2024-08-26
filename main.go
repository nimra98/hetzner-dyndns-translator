package main

import (
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/nimra98/hetzner-dyndns-translator/hetzner_dns"
)

const VERSION = "1.3.0"

func main() {
	// Set release mode for Gin
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.GET("/dyndns/:recordName/:zoneName/:token/:value", func(c *gin.Context) {
		token, _ := c.Params.Get("token")
		zoneName, _ := c.Params.Get("zoneName")
		recordName, _ := c.Params.Get("recordName")
		value, _ := c.Params.Get("value")
		if len(token) == 0 || len(zoneName) == 0 || len(recordName) == 0 || len(value) == 0 {
			c.String(http.StatusBadRequest, "badreq")
			log.Printf("badreq - token: %s, zoneName: %s, recordName: %s, value: %s", token, zoneName, recordName, value)
			return
		}

		log.Printf("%s-%s-%s", zoneName, recordName, value)
		dns := hetzner_dns.NewHetznerDNS(token)
		err := dns.PatchRecord(zoneName, recordName, value)
		if err != nil {
			c.String(http.StatusBadRequest, "err")
			log.Printf("error - error: %s, token: %s, zoneName: %s, recordName: %s, value: %s", err.Error(), token, zoneName, recordName, value)
			return
		}

		c.String(http.StatusOK, "OK")
		log.Printf("ok - token: %s, zoneName: %s, recordName: %s, value: %s", token, zoneName, recordName, value)
	})

	log.Printf("Starting server ...")
	log.Printf("Translator Version: %s", VERSION)

	r.Run(":3000")
}
