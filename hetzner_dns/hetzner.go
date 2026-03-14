package hetzner_dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// DNSProvider defines the interface for interacting with various Hetzner DNS APIs.
type DNSProvider interface {
	PatchRecord(zoneName, recordName, value string) error
}

// NewHetznerDNS creates a new DNSProvider based on the HETZNER_API_VERSION environment variable.
func NewHetznerDNS(accessToken string) DNSProvider {
	apiVersion := os.Getenv("HETZNER_API_VERSION")
	if apiVersion == "cloud" {
		return NewCloudHetznerDNS(accessToken)
	}
	return NewLegacyHetznerDNS(accessToken)
}

// LegacyHetznerDNS represents a client for interacting with the legacy Hetzner DNS API.
type LegacyHetznerDNS struct {
	accessToken string
	client      http.Client
}

// NewLegacyHetznerDNS creates a new instance of LegacyHetznerDNS.
func NewLegacyHetznerDNS(accessToken string) *LegacyHetznerDNS {
	return &LegacyHetznerDNS{
		accessToken: accessToken,
		client:      http.Client{},
	}
}

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
	respBody, _ := io.ReadAll(resp.Body)

	var zones Zones
	err = json.Unmarshal(respBody, &zones)
	if err != nil {
		return nil, err
	}

	for _, zone := range zones.Zones {
		if zone.Name == zoneName {
			return &zone, nil
		}
	}

	return nil, fmt.Errorf("Zone not found")
}

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
	respBody, _ := io.ReadAll(resp.Body)

	var records Records
	err = json.Unmarshal(respBody, &records)
	if err != nil {
		return nil, err
	}

	for _, record := range records.Records {
		if record.Name == recordName {
			return &record, nil
		}
	}

	return nil, fmt.Errorf("Record not found")
}

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

func (h *LegacyHetznerDNS) PatchRecord(zoneName, recordName, value string) error {
	zone, err := h.findZone(zoneName)
	if err != nil {
		return err
	}

	record, err := h.findRecord(zone.Id, recordName)
	if err != nil {
		return err
	}

	record.Value = value
	return h.updateRecord(*record)
}
