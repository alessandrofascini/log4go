package tests

import (
	"archive/tar"
	"os"
	"testing"
)

func TestTarGz(t *testing.T) {
	file, _ := os.Open("./readme.txt")
	defer file.Close()

	tw := tar.NewWriter(file)
	defer tw.Close()

	info, _ := file.Stat()

	hdr := tar.Header{
		Name: info.Name(),
		Size: info.Size(),
		Mode: 0600,
	}
	if err := tw.WriteHeader(&hdr); err != nil {
		panic(err)
	}
	//if _, err := tw.Write(); err != nil {
	//	panic(err)
	//}
}
