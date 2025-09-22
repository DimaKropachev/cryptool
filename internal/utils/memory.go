package utils

import (
	"fmt"

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
