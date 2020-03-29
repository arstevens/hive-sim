package controller

import (
	"fmt"
	"testing"
)

func TestBasicNode(t *testing.T) {
	bn := NewRandomBasicNode(nil, nil)
	id := bn.Id()
	if id != "" {
		fmt.Println(id)
		t.Log(id)
		return
	}
	t.Errorf("No id generated")
}
