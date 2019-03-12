package smbios

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// A MemoryDeviceStructure is an SMBIOS structure.
type MemoryDeviceStructure struct {
	Header                                  Header
	PhysicalMemoryArrayHandle               uint16
	MemoryErrorInformationHandle            uint16
	TotalWidth                              uint16 /* 2.1+ */
	DataWidth                               uint16
	Size                                    int
	FormFactor                              string
	DeviceSet                               string
	DeviceLocator                           string
	BankLocator                             string
	MemoryType                              string
	TypeDetail                              string
	Speed                                   uint16 /* 2.3+ */
	Manufacturer                            string
	SerialNumber                            string
	AssetTag                                string
	PartNumber                              string
	Attributes                              uint8  /* 2.6+ */
	ExtendedSize                            uint32 /* 2.7+ */
	ConfiguredMemoryClockSpeed              uint16
	MinimumVoltage                          uint16 /* 2.8+ */
	MaximumVoltage                          uint16
	ConfiguredVoltage                       uint16
	MemoryTechnology                        string /* 3.2+ */
	MemoryOperatingModeCapability           string
	FirmwareVersion                         string
	ModuleManufacturerID                    uint16
	ModuleProductID                         uint16
	MemorySubsystemControllerManufacturerID uint16
	MemorySubsystemControllerProductID      uint16
	NonVolatileSize                         uint32
	VolatileSize                            uint32
	CacheSize                               uint32
	LogicalSize                             uint32
}

// Memory_Device_Factor 设备接口类型？
var Memory_Device_Factor = []string{ /* 7.18.1 */
	"Unknown", /* 0x00 为了方便处理*/
	"Other",   /* 0x01 */
	"Unknown",
	"SIMM",
	"SIP",
	"Chip",
	"DIP",
	"ZIP",
	"Proprietary Card",
	"DIMM",
	"TSOP",
	"Row Of Chips",
	"RIMM",
	"SODIMM",
	"SRIMM",
	"FB-DIMM", /* 0x0F */
}

// getMemoryDeviceFactor 获取内存接口类型
func getMemoryDeviceFactor(n int) string {
	if n < 0 || n > len(Memory_Device_Factor) {
		return Memory_Device_Factor[0]
	}
	return Memory_Device_Factor[n]
}

// Memory_Device_Type 内存类型
var Memory_Device_Type = []string{ /* 7.18.2 */
	"Unknown", /* 0x00 为了方便处理*/
	"Other",   /* 0x01 */
	"Unknown",
	"DRAM",
	"EDRAM",
	"VRAM",
	"SRAM",
	"RAM",
	"ROM",
	"Flash",
	"EEPROM",
	"FEPROM",
	"EPROM",
	"CDRAM",
	"3DRAM",
	"SDRAM",
	"SGRAM",
	"RDRAM",
	"DDR",
	"DDR2",
	"DDR2 FB-DIMM",
	"Reserved",
	"Reserved",
	"Reserved",
	"DDR3",
	"FBD2",
	"DDR4",
	"LPDDR",
	"LPDDR2",
	"LPDDR3",
	"LPDDR4",
	"Logical non-volatile device", /* 0x1F */
}

// getMemoryDeviceType 获取内存类型
func getMemoryDeviceType(n int) string {
	if n < 0 || n > len(Memory_Device_Type) {
		return Memory_Device_Type[0]
	}
	return Memory_Device_Type[n]
}

// Memory_Device_Detail 内存信息明细
var Memory_Device_Detail = map[int]string{ /* 7.18.3 */
	0:     "Unknown",
	1:     "Other", /* 1 */
	2:     "Unknown",
	4:     "Fast-paged",
	8:     "Static Column",
	16:    "Pseudo-static",
	32:    "RAMBus",
	64:    "Synchronous",
	128:   "CMOS",
	256:   "EDO",
	512:   "Window DRAM",
	1024:  "Cache DRAM",
	2048:  "Non-Volatile",
	4096:  "Registered (Buffered)",
	8192:  "Unbuffered (Unregistered)",
	16384: "LRDIMM", /* 15 */
}

// Memory_Device_Technology 设备技术 /* 7.18.6 */
var Memory_Device_Technology = []string{
	"Unknown",
	"Other", /* 0x01 */
	"Unknown",
	"DRAM",
	"NVDIMM-N",
	"NVDIMM-F",
	"NVDIMM-P",
	"Intel persistent memory", /* 0x07 */
}

var Memory_Device_Operating_Mode_Capability = []string{ /* 7.18.7 */
	"Other", /* 1 */
	"Unknown",
	"Volatile memory",
	"Byte-accessible persistent memory",
	"Block-accessible persistent memory", /* 5 */
}

// getMemoryDeviceDetail 获取内存设备信息详情
func getMemoryDeviceDetail(n int) string {
	out, exists := Memory_Device_Detail[n]
	if !exists {
		return Memory_Device_Detail[0]
	}
	return out
}

