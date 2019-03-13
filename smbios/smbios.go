package smbios

import (
	"fmt"
	"log"
)

type SMBIOS struct {
	Major                 int
	Minor                 int
	Revision              int
	BIOSInformation       *BIOSInformation         // type 0
	SystemInformation     *SystemInformation       // type 1
	BaseboardInformations []*BaseboardInformation  // type 2
	SystemEnclosures      []*SystemEnclosure       // type 3
	ProcessorInformations []*ProcessorInformation  // type 4
	MemoryDevices         []*MemoryDeviceStructure // type 17
}

func GetSMBIOS() *SMBIOS {
	// Find SMBIOS data in operating system-specific location.
	rc, ep, err := Stream()
	if err != nil {
		log.Fatalf("failed to open stream: %v", err)
	}
	// Be sure to close the stream!
	defer rc.Close()

	// Decode SMBIOS structures from the stream.
	d := NewDecoder(rc)
	ss, err := d.Decode()
	if err != nil {
		log.Fatalf("failed to decode structures: %v", err)
	}

	major, minor, rev := ep.Version()
	ret := &SMBIOS{}
	ret.Major = major
	ret.Minor = minor
	ret.Revision = rev

	for _, s := range ss {
		// Code based on: https://www.dmtf.org/sites/default/files/standards/documents/DSP0134_3.1.1.pdf.

		if s.Header.Type == 0 {
			out, err := ParseBIOSInformation(s)
			if err != nil {
				fmt.Print(err)
			}
			ret.BIOSInformation = out
		}
		if s.Header.Type == 1 {
			out, err := ParseSystemInformation(s)
			if err != nil {
				fmt.Print(err)
			}
			ret.SystemInformation = out
		}
		if s.Header.Type == 2 {
			out, err := ParseBaseboardInformation(s)
			if err != nil {
				fmt.Print(err)
			}
			ret.BaseboardInformations = append(ret.BaseboardInformations, out)
		}
		if s.Header.Type == 3 {
			out, err := ParseSystemEnclosure(s)
			if err != nil {
				fmt.Print(err)
			}

			ret.SystemEnclosures = append(ret.SystemEnclosures, out)
		}
		if s.Header.Type == 4 {
			out, err := ParseProcessorInformation(s)
			if err != nil {
				fmt.Print(err)
			}
			ret.ProcessorInformations = append(ret.ProcessorInformations, out)
		}

		if s.Header.Type == 17 {
			out, err := ParseMemoryDevice(s)
			if err != nil {
				fmt.Print(err)
			}
			ret.MemoryDevices = append(ret.MemoryDevices, out)
		}

	}
	return ret
}
