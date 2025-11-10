package memory

import (
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

func GetFreeRAM() (uint64, error) {
	memStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, fmt.Errorf("error getting the amount of free RAM: %w", err)
	}

	return memStat.Available, nil
}

func GetFreeDiskMemory() (uint64, error) {
	diskMemStat, err := disk.Usage("C:\\")
	if err != nil {
		return 0, fmt.Errorf("error getting the amount of free disk memory: %w", err)
	}

	return diskMemStat.Free, nil
}

type MemoryUsed struct {
	Quantity float64
	Units    string
}

func MemoryMeasurement(oper func()) *MemoryUsed {
	runtime.GC()

	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)

	oper()

	runtime.ReadMemStats(&after)

	allocated := float64(after.TotalAlloc - before.TotalAlloc)

	const (
		kb = 1024
		mb = 1024 * kb
		gb = 1024 * mb
	)

	memUsed := &MemoryUsed{}

	switch {
	case allocated >= 0.5*gb:
		memUsed.Quantity = allocated / gb
		memUsed.Units = "GB"
	case allocated >= 0.5*mb:
		memUsed.Quantity = allocated / mb
		memUsed.Units = "MB"
	case allocated >= 0.5*kb:
		memUsed.Quantity = allocated / kb
		memUsed.Units = "KB"
	default:
		memUsed.Quantity = allocated
		memUsed.Units = "B"
	}

	return memUsed
}

func FormatBytes(bytes float64) string {
	const (
		kb = 1024
		mb = 1024 * kb
		gb = 1024 * mb
	)

	switch {
	case bytes >= 0.5*gb:
		return fmt.Sprintf("%.2f GB", bytes / gb)
	case bytes >= 0.5*mb:
		return fmt.Sprintf("%.2f MB", bytes / mb)
	case bytes >= 0.5*kb:
		return fmt.Sprintf("%.2f KB", bytes / kb)
	default:
  	return fmt.Sprintf("%.2f B", bytes)
	}	
}

func FormatTime(t time.Duration) string {
	switch {
	case t.Minutes() > 2:
		return fmt.Sprintf("%.2f m", t.Minutes())
	default:
		return fmt.Sprintf("%.2f s", t.Seconds())
	}
}

func CalculateThroughput(fileSize int64, duration time.Duration) float64 {
	fileSizeMB := float64(fileSize) / (1024 * 1024)
	timeSec := duration.Seconds()

	return fileSizeMB / timeSec
}
