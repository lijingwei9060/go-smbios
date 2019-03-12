package smbios

import (
	"encoding/binary"
	"fmt"
)

// TODO: 涉及到的参数太复杂，现在先不做太详细，很多细节没处理
type ProcessorInformation struct { // 7.5 type 4
	Header                   Header
	SocketDesignation        string   // 4 String number
	ProcessorType            string   // 5 7.5.1
	ProcessorFamily          string   // 6 7.5.2
	ProcessorManufacturer    string   // 7 String number
	ProcessorID              string   // 8-15 7.5.3
	ProcessorVersion         string   // 16 String number
	Voltage                  string   // 17 7.5.4
	ExternalClock            uint16   // 18-19
	MaxSpeed                 int      // 20-21 uint:MHz
	CurrentSpeed             int      // 22-23 uint:MHz
	Status                   uint8    // 24
	ProcessorUpgrade         string   // 25 7.5.5
	L1CacheHandle            uint16   // 26-27 2.1+
	L2CacheHandle            uint16   // 28-29
	L3CacheHandle            uint16   // 30-31 2.3+
	SerialNumber             string   // 32 String number
	AssetTag                 string   // 33 String number
	PartNumber               string   // 34 String number
	CoreCount                int      // 35 7.5.6 2.5+
	CoreEnabled              int      // 36 7.5.7
	ThreadCount              int      // 37 7.5.8
	ProcessorCharacteristics []string // 38-39  7.5.9
	ProcessorFamily2         string   // 40-41 7.5.2 2.6+
	CoreCount2               int      // 42-43 7.5.6 3.0+
	CoreEnabled2             int      // 44-45 7.5.7
	ThreadCount2             int      // 46-47 7.5.8
}

func ParseProcessorInformation(s *Structure) (*ProcessorInformation, error) {
	if s == nil {
		return nil, fmt.Errorf("structure s is null")
	}
	if s.Header.Type != 4 {
		return nil, fmt.Errorf("structure s is not type 3, but %d", s.Header.Type)
	}
	ret := &ProcessorInformation{}
	var t int

	ret.Header = s.Header
	if len(s.Strings) >= 1 {
		ret.SocketDesignation = s.Strings[0]
	}

	t = int(s.Formatted[1])
	if t > 0 && t <= len(Processor_Type) {
		ret.ProcessorType = Processor_Type[t-1]
	}

	// TODO: Processor_Family 有一些型号比较模糊，需要具体判断一下
	t = int(s.Formatted[2])
	pt, exists := Processor_Family[t]
	if !exists {
		ret.ProcessorFamily = Processor_Family[2]
	} else {
		ret.ProcessorFamily = pt
	}

	if len(s.Strings) >= 2 {
		ret.ProcessorManufacturer = s.Strings[1]
	}

	pid := s.Formatted[4:12]
	ret.ProcessorID = fmt.Sprintf("%02X %02X %02X %02X %02X %02X %02X %02X", pid[0], pid[1], pid[2], pid[3], pid[4], pid[5], pid[6], pid[7])

	if len(s.Strings) >= 3 {
		ret.ProcessorVersion = s.Strings[2]
	}

	t = int(s.Formatted[13])
	if t >= 0 && t < len(Processor_Voltage) {
		ret.Voltage = Processor_Voltage[t]
	} else {
		ret.Voltage = "Unknown"
	}

	ret.ExternalClock = binary.LittleEndian.Uint16(s.Formatted[14:16])
	ret.MaxSpeed = int(binary.LittleEndian.Uint16(s.Formatted[16:18]))
	ret.CurrentSpeed = int(binary.LittleEndian.Uint16(s.Formatted[18:20]))
	ret.Status = uint8(s.Formatted[20])

	t = int(s.Formatted[21])
	if t > 0 && t <= len(Processor_Upgrade) {
		ret.ProcessorUpgrade = Processor_Upgrade[t-1]
	} else {
		ret.ProcessorUpgrade = Processor_Upgrade[1]
	}

	ret.L1CacheHandle = binary.LittleEndian.Uint16(s.Formatted[22:24])
	ret.L2CacheHandle = binary.LittleEndian.Uint16(s.Formatted[24:26])
	ret.L3CacheHandle = binary.LittleEndian.Uint16(s.Formatted[26:28])

	if len(s.Strings) >= 4 {
		ret.SerialNumber = s.Strings[3]
	}
	if len(s.Strings) >= 5 {
		ret.AssetTag = s.Strings[4]
	}
	if len(s.Strings) >= 6 {
		ret.PartNumber = s.Strings[5]
	}

	cc := uint8(s.Formatted[31])
	if cc != 0xFF {
		ret.CoreCount = int(cc)
	} else {
		ret.CoreCount = int(binary.LittleEndian.Uint16(s.Formatted[38:40])) //count2

	}

	ce := uint8(s.Formatted[32])
	if cc != 0xFF {
		ret.CoreEnabled = int(ce)
	} else {
		ret.CoreEnabled = int(binary.LittleEndian.Uint16(s.Formatted[40:42]))

	}

	tc := uint8(s.Formatted[33])
	if tc != 0xFF {
		ret.ThreadCount = int(tc)
	} else {
		ret.ThreadCount = int(binary.LittleEndian.Uint16(s.Formatted[42:44]))
	}

	t = int(binary.LittleEndian.Uint16(s.Formatted[34:36]))
	for i := uint(0); i < 8; i++ {
		if (t>>i)&0x01 == 0x01 {
			ret.ProcessorCharacteristics = append(ret.ProcessorCharacteristics, Processor_Characteristics[i])
		}
	}

	t = int(binary.LittleEndian.Uint16(s.Formatted[36:38]))
	pf2, exists := Processor_Family[t]
	if !exists {
		ret.ProcessorFamily2 = Processor_Family[2]
	} else {
		ret.ProcessorFamily2 = pf2
	}

	// 3.0+
	if s.Header.Length > 42 {
		ret.CoreCount2 = int(binary.LittleEndian.Uint16(s.Formatted[38:40])) //count2
		ret.CoreEnabled2 = int(binary.LittleEndian.Uint16(s.Formatted[40:42]))
		ret.ThreadCount2 = int(binary.LittleEndian.Uint16(s.Formatted[42:44]))
	}

	return ret, nil
}

