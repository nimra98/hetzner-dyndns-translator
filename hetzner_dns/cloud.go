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

func (h *CloudHetznerDNS) findRRset(zoneId, recordName string, value string) (*RRset, error) {
	// Determine type based on IP value
	recordType := "A"
	if strings.Contains(value, ":") {
		recordType = "AAAA"
	}

	url := fmt.Sprintf("https://api.hetzner.cloud/v1/zones/%s/rrsets", zoneId)
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
	var rrsets RRsets
	err = json.Unmarshal(body, &rrsets)
	if err != nil {
		return nil, err
	}

	// In Cloud API, root record is often "@"
	searchName := recordName
	if recordName == "" || recordName == "@" {
		searchName = "@"
	}

	for _, rrset := range rrsets.RRsets {
		if rrset.Name == searchName && rrset.Type == recordType {
			return &rrset, nil
		}
	}

	return nil, fmt.Errorf("RRset not found (name: %s, type: %s)", searchName, recordType)
}

func (h *CloudHetznerDNS) updateRRset(zoneId string, rrset RRset) error {
	// Create update request with only records field as required by set_records action
	update := struct {
		Records []Value `json:"records"`
	}{
		Records: rrset.Records,
	}

	data, err := json.Marshal(update)
	if err != nil {
		return err
	}

	// Correct endpoint for setting records in an RRset
	url := fmt.Sprintf("https://api.hetzner.cloud/v1/zones/%s/rrsets/%s/%s/actions/set_records",
		zoneId, rrset.Name, rrset.Type)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update RRset, status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (h *CloudHetznerDNS) PatchRecord(zoneName, recordName, value string) error {
	zone, err := h.findZone(zoneName)
	if err != nil {
		return err
	}

	rrset, err := h.findRRset(zone.GetId(), recordName, value)
	if err != nil {
		return err
	}

	// Update records in RRset
	rrset.Records = []Value{{Value: value}}

	return h.updateRRset(zone.GetId(), *rrset)
}
