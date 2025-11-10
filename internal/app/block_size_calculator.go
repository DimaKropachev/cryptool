package app

import (
	mem "github.com/DimaKropachev/cryptool/pkg/memory"
)

func CalculateOptimalBlockSize(fileSize int) (int, error) {
	freeRAM, err := mem.GetFreeRAM()
	if err != nil {
		return 0, err
	}

	if fileSize < int(freeRAM)/2 {
		return fileSize, nil
	} else {
		quarterRAM := freeRAM / 4
		if quarterRAM > 2*1024*1024*1024 {
			return 2 * 1024 * 1024 * 1024, nil
		}
		return int(quarterRAM), nil
	}
}
