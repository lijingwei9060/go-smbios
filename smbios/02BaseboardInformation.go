package smbios

import (
	"encoding/binary"
	"fmt"
)

type BaseboardInformation struct { // 7.3 type 2
	Manufacturer                   string   // 4 String number
	ProductName                    string   // 5 String number
	Version                        string   // 6 String number
	SerialNumber                   string   // 7 String number
	AssetTag                       string   // 8 String number
	FeatureFlags                   []string //9
	LocationInChassis              string   // 10 String number
	ChassisHandle                  uint16   // 11-12
	BoardType                      string   // 13
	NumberOfContainedObjectHandles uint8    // 14
	//ContainedObjectHandles
}

func ParseBaseboardInformation(s *Structure) (*BaseboardInformation, error) { // type 2
	if s == nil {
		return nil, fmt.Errorf("structure s is null")
	}
	if s.Header.Type != 2 {
		return nil, fmt.Errorf("structure s is not type 2, but %d", s.Header.Type)
	}

	ret := &BaseboardInformation{}
	if len(s.Strings) >= 1 {
		ret.Manufacturer = s.Strings[0]
	}
	if len(s.Strings) >= 2 {
		ret.ProductName = s.Strings[1]
	}

	if len(s.Strings) >= 3 {
		ret.Version = s.Strings[2]
	}
	if len(s.Strings) >= 4 {
		ret.SerialNumber = s.Strings[3]
	}
	if len(s.Strings) >= 5 {
		ret.AssetTag = s.Strings[4]
	}
	// features
	featureflag := uint8(s.Formatted[9])
	bit := uint8(0x01)
	for i := uint8(0); i < uint8(len(bifeatures)); i++ {
		if (featureflag>>i)&bit == bit {
			ret.FeatureFlags = append(ret.FeatureFlags, bifeatures[i])
		}
	}
	if len(s.Strings) >= 6 {
		ret.LocationInChassis = s.Strings[5]
	}
	ret.ChassisHandle = binary.LittleEndian.Uint16(s.Formatted[7:9])
	b := int(s.Formatted[9])
	if b > 0 && b < len(bitypes) {
		ret.BoardType = bitypes[b]
	} else {
		ret.BoardType = bitypes[0]
	}
	ret.NumberOfContainedObjectHandles = uint8(s.Formatted[10])

	return ret, nil
}

/* 7.3.1 */
var bifeatures = []string{
	"Board is a hosting board", /* bit 0 */
	"Board requires at least one daughter board",
	"Board is removable",
	"Board is replaceable",
	"Board is hot swappable", /* bit 4 */
}

/* 7.3.2 */
var bitypes = []string{
	"Unknown", /* 0x00 */
	"Unknown", /* 0x01 */
	"Other",
	"Server Blade",
	"Connectivity Switch",
	"System Management Module",
	"Processor Module",
	"I/O Module",
	"Memory Module",
	"Daughter Board",
	"Motherboard",
	"Processor+Memory Module",
	"Processor+I/O Module",
	"Interconnect Board", /* 0x0D */
}
