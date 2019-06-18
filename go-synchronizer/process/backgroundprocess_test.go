package process

import (
	"testing"
)

func TestBackgroundProcessNew(t *testing.T) {
	process := BackgroundProcessNew(nil)
	if process == nil {
		t.Fatal("create background process returned nil")
	}
}
