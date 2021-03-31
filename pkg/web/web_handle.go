package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"minecraft-richnet-addin/pkg/global"
	"minecraft-richnet-addin/pkg/global/googleOauth2"
	"minecraft-richnet-addin/pkg/web/encrypt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

//認證 使用者是否在gooogle email 名單內
func IsAllowEmail(email string) bool {
	for _, allowEmail := range global.Oauth2EmailList {
		if allowEmail == email {
			return true
		}
	}
	return false
}

// func HandleDebug(ginCTX *gin.Context) {

// 	log.Printf("Connected with > %v ,return the index Page\n", ginCTX.Request.RemoteAddr)
// 	ginCTX.Writer.WriteString("<a href=\"http://localhost:8080/v1/login?state=test\">http://localhost:8080/v1/login?state=test</a><br>")
// 	ginCTX.Writer.WriteString("<a href=\"http://localhost:8080/v1/login?state=loginTest\">http://localhost:8080/v1/login?state=loginTest</a><br>")
// 	ginCTX.Writer.WriteString("<a href=\"http://localhost:8080/v1/login?state=upload\">http://localhost:8080/v1/login?state=upload</a><br>")

// 	return

// }

/*
URN:/
Response: 307
說明: 將根目錄導向首頁選單
*/
func HandleIndexMenu(ginCTX *gin.Context) {

	ginCTX.Redirect(http.StatusTemporaryRedirect, global.URNPath_Login[1:]+"?state="+global.CallbackType_IndexMenu)
	ginCTX.Abort() //首頁

}

/*
URN:/login
Response: 307
說明: 依據Query:state=值,重新導向登陸googleOauto2進行認證
*/
func HandleLogin(ginCTX *gin.Context) {

	ginCTX.Request.ParseForm()
	//fmt.Println(ginCTX.Request.PostForm)
	buffer := bytes.NewBuffer([]byte{})
	for key, value := range ginCTX.Request.PostForm {
		buffer.WriteString(key + "\t" + value[len(value)-1] + "\n")
	} //post from data

	var state string
	if buffer.Len() == 0 {
		buffer = nil //CG優化
		state = ginCTX.Query("state")
	} else {
		cacheStr := buffer.String()
		encode, err := encrypt.Encode(cacheStr[:len(cacheStr)-1])
		if err != nil {
			log.Printf("Post編碼錯誤 : %v (%v)", ginCTX.Request.PostForm, err)
		}
		state = ginCTX.Query("state") + "~" + encode
	} //Post From Data 將會存入state當中,例如:state=test~3128937

	//state為googleOauth2准許自定義的query的自訂值
	googleOauthURL := googleOauth2.OauthConfig.AuthCodeURL(state) //URI Query : state
	ginCTX.Redirect(http.StatusTemporaryRedirect, googleOauthURL)
	return

}

/*
URN:/callback
Response: html-template
說明: 依據client回傳Query:state=值,由server對googleOauto2進行驗證,驗證正確後即回傳模版
*/
func HandleCallback(ginCTX *gin.Context) {
	//http://localhost:8080/v1/login?state=test&code=4%2F0AY0e-g42CDKUTrW7IG_0k4nI6tCxjILz776Z9zizcAcU4x0BTf1RfmbZWfTrbV0D3_UILQ&scope=email+openid+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email&authuser=0&prompt=none
	token, err := googleOauth2.OauthConfig.Exchange(oauth2.NoContext, ginCTX.Query("code"))
	if err != nil {
		state := ginCTX.Query("state")
		ginCTX.Redirect(http.StatusTemporaryRedirect, global.URNPath_Login[1:]+"?state="+state)
		return
	} //將client接收到的callback資訊,回傳給google伺服器作認證,若認證失敗則要求client重新登陸

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		ginCTX.Writer.WriteString(fmt.Sprintf("Failed getting user info: %s", err.Error()))
		return
	} //googleOauth2認證錯誤

	defer response.Body.Close()

	var account googleOauth2.GoogleAcc
	json.NewDecoder(response.Body).Decode(&account) //save the google account to struct
	if err != nil {
		ginCTX.Writer.WriteString(fmt.Sprintf("Failed reading response body: %s", err.Error()))
		return
	} //googleOauth2回傳格式為json

	if !IsAllowEmail(account.Email) {
		ginCTX.Writer.WriteString("Account fail.")
		log.Printf("有個傢伙亂登入伺服器注意一下: %v", account.Email)
		return
	}

	state := ginCTX.Query("state")
	template, err := NewTemplate(state)
	if err != nil {
		ginCTX.Writer.WriteString(fmt.Sprintf("Somthing error : %v)", err))
		return
	} //依據query內的state參數設置,轉發給api_model作實際的業務處理

	template.Action(ginCTX, &account)
	return
}

/*
URN:/upload
Response: html-template
說明: 依據client回傳Query:state=值,由server對googleOauto2進行驗證,驗證正確後即產生Token並回傳模版
*/
func HandleUploadFile(ginCTX *gin.Context) {

	token := ginCTX.PostForm(global.H5TAGName_Token)
	var errList []error
	decode, err := encrypt.Decode(token)
	if err != nil {
		errList = append(errList, err)
	}

	cache := strings.Split(decode, ":")
	if len(cache) == 2 {
		tokenTime, _ := strconv.ParseInt(cache[0], 10, 64)
		if (time.Now().Unix() - tokenTime) > 300 { //token超過5分鐘
			errList = append(errList, fmt.Errorf("token超過規定時間"))
		} else if !IsAllowEmail(cache[1]) {
			errList = append(errList, fmt.Errorf("Gmail帳號不准許使用！"))
		}
	}

	if len(errList) > 0 {
		log.Println("注意一下token解析異常,留意可能的惡意使用者！")
		ginCTX.Writer.WriteString(fmt.Sprintf("token error : %v)", errList))
		return
	}

	uploadType := ginCTX.PostForm(global.H5TAGName_Upload_Type)
	fileUpdater, err := NewFileUploader(uploadType)
	if err != nil {
		ginCTX.Writer.WriteString(fmt.Sprintf("Somthing error : %v)", err))
		return
	}

	fileUpdater.Action(ginCTX)
	return

}