var Processor_Type = []string{ /* 7.5.1 */
	"Other", /* 0x01 */
	"Unknown",
	"Central Processor",
	"Math Processor",
	"DSP Processor",
	"Video Processor", /* 0x06 */
}

var Processor_Family = map[int]string{ /* 7.5.2 */
	0x01: "Other",
	0x02: "Unknown",
	0x03: "8086",
	0x04: "80286",
	0x05: "80386",
	0x06: "80486",
	0x07: "8087",
	0x08: "80287",
	0x09: "80387",
	0x0A: "80487",
	0x0B: "Pentium",
	0x0C: "Pentium Pro",
	0x0D: "Pentium II",
	0x0E: "Pentium MMX",
	0x0F: "Celeron",
	0x10: "Pentium II Xeon",
	0x11: "Pentium III",
	0x12: "M1",
	0x13: "M2",
	0x14: "Celeron M",
	0x15: "Pentium 4 HT",
	0x18: "Duron",
	0x19: "K5",
	0x1A: "K6",
	0x1B: "K6-2",
	0x1C: "K6-3",
	0x1D: "Athlon",
	0x1E: "AMD29000",
	0x1F: "K6-2+",
	0x20: "Power PC",
	0x21: "Power PC 601",
	0x22: "Power PC 603",
	0x23: "Power PC 603+",
	0x24: "Power PC 604",
	0x25: "Power PC 620",
	0x26: "Power PC x704",
	0x27: "Power PC 750",
	0x28: "Core Duo",
	0x29: "Core Duo Mobile",
	0x2A: "Core Solo Mobile",
	0x2B: "Atom",
	0x2C: "Core M",
	0x2D: "Core m3",
	0x2E: "Core m5",
	0x2F: "Core m7",
	0x30: "Alpha",
	0x31: "Alpha 21064",
	0x32: "Alpha 21066",
	0x33: "Alpha 21164",
	0x34: "Alpha 21164PC",
	0x35: "Alpha 21164a",
	0x36: "Alpha 21264",
	0x37: "Alpha 21364",
	0x38: "Turion II Ultra Dual-Core Mobile M",
	0x39: "Turion II Dual-Core Mobile M",
	0x3A: "Athlon II Dual-Core M",
	0x3B: "Opteron 6100",
	0x3C: "Opteron 4100",
	0x3D: "Opteron 6200",
	0x3E: "Opteron 4200",
	0x3F: "FX",
	0x40: "MIPS",
	0x41: "MIPS R4000",
	0x42: "MIPS R4200",
	0x43: "MIPS R4400",
	0x44: "MIPS R4600",
	0x45: "MIPS R10000",
	0x46: "C-Series",
	0x47: "E-Series",
	0x48: "A-Series",
	0x49: "G-Series",
	0x4A: "Z-Series",
	0x4B: "R-Series",
	0x4C: "Opteron 4300",
	0x4D: "Opteron 6300",
	0x4E: "Opteron 3300",
	0x4F: "FirePro",
	0x50: "SPARC",
	0x51: "SuperSPARC",
	0x52: "MicroSPARC II",
	0x53: "MicroSPARC IIep",
	0x54: "UltraSPARC",
	0x55: "UltraSPARC II",
	0x56: "UltraSPARC IIi",
	0x57: "UltraSPARC III",
	0x58: "UltraSPARC IIIi",
	0x60: "68040",
	0x61: "68xxx",
	0x62: "68000",
	0x63: "68010",
	0x64: "68020",
	0x65: "68030",
	0x66: "Athlon X4",
	0x67: "Opteron X1000",
	0x68: "Opteron X2000",
	0x69: "Opteron A-Series",
	0x6A: "Opteron X3000",
	0x6B: "Zen",
	0x70: "Hobbit",
	0x78: "Crusoe TM5000",
	0x79: "Crusoe TM3000",
	0x7A: "Efficeon TM8000",
	0x80: "Weitek",
	0x82: "Itanium",
	0x83: "Athlon 64",
	0x84: "Opteron",
	0x85: "Sempron",
	0x86: "Turion 64",
	0x87: "Dual-Core Opteron",
	0x88: "Athlon 64 X2",
	0x89: "Turion 64 X2",
	0x8A: "Quad-Core Opteron",
	0x8B: "Third-Generation Opteron",
	0x8C: "Phenom FX",
	0x8D: "Phenom X4",
	0x8E: "Phenom X2",
	0x8F: "Athlon X2",
	0x90: "PA-RISC",
	0x91: "PA-RISC 8500",
	0x92: "PA-RISC 8000",
	0x93: "PA-RISC 7300LC",
	0x94: "PA-RISC 7200",
	0x95: "PA-RISC 7100LC",
	0x96: "PA-RISC 7100",
	0xA0: "V30",
	0xA1: "Quad-Core Xeon 3200",
	0xA2: "Dual-Core Xeon 3000",
	0xA3: "Quad-Core Xeon 5300",
	0xA4: "Dual-Core Xeon 5100",
	0xA5: "Dual-Core Xeon 5000",
	0xA6: "Dual-Core Xeon LV",
	0xA7: "Dual-Core Xeon ULV",
	0xA8: "Dual-Core Xeon 7100",
	0xA9: "Quad-Core Xeon 5400",
	0xAA: "Quad-Core Xeon",
	0xAB: "Dual-Core Xeon 5200",
	0xAC: "Dual-Core Xeon 7200",
	0xAD: "Quad-Core Xeon 7300",
	0xAE: "Quad-Core Xeon 7400",
	0xAF: "Multi-Core Xeon 7400",
	0xB0: "Pentium III Xeon",
	0xB1: "Pentium III Speedstep",
	0xB2: "Pentium 4",
	0xB3: "Xeon",
	0xB4: "AS400",
	0xB5: "Xeon MP",
	0xB6: "Athlon XP",
	0xB7: "Athlon MP",
	0xB8: "Itanium 2",
	0xB9: "Pentium M",
	0xBA: "Celeron D",
	0xBB: "Pentium D",
	0xBC: "Pentium EE",
	0xBD: "Core Solo",
	/* 0x:E handled as a special case */
	0xBF:  "Core 2 Duo",
	0xC0:  "Core 2 Solo",
	0xC1:  "Core 2 Extreme",
	0xC2:  "Core 2 Quad",
	0xC3:  "Core 2 Extreme Mobile",
	0xC4:  "Core 2 Duo Mobile",
	0xC5:  "Core 2 Solo Mobile",
	0xC6:  "Core i7",
	0xC7:  "Dual-Core Celeron",
	0xC8:  "IBM390",
	0xC9:  "G4",
	0xCA:  "G5",
	0xCB:  "ESA/390 G6",
	0xCC:  "z/Architecture",
	0xCD:  "Core i5",
	0xCE:  "Core i3",
	0xCF:  "Core i9",
	0xD2:  "C7-M",
	0xD3:  "C7-D",
	0xD4:  "C7",
	0xD5:  "Eden",
	0xD6:  "Multi-Core Xeon",
	0xD7:  "Dual-Core Xeon 3xxx",
	0xD8:  "Quad-Core Xeon 3xxx",
	0xD9:  "Nano",
	0xDA:  "Dual-Core Xeon 5xxx",
	0xDB:  "Quad-Core Xeon 5xxx",
	0xDD:  "Dual-Core Xeon 7xxx",
	0xDE:  "Quad-Core Xeon 7xxx",
	0xDF:  "Multi-Core Xeon 7xxx",
	0xE0:  "Multi-Core Xeon 3400",
	0xE4:  "Opteron 3000",
	0xE5:  "Sempron II",
	0xE6:  "Embedded Opteron Quad-Core",
	0xE7:  "Phenom Triple-Core",
	0xE8:  "Turion Ultra Dual-Core Mobile",
	0xE9:  "Turion Dual-Core Mobile",
	0xEA:  "Athlon Dual-Core",
	0xEB:  "Sempron SI",
	0xEC:  "Phenom II",
	0xED:  "Athlon II",
	0xEE:  "Six-Core Opteron",
	0xEF:  "Sempron M",
	0xFA:  "i860",
	0xFB:  "i960",
	0x100: "ARMv7",
	0x101: "ARMv8",
	0x104: "SH-3",
	0x105: "SH-4",
	0x118: "ARM",
	0x119: "StrongARM",
	0x12C: "6x86",
	0x12D: "MediaGX",
	0x12E: "MII",
	0x140: "WinChip",
	0x15E: "DSP",
	0x1F4: "Video Processor",
}

