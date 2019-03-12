package smbios

import (
	"encoding/binary"
	"fmt"
)

type SystemEnclosure struct { // 7.4 type 3
	Header                       Header
	Manufacturer                 string // 4 String number
	Type                         string // 5
	Version                      string // 6 String number
	SerialNumber                 string // 7 String number
	AssetTag                     string // 8 String number
	BootUpState                  string // 9 7.4.2 2.1+
	PowerSupplyState             string // 10 7.4.2
	ThermalState                 string // 11 7.4.2
	SecurityStatus               string // 12 7.4.3
	OEMDefined                   uint32 // 13-16 2.3+
	Height                       uint8  // 17 uint:U 1.75inches or 4.445 cm
	NumberOfPowerCords           uint8  // 18
	ContainedElementCount        uint8  // 19 n
	ContainedElementRecordLength uint8  // 20 m
	ContainedElements            int    // n*m 7.4.4
	SKUNumber                    string // String number 2.7+
}

func parseSystemEnclosure(s *Structure) (*SystemEnclosure, error) {
	if s == nil {
		return nil, fmt.Errorf("structure s is null")
	}
	if s.Header.Type != 3 {
		return nil, fmt.Errorf("structure s is not type 3, but %d", s.Header.Type)
	}

	ret := &SystemEnclosure{}
	ret.Header = s.Header
	if len(s.Strings) >= 1 {
		ret.Manufacturer = s.Strings[0]
	}
	var t int // 临时变量
	t = int(s.Formatted[1])
	if t > 0 && t <= len(Chassis_Type) {
		ret.Type = Chassis_Type[t-1]
	} else {
		ret.Type = Chassis_Type[1]
	}

	if len(s.Strings) >= 2 {
		ret.Version = s.Strings[1]
	}
	if len(s.Strings) >= 3 {
		ret.SerialNumber = s.Strings[2]
	}
	if len(s.Strings) >= 4 {
		ret.AssetTag = s.Strings[3]
	}

	t = int(s.Formatted[5])
	if t > 0 && t <= len(Chassis_State) {
		ret.BootUpState = Chassis_State[t-1]
	} else {
		ret.BootUpState = Chassis_State[1]
	}

	t = int(s.Formatted[6])
	if t > 0 && t <= len(Chassis_State) {
		ret.PowerSupplyState = Chassis_State[t-1]
	} else {
		ret.PowerSupplyState = Chassis_State[1]
	}

	t = int(s.Formatted[7])
	if t > 0 && t <= len(Chassis_State) {
		ret.ThermalState = Chassis_State[t-1]
	} else {
		ret.ThermalState = Chassis_State[1]
	}

	t = int(s.Formatted[8])
	if t > 0 && t <= len(Chassis_Security_State) {
		ret.SecurityStatus = Chassis_Security_State[t-1]
	} else {
		ret.SecurityStatus = Chassis_Security_State[1]
	}
	ret.OEMDefined = binary.LittleEndian.Uint32(s.Formatted[13:17])
	ret.Height = uint8(s.Formatted[17])
	ret.NumberOfPowerCords = uint8(s.Formatted[18])
	ret.ContainedElementCount = uint8(s.Formatted[19])
	ret.ContainedElementRecordLength = uint8(s.Formatted[20])
	if len(s.Strings) >= 5 {
		ret.SKUNumber = s.Strings[4]
	}
	return ret, nil
}

var Chassis_Type = []string{
	"Other", /* 0x01 */
	"Unknown",
	"Desktop",
	"Low Profile Desktop",
	"Pizza Box",
	"Mini Tower",
	"Tower",
	"Portable",
	"Laptop",
	"Notebook",
	"Hand Held",
	"Docking Station",
	"All In One",
	"Sub Notebook",
	"Space-saving",
	"Lunch Box",
	"Main Server Chassis", /* CIM_Chassis.ChassisPackageType says "Main System Chassis" */
	"Expansion Chassis",
	"Sub Chassis",
	"Bus Expansion Chassis",
	"Peripheral Chassis",
	"RAID Chassis",
	"Rack Mount Chassis",
	"Sealed-case PC",
	"Multi-system",
	"CompactPCI",
	"AdvancedTCA",
	"Blade",
	"Blade Enclosing",
	"Tablet",
	"Convertible",
	"Detachable",
	"IoT Gateway",
	"Embedded PC",
	"Mini PC",
	"Stick PC", /* 0x24 */
}

var Chassis_State = []string{ // 7.4.2
	"Other", /* 0x01 */
	"Unknown",
	"Safe",
	"Warning",
	"Critical",
	"Non-recoverable", /* 0x06 */
}

var Chassis_Security_State = []string{ // 7.4.3
	"Other", /* 0x01 */
	"Unknown",
	"None",
	"External Interface Locked Out",
	"External Interface Enabled", /* 0x05 */
}
