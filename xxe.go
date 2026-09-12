package main

import "encoding/xml"

// XXE vulnerability: xml.Unmarshal on untrusted input
func parseXMLPayload(data []byte) error {
	var result struct {
		XMLName xml.Name `xml:"root"`
		Data    string   `xml:",chardata"`
	}
	return xml.Unmarshal(data, &result)
}
