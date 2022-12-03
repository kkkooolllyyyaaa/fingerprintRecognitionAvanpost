package app

import (
	"context"
	"fingerprintRecognitionAvanpost/internal/file"
	"fingerprintRecognitionAvanpost/internal/services"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"path/filepath"
)

func Run(ctx context.Context, erg *errgroup.Group) error {
	// Сюда нужно поместить 1 файл
	fileRoot := "files/train/SOCOFing/OneImage/"
	//fileRoot := "files/train/SOCOFing/TenPeople/"
	//fileRoot := "files/train/SOCOFing/Real/"
	//fileRoot := "files/train/SOCOFing/Altered/Altered-Hard/"

	filepath.Walk(fileRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}

		if info.IsDir() {
			return nil
		}
		fileWorker := file.New(fileRoot)
		bitset, err := fileWorker.ReadToBitset(info.Name())
		bitset.Print()
		services.Skeleton(bitset.Bin)
		if err != nil {
			log.Fatalf(err.Error())
		}
		bitset.Print()
		return nil
	})
	return nil
}
