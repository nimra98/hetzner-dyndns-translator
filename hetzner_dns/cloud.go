package hetzner_dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	err = json.Unmarshal([]byte{}, &zones) // Will be populated from body
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

func (h *CloudHetznerDNS) findRRset(zoneId, recordName string) (*RRset, error) {
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

	for _, rrset := range rrsets.RRsets {
		// In Cloud API, recordName might be just the subdomain or the full name.
		// Usually it's just the subdomain.
		if rrset.Name == recordName {
			return &rrset, nil
		}
	}

	return nil, fmt.Errorf("RRset not found")
}

func (h *CloudHetznerDNS) updateRRset(zoneId string, rrset RRset) error {
	data, err := json.Marshal(rrset)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.hetzner.cloud/v1/zones/%s/rrsets/%s", zoneId, rrset.GetId())
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
		return fmt.Errorf("failed to update RRset, status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (h *CloudHetznerDNS) PatchRecord(zoneName, recordName, value string) error {
	zone, err := h.findZone(zoneName)
	if err != nil {
		return err
	}

	rrset, err := h.findRRset(zone.GetId(), recordName)
	if err != nil {
		return err
	}

	// For DynDNS, we typically have only one record in the set.
	// We replace all records in the set with the new IP.
	rrset.Records = []Value{{Value: value}}

	return h.updateRRset(zone.GetId(), *rrset)
}
