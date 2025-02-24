package fileutil

import "fmt"

var (
	KiloByte int64 = 1024
	MegaByte       = 1024 * KiloByte
	GigaByte       = 1024 * MegaByte
)

func ByteToAppropriateUnit(byteUnit int64) string {
	switch {
	case byteUnit >= GigaByte:
		return fmt.Sprintf("%.1fGB", float64(byteUnit)/float64(GigaByte))
	case byteUnit >= MegaByte:
		return fmt.Sprintf("%.1fMB", float64(byteUnit)/float64(MegaByte))
	case byteUnit >= KiloByte:
		return fmt.Sprintf("%.1fKB", float64(byteUnit)/float64(KiloByte))
	default:
		return fmt.Sprintf("%dB", byteUnit)
	}
}
