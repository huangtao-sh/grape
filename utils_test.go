package grape

import (
	"testing"
)

func TestExtractPos(t *testing.T) {
	if ExtractPos(`(\d{3})`, "fsdab134fsda", 1) != "134" {
		t.Error("test ExtractPos failed")
	}
}
