package file

import (
	"os"

	"github.com/DimaKropachev/cryptool/internal/utils"
)

type BlockCalculator struct {
	
}

func CalculateOptimalBlockSize(path string) (int, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	fileSize := int(fileInfo.Size())

	freeRAM, err := utils.GetFreeRAM()
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