var Processor_Voltage = []string{ /* 7.5.4 */
	"5.0 V", /* 0 */
	"3.3 V",
	"2.9 V", /* 2 */
}

var Processor_Status = []string{
	"Unknown", /* 0x00 */
	"Enabled",
	"Disabled By User",
	"Disabled By BIOS",
	"Idle", /* 0x04 */
	"Unknown",
	"Unknown",
	"Other", /* 0x07 */
}

var Processor_Upgrade = []string{ /* 7.5.5 */
	"Other", /* 0x01 */
	"Unknown",
	"Daughter Board",
	"ZIF Socket",
	"Replaceable Piggy Back",
	"None",
	"LIF Socket",
	"Slot 1",
	"Slot 2",
	"370-pin Socket",
	"Slot A",
	"Slot M",
	"Socket 423",
	"Socket A (Socket 462)",
	"Socket 478",
	"Socket 754",
	"Socket 940",
	"Socket 939",
	"Socket mPGA604",
	"Socket LGA771",
	"Socket LGA775",
	"Socket S1",
	"Socket AM2",
	"Socket F (1207)",
	"Socket LGA1366",
	"Socket G34",
	"Socket AM3",
	"Socket C32",
	"Socket LGA1156",
	"Socket LGA1567",
	"Socket PGA988A",
	"Socket BGA1288",
	"Socket rPGA988B",
	"Socket BGA1023",
	"Socket BGA1224",
	"Socket BGA1155",
	"Socket LGA1356",
	"Socket LGA2011",
	"Socket FS1",
	"Socket FS2",
	"Socket FM1",
	"Socket FM2",
	"Socket LGA2011-3",
	"Socket LGA1356-3",
	"Socket LGA1150",
	"Socket BGA1168",
	"Socket BGA1234",
	"Socket BGA1364",
	"Socket AM4",
	"Socket LGA1151",
	"Socket BGA1356",
	"Socket BGA1440",
	"Socket BGA1515",
	"Socket LGA3647-1",
	"Socket SP3",
	"Socket SP3r2",
	"Socket LGA2066",
	"Socket BGA1392",
	"Socket BGA1510",
	"Socket BGA1528", /* 0x3C */
}

var Processor_Characteristics = []string{ /* 7.5.9 */
	"Reserved",       /* bit 0 */
	"Unknown",        /* bit 1 */
	"64-bit capable", /* bit 2 */
	"Multi-Core",
	"Hardware Thread",
	"Execute Protection",
	"Enhanced Virtualization",
	"Power/Performance Control", /* bit 7 */
}
