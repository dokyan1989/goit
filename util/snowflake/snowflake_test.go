package snowflake

import (
	"testing"
)

func Test_createNodeID(t *testing.T) {
	got := createNodeID()

	if got < 0 {
		t.Errorf("createNodeID() = %v, want > 0", got)
		return
	}

	if got > 1023 {
		t.Errorf("createNodeID() = %v, want < 1023", got)
	}
}

func Test_NextID(t *testing.T) {
	got, err := NextID()

	if err != nil {
		t.Errorf("NextID() has error but want nil")
		return
	}

	if got < 0 {
		t.Errorf("NextID() = %v, want > 0", got)
	}
}
