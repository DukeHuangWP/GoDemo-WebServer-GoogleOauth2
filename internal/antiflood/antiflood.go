package antiflood

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

var antFlood_Count uint8 = 30           //短時間連線次數上限
var antFlood_BetSecs int64 = 5          //短時間連線秒數計算
var antFlood_TimerClearSecs int64 = 600 //多久後解封ClientIP時間
var antFlood_TimerBetSecs int = 1800    //ClientIP監控執行周期

type clientIPType struct {
	ConnectCount uint8 //短時間連線次數
	LastTime     int64 //客戶端最後連線時間戳
}

//[clientIP]資訊
var antFlood_ClientIPs = make(map[string]*clientIPType)

const banPageHTML =  `<!DOCTYPE html><html><body><img src="https://stickershop.line-scdn.net/stickershop/v1/sticker/34521489/IOS/sticker.png"></body></html>`

//啟動計時goroutine
func StartTimer() {
	//Timer定時器,定時清除閒置clientIP
	antFlood_TimerTicker := time.NewTicker(time.Duration(antFlood_TimerBetSecs) * time.Second)
	defer antFlood_TimerTicker.Stop()
	go func(ticker *time.Ticker) {
		for {
			<-ticker.C
			if antFlood_ClientIPs != nil {
				clearClientIPList(antFlood_ClientIPs)
			}
		}
	}(antFlood_TimerTicker)
}

//啟用反SYN-Flood機制
func clearClientIPList(antFlood_ClientIPs map[string]*clientIPType) {
	for index, value := range antFlood_ClientIPs {
		if time.Now().Unix()-value.LastTime > antFlood_TimerClearSecs {
			log.Printf("ClientIP 監控佔存清除 > %v (Count:'%v', LastTime:'%v')", index, value.ConnectCount, value.LastTime)
			delete(antFlood_ClientIPs, index)
		}
	}
}

//觸發Handler(計算連線次數)
func TriggerHandler(ginCTX *gin.Context) {
	vs_ClientIP := ginCTX.ClientIP()
	if antFlood_ClientIPs[vs_ClientIP] == nil {

		antFlood_ClientIPs[vs_ClientIP] = &clientIPType{ConnectCount: 1, LastTime: time.Now().Unix()}

	} else {

		if time.Now().Unix()-antFlood_ClientIPs[vs_ClientIP].LastTime < antFlood_BetSecs {
			antFlood_ClientIPs[vs_ClientIP].ConnectCount++
		} else if antFlood_ClientIPs[vs_ClientIP].ConnectCount > 0 {
			antFlood_ClientIPs[vs_ClientIP].ConnectCount--
		}
		antFlood_ClientIPs[ginCTX.ClientIP()].LastTime = time.Now().Unix()

		if antFlood_ClientIPs[vs_ClientIP].ConnectCount > antFlood_Count {

			ginCTX.Writer.WriteString(banPageHTML)
			ginCTX.Abort()
			return
		}
	}
	//fmt.Println(antFlood_ClientIPs[vs_ClientIP])
}
