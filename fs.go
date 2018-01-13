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

type FileContainer interface {
	SetDataAsString(data string)
	GetDataAsString() string
	GetData() []byte
	SetData(data []byte)
	SetPath(dirpath string)
	GetPath() string
	SetFilename(filename string)
	GetFilename() string
	Write()
	Read()
}

func NewFileContainer() FileContainer {
	return new(FileContainerImpl)
}

type FileContainerImpl struct {
	data              []byte
	dirpath, filename string
}

func (f *FileContainerImpl) SetData(data []byte) {
	f.data = data
}

func (f *FileContainerImpl) GetData() []byte {
	return f.data
}

func (f *FileContainerImpl) SetDataAsString(data string) {
	f.data = []byte(data)
}

func (f *FileContainerImpl) GetDataAsString() string {
	return string(f.data)
}

func (f *FileContainerImpl) SetPath(dirpath string) {
	f.dirpath = dirpath
}

func (f *FileContainerImpl) GetPath() string {
	return f.dirpath
}

func (f *FileContainerImpl) SetFilename(filename string) {
	f.filename = filename
}

func (f *FileContainerImpl) GetFilename() string {
	return f.filename
}

func (f *FileContainerImpl) Write() {
	WriteFromFileContainer(f)
}

func (f *FileContainerImpl) Read() {
	ReadFromFileContainer(f)
}

/* util fns */

func ReadDirEntriesEndingWith(path string, ending ...string) []string {
	fileInfos := getDirContentInfos(path)
	names := []string{}
	for _, f := range fileInfos {
		for _, e := range ending {
			if strings.HasSuffix(f.Name(), e) {
				names = append(names, f.Name())
			}
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

func WriteFromFileContainer(f FileContainer) {
	path := f.GetPath()
	filename := f.GetFilename()
	log.Println("path", path)
	log.Println("filename", filename)
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	WriteStringToFS(path, filename, f.GetDataAsString())
}

func WriteStringToFS(path, filename, content string) {
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

func ReadFromFileContainer(f FileContainer) {
	data := ReadFileAsString(f.GetPath() + f.GetFilename())
	f.SetDataAsString(data)
}

func ReadFileAsString(path string) string {
	content := string(ReadByteArrayFromFile(path))
	return content
}

func ReadByteArrayFromFile(path string) []byte {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err.Error(), path)
	}
	return raw
}

func IsValidPathTo(path string, suffixes ...string) bool {
	if _, err := os.Stat(path); err == nil {
		for _, s := range suffixes {
			if strings.HasSuffix(path, s) {
				return true
			}
		}
	}
	fmt.Println("This path doesn't lead to any file with an ending like", strings.Join(suffixes, ", "))
	return false
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

func GetPathWithoutFilename(path string) string {
	if _, err := os.Stat(path); err == nil {
		parts := strings.Split(path, "/")
		newpath := strings.Join(parts[:len(parts)-1], "/")
		return newpath + "/"
	}
	log.Fatalln("Not a valid path", path)
	return ""
}

func GetFilenameFromPath(path string) string {
	if _, err := os.Stat(path); err == nil {
		parts := strings.Split(path, "/")
		return parts[len(parts)-1]
	}
	log.Fatalln("Not a valid path", path)
	return ""
}
