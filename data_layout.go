package irutil

import (
	"fmt"
	"strconv"
	"strings"
)

// DataLayout is a structural representation of the datalayout
// string from https://llvm.org/docs/LangRef.html#data-layout
type DataLayout struct {
	isBigEndian               bool                                   // default - True
	naturalStackAlignment     uint64                                 // bits in multiple of 8, default: 0
	programMemoryAddressSpace uint64                                 // default - 0
	globalVarAddressSpace     uint64                                 // default - 0
	allocaAddressSpace        uint64                                 // default - 0
	ptrSizeAlign              map[uint64]*PointerSizeAlignment       // Key is addressspace
	intSizeAlign              map[uint64]*IntegerSizeAlignment       // Key is size
	vecSizeAlign              map[uint64]*VectorSizeAlignment        // Key is size
	floatSizeAlign            map[uint64]*FloatingPointSizeAlignment // Key is size
	aggAlign                  *AggregateAlignment                    // default - 0:64
	funcPtrAlign              *FunctionPointerAlignment              // default - Not specified in LLVM Ref doc
	mangling                  ManglingStyle                          // depends on os - "e" for linux, "o" for mac
	nativeIntBitWidths        []uint64                               // default - depends on arch "8:16:32:64" for x86-64
	nonIntegralPointerTypes   []string                               // Unsure of the datatype, so leaving it as string to accommodate anything.
}

// NewDataLayout returns an instance of DataLayout with default values
// taken from https://llvm.org/docs/LangRef.html#data-layout
func NewDataLayout(os string, arch string) *DataLayout {
	dl := &DataLayout{
		isBigEndian:               true,
		naturalStackAlignment:     0,
		programMemoryAddressSpace: 0,
		globalVarAddressSpace:     0,
		allocaAddressSpace:        0,
		ptrSizeAlign:              getDefaultPtrAligns(),
		intSizeAlign:              getDefaultIntAligns(),
		vecSizeAlign:              getDefaultVectorAligns(),
		floatSizeAlign:            getDefaultFloatAligns(),
		aggAlign:                  NewAggregateAlignment(0, 64),
	}
	switch arch {
	case "x86-64":
		dl.nativeIntBitWidths = []uint64{8, 16, 32, 64}
	}

	if os == "linux" {
		dl.mangling = ELF
	} else if strings.HasPrefix(os, "darwin") || strings.HasPrefix(os, "macosx") {
		dl.mangling = MachO
	}
	return dl
}

// NewDataLayoutFromString parses a datalayout string and constructs a DataLayout object.
// Default values whereever applicable will be added if options are ommitted.
//
// This API expects the string to be valid as per LLVM spec, so there aren't many validations.
func NewDataLayoutFromString(layoutString, os, arch string) *DataLayout {
	dl := NewDataLayout(os, arch)
	specs := strings.Split(layoutString, "-")
	for _, spec := range specs {
		if spec == "e" {
			dl.isBigEndian = false
		} else if spec == "E" {
			dl.isBigEndian = true
		} else if strings.HasPrefix(spec, "S") {
			dl.naturalStackAlignment = getUInt64(strings.TrimPrefix(spec, "S"))
		} else if strings.HasPrefix(spec, "P") {
			dl.programMemoryAddressSpace = getUInt64(strings.TrimPrefix(spec, "P"))
		} else if strings.HasPrefix(spec, "G") {
			dl.globalVarAddressSpace = getUInt64(strings.TrimPrefix(spec, "G"))
		} else if strings.HasPrefix(spec, "A") {
			dl.allocaAddressSpace = getUInt64(strings.TrimPrefix(spec, "A"))
		} else if strings.HasPrefix(spec, "p") {
			addPtrSizeAlignFromString(spec, dl)
		} else if strings.HasPrefix(spec, "i") {
			addIntSizeAlignFromString(spec, dl)
		} else if strings.HasPrefix(spec, "v") {
			addVectorSizeAlignFromString(spec, dl)
		} else if strings.HasPrefix(spec, "f") {
			addFloatSizeAlignFromString(spec, dl)
		} else if strings.HasPrefix(spec, "a") {
			addAggAlignFromString(spec, dl)
		} else if strings.HasPrefix(spec, "F") {
			addFuncPtrAlignFromString(spec, dl)
		} else if strings.HasPrefix(spec, "m") {
			dl.mangling = ManglingStyle(strings.Split(spec, ":")[1])
		} else if strings.HasPrefix(spec, "n") {
			addNativeIntBitWidths(spec, dl)
		} else if strings.HasPrefix(spec, "ni") {
			dl.nonIntegralPointerTypes = strings.Split(strings.TrimPrefix(spec, "ni"), ":")
		}
	}
	return dl
}

