package fs

import "testing"

func TestGetImageDimensions(t *testing.T) {
	pngPath := "/Users/drewing/Desktop/devabo_de_uploads/colour_tutorial.png"
	w, h := GetImageDimensions(pngPath)
	if w != 800 || h != 1645 {
		t.Errorf("expected %d to be 800 and %d to be 1645", w, h)
	}
}
