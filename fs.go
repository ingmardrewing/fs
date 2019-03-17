package fs

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
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
	f.dirpath = filepath.FromSlash(dirpath)
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
	exists, err := PathExists(absPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if !exists {
		os.MkdirAll(absPath, 0755)
	}
}

func CreateDir(absPath string) error {
	exists, err := PathExists(absPath)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("Path already exists")
	}
	os.MkdirAll(absPath, 0755)
	return nil
}

func RemoveDir(absPath string) error {
	exists, err := PathExists(absPath)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Dir doesn't exist at" + absPath)
	}
	fi, err := os.Stat(absPath)
	if err != nil {
		return err
	}
	mode := fi.Mode()
	if !mode.IsDir() {
		return errors.New("Is not a directory:" + absPath)
	}
	os.Remove(absPath)
	return nil
}

func PathExists(pth string) (bool, error) {
	pth = filepath.FromSlash(pth)
	_, err := os.Stat(pth)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func WriteFromFileContainer(f FileContainer) {
	WriteStringToFS(f.GetPath(), f.GetFilename(), f.GetDataAsString())
}

func WriteStringToFS(fpath, filename, content string) {
	pth := filepath.FromSlash(fpath)
	pathExists, _ := PathExists(pth)
	if !pathExists {
		createPath(pth)
	}
	completepath := path.Join(pth, filename)
	err := ioutil.WriteFile(completepath, []byte(content), 0644)
	if err != nil {
		fmt.Println("fs.WriteStringToFs", err.Error())
		os.Exit(1)
	}
}

func ReadFromFileContainer(f FileContainer) {
	path := f.GetPath()
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	data := ReadFileAsString(path + f.GetFilename())
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
	if len(raw) == 0 {
		log.Println("Empty file:", path)
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

func RemoveDirContents(basedir string) error {
	d, err := os.Open(basedir)
	if err != nil {
		return err
	}
	defer d.Close()
	subdirnames, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, subdirname := range subdirnames {
		err = os.RemoveAll(filepath.Join(basedir, subdirname))
		if err != nil {
			return err
		}
	}
	return nil
}

func RemoveFile(p, filename string) {
	err := os.Remove(path.Join(p, filename))
	if err != nil {
		log.Fatalln(err)
	}
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

func Pwd() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}

	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