func addPtrSizeAlignFromString(spec string, dl *DataLayout) {
	ptrVals := strings.Split(spec, ":")
	ptrValsLen := len(ptrVals)
	ptrAddSpace := uint64(0)
	if len(ptrVals[0]) > 1 {
		ptrAddSpace = getUInt64(strings.TrimPrefix(ptrVals[0], "p"))
	}
	size := getUInt64(ptrVals[1])
	abi := getUInt64(ptrVals[2])
	pref := abi
	if ptrValsLen > 3 {
		pref = getUInt64(ptrVals[3])
	}
	ind := size
	if ptrValsLen == 5 {
		ind = getUInt64(ptrVals[4])
	}
	dl.ptrSizeAlign[ptrAddSpace] = NewPointerSizeAlignment(ptrAddSpace, size, abi, pref, ind)
}

func addIntSizeAlignFromString(spec string, dl *DataLayout) {
	intVals := strings.Split(spec, ":")
	intValsLen := len(intVals)
	size := getUInt64(strings.TrimPrefix(intVals[0], "i"))
	abi := getUInt64(intVals[1])
	pref := abi
	if intValsLen == 3 {
		pref = getUInt64(intVals[2])
	}
	dl.intSizeAlign[size] = NewIntegerSizeAlignment(size, abi, pref)
}

func addVectorSizeAlignFromString(spec string, dl *DataLayout) {
	vectorVals := strings.Split(spec, ":")
	vectorValsLen := len(vectorVals)
	size := getUInt64(strings.TrimPrefix(vectorVals[0], "v"))
	abi := getUInt64(vectorVals[1])
	pref := abi
	if vectorValsLen == 3 {
		pref = getUInt64(vectorVals[2])
	}
	dl.vecSizeAlign[size] = NewVectorSizeAlignment(size, abi, pref)
}

func addFloatSizeAlignFromString(spec string, dl *DataLayout) {
	floatVals := strings.Split(spec, ":")
	floatValsLen := len(floatVals)
	size := getUInt64(strings.TrimPrefix(floatVals[0], "f"))
	abi := getUInt64(floatVals[1])
	pref := abi
	if floatValsLen == 3 {
		pref = getUInt64(floatVals[2])
	}
	dl.floatSizeAlign[size] = NewFloatingPointSizeAlignment(size, abi, pref)
}

func addAggAlignFromString(spec string, dl *DataLayout) {
	aggVals := strings.Split(spec, ":")
	aggValsLen := len(aggVals)
	abi := getUInt64(aggVals[1])
	pref := abi
	if aggValsLen == 3 {
		pref = getUInt64(aggVals[2])
	}
	dl.aggAlign = NewAggregateAlignment(abi, pref)
}

func addFuncPtrAlignFromString(spec string, dl *DataLayout) {
	val := strings.TrimPrefix(spec, "F")
	var typ string
	if strings.HasPrefix(val, "i") {
		typ = "i"
	} else if strings.HasPrefix(val, "n") {
		typ = "n"
	} else {
		panic("Invalid Function pointer alignment type, possible options are \"i\", \"n\".")
	}
	abi := getUInt64(strings.TrimPrefix(val, typ))
	dl.funcPtrAlign = NewFunctionPointerAlignment(typ == "i", abi)
}

func addNativeIntBitWidths(spec string, dl *DataLayout) {
	widthContents := strings.Split(strings.TrimPrefix(spec, "n"), ":")
	widths := make([]uint64, len(widthContents))
	for i, width := range widthContents {
		widths[i] = getUInt64(width)
	}
	dl.nativeIntBitWidths = widths
}

func getUInt64(str string) uint64 {
	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		panic(err.Error())
	}
	return val
}

func (dl *DataLayout) SetIsBigEndian(isBigEndian bool) {
	dl.isBigEndian = isBigEndian
}

func (dl *DataLayout) SetNaturalStackAlignment(al uint64) {
	dl.naturalStackAlignment = al
}

func (dl *DataLayout) SetProgramMemoryAddressSpace(as uint64) {
	dl.programMemoryAddressSpace = as
}

func (dl *DataLayout) SetGlobalVarAddressSpace(as uint64) {
	dl.globalVarAddressSpace = as
}

func (dl *DataLayout) SetAllocaAddressSpace(as uint64) {
	dl.allocaAddressSpace = as
}

