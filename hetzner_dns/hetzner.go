package hetzner_dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HetznerDNS struct {
	accessToken string
	client      http.Client
}

func NewHetznerDNS(accessToken string) *HetznerDNS {
	client := &http.Client{}
	return &HetznerDNS{
		accessToken: accessToken,
		client:      *client,
	}
}

func (h *HetznerDNS) findZone(zoneName string) (*Zone, error) {
	// Get Record (GET https://dns.hetzner.com/api/v1/zones)

	// Create request
	req, err := http.NewRequest("GET", "https://dns.hetzner.com/api/v1/zones", nil)
	if err != nil {
		return nil, err
	}

	// Headers
	req.Header.Add("Auth-API-Token", h.accessToken)

	// Fetch Request
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

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

func (h *HetznerDNS) findRecord(zoneId, recordName string) (*Record, error) {
	// Get Record (GET https://dns.hetzner.com/api/v1/records?zone_id=)

	// Create request
	url := fmt.Sprintf("https://dns.hetzner.com/api/v1/records?zone_id=%s", zoneId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Headers
	req.Header.Add("Auth-API-Token", h.accessToken)

	// Fetch Request
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

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

func (h *HetznerDNS) updateRecord(record Record) error {
	// Update Record (PUT https://dns.hetzner.com/api/v1/records/{RecordID})

	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	// Create request
	url := fmt.Sprintf("https://dns.hetzner.com/api/v1/records/%s", record.Id)
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return err
	}

	// Headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-API-Token", h.accessToken)

	// Fetch Request
	_, err = h.client.Do(req)

	if err != nil {
		return err
	}

	return nil
}

func (h *HetznerDNS) PatchRecord(zoneName, recordName, value string) error {
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
