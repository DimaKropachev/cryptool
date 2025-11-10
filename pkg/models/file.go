package models

import (
	"os"

	"github.com/DimaKropachev/cryptool/pkg/progressbar"
)

type File struct {
	Name string
	Info os.FileInfo
	Path string
	PB *progressbar.ProgressBar
}