func (dl *DataLayout) AddPointerSizeAlignment(instance *PointerSizeAlignment) {
	dl.ptrSizeAlign[instance.addressSpace] = instance
}

func (dl *DataLayout) AddIntegerSizeAlignment(instance *IntegerSizeAlignment) {
	dl.intSizeAlign[instance.size] = instance
}

func (dl *DataLayout) AddVectorSizeAlignment(instance *VectorSizeAlignment) {
	dl.vecSizeAlign[instance.size] = instance
}

func (dl *DataLayout) AddFloatingPointSizeAlignment(instance *FloatingPointSizeAlignment) {
	dl.floatSizeAlign[instance.size] = instance
}

func (dl *DataLayout) SetAggregateAlignment(instance *AggregateAlignment) {
	dl.aggAlign = instance
}

func (dl *DataLayout) SetFunctionPointerAlignment(instance *FunctionPointerAlignment) {
	dl.funcPtrAlign = instance
}

func (dl *DataLayout) SetManglingStyle(style ManglingStyle) {
	dl.mangling = style
}

func (dl *DataLayout) SetNativeIntBitWidths(bitWidths []uint64) {
	dl.nativeIntBitWidths = bitWidths
}

func (dl *DataLayout) GetLayoutString() string {
	layout := ""
	if dl.isBigEndian {
		layout += "E"
	} else {
		layout += "e"
	}
	if dl.naturalStackAlignment != 0 {
		layout += "-S" + fmt.Sprint(dl.naturalStackAlignment)
	}
	if dl.programMemoryAddressSpace != 0 {
		layout += "-P" + fmt.Sprint(dl.programMemoryAddressSpace)
	}
	if dl.globalVarAddressSpace != 0 {
		layout += "-G" + fmt.Sprint(dl.globalVarAddressSpace)
	}
	if dl.allocaAddressSpace != 0 {
		layout += "-A" + fmt.Sprint(dl.allocaAddressSpace)
	}
	if len(dl.ptrSizeAlign) > 0 {
		for k, v := range dl.ptrSizeAlign {
			if k == 0 {
				layout += "-p"
			} else {
				layout += "-p" + fmt.Sprint(k)
			}
			layout += ":" + fmt.Sprint(v.size) + ":" + fmt.Sprint(v.abiAlignment) + ":" +
				fmt.Sprint(v.preferredAlignment) + ":" + fmt.Sprint(v.addressCalculationIndex)
		}
	}
	if len(dl.intSizeAlign) > 0 {
		for k, v := range dl.intSizeAlign {
			layout += "-i" + fmt.Sprint(k) + ":" + fmt.Sprint(v.abiAlignment) + ":" + fmt.Sprint(v.preferredAlignment)
		}
	}
	if len(dl.vecSizeAlign) > 0 {
		for k, v := range dl.vecSizeAlign {
			layout += "-v" + fmt.Sprint(k) + ":" + fmt.Sprint(v.abiAlignment) + ":" + fmt.Sprint(v.preferredAlignment)
		}
	}
	if len(dl.floatSizeAlign) > 0 {
		for k, v := range dl.floatSizeAlign {
			layout += "-f" + fmt.Sprint(k) + ":" + fmt.Sprint(v.abiAlignment) + ":" + fmt.Sprint(v.preferredAlignment)
		}
	}
	if dl.aggAlign != nil {
		layout += "-a:" + fmt.Sprint(dl.aggAlign.abiAlignment) + ":" + fmt.Sprint(dl.aggAlign.preferredAlignment)
	}
	if dl.funcPtrAlign != nil {
		if dl.funcPtrAlign.isIndependant {
			layout += "-Fi" + fmt.Sprint(dl.funcPtrAlign.abiAlignment)
		} else {
			layout += "-Fn" + fmt.Sprint(dl.funcPtrAlign.abiAlignment)
		}
	}
	layout += "-m:" + string(dl.mangling)
	if bitWidths := len(dl.nativeIntBitWidths); bitWidths > 0 {
		layout += "-n" + fmt.Sprint(dl.nativeIntBitWidths[0])
		for i := 1; i < bitWidths; i++ {
			layout += ":" + fmt.Sprint(dl.nativeIntBitWidths[i])
		}
	}
	if len(dl.nonIntegralPointerTypes) > 0 {
		layout += "-ni"
		for _, v := range dl.nonIntegralPointerTypes {
			layout += ":" + fmt.Sprint(v)
		}
	}
	return layout
}

