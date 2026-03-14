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

// CloudRecord represents a DNS record in the new Hetzner Cloud DNS API.
type CloudRecord struct {
	Id    json.RawMessage `json:"id,omitempty"`
	Name  string          `json:"name"`
	Type  string          `json:"type"`
	Value string          `json:"value"`
	TTL   int             `json:"ttl"`
}

// GetId returns the ID as a string.
func (r *CloudRecord) GetId() string {
	var s string
	if err := json.Unmarshal(r.Id, &s); err == nil {
		return s
	}
	var i int64
	if err := json.Unmarshal(r.Id, &i); err == nil {
		return json.Number(string(r.Id)).String()
	}
	return string(r.Id)
}

// CloudRecords represents a collection of DNS records in the Cloud API.
type CloudRecords struct {
	Records []CloudRecord `json:"records"`
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
	Id      json.RawMessage `json:"id,omitempty"`
	Name    string          `json:"name"`
	Type    string          `json:"type"`
	TTL     int             `json:"ttl"`
	Records []Value         `json:"records"`
}

// GetId returns the ID as a string.
func (r *RRset) GetId() string {
	var s string
	if err := json.Unmarshal(r.Id, &s); err == nil {
		return s
	}
	return string(r.Id)
}

// Value represents a single value within an RRset.
type Value struct {
	Value string `json:"value"`
}

// RRsets represents a collection of RRsets.
type RRsets struct {
	RRsets []RRset `json:"rrsets"`
}
