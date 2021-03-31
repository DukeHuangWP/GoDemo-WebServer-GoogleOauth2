package common

import (
	"testing"
)

func Test_CatchPanic(t *testing.T) {
	defer CatchPanic("單元測試")
	panic("這是單元測試") //此panic會遭攔截
}

func Test_CreatFolderAndRemove(t *testing.T) {
	if err := CreatFolderIfNotExist("./UNITEST"); err != nil {
		t.Fatalf("%v", err)
	}

	if err := RemoveFolderPath("./UNITEST"); err != nil {
		t.Fatalf("%v", err)
	}

}
