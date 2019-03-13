package smbios

import (
	"fmt"
	"strings"
)

type SystemInformation struct {
	Manufacturer string // 4 String number
	ProductName  string // 5 String number
	Version      string // 6 String number
	SerialNumber string // 7 String number
	UUID         string // 8-23 UUID 16bytes 2.1+
	WakeUpType   string // 24
	SKUNumber    string // 25 String number 2.4+
	Family       string //26 Stringnumber
}

func ParseSystemInformation(s *Structure) (*SystemInformation, error) { // type 1
	if s == nil {
		return nil, fmt.Errorf("structure s is null")
	}
	if s.Header.Type != 1 {
		return nil, fmt.Errorf("structure s is not type 1, but %d", s.Header.Type)
	}

	ret := &SystemInformation{}
	if len(s.Strings) >= 1 {
		ret.Manufacturer = strings.TrimSpace(s.Strings[0])
	}
	if len(s.Strings) >= 2 {
		ret.ProductName = strings.TrimSpace(s.Strings[1])
	}

	if len(s.Strings) >= 3 {
		ret.Version = strings.TrimSpace(s.Strings[2])
	}
	if len(s.Strings) >= 4 {
		ret.SerialNumber = strings.TrimSpace(s.Strings[3])
	}
	/*
	 * As of version 2.6 of the SMBIOS specification, the first 3
	 * fields of the UUID are supposed to be encoded on little-endian.
	 * The specification says that this is the defacto standard,
	 * however I've seen systems following RFC 4122 instead and use
	 * network byte order, so I am reluctant to apply the byte-swapping
	 * for older versions.
	 */
	p := []uint8(s.Formatted[4:20])
	ret.UUID = fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x", p[3], p[2], p[1], p[0], p[5], p[4], p[7], p[6],
		p[8], p[9], p[10], p[11], p[12], p[13], p[14], p[15])

	w := int(s.Formatted[20])
	if w > 0 && w < len(wakeUpType) {
		ret.WakeUpType = wakeUpType[w]
	} else {
		ret.WakeUpType = wakeUpType[2]
	}

	if len(s.Strings) >= 5 {
		ret.SKUNumber = strings.TrimSpace(s.Strings[4])
	}
	if len(s.Strings) >= 6 {
		ret.Family = strings.TrimSpace(s.Strings[5])
	}
	return ret, nil
}

var wakeUpType = []string{
	"Reserved", /* 0x00 */
	"Other",
	"Unknown",
	"APM Timer",
	"Modem Ring",
	"LAN Remote",
	"Power Switch",
	"PCI PME#",
	"AC Power Restored", /* 0x08 */
}
