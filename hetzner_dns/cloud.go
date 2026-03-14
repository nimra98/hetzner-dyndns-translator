package hetzner_dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// CloudHetznerDNS implements the DNSProvider interface for the new Hetzner Cloud API (api.hetzner.cloud).
// This API uses Resource Record Sets (RRsets) to group records of the same name and type.
type CloudHetznerDNS struct {
	accessToken string
	client      http.Client
}

// NewCloudHetznerDNS creates a new client for the Hetzner Cloud DNS API.
func NewCloudHetznerDNS(accessToken string) *CloudHetznerDNS {
	return &CloudHetznerDNS{
		accessToken: accessToken,
		client:      http.Client{},
	}
}

// findZone retrieves zone information by its name from the Cloud API.
func (h *CloudHetznerDNS) findZone(zoneName string) (*Zone, error) {
	req, err := http.NewRequest("GET", "https://api.hetzner.cloud/v1/zones", nil)
	if err != nil {
		return nil, err
	}
	// Cloud API uses Bearer authentication
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

// findRRset searches for a Resource Record Set by name and automatically determines the type (A vs AAAA).
func (h *CloudHetznerDNS) findRRset(zoneId, recordName string, value string) (*RRset, error) {
	// Determine DNS record type based on the IP address format
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

	// The root domain is often represented as "@" in DNS APIs
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

// updateRRset updates the records of an existing RRset using the 'set_records' action.
func (h *CloudHetznerDNS) updateRRset(zoneId string, rrset RRset) error {
	// The 'set_records' action only expects the 'records' array in the body
	update := struct {
		Records []Value `json:"records"`
	}{
		Records: rrset.Records,
	}

	data, err := json.Marshal(update)
	if err != nil {
		return err
	}

	// Action-based endpoint for targeted record updates
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

	// 201 Created is the standard response for successful action execution
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update RRset, status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// PatchRecord implements the DNSProvider interface for the Cloud API.
func (h *CloudHetznerDNS) PatchRecord(zoneName, recordName, value string) error {
	zone, err := h.findZone(zoneName)
	if err != nil {
		return err
	}

	rrset, err := h.findRRset(zone.GetId(), recordName, value)
	if err != nil {
		return err
	}

	// For DynDNS, we replace all existing records in the set with the single new IP
	rrset.Records = []Value{{Value: value}}

	return h.updateRRset(zone.GetId(), *rrset)
}
