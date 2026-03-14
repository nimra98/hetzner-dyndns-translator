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

// Records represents a collection of DNS records (legacy).
type Records struct {
	Records []Record `json:"records"`
}

// Zone represents a DNS zone in both APIs.
type Zone struct {
	Id   json.RawMessage `json:"id"`
	Name string          `json:"name"`
}

// GetId returns the ID as a string, regardless of whether it was a number or string in JSON.
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
type RRset struct {
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	TTL     int     `json:"ttl"`
	Records []Value `json:"records"`
}

// Value represents a single value within an RRset.
type Value struct {
	Value string `json:"value"`
}

// RRsets represents a collection of RRsets.
type RRsets struct {
	RRsets []RRset `json:"rrsets"`
}
