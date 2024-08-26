package hetzner_dns

// Record represents a DNS record in the Hetzner DNS API.
//
// Fields:
// - Id: The unique identifier of the DNS record.
// - Name: The name of the DNS record.
// - TTL: The time-to-live value of the DNS record.
// - Type: The type of the DNS record (e.g., A, AAAA, CNAME).
// - Value: The value of the DNS record.
// - ZoneID: The ID of the zone to which the DNS record belongs.
type Record struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	TTL    int    `json:"ttl"`
	Type   string `json:"type"`
	Value  string `json:"value"`
	ZoneID string `json:"zone_id"`
}

// Records represents a collection of DNS records.
//
// Fields:
// - Records: A slice of Record structs.
type Records struct {
	Records []Record `json:"records"`
}

// Zone represents a DNS zone in the Hetzner DNS API.
//
// Fields:
// - Id: The unique identifier of the DNS zone.
// - Name: The name of the DNS zone.
type Zone struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Zones represents a collection of DNS zones.
//
// Fields:
// - Zones: A slice of Zone structs.
type Zones struct {
	Zones []Zone `json:"zones"`
}
