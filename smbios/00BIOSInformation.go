package smbios

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type BIOSInformation struct { // type 0
	Vendor                     string // 4 String number
	BIOSVersion                string // 5 String number
	BIOSStartingAddressSegment uint16 // 6-7
	BIOSReleaseDate            string // 8 String number
	// Size (n) where 64K * (n+1) is the size of the physical device containing the BIOS, in bytes.
	// FFh - size is 16MB or greater, see Extended BIOS ROM Size for actual size
	BIOSROMSize                            int      // 9
	BIOSCharacteristics                    []string // 10-17
	BIOSCharacteristicsExtensionBytes      uint16   // 18-19 2.4+
	SystemBIOSMajorRelease                 uint8    // 20
	SystemBIOSMinorRelease                 uint8    // 21
	EmbeddedControllerFirmwareMajorRelease uint8    // 22
	EmbeddedControllerFirmwareMinorRelease uint8    // 23
	ExtendedBIOSROMSize                    uint16   // 24
}

// ParseBIOSInformation 解析structure
func ParseBIOSInformation(s *Structure) (*BIOSInformation, error) {
	if s == nil {
		return nil, fmt.Errorf("structure s is null")
	}

	if s.Header.Type != 0 {
		return nil, fmt.Errorf("sturcture s is not a BiosInformation type 0 ,but %d", s.Header.Type)
	}

	ret := &BIOSInformation{}
	if len(s.Strings) >= 1 {
		ret.Vendor = strings.TrimSpace(s.Strings[0])
	}

	if len(s.Strings) >= 2 {
		ret.BIOSVersion = strings.TrimSpace(s.Strings[1])
	}
	ret.BIOSStartingAddressSegment = binary.LittleEndian.Uint16(s.Formatted[2:4])

	if len(s.Strings) >= 3 {
		ret.BIOSReleaseDate = strings.TrimSpace(s.Strings[2])
	}

	// unit :KB
	BIOSROMSize := uint8(s.Formatted[5])
	if BIOSROMSize != 0xff || s.Header.Length <= 24 { // 小于版本3.1或者小于16M
		ret.BIOSROMSize = (int(s.Formatted[5]) + 1) * 64
	} else {
		ret.BIOSROMSize = (int(binary.LittleEndian.Uint16(s.Formatted[20:22])) + 1) * 64
	}

	BIOSCharacteristics := binary.LittleEndian.Uint64(s.Formatted[6:14])
	bit := uint64(0x01)
	for i := uint(0); i <= 31; i++ {
		if (BIOSCharacteristics>>i)&bit == bit {
			ret.BIOSCharacteristics = append(ret.BIOSCharacteristics, BIOSInformationCharacter[i])
		}
	}
	if s.Header.Length > 18 { // 2.4+
		ceb := binary.LittleEndian.Uint16(s.Formatted[14:16])
		bit16 := uint16(0x01)
		for i := uint(0); i <= 12; i++ {
			if (ceb>>i)&bit16 == bit16 {
				ret.BIOSCharacteristics = append(ret.BIOSCharacteristics, BIOSCharacteristicsExtensionBytes[i])
			}
		}
		ret.BIOSCharacteristicsExtensionBytes = binary.LittleEndian.Uint16(s.Formatted[14:16])
		ret.SystemBIOSMajorRelease = uint8(s.Formatted[16])
		ret.SystemBIOSMinorRelease = uint8(s.Formatted[17])
		ret.EmbeddedControllerFirmwareMajorRelease = uint8(s.Formatted[18])
		ret.EmbeddedControllerFirmwareMinorRelease = uint8(s.Formatted[19])
	}
	if s.Header.Length > 24 { // 3.1+
		ret.ExtendedBIOSROMSize = binary.LittleEndian.Uint16(s.Formatted[20:22])
	}

	return ret, nil
}

var BIOSInformationCharacter = []string{
	"Reserved",                           /* bit 0 */
	"Reserved",                           /* bit 1 */
	"Unknown",                            /* bit 2 */
	"BIOS characteristics not supported", /* bit 3 */
	"ISA is supported",
	"MCA is supported",
	"EISA is supported",
	"PCI is supported",
	"PC Card (PCMCIA) is supported",
	"PNP is supported",
	"APM is supported",
	"BIOS is upgradeable",
	"BIOS shadowing is allowed",
	"VLB is supported",
	"ESCD support is available",
	"Boot from CD is supported",
	"Selectable boot is supported",
	"BIOS ROM is socketed",
	"Boot from PC Card (PCMCIA) is supported",
	"EDD is supported",
	"Japanese floppy for NEC 9800 1.2 MB is supported (int 13h)",
	"Japanese floppy for Toshiba 1.2 MB is supported (int 13h)",
	"5.25\"/360 kB floppy services are supported (int 13h)",
	"5.25\"/1.2 MB floppy services are supported (int 13h)",
	"3.5\"/720 kB floppy services are supported (int 13h)",
	"3.5\"/2.88 MB floppy services are supported (int 13h)",
	"Print screen service is supported (int 5h)",
	"8042 keyboard services are supported (int 9h)",
	"Serial services are supported (int 14h)",
	"Printer services are supported (int 17h)",
	"CGA/mono video services are supported (int 10h)",
	"NEC PC-98", /* 31 */
}

var BIOSCharacteristicsExtensionBytes = []string{
	"ACPI is supported", /* bit 0 */
	"USB legacy is supported",
	"AGP is supported",
	"I2O boot is supported",
	"LS-120 boot is supported",
	"ATAPI Zip drive boot is supported",
	"IEEE 1394 boot is supported",
	"Smart battery is supported",           /* bit 7 */
	"BIOS boot specification is supported", /* bit 8  */
	"Function key-initiated network boot is supported",
	"Targeted content distribution is supported",
	"UEFI is supported",
	"System is a virtual machine", /* bit 12 */
}
