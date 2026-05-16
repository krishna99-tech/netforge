package snmp

import (
	"fmt"
	"time"

	gosnmp "github.com/gosnmp/gosnmp"
)

// WalkResult holds a single OID result from an SNMP walk.
type WalkResult struct {
	OID   string `json:"oid"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// WalkOptions configures the SNMP session.
type WalkOptions struct {
	Target    string
	Port      uint16
	Community string
	Version   int
	Timeout   time.Duration
	Retries   int
	RootOID   string
}

func buildClient(opts WalkOptions) *gosnmp.GoSNMP {
	version := gosnmp.Version2c
	if opts.Version == 1 {
		version = gosnmp.Version1
	}

	return &gosnmp.GoSNMP{
		Target:    opts.Target,
		Port:      opts.Port,
		Community: opts.Community,
		Version:   version,
		Timeout:   opts.Timeout,
		Retries:   opts.Retries,
	}
}

// Walk performs a BulkWalk from the given root OID.
func Walk(opts WalkOptions) ([]WalkResult, error) {
	client := buildClient(opts)
	if err := client.Connect(); err != nil {
		return nil, fmt.Errorf("SNMP connect failed: %v", err)
	}
	defer client.Conn.Close()

	var results []WalkResult
	err := client.BulkWalk(opts.RootOID, func(pdu gosnmp.SnmpPDU) error {
		results = append(results, WalkResult{
			OID:   pdu.Name,
			Type:  pdu.Type.String(),
			Value: fmt.Sprintf("%v", gosnmp.ToBigInt(pdu.Value)),
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("SNMP walk failed: %v", err)
	}
	return results, nil
}

// Get retrieves a single OID value.
func Get(opts WalkOptions, oid string) (*WalkResult, error) {
	client := buildClient(opts)
	if err := client.Connect(); err != nil {
		return nil, fmt.Errorf("SNMP connect failed: %v", err)
	}
	defer client.Conn.Close()

	result, err := client.Get([]string{oid})
	if err != nil {
		return nil, fmt.Errorf("SNMP get failed: %v", err)
	}
	if len(result.Variables) == 0 {
		return nil, fmt.Errorf("no result returned for OID %s", oid)
	}

	pdu := result.Variables[0]
	return &WalkResult{
		OID:   pdu.Name,
		Type:  pdu.Type.String(),
		Value: fmt.Sprintf("%v", gosnmp.ToBigInt(pdu.Value)),
	}, nil
}
