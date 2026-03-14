package hetzner_dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// DNSProvider defines the common interface for different Hetzner DNS API implementations.
// It allows the main application to remain agnostic of the underlying API version (Legacy vs. Cloud).
type DNSProvider interface {
	// PatchRecord finds a DNS record by name in a given zone and updates its value (IP).
	PatchRecord(zoneName, recordName, value string) error
}

// NewHetznerDNS is a factory function that returns the appropriate DNSProvider implementation.
// It switches between 'legacy' and 'cloud' based on the HETZNER_API_VERSION environment variable.
func NewHetznerDNS(accessToken string) DNSProvider {
	if os.Getenv("HETZNER_API_VERSION") == "cloud" {
		return NewCloudHetznerDNS(accessToken)
	}
	return NewLegacyHetznerDNS(accessToken)
}

// LegacyHetznerDNS implements the DNSProvider interface for the deprecated Hetzner DNS API (dns.hetzner.com).
type LegacyHetznerDNS struct {
	accessToken string
	client      http.Client
}

// NewLegacyHetznerDNS creates a new client for the legacy DNS API.
func NewLegacyHetznerDNS(accessToken string) *LegacyHetznerDNS {
	return &LegacyHetznerDNS{
		accessToken: accessToken,
		client:      http.Client{},
	}
}

// findZone retrieves zone information by its name from the legacy API.
func (h *LegacyHetznerDNS) findZone(zoneName string) (*Zone, error) {
	req, err := http.NewRequest("GET", "https://dns.hetzner.com/api/v1/zones", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Auth-API-Token", h.accessToken)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var zones Zones
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &zones); err != nil {
		return nil, err
	}

	for _, zone := range zones.Zones {
		if zone.Name == zoneName {
			return &zone, nil
		}
	}

	return nil, fmt.Errorf("Zone not found")
}

// findRecord retrieves a specific DNS record by name within a zone using the legacy API.
func (h *LegacyHetznerDNS) findRecord(zoneId, recordName string) (*Record, error) {
	url := fmt.Sprintf("https://dns.hetzner.com/api/v1/records?zone_id=%s", zoneId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Auth-API-Token", h.accessToken)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var records Records
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &records); err != nil {
		return nil, err
	}

	for _, record := range records.Records {
		if record.Name == recordName {
			return &record, nil
		}
	}

	return nil, fmt.Errorf("Record not found")
}

// updateRecord sends a PUT request to update an existing DNS record in the legacy API.
func (h *LegacyHetznerDNS) updateRecord(record Record) error {
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://dns.hetzner.com/api/v1/records/%s", record.Id)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-API-Token", h.accessToken)

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update record, status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// PatchRecord implements the DNSProvider interface for the legacy API.
func (h *LegacyHetznerDNS) PatchRecord(zoneName, recordName, value string) error {
	zone, err := h.findZone(zoneName)
	if err != nil {
		return err
	}

	record, err := h.findRecord(zone.GetId(), recordName)
	if err != nil {
		return err
	}

	record.Value = value
	return h.updateRecord(*record)
}