// ParseMemoryDevice 解析MemoryDevice
func ParseMemoryDevice(s *Structure) (*MemoryDeviceStructure, error) {
	if s == nil {
		return nil, fmt.Errorf("parameter s is null")
	}
	if s.Header.Type != 17 {
		return nil, fmt.Errorf("parameter is not a memory device")
	}

	ret := &MemoryDeviceStructure{}
	ret.Header = s.Header                                                           // 0-3
	ret.PhysicalMemoryArrayHandle = binary.LittleEndian.Uint16(s.Formatted[:2])     // 0-1
	ret.MemoryErrorInformationHandle = binary.LittleEndian.Uint16(s.Formatted[2:4]) // 2-3
	ret.TotalWidth = binary.LittleEndian.Uint16(s.Formatted[4:6])                   // 4-5
	ret.DataWidth = binary.LittleEndian.Uint16(s.Formatted[6:8])                    // 6-7

	// Only parse the DIMM size.
	dimmSize := int(binary.LittleEndian.Uint16(s.Formatted[8:10]))
	//If the DIMM size is 32GB or greater, we need to parse the extended field.
	// Spec says 0x7fff in regular size field means we should parse the extended.
	if dimmSize == 0x7fff {
		dimmSize = int(binary.LittleEndian.Uint32(s.Formatted[24:28]))
	}

	// The granularity in which the value is specified
	// depends on the setting of the most-significant bit (bit
	// 15). If the bit is 0, the value is specified in megabyte
	// units; if the bit is 1, the value is specified in kilobyte
	// units.
	//
	// Little endian MSB for uint16 is in second byte.

	if s.Formatted[9]&0x80 == 0 {
		dimmSize = dimmSize * 1024
	}
	ret.Size = dimmSize // 8-9

	formFactor := int(s.Formatted[10]) //10
	if formFactor < 0 || formFactor > len(Memory_Device_Factor) {
		ret.FormFactor = Memory_Device_Factor[0]
	} else {
		ret.FormFactor = Memory_Device_Factor[formFactor]
	}

	if len(s.Strings) >= 1 {
		ret.DeviceSet = strings.TrimSpace(s.Strings[0]) // 11 String number
	}
	if len(s.Strings) >= 2 {
		ret.DeviceLocator = strings.TrimSpace(s.Strings[1]) // 12 String number
	}
	if len(s.Strings) >= 3 {
		ret.BankLocator = strings.TrimSpace(s.Strings[2]) // 13 String number
	}

	memoryType := int(s.Formatted[14]) // 14
	if memoryType > 0 || memoryType > len(Memory_Device_Type) {
		ret.MemoryType = Memory_Device_Type[0]
	} else {
		ret.MemoryType = Memory_Device_Type[memoryType]
	}

	ret.TypeDetail = Memory_Device_Detail[int(binary.LittleEndian.Uint16(s.Formatted[15:17]))] // 15-16
	ret.Speed = binary.LittleEndian.Uint16(s.Formatted[17:19])                                 // 17-18
	if len(s.Strings) >= 4 {
		ret.Manufacturer = strings.TrimSpace(s.Strings[3]) // 19 String number  2.3+

	}
	if len(s.Strings) >= 5 {
		ret.SerialNumber = strings.TrimSpace(s.Strings[4]) // 20 String number
	}
	if len(s.Strings) >= 6 {
		ret.AssetTag = strings.TrimSpace(s.Strings[5]) // 21 String number
	}
	if len(s.Strings) >= 7 {
		ret.PartNumber = strings.TrimSpace(s.Strings[6]) // 22 String number

	}
	// 2.6+ length > 27
	if s.Header.Length > 27 {
		ret.Attributes = uint8(s.Formatted[23]) // 23 2.6+
	}
	// 2.7+ length > 28
	if s.Header.Length > 28 {
		ret.ExtendedSize = binary.LittleEndian.Uint32(s.Formatted[24:28])               // 24-27 2.7+
		ret.ConfiguredMemoryClockSpeed = binary.LittleEndian.Uint16(s.Formatted[28:30]) //28-29
	}
	// 2.8+ length >34
	if s.Header.Length > 34 {
		ret.MinimumVoltage = binary.LittleEndian.Uint16(s.Formatted[30:32])    // 30-31 2.8+
		ret.MaximumVoltage = binary.LittleEndian.Uint16(s.Formatted[32:34])    // 32-33
		ret.ConfiguredVoltage = binary.LittleEndian.Uint16(s.Formatted[34:36]) // 34-35
	}
	// 3.2+ length >40
	if s.Header.Length > 40 {

		memoryTechnology := int(s.Formatted[36]) // 36 3.2+
		if memoryTechnology < 0 || memoryTechnology > len(Memory_Device_Technology) {
			ret.MemoryTechnology = Memory_Device_Technology[0]
		} else {
			ret.MemoryTechnology = Memory_Device_Technology[memoryTechnology]
		}

		MemoryOperatingModeCapability := int(binary.LittleEndian.Uint16(s.Formatted[37:39])) // 37-38
		if MemoryOperatingModeCapability < 0 || MemoryOperatingModeCapability > len(Memory_Device_Operating_Mode_Capability) {
			ret.MemoryOperatingModeCapability = Memory_Device_Operating_Mode_Capability[0]
		} else {
			ret.MemoryOperatingModeCapability = Memory_Device_Operating_Mode_Capability[MemoryOperatingModeCapability]
		}
		if len(s.Strings) >= 8 {
			ret.FirmwareVersion = strings.TrimSpace(s.Strings[7]) // 39 String number
		}
		ret.ModuleManufacturerID = binary.LittleEndian.Uint16(s.Formatted[40:42])                    // 40-41
		ret.ModuleProductID = binary.LittleEndian.Uint16(s.Formatted[42:44])                         // 42-43
		ret.MemorySubsystemControllerManufacturerID = binary.LittleEndian.Uint16(s.Formatted[44:46]) // 44-45
		ret.MemorySubsystemControllerProductID = binary.LittleEndian.Uint16(s.Formatted[46:48])      //46-47
		ret.NonVolatileSize = binary.LittleEndian.Uint32(s.Formatted[48:52])                         // 48-51
		ret.VolatileSize = binary.LittleEndian.Uint32(s.Formatted[52:56])                            // 52-55
		ret.CacheSize = binary.LittleEndian.Uint32(s.Formatted[56:60])                               // 56-59
		ret.LogicalSize = binary.LittleEndian.Uint32(s.Formatted[60:64])                             // 60-63
	}

	return ret, nil
}
