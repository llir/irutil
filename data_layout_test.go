package irutil

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/llir/llvm/asm"
)

func TestDataLayoutString(t *testing.T) {
	opSys := "linux"
	arch := "x86-64"
	folderPath := "testdata/coreutils/test"
	os.Open("testdata")
	files, err := os.ReadDir(folderPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			t.Skip("testdata not found, clone it from https://github.com/llir/testdata")
		}
		t.Fatal(err.Error())
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".golden") {
			continue
		}
		filePath := folderPath + "/" + f.Name()
		t.Logf("=== [ %s ] ===", filePath)
		m, err := asm.ParseFile(filePath)
		if err != nil {
			t.Fatal(err.Error())
		}
		origDL, err := NewDataLayoutFromString(m.DataLayout, opSys, arch)
		if err != nil {
			t.Error(err.Error())
		}
		newDL, err := NewDataLayoutFromString(origDL.LLString(), opSys, arch)
		if err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(origDL, newDL) {
			t.Errorf("data layout mismatch, expected %q, got %q", origDL.LLString(), newDL.LLString())
		}
	}
}
