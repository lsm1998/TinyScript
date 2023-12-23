package native

import "testing"

func Test_nativeFile_WriteFile(t *testing.T) {
	file := &nativeFile{}

	err := file.WriteFile("test.txt", "hello world", "append")
	if err != nil {
		t.Fatal(err)
	}
}
