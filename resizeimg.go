package main

import (
	"bytes"
	"image"
	"image/png"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/nfnt/resize"
)

// ResizeImgByDir(`ord_dtl`, 10000)

func ResizeImgByDir(dirpath string, limit int64) {
	ReadDir(dirpath, func(fileinfo fs.DirEntry) {
		imgpath := dirpath + "/" + fileinfo.Name()
		savepath := strings.Replace(imgpath, `.png`, `_x10.png`, 1)
		fino, err := fileinfo.Info()
		if err != nil {
			panic(err)
		}
		size := fino.Size()
		if size > limit {
			ResizeJpegByPath(imgpath, savepath, 150, 0)
			log.Printf("--SUCCESS-%s----size: %d------\n", imgpath, size)

		} else {

			ResizeJpegByPath(imgpath, savepath, 150, 0)
			log.Printf("--SKIP----%s----size: %d------\n", imgpath, size)
		}
	})
}

func resizeJpeg(data []byte, width, height uint) []byte {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return data
	}
	buf := bytes.Buffer{}
	m := resize.Resize(width, height, img, resize.Lanczos3)
	err = png.Encode(&buf, m) // &jpeg.Options{Quality: 40}
	if err != nil {
		return data
	}
	if buf.Len() > len(data) {
		return data
	}
	return buf.Bytes()
}

func ResizeJpegByPath(path, savepath string, width, height uint) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	compressBytes := resizeJpeg(data, width, height)
	os.WriteFile(savepath, compressBytes, 0666)
	f, _ := os.Create(savepath)
	f.Write(compressBytes)
}

func ReadDir(path string, callback func(fileinfo fs.DirEntry)) error {
	filelist, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
		return err
	}

	for _, fileinfo := range filelist {
		if fileinfo.Type().IsRegular() {
			callback(fileinfo)
		}
	}
	return nil
}
