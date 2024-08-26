package hetzner_dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HetznerDNS represents a client for interacting with the Hetzner DNS API.
//
// Fields:
// - accessToken: A string containing the access token for authenticating with the Hetzner DNS API.
// - client: An http.Client used to send requests to the Hetzner DNS API.
//
// This struct is used to encapsulate the necessary information and methods for making API requests to Hetzner's DNS service.
type HetznerDNS struct {
	accessToken string
	client      http.Client
}

// NewHetznerDNS creates a new instance of HetznerDNS with the provided access token.
//
// Parameters:
// - accessToken: A string containing the access token for authenticating with the Hetzner DNS API.
//
// Returns:
// - *HetznerDNS: A pointer to a new HetznerDNS instance initialized with the provided access token and an HTTP client.
//
// This function initializes an HTTP client and returns a HetznerDNS struct with the provided access token and the initialized client.
func NewHetznerDNS(accessToken string) *HetznerDNS {
	client := &http.Client{}
	return &HetznerDNS{
		accessToken: accessToken,
		client:      *client,
	}
}

// findZone searches for a DNS zone by its name using the Hetzner DNS API.
//
// Parameters:
// - zoneName: The name of the zone to search for.
//
// Returns:
// - *Zone: A pointer to the found Zone struct, if any.
// - error: An error if the zone is not found or another error occurs.
//
// This function sends a GET request to the Hetzner DNS API to retrieve all zones.
// It then searches the response to find the zone with the specified name.
// If the zone is found, it returns a pointer to the Zone struct; otherwise, it returns an error.
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

// findRecord searches for a DNS record in a specific zone.
//
// Parameters:
// - zoneId: The ID of the zone in which to search for the record.
// - recordName: The name of the record to search for.
//
// Returns:
// - *Record: A pointer to the found record, if any.
// - error: An error if the record is not found or another error occurs.
//
// This function sends a GET request to the Hetzner DNS API to retrieve all records in the specified zone.
// It then searches the response to find the record with the specified name.
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

// updateRecord updates a DNS record using the Hetzner DNS API.
//
// Parameters:
// - record: The Record struct containing the updated information for the DNS record.
//
// Returns:
// - error: An error if the update operation fails.
//
// This function sends a PUT request to the Hetzner DNS API to update the specified DNS record.
// It marshals the record into JSON, creates the request with the appropriate headers, and sends it.
// If any step fails, it returns an error.
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

// PatchRecord updates the value of a DNS record in a specified zone.
//
// Parameters:
// - zoneName: The name of the zone containing the DNS record.
// - recordName: The name of the DNS record to update.
// - value: The new value to set for the DNS record.
//
// Returns:
// - error: An error if the zone or record is not found, or if the update operation fails.
//
// This function first finds the zone by its name using the findZone method.
// It then finds the record within that zone using the findRecord method.
// The record's value is updated to the provided value, and the updateRecord method is called to apply the change.
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
