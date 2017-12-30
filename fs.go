package fs

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

func ReadDirEntriesEndingWith(path string, ending string) []string {
	log.Println("Opening files inside " + path + " ending with " + ending)
	fileInfos := getDirContentInfos(path)
	names := []string{}
	for _, f := range fileInfos {
		if strings.HasSuffix(f.Name(), ending) {
			names = append(names, f.Name())
		}
	}
	sort.Strings(names)
	return names
}

func createPath(absPath string) {
	exists, err := pathExists(absPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if !exists {
		log.Println("creating path", absPath)
		os.MkdirAll(absPath, 0755)
	}
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func WriteStringToFS(path, filename, content string) {
	log.Println("Writing to " + path + filename)
	pathExists, _ := pathExists(path)
	if !pathExists {
		createPath(path)
	}
	if strings.LastIndex(path, "/") == -1 {
		path += "/"
	}
	err := ioutil.WriteFile(path+filename, []byte(content), 0644)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func ReadFileAsString(path string) string {
	return string(ReadByteArrayFromFile(path))
}

func ReadByteArrayFromFile(path string) []byte {
	log.Println("Reading file " + path)
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err.Error(), path)
	}
	return raw
}

func ReadDirEntries(path string, beingDir bool) []string {
	log.Println("Reading dir entries " + path)
	fileInfos := getDirContentInfos(path)
	names := []string{}
	for _, f := range fileInfos {
		if beingDir == f.IsDir() {
			names = append(names, f.Name())
		}
	}
	return names
}

func getDirContentInfos(path string) []os.FileInfo {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return fileInfos
}

func GetBase64FromPngFile(path string) (string, int, int) {
	imgFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer imgFile.Close()

	fInfo, _ := imgFile.Stat()
	var size int64 = fInfo.Size()
	buf := make([]byte, size)

	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)
	b := base64.StdEncoding.EncodeToString(buf)

	w, h := GetImageDimensions(path)

	return b, w, h
}

func GetImageDimensions(path string) (int, int) {
	img := GetImageConfig(path)
	return img.Width, img.Height
}

func GetImageConfig(path string) image.Config {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Println(err)
	}
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Println(err)
	}
	return img
}
