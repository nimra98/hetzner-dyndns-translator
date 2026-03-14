package hetzner_dns

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
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Zones represents a collection of DNS zones.
type Zones struct {
	Zones []Zone `json:"zones"`
}

// RRset represents a Resource Record Set in the new Hetzner Cloud DNS API.
type RRset struct {
	Id      string   `json:"id,omitempty"`
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	TTL     int      `json:"ttl"`
	Records []Value  `json:"records"`
}

// Value represents a single value within an RRset.
type Value struct {
	Value string `json:"value"`
}

// RRsets represents a collection of RRsets.
type RRsets struct {
	RRsets []RRset `json:"rrsets"`
}
