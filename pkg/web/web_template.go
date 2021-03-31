package web

import (
	"bytes"
	"fmt"
	"log"
	"minecraft-richnet-addin/pkg/common/fileInfo"
	"minecraft-richnet-addin/pkg/global"
	"minecraft-richnet-addin/pkg/global/googleOauth2"
	"minecraft-richnet-addin/pkg/web/encrypt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//Design Pattern : Factory
// type Handlers interface {
// 	Action(account *googleOauth2.GoogleAcc, w http.ResponseWriter, r *http.Request)
// }

type Templates interface {
	Action(ginCTX *gin.Context, account *googleOauth2.GoogleAcc)
}

func NewTemplate(stateType string) (Templates, error) {

	postFrom := make(map[string]string)
	if strings.Count(stateType, "~") == 1 {
		decode, err := encrypt.Decode(stateType[strings.Index(stateType, "~")+1:])
		if err != nil {
			log.Printf("Post解碼錯誤 : %v (%v)", decode, err)
		} else {
			for _, lineStr := range strings.Split(decode, "\n") {
				if strings.Count(lineStr, "\t") != 1 {
					log.Printf("Post解碼異常 : %v (%v)", decode, lineStr)
					break
				}
				postFrom[lineStr[:strings.Index(lineStr, "\t")]] = lineStr[strings.Index(lineStr, "\t")+1:]
			}
		}
		stateType = stateType[:strings.Index(stateType, "~")]
	}

	// fmt.Println(postFrom)
	// fmt.Println(stateType)

	switch stateType {

	case global.CallbackType_IndexMenu: //轉入首頁選單
		return &templIndexMenu{}, nil

	case global.CallbackType_Test: //測試callback認證功能
		return &templTest{}, nil

	case global.CallbackType_PostTest: //執行測試腳本
		return &templPostTest{fromData: postFrom}, nil
	case global.CallbackType_RunSh: //執行測試腳本
		return &templRunSh{}, nil
	case global.CallbackType_Download: //Mincraft伺服器檔案 - 下載檔案
		return &templDownloadList{}, nil
	case global.CallbackType_Upload: //Mincraft伺服器檔案 - 上傳檔案
		return &templUploadMenu{}, nil

	default:
		return nil, fmt.Errorf("get url query state 設置錯誤！")
	}

}

type templIndexMenu struct{}

func (_ *templIndexMenu) Action(ginCTX *gin.Context, _ *googleOauth2.GoogleAcc) {
	ginCTX.Redirect(http.StatusFound, "./"+global.FolderPath_Pages+"/"+global.Template_IndexMenu)
}

type templTest struct{}

func (_ *templTest) Action(ginCTX *gin.Context, account *googleOauth2.GoogleAcc) {
	ginCTX.Writer.WriteString(fmt.Sprintf("Your Email: %v\n", account.Email))
}

type templPostTest struct {
	fromData map[string]string
}

func (templ *templPostTest) Action(ginCTX *gin.Context, account *googleOauth2.GoogleAcc) {

	if len(templ.fromData) > 0 {
		ginCTX.Writer.WriteString(fmt.Sprintf("Your Email: %v\n", account.Email))
		ginCTX.Writer.WriteString(fmt.Sprintf("POST FromData: %v\n", templ.fromData))
		ginCTX.Abort()
		return
	}

	ginCTX.HTML(http.StatusOK, global.Template_PostTest, gin.H{
		global.H5Value_LastUpdate: global.LastOnLineTime,
		global.H5Value_Title:      "Post 頁面測試",
		global.H5Page_Menu:        "./" + global.FolderPath_Pages + "/" + global.Template_IndexMenu,
		global.H5Value_Login:      "." + global.URNPath_Login + "?state=" + global.CallbackType_PostTest,
	})

}

type templRunSh struct{}

func (_ *templRunSh) Action(ginCTX *gin.Context, _ *googleOauth2.GoogleAcc) {
	outBytes, err := exec.Command("/bin/sh", global.Model_RunShPath).Output()
	outputMessage := string(outBytes)
	if err != nil {
		outputMessage = "Script path : " + global.Model_RunShPath + "\r\n" + err.Error()
	}

	ginCTX.HTML(http.StatusOK, global.Template_RunSh, gin.H{
		global.H5Value_LastUpdate: global.LastOnLineTime,
		global.H5Value_Title:      "RunSh 腳本執行測試",
		global.H5Page_Menu:        "./" + global.FolderPath_Pages + "/" + global.Template_IndexMenu,
		"HTMLValue_RunScript":     "." + global.URNPath_Login + "?state=" + global.CallbackType_RunSh,
		"HTMLValue_OutputMessage": outputMessage,
		"HTMLValue_OutputReturn":  "執行腳本時間:" + time.Now().Format("2006-01-02_15:04:05"),
	})

}

type templDownloadList struct{}

func (_ *templDownloadList) Action(ginCTX *gin.Context, _ *googleOauth2.GoogleAcc) {

	fileList, _ := fileInfo.FileTree(global.FolderPath_Public + "/" + global.FolderPath_Downloads)
	pathStr := global.FolderPath_Public[strings.Index(global.FolderPath_Public, "./")+2:] + "/" + global.FolderPath_Downloads
	var downloadLinks string
	buffer := bytes.NewBuffer([]byte{})
	for _, value := range fileList {
		subPath := value.FilePath[strings.LastIndex(value.FilePath, pathStr)+len(pathStr):]
		buffer.WriteString(fmt.Sprintf(`<a href="%v" download="%v" >%v</a><br>`, "./"+global.FolderPath_Downloads+"/"+subPath, value.FileName, subPath))
	}
	downloadLinks = buffer.String()

	ginCTX.HTML(http.StatusOK, global.Template_DownloadList, gin.H{
		global.H5Value_LastUpdate:       global.LastOnLineTime,
		global.H5Value_Title:            "下載清單",
		global.H5Page_Menu:              "./" + global.FolderPath_Pages + "/" + global.Template_IndexMenu,
		"HTMLValue_OutputDownloadLinks": downloadLinks,
		"HTMLValue_OutputReturn":        "檔案掃描時間:" + time.Now().Format("2006-01-02_15:04:05"),
	})

}

type templUploadMenu struct{}

func (_ *templUploadMenu) Action(ginCTX *gin.Context, account *googleOauth2.GoogleAcc) {

	token, err := encrypt.Encode(fmt.Sprint(time.Now().Unix()) + ":" + account.Email)
	if err != nil {
		ginCTX.HTML(http.StatusOK, global.Template_UploadMenu, gin.H{
			global.H5Value_LastUpdate: global.LastOnLineTime,
			global.H5Value_Title:      "上傳檔案選單:" + err.Error(),
			global.H5Page_Menu:        "./" + global.FolderPath_Pages + "/" + global.Template_IndexMenu,
			global.H5Value_Upload:     "." + global.URNPath_Upload + "?state=" + global.CallbackType_PostTest,
			"HTMLValue_OutputReturn":  "處理時間:" + time.Now().Format("2006-01-02_15:04:05"),
		})
		return
	}

	ginCTX.HTML(http.StatusOK, global.Template_UploadMenu, gin.H{
		global.H5Value_LastUpdate:    global.LastOnLineTime,
		global.H5Value_Title:         "上傳檔案選單:",
		global.H5Page_Menu:           "./" + global.FolderPath_Pages + "/" + global.Template_IndexMenu,
		global.H5Value_Upload:        "." + global.URNPath_Upload + "?state=" + global.CallbackType_PostTest,
		global.H5TAGName_Token:       global.H5TAGName_Token,
		global.H5TAGName_Upload_Type: global.H5TAGName_Upload_Type,
		global.H5Value_Token:         token,
		global.H5Value_UploadType:    global.UploadType_Test,
		"HTMLValue_OutputReturn":     "",
	})
	return
}
