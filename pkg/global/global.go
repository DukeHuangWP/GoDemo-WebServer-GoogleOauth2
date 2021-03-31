package global

import (
	"fmt"
	"log"
	"minecraft-richnet-addin/pkg/common"
	"minecraft-richnet-addin/pkg/global/config"
	"minecraft-richnet-addin/pkg/global/customVar"
	"minecraft-richnet-addin/pkg/global/googleOauth2"
	"strings"
)

const (
	ConfigFilePath = "./configs/webstie.conf" //環境變數設定檔

	WebURLRoot           = "minecraft" //網頁根目錄 http://localhost:8080/xxx
	FolderPath_Public    = "./website" //網頁公用檔案目錄
	FolderPath_Images    = "images"    //網頁公用目錄名稱 - 圖片
	FolderPath_Pages     = "pages"     //網頁公用目錄名稱 - 純html頁面
	FolderPath_Templs    = "templates" //網頁公用目錄名稱 - 版模
	FolderPath_Downloads = "downloads" //網頁公用目錄名稱 - 下載目錄
	FolderPath_Uploads   = "uploads"   //網頁公用目錄名稱 - 上傳目錄
	FolderPath_Others    = "others"    //網頁公用目錄名稱 - 其他檔案分類

	URNPath_Login  = "/login"  //登陸Path名稱
	URNPath_Upload = "/upload" //上傳檔案Path名稱

	CallbackType_IndexMenu = "indexMenu" //首頁
	CallbackType_Test      = "test"
	CallbackType_PostTest  = "postTest"
	CallbackType_RunSh     = "runsh"
	CallbackType_Download  = "dowlnoad"
	CallbackType_Upload    = "upload"

	UploadType_Test = "uploadTest" //上傳檔案Post 標籤名稱

	Template_IndexMenu    = "index-Menu.html"
	Template_PostTest     = "Post.html"
	Template_RunSh        = "RunSh.html"
	Template_DownloadList = "Download-List.html"
	Template_UploadMenu   = "Upload-Menu.html"

	H5Page_Menu           = "HTMLPage_Menu"
	H5Value_Login         = "HTMLValue_ActionLogin"
	H5Value_Upload        = "HTMLValue_ActionUpload"
	H5Value_Title         = "HTMLValue_Title"
	H5Value_LastUpdate    = "HTMLValue_LastUpdate"
	H5Value_Token         = "HTMLValue_Token"
	H5Value_UploadType    = "HTMLValue_Upload_Type"
	H5TAGName_Token       = "HTMLOutput_Token"
	H5TAGName_Upload_Type = "HTMLOutput_Upload_Type"
)

var (
	LastOnLineTime string //伺服器啟動時間

	WebURLPage = WebURLRoot + "/" + FolderPath_Pages //直接顯示頁面

	Oauth2ClientID     string   //= "32054680383-r6adfs24i84edvqtfr3eesskums6l2n1.apps.googleusercontent.com"
	Oauth2ClientSecret string   //= "_cha_wE8IT-pMGeGUDrbcn9j"
	Oauth2CallbackName string   //= "callback"
	Oauth2EmailList    []string //准許的google email 清單

	HostScheme     string //= "http://"
	HostDomainName string //= "localhost"
	HostPort       string //= "8080"

	Model_RunShPath string //重啟伺服器腳本 //= "./configs/restart_script.sh"

	Minecraft_DownloadList_Excludes map[string]struct{}    //下載清單中要忽略的檔案名稱
	Minecraft_DownloadList_Settings map[string]interface{} //下載清單設定
	SettingMap                      = config.SettingMap    //參照package config

	//GoogleOauth2Config googleOauth2.Config
)

func init() {

	//環境變數設定值解析 SettingMap["XXX"] = &config.AddSetings{初始值, 輸出變數指標, 自訂設定值類型}
	SettingMap["host_scheme"] = &config.AddSetings{DefaultValue: "http://", OutputPointer: &HostScheme, Custom: &customVar.StringType{}}   //設定變數google Callback host_scheme ,預設值: "http://"
	SettingMap["host_domain_name"] = &config.AddSetings{DefaultValue: "", OutputPointer: &HostDomainName, Custom: &customVar.StringType{}} //設定變數 ,預設值: ""
	SettingMap["host_port"] = &config.AddSetings{DefaultValue: "8080", OutputPointer: &HostPort, Custom: &customVar.Uint16Type{}}          //設定變數 ,預設值: "8080"

	SettingMap["google_oauth2_clientID"] = &config.AddSetings{DefaultValue: "", OutputPointer: &Oauth2ClientID, Custom: &customVar.StringType{}}                  //設定變數 google_oauth2_clientID,預設值: ""
	SettingMap["google_oauth2_secret_code"] = &config.AddSetings{DefaultValue: "", OutputPointer: &Oauth2ClientSecret, Custom: &customVar.StringType{}}           //設定變數 google_oauth2_secret_code ,預設值: ""
	SettingMap["google_oauth2_callback_path"] = &config.AddSetings{DefaultValue: "callback", OutputPointer: &Oauth2CallbackName, Custom: &customVar.StringType{}} //設定變數 google_oauth2_callback_path ,預設值: "callback"
	SettingMap["add_google_email"] = &config.AddSetings{DefaultValue: "", OutputPointer: &Oauth2EmailList, Custom: &customVar.AddStringSlice{}}                   //設定變數 准許的google E-Mail地址,預設值: ""

	SettingMap["path_script"] = &config.AddSetings{DefaultValue: "./configs/script.sh", OutputPointer: &Model_RunShPath, Custom: &customVar.StringType{}} //設定變數 ,預設值: "./configs/script.sh"

	config.InitLoad(ConfigFilePath) //載入環境設定檔案
	log.Printf("host_scheme=%v\n", HostScheme)
	log.Printf("host_domain_name=%v\n", HostDomainName)
	log.Printf("host_port=%v\n", HostPort)
	log.Printf("google_oauth2_clientID=%v\n", Oauth2ClientID)
	log.Printf("google_oauth2_secret_code=%v\n", Oauth2ClientSecret)
	log.Printf("google_oauth2_callback_path=%v\n", Oauth2CallbackName)
	log.Printf("add_google_email=%v\n", Oauth2EmailList)
	log.Printf("path_script=%v\n", Model_RunShPath)

	if err := common.CreatFolderIfNotExist(FolderPath_Public + "/" + FolderPath_Uploads); err != nil {
		log.Fatal(err)
	} //檢查並建立webserver上傳目錄

	redirectURL := HostScheme + HostDomainName + ":" + HostPort + "/" + WebURLRoot + "/" + Oauth2CallbackName
	googleOauth2.InitConfig(redirectURL, Oauth2ClientID, Oauth2ClientSecret) //set google credentials

}

//環境設定轉成自定義格式輸出到變數 XXX=key|value
type filesMapType struct{}

func (_ filesMapType) GetValue(inputValue string) (output interface{}, err error) {

	cacheSlice := strings.Split(inputValue, "|")
	if len(cacheSlice) < 1 {
		return nil, fmt.Errorf("'%v' > ' | ' 分隔號數量不可低於1", inputValue)
	} else if cacheSlice[0] == "" || cacheSlice[1] == "" {
		return nil, fmt.Errorf("'%v' > 部份數值不可為空", inputValue)
	}

	cacheMap := make(map[string]interface{})
	cacheMap[cacheSlice[0]] = cacheSlice[1]
	return cacheMap, nil
}
