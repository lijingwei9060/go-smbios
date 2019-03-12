// Copyright 2017-2018 DigitalOcean.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Command lsdimms lists memory DIMM information from SMBIOS.
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/lijingwei9060/go-smbios/smbios"
)

func main() {
	// Find SMBIOS data in operating system-specific location.
	rc, ep, err := smbios.Stream()
	if err != nil {
		log.Fatalf("failed to open stream: %v", err)
	}
	// Be sure to close the stream!
	defer rc.Close()

	// Decode SMBIOS structures from the stream.
	d := smbios.NewDecoder(rc)
	ss, err := d.Decode()
	if err != nil {
		log.Fatalf("failed to decode structures: %v", err)
	}

	major, minor, rev := ep.Version()
	fmt.Printf("SMBIOS %d.%d.%d\n", major, minor, rev)

	for _, s := range ss {
		// Only look at memory devices.
		// if s.Header.Type == 17 {
		// 	out, err := smbios.ParseMemoryDevice(s)
		// 	if err != nil {
		// 		fmt.Print(err)
		// 	}
		// 	str, _ := json.Marshal(out)
		// 	fmt.Print(string(str))
		// }

		// Code based on: https://www.dmtf.org/sites/default/files/standards/documents/DSP0134_3.1.1.pdf.

		if s.Header.Type == 0 {
			out, err := smbios.ParseBIOSInformation(s)
			if err != nil {
				fmt.Print(err)
			}
			str, _ := json.Marshal(out)
			fmt.Print(string(str))
		}
		if s.Header.Type == 1 {
			out, err := smbios.ParseSystemInformation(s)
			if err != nil {
				fmt.Print(err)
			}
			str, _ := json.Marshal(out)
			fmt.Print(string(str))
		}
		if s.Header.Type == 2 {
			out, err := smbios.ParseBaseboardInformation(s)
			if err != nil {
				fmt.Print(err)
			}
			str, _ := json.Marshal(out)
			fmt.Print(string(str))
		}
		if s.Header.Type == 3 {
			out, err := smbios.ParseSystemEnclosure(s)
			if err != nil {
				fmt.Print(err)
			}
			str, _ := json.Marshal(out)
			fmt.Print(string(str))
		}
		if s.Header.Type == 4 {
			out, err := smbios.ParseProcessorInformation(s)
			if err != nil {
				fmt.Print(err)
			}
			str, _ := json.Marshal(out)
			fmt.Print(string(str))
		}

	}
}
