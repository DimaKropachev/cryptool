package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/DimaKropachev/cryptool/pkg/models"
)

type DirectoryScanner struct {
	BasePath string
	Files    []*models.File
}

func newDirScanner(basePath string) *DirectoryScanner {
	return &DirectoryScanner{
		BasePath: basePath + string(filepath.Separator),
		Files:    []*models.File{},
	}
}

func ReadDirectory(dirPath string) ([]*models.File, error) {
	ds := newDirScanner(dirPath)

	err := ds.readDirectory(dirPath)
	if err != nil {
		return nil, pathError(ActionScanning, dirPath, err)
	}
	return ds.Files, nil
}

func (ds *DirectoryScanner) readDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed get information from <%s>", path)
	}

	if info.IsDir() {
		d, err := os.ReadDir(path)
		if err != nil {
			return fmt.Errorf("failed opening directory <%s>", path)
		}

		for _, obj := range d {
			currPath := path + string(filepath.Separator) + obj.Name()
			currName := strings.TrimPrefix(currPath, ds.BasePath)

			if obj.IsDir() {
				err = ds.readDirectory(currPath)
				if err != nil {
					return err
				}
			} else {
				info, err := obj.Info()
				if err != nil {
					return fmt.Errorf("failed get information from <%s>", currPath)
				}

				ds.Files = append(ds.Files, &models.File{
					Name: currName,
					Info: info,
					Path: currPath,
				})
			}
		}
	} else {
		ds.Files = append(ds.Files, &models.File{
			Name: strings.TrimPrefix(path, ds.BasePath),
			Info: info,
		})
	}

	return nil
}
