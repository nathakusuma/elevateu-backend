package fileutil

import "fmt"

var (
	KiloByte int64 = 1024
	MegaByte       = 1024 * KiloByte
	GigaByte       = 1024 * MegaByte
)

func ByteToAppropriateUnit(byte int64) string {
	if byte >= GigaByte {
		return fmt.Sprintf("%.1fGB", float64(byte)/float64(GigaByte))
	} else if byte >= MegaByte {
		return fmt.Sprintf("%.1fMB", float64(byte)/float64(MegaByte))
	} else if byte >= KiloByte {
		return fmt.Sprintf("%.1fKB", float64(byte)/float64(KiloByte))
	}
	return fmt.Sprintf("%dB", byte)
}