type PointerSizeAlignment struct { // All sizes are in bits.
	addressSpace            uint64
	size                    uint64
	abiAlignment            uint64
	preferredAlignment      uint64 // optional and defaults to abi
	addressCalculationIndex uint64 // default 0
}

func NewPointerSizeAlignment(addSp, size, abiAl, prefAl, addCalInd uint64) *PointerSizeAlignment {
	return &PointerSizeAlignment{addressSpace: addSp, size: size, abiAlignment: abiAl, preferredAlignment: prefAl, addressCalculationIndex: addCalInd}
}

func getDefaultPtrAligns() map[uint64]*PointerSizeAlignment {
	ptrAlignMap := make(map[uint64]*PointerSizeAlignment)
	ptrAlignMap[0] = NewPointerSizeAlignment(0, 64, 64, 64, 64)
	return ptrAlignMap
}

type IntegerSizeAlignment struct { // All sizes are in bits.
	size               uint64
	abiAlignment       uint64
	preferredAlignment uint64 // optional and defaults to abi
}

func NewIntegerSizeAlignment(size, abiAl, prefAl uint64) *IntegerSizeAlignment {
	return &IntegerSizeAlignment{size: size, abiAlignment: abiAl, preferredAlignment: prefAl}
}

func getDefaultIntAligns() map[uint64]*IntegerSizeAlignment {
	intAlignMap := make(map[uint64]*IntegerSizeAlignment)
	intAlignMap[1] = NewIntegerSizeAlignment(1, 8, 8)
	intAlignMap[8] = NewIntegerSizeAlignment(8, 8, 8)
	intAlignMap[16] = NewIntegerSizeAlignment(16, 16, 16)
	intAlignMap[32] = NewIntegerSizeAlignment(32, 32, 32)
	intAlignMap[64] = NewIntegerSizeAlignment(64, 32, 64)
	return intAlignMap
}

type VectorSizeAlignment struct { // All sizes are in bits.
	size               uint64
	abiAlignment       uint64
	preferredAlignment uint64 // optional and defaults to abi
}

func NewVectorSizeAlignment(size, abiAl, prefAl uint64) *VectorSizeAlignment {
	return &VectorSizeAlignment{size: size, abiAlignment: abiAl, preferredAlignment: prefAl}
}

func getDefaultVectorAligns() map[uint64]*VectorSizeAlignment {
	vectorAlignMap := make(map[uint64]*VectorSizeAlignment)
	vectorAlignMap[64] = NewVectorSizeAlignment(64, 64, 64)
	vectorAlignMap[128] = NewVectorSizeAlignment(128, 128, 128)
	return vectorAlignMap
}

type FloatingPointSizeAlignment struct { // All sizes are in bits.
	size               uint64
	abiAlignment       uint64
	preferredAlignment uint64 // optional and defaults to abi
}

func NewFloatingPointSizeAlignment(size, abiAl, prefAl uint64) *FloatingPointSizeAlignment {
	return &FloatingPointSizeAlignment{size: size, abiAlignment: abiAl, preferredAlignment: prefAl}
}

func getDefaultFloatAligns() map[uint64]*FloatingPointSizeAlignment {
	floatAlignMap := make(map[uint64]*FloatingPointSizeAlignment)
	floatAlignMap[16] = NewFloatingPointSizeAlignment(16, 16, 16)
	floatAlignMap[32] = NewFloatingPointSizeAlignment(32, 32, 32)
	floatAlignMap[64] = NewFloatingPointSizeAlignment(64, 64, 64)
	floatAlignMap[128] = NewFloatingPointSizeAlignment(128, 128, 128)
	return floatAlignMap
}

type AggregateAlignment struct {
	abiAlignment       uint64
	preferredAlignment uint64 // optional and defaults to abi
}

func NewAggregateAlignment(abiAl, prefAl uint64) *AggregateAlignment {
	return &AggregateAlignment{abiAlignment: abiAl, preferredAlignment: prefAl}
}

type FunctionPointerAlignment struct {
	isIndependant bool
	abiAlignment  uint64
}

func NewFunctionPointerAlignment(isInd bool, abiAl uint64) *FunctionPointerAlignment {
	return &FunctionPointerAlignment{isIndependant: isInd, abiAlignment: abiAl}
}

type ManglingStyle string

const (
	ELF        ManglingStyle = "e"
	Mips       ManglingStyle = "m"
	MachO      ManglingStyle = "o"
	WinX86COFF ManglingStyle = "x"
	WinCOFF    ManglingStyle = "w"
	XCOFF      ManglingStyle = "a"
)
