package fileInfo

import (
	"os"
	"path/filepath"
	"strings"
)

type FileInfo struct {
	FilePath      string
	FileName      string
	FileExtName   string
	FileFolder    string
	FileDirectory string
}

func FileTree(path string) (fileList []*FileInfo,err error)  {

	exclude := make(map[string]struct{})
	exclude[".DS_Store"] = struct{}{}
	exclude["Thumbs.db"] = struct{}{}
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil // not a file.  ignore.
		}

		fileName := info.Name()
		if _, isExsit := exclude[fileName]; isExsit {
			return nil
		}

		fileExtName := ""
		if index := strings.LastIndex(fileName, "."); index > 0 {
			fileExtName = fileName[index+1:]
		}

		fileDirectory := "/"
		if len(path) != len(fileName) {
			fileDirectory = path[:len(path)-len(fileName)-1]
		}

		fileFolder := ""
		if fileDirectory != "/" {
			fileFolder = fileDirectory[strings.LastIndex(fileDirectory, "/")+1:]
		}

		fileList = append(fileList, &FileInfo{FilePath: path, FileName: fileName, FileExtName: fileExtName, FileDirectory: fileDirectory, FileFolder: fileFolder})
		return nil
	})
	return

}
