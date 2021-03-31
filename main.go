package main

import (
	"fmt"
	"log"
	"minecraft-richnet-addin/internal/antiflood"
	"minecraft-richnet-addin/pkg/global"
	"minecraft-richnet-addin/pkg/web"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	//gin.DisableConsoleColor()

	ginRouter := gin.Default()
	global.LastOnLineTime = fmt.Sprint(time.Now().Format("2006-01-02_15:04:05")) //紀錄Web Server 啟動時間

	//ginRouter.Use中每個路徑都會先加載func
	ginRouter.Use(func(ginCTX *gin.Context) {
		antiflood.TriggerHandler(ginCTX) //反SYN-Flood
	})

	ginRouter.LoadHTMLGlob(global.FolderPath_Public + "/" + global.FolderPath_Templs + "/*.html")         //.LoadHTMLGlob()僅會在板模後傳給客戶端，該html檔案並不會被公開
	ginRouter.StaticFS(global.WebURLPage, http.Dir(global.FolderPath_Public+"/"+global.FolderPath_Pages)) //公開靜態頁面

	ginGroup1 := ginRouter.Group("/" + global.WebURLRoot)
	{
		//StaticFS()可指定目錄內哪些檔案可以被公開
		//ginRouter.Static("BoswerPath", "./Public") //Static()會將目錄內檔案公開出去，較不安全
		//ginRouter.StaticFS("BoswerPath", http.Dir("./Public")) //StaticFS()可指定目錄內哪些檔案可以被公開
		//ginRouter.StaticFile(global.FolderPath_Images, "./Public/gabo.png") //只能公開一個檔案
		ginGroup1.StaticFS(global.FolderPath_Images, http.Dir(global.FolderPath_Public+"/"+global.FolderPath_Images))
		ginGroup1.StaticFS(global.FolderPath_Others, http.Dir(global.FolderPath_Public+"/"+global.FolderPath_Others))
		ginGroup1.StaticFS(global.FolderPath_Downloads, http.Dir(global.FolderPath_Public+"/"+global.FolderPath_Downloads))

		//ginGroup1.GET("/debug", api.HandleDebug)
		ginGroup1.GET("/", web.HandleIndexMenu)                      //根目錄
		ginGroup1.GET(global.URNPath_Login, web.HandleLogin)         //GET 登陸api接口,包含所有服務需通過接口
		ginGroup1.POST(global.URNPath_Login, web.HandleLogin)        //POST 登陸api接口,包含所有服務接需通過接口
		ginGroup1.GET(global.Oauth2CallbackName, web.HandleCallback) //GET 驗證客戶端獲得的GoogleOauth2認證碼,並藉由Query 'state'標籤後開始後續使用服務
		ginGroup1.POST(global.URNPath_Upload, web.HandleUploadFile)  //POST 上傳檔案接口,上傳過程包含Token驗證

	}

	antiflood.StartTimer() //
	if err := ginRouter.Run(":" + global.HostPort); err != nil {
		log.Fatal("HTTP端口監聽失敗 > ", err.Error())
	}

}
