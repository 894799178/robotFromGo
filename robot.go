package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)


var(

	handle int
	token string = "c3244957-fadc-4114-b1a1-f88d8fd82205"
	bodyContent []byte
	url = "http://localhost:8080/getAutoMessage"
	moduser32 = syscall.NewLazyDLL("user32.dll")
	procSendMessage = moduser32.NewProc("SendMessageA")

)
func main() {
	fmt.Print("请输入句柄和token值(通过空格隔开):--->")
	fmt.Scanln(&handle,&token)

	fmt.Println("句柄:%v token:%v",handle,token)
	song := make(map[string]interface{})
	song["token"] = token
	var data []map[string]string
	for {
		request := sendPostRequest(url, song)
		if len(request) > 2{
			if err := json.Unmarshal([]byte(request), &data); err == nil {
				for _, value := range data {
					fmt.Println("json数据-->", value)
					x,_ := strconv.Atoi(value["pointX"])
					y,_ := strconv.Atoi(value["pointY"])
					mouseClick(uintptr(x),uintptr(y))
				}
			} else {
				fmt.Println("error", err)
			}
		}
	}
}


func mouseClick(x uintptr,y uintptr){
	//延迟1s
 	time.Sleep(time.Duration(1)*time.Second)
	SendMessage(handle,513,0,x+y*65536)//左键按下
	SendMessage(handle,514,0,x+y*65536)//左键弹起
}



func  sendPostRequest(url string,song map[string]interface{}) string {

	bytesData, err := json.Marshal(song)
	if err != nil {
		fmt.Println(err.Error() )
		return ""
	}
	reader := bytes.NewReader(bytesData)

	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	//byte数组直接转成string，优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	return *str
}

func sendGetRequest(url string){

}

func SendMessage(hwnd int, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := procSendMessage.Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam)
	return ret
}