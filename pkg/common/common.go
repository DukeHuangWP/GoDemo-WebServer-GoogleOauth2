package common

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

//攔截Panic，防止非預期錯誤導致當機
func CatchPanic(errorTitle string) {
	if err := recover(); err != nil {
		var logText string
		for index := 0; index < 8; index++ { //最多捕捉10層log
			ptr, filename, line, ok := runtime.Caller(index)
			if !ok {
				break
			}
			logText = logText + fmt.Sprintf(" %v:%d,%v > %v\n", filename, line, runtime.FuncForPC(ptr).Name(), err)
		}
		log.Printf("%v : 發生嚴重錯誤 %v\n", errorTitle, logText)
		return
	}
}

//檢查資料夾是否存在否則建立
func CreatFolderIfNotExist(folderPath string) error {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if err := os.Mkdir(folderPath, 0755); err != nil {
			return err
		}
	}
	return nil
}

//刪除整個資料夾路徑
func RemoveFolderPath(folderPath string) error {
	if err := os.RemoveAll(folderPath); err != nil {
		return err
	}
	return nil
}
