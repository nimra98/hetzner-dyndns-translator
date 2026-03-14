package hetzner_dns

import "encoding/json"

// Record represents a DNS record in the legacy Hetzner DNS API.
type Record struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	TTL    int    `json:"ttl"`
	Type   string `json:"type"`
	Value  string `json:"value"`
	ZoneID string `json:"zone_id"`
}

// Records represents a collection of DNS records returned by the legacy API.
type Records struct {
	Records []Record `json:"records"`
}

// Zone represents a DNS zone in both Hetzner DNS APIs.
type Zone struct {
	Id   json.RawMessage `json:"id"`
	Name string          `json:"name"`
}

// GetId returns the Zone ID as a string.
// Since numeric IDs (Cloud API) and string IDs (Legacy API) are used,
// json.RawMessage is used to handle both formats transparently.
func (z *Zone) GetId() string {
	var s string
	if err := json.Unmarshal(z.Id, &s); err == nil {
		return s
	}
	var i int64
	if err := json.Unmarshal(z.Id, &i); err == nil {
		return json.Number(string(z.Id)).String()
	}
	return string(z.Id)
}

// Zones represents a collection of DNS zones.
type Zones struct {
	Zones []Zone `json:"zones"`
}

// RRset represents a Resource Record Set in the new Hetzner Cloud DNS API.
// It groups records of the same name and type together.
type RRset struct {
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	TTL     int     `json:"ttl"`
	Records []Value `json:"records"`
}

// Value represents a single DNS record value within an RRset.
type Value struct {
	Value string `json:"value"`
}

// RRsets represents a collection of RRsets returned by the Cloud API.
type RRsets struct {
	RRsets []RRset `json:"rrsets"`
}
