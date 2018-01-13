package fs

import "testing"

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
