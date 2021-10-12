package hetzner_dns

type Record struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	TTL    int    `json:"ttl"`
	Type   string `json:"type"`
	Value  string `json:"value"`
	ZoneID string `json:"zone_id"`
}

type Records struct {
	Records []Record `json:"records"`
}

type Zone struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Zones struct {
	Zones []Zone `json:"zones"`
}
