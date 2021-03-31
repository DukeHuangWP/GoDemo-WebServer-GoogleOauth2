package fileInfo

import (
	"os"
	"testing"
)

func Test_CreatFolderAndRemove(t *testing.T) {

	path := "empty.txt"
	emptyFile, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	emptyFile.Close()

	exclude_macos, err := os.Create(".DS_Store")
	if err != nil {
		t.Fatal(err)
	}
	exclude_macos.Close()

	exclude_win, err := os.Create("Thumbs.db")
	if err != nil {
		t.Fatal(err)
	}
	exclude_win.Close()

	if fileList, err := FileTree("./"); err != nil {
		if len(fileList) < 0 {
			t.Fatal("目錄內檔案不應為空")
		}

		for _, value := range fileList {
			if value.FileName == ".DS_Store" || value.FileName == "Thumbs.db" {
				t.Fatal("不應出現OS暫存檔案")
			}
		}
	}

	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}

	if err := os.Remove(".DS_Store"); err != nil {
		t.Fatal(err)
	}

	if err := os.Remove("Thumbs.db"); err != nil {
		t.Fatal(err)
	}
}
