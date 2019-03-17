package fs

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func tearDown() {
	dir := path.Join(getTestFileDirPath(), "testResources/db")
	RemoveDirContents(dir)
}

func getTestFileDirPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func TestPathExists(t *testing.T) {
	p := path.Join(getTestFileDirPath(), "testResources/not-here")

	expected := false
	actual, _ := PathExists(p)

	if actual != expected {
		t.Error("Expected pathExists to return", expected, "but it returned", actual)
	}

	p = path.Join(getTestFileDirPath(), "testResources")
	files, err := ioutil.ReadDir(p)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Contents of " + p)
	for _, f := range files {
		fmt.Println(f.Name())
	}

	p = path.Join(getTestFileDirPath(), "testResources/db")
	expected = true
	actual, _ = PathExists(p)

	if actual != expected {
		t.Error("Expected pathExists to return", expected, "for path", p, "but it returned", actual)
	}
}

func TestGetImageDimensions(t *testing.T) {
	pngPath := "testResources/add/ClassicPen.png"
	w, h := GetImageDimensions(pngPath)
	if w != 800 || h != 1498 {
		t.Errorf("expected %d to be 800 and %d to be 1498", w, h)
	}
}

func TestGetPathWithoutFilename(t *testing.T) {
	path := "testResources/add/ClassicPen.png"
	actual := GetPathWithoutFilename(path)
	expected := "testResources/add/"
	if actual != expected {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestGetFilenameFromPath(t *testing.T) {
	path := "testResources/add/ClassicPen.png"
	actual := GetFilenameFromPath(path)
	expected := "ClassicPen.png"
	if actual != expected {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestCreateDir(t *testing.T) {
	p := path.Join(getTestFileDirPath(),
		"testResources/db/abouttobecreated")
	dirExists, _ := PathExists(p)

	if dirExists {
		t.Error("Test expects this path to be non existent: ", p)
	}

	err := CreateDir(p)
	if err != nil {
		t.Error(err)
	}

	err = CreateDir(p)
	if err == nil {
		t.Error("CreateDir should return an error when given a path that already exits")
	}
}
