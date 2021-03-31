package web

import (
	"bytes"
	"fmt"
	"log"
	"minecraft-richnet-addin/pkg/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FileUploader interface {
	Action(ginCTX *gin.Context)
}

func NewFileUploader(uploadType string) (FileUploader, error) {

	switch uploadType {

	case global.UploadType_Test:
		return &FileUpdaterAllowList{}, nil
	default:
		return nil, fmt.Errorf("get url query uploadType 設置錯誤！")
	}

}

type FileUpdaterAllowList struct{}

//func (_ *FileUpdaterAllowList) Action(ginCTX *gin.Context) {}
func (_ *FileUpdaterAllowList) Action(ginCTX *gin.Context) {

	//filesMap, foldersMap := minecraft.GetUploadlist()
	muiltFilesForm, _ := ginCTX.MultipartForm()
	savePath := global.FolderPath_Public + "/" + global.FolderPath_Uploads
	var getFileList []string
	for fromName, muiltFilesHeaders := range muiltFilesForm.File {
		for _, headerCache := range muiltFilesHeaders {
			log.Printf("接收到上傳檔案 : %v(%v)(%v)", headerCache.Filename, fromName, headerCache.Size)

			//if fileMap, isExsit := filesMap[fromName]; isExsit {
			//for filePath := range fileMap {
			getFileList = append(getFileList, headerCache.Filename)
			filePath := savePath + "/" + headerCache.Filename
			ginCTX.SaveUploadedFile(headerCache, filePath) // 接收上傳單檔
			log.Printf("成功接收上傳檔案%v", filePath)
			//}
			//}

		}
	}

	var outputReturn bytes.Buffer
	for _, fileName := range getFileList {
		outputReturn.WriteString(fileName + "\n")
	}
	outputReturn.WriteString(fmt.Sprintf("%d files uploaded!", len(getFileList)))
	ginCTX.String(http.StatusOK, outputReturn.String())
	//minecraft.CopyToDownload()

	return
}
