package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nimra98/hetzner-dyndns-translator/hetzner_dns"
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
			log.Debug().Msgf("badreq - token: %s, zoneName: %s, recordName: %s, value: %s", token, zoneName, recordName, value)
			return
		}

		log.Debug().Msgf("%s-%s-%s", zoneName, recordName, value)

		dns := hetzner_dns.NewHetznerDNS(token)
		err := dns.PatchRecord(zoneName, recordName, value)
		if err != nil {
			c.String(http.StatusBadRequest, "err")
			log.Debug().Msgf("error - error: %s, token: %s, zoneName: %s, recordName: %s, value: %s", err.Error(), token, zoneName, recordName, value)
			return
		}

		c.String(http.StatusOK, "OK")
		log.Debug().Msgf("ok - token: %s, zoneName: %s, recordName: %s, value: %s", token, zoneName, recordName, value)
	})

	log.Debug().Msg("Starting server ...")

	r.Run(":3000")
}
