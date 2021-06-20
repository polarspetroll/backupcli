package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/polarspetroll/gocolor"
)

const BackupDir = ""

func main() {
	arg := os.Args
	if len(arg) < 2 {
		fmt.Println(gocolor.ColorString("Usage: \ngbackup [Directory]", "cyan", "bold"))
		return
	}
	tarname := ArchiveTar(os.Args[1])
	Compress(tarname)

}

func ArchiveTar(src string) (dst string) {
	fmt.Println(gocolor.ColorString("Generating tar...", "green", "italic"))
	y, m, d := time.Now().Date()
	h, min, sec := time.Now().Clock()
	dst = fmt.Sprintf("backup_%v-%v-%v_%v-%v-%v.tar", d, m, y, h, min, sec)
	dir, err := os.Open(src)
	CheckErr(err)
	defer dir.Close()
	files, err := dir.Readdir(0)
	CheckErr(err)
	t, err := os.Create("/tmp/" + dst)
	CheckErr(err)
	defer t.Close()
	var fw io.WriteCloser = t
	tarw := tar.NewWriter(fw)
	defer tarw.Close()
	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			continue
		}
		file, err := os.Open(dir.Name() + string(filepath.Separator) + fileInfo.Name())
		CheckErr(err)
		defer file.Close()
		header := new(tar.Header)
		header.Name = file.Name()
		header.Size = fileInfo.Size()
		header.Mode = int64(fileInfo.Mode())
		header.ModTime = fileInfo.ModTime()
		err = tarw.WriteHeader(header)
		CheckErr(err)
		_, err = io.Copy(tarw, file)
		CheckErr(err)
	}
	return dst
}

func Compress(fname string) {
	fmt.Println(gocolor.ColorString("Compressing...", "red", "bold"))
	f, err := os.Open("/tmp/" + fname)
	CheckErr(err)
	defer f.Close()
	r, err := ioutil.ReadAll(f)
	CheckErr(err)
	out, err := os.Create(BackupDir + "/" + fname + ".gz")
	CheckErr(err)
	defer out.Close()
	w := gzip.NewWriter(out)
	defer w.Close()
	w.Write(r)
	os.Remove("/tmp/" + fname)
	fmt.Println(gocolor.ColorString(fmt.Sprintf("Done!\nBackup saved at %v/%v", BackupDir, fname), "green", "bold"))
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
