package hetzner_dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// CloudHetznerDNS represents a client for the new Hetzner Cloud DNS API.
type CloudHetznerDNS struct {
	accessToken string
	client      http.Client
}

func NewCloudHetznerDNS(accessToken string) *CloudHetznerDNS {
	return &CloudHetznerDNS{
		accessToken: accessToken,
		client:      http.Client{},
	}
}

func (h *CloudHetznerDNS) findZone(zoneName string) (*Zone, error) {
	req, err := http.NewRequest("GET", "https://api.hetzner.cloud/v1/zones", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+h.accessToken)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch zones, status: %d", resp.StatusCode)
	}

	var zones Zones
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &zones)
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

func (h *CloudHetznerDNS) findRecord(zoneId, recordName string, value string) (*CloudRecord, error) {
	// Determine type based on IP value
	recordType := "A"
	if strings.Contains(value, ":") {
		recordType = "AAAA"
	}

	url := fmt.Sprintf("https://api.hetzner.cloud/v1/zones/%s/records", zoneId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+h.accessToken)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var records CloudRecords
	err = json.Unmarshal(body, &records)
	if err != nil {
		return nil, err
	}

	// In Cloud API, root record is often "@"
	searchName := recordName
	if recordName == "" || recordName == "@" {
		searchName = "@"
	}

	for _, record := range records.Records {
		if record.Name == searchName && record.Type == recordType {
			return &record, nil
		}
	}

	return nil, fmt.Errorf("Record not found (name: %s, type: %s)", searchName, recordType)
}

func (h *CloudHetznerDNS) updateRecord(zoneId string, record CloudRecord) error {
	// Create update request (Hetzner Cloud API uses the same format for PUT /records/{id})
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.hetzner.cloud/v1/zones/%s/records/%s", zoneId, record.GetId())
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+h.accessToken)

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

func (h *CloudHetznerDNS) PatchRecord(zoneName, recordName, value string) error {
	zone, err := h.findZone(zoneName)
	if err != nil {
		return err
	}

	record, err := h.findRecord(zone.GetId(), recordName, value)
	if err != nil {
		return err
	}

	// Update value
	record.Value = value

	return h.updateRecord(zone.GetId(), *record)
}
