package x264

// H.264 NAL unit types.
const (
	NALTypeNonIDR = 1 // Coded slice of a non-IDR picture
	NALTypeIDR    = 5 // Coded slice of an IDR picture
	NALTypeSEI    = 6 // Supplemental Enhancement Information
	NALTypeSPS    = 7 // Sequence Parameter Set
	NALTypePPS    = 8 // Picture Parameter Set
	NALTypeAUD    = 9 // Access Unit Delimiter
)

// readNALStartCode checks whether data at position pos begins with a H.264
// Annex B start code prefix. H.264 allows two start code forms:
//
//	4-byte: 0x00 0x00 0x00 0x01
//	3-byte: 0x00 0x00 0x01
//
// as specified in ITU-T H.264 (2021) Section 7.4.1 ("NAL unit semantics").
// The NAL unit type is encoded in the first byte after the start code:
//
//	nal_unit_type = first_byte & 0x1f
//
// It returns the NAL unit type, the start code prefix length (3 or 4), and
// true if a valid start code is found at position pos.
func readNALStartCode(data []byte, pos int) (nalType uint8, prefixLen int, ok bool) {
	if pos+5 <= len(data) && data[pos] == 0 && data[pos+1] == 0 && data[pos+2] == 0 && data[pos+3] == 1 {
		return data[pos+4] & 0x1f, 4, true
	}
	if pos+4 <= len(data) && data[pos] == 0 && data[pos+1] == 0 && data[pos+2] == 1 {
		return data[pos+3] & 0x1f, 3, true
	}
	return 0, 0, false
}

// NALTypes extracts all H.264 NAL unit types from encoded data.
func NALTypes(data []byte) []uint8 {
	var types []uint8
	pos := 0
	for pos < len(data) {
		if t, n, ok := readNALStartCode(data, pos); ok {
			types = append(types, t)
			pos += n
		} else {
			pos++
		}
	}
	return types
}

// HasNALType reports whether the encoded data contains a NAL unit of the given type.
func HasNALType(data []byte, nalType uint8) bool {
	pos := 0
	for pos < len(data) {
		if t, n, ok := readNALStartCode(data, pos); ok {
			if t == nalType {
				return true
			}
			pos += n
		} else {
			pos++
		}
	}
	return false
}
