package process

import (
	"testing"
)

func TestBackgroundProcessNew(t *testing.T) {
	process := BackgroundProcessNew(nil, nil, nil, nil, nil)
	if process == nil {
		t.Fatal("create background process returned nil")
	}
}
