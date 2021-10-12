package main

import (
	"net/http"

	"codeberg.org/anbraten/hetzner-dyndns-translator/hetzner_dns"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	r := gin.Default()
	r.GET("/dyndns/:recordName/:zoneName/:token/:value", func(c *gin.Context) {
		token, _ := c.Params.Get("token")
		zoneName, _ := c.Params.Get("zoneName")
		recordName, _ := c.Params.Get("recordName")
		value, _ := c.Params.Get("value")
		if len(token) == 0 || len(zoneName) == 0 || len(recordName) == 0 || len(value) == 0 {
			c.String(http.StatusBadRequest, "badreq")
			return
		}

		log.Debug().Msgf("%s-%s-%s", zoneName, recordName, value)

		dns := hetzner_dns.NewHetznerDNS(token)
		err := dns.PatchRecord(zoneName, recordName, value)
		if err != nil {
			log.Error().Msg(err.Error())
			c.String(http.StatusBadRequest, "err")
			return
		}

		c.String(http.StatusAccepted, "good")
	})

	log.Debug().Msg("Starting server ...")

	r.Run(":3000")
}
