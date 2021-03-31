package antiflood

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_TryAntiSYNFlood(t *testing.T) {

	ginRouter := gin.Default()
	ginRouter.Use(func(ginCTX *gin.Context) {
		TriggerHandler(ginCTX) //反SYN-Flood
	})

	ginRouter.Any("/", func(ginCTX *gin.Context) {
	})

	StartTimer()

	go func() {
		if err := ginRouter.Run(":8080"); err != nil {
			log.Fatal("請檢查localhost:8080是否正常開啟！")
		}
	}()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8080", nil)

	for index := uint8(0); index < antFlood_Count+10; index++ {
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("測試連線過程發生錯誤: %v\n", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Get response 值非 200: %v\n", resp.Status)
		}

		if index > antFlood_Count+10 {
			respBody, _ := ioutil.ReadAll(resp.Body)
			if string(respBody) != banPageHTML {
				t.Error("Anti SYN-Flood didnt work!\n")
			}
			log.Println("Anti SYN-Flood is working")
		}

	}
	log.Println("--------測試結束-------")
}
