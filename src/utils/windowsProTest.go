package utils

import (
	"fmt"
	"github.com/axgle/mahonia"
	"strings"
	"syscall"
	"unsafe"
)

const (
	GW_HWNDFIRST = 0
	GW_HWNDLAST  = 1
	GW_HWNDNEXT  = 2
	HWNDPREV     = 3
	GW_OWNER     = 4
	GW_CHILD     = 5
)

var (
	handle1      uintptr
	moduserver32 = syscall.NewLazyDLL("user32.dll")
	//procSendMessage = moduser32.NewProc("SendMessageA")
	findWindowA      = moduserver32.NewProc("FindWindowA")
	getWindowTextA   = moduserver32.NewProc("GetWindowTextA")
	getWindow        = moduserver32.NewProc("GetWindow")
	enumWindows      = moduserver32.NewProc("EnumWindows")
	getDesktopWindow = moduserver32.NewProc("GetDesktopWindow")
)

//
//func main() {
//	//handle1 = FindWindow("", "CocosCreator | sqzz - Google Chrome")
//	//fmt.Println(handle1)
//	//var str = make([]byte, 256)
//	//GetWindowText(handle1, str, 256)
//	//fmt.Println(string(str))
//	//window := GetWindow(handle1, GW_CHILD)
//	//fmt.Println(window)
//	handers := EnumWindowsByTitle("CocosCreator | sqzz - Google Chrome")
//	for _,value := range handers{
//		fmt.Println(value)
//	}
//}

/**
获取一个SQZZ的句柄,需要通过持久化的数据进行排他处理
*/
func getSQZZHandle() {

}

/**
查找顶级窗口
输入空字符串,默认匹配所有
*/
func FindWindow(className string, windowsName string) uintptr {
	ret, _, _ := findWindowA.Call(
		stringToUintptr(className), stringToUintptr(windowsName),
	)
	return ret
}

/**
 * 转换string类型为指针类型的数据(无符号16为数值 uint16)
 */
func stringToUintptr(str string) uintptr {
	if str == "" {
		return 0
	}
	enc := mahonia.NewEncoder("GBK")
	str1 := enc.ConvertString(str) //字符转换
	a1 := []byte(str1)
	p1 := &a1[0] //把字符串转字节指针

	return uintptr(unsafe.Pointer(p1))
}

func GetWindowText(hander uintptr, str []byte, MaxCount uintptr) uintptr {
	ret, _, _ := getWindowTextA.Call(
		hander, uintptr(unsafe.Pointer(&str[0])), MaxCount,
	)
	return ret
}

/**
返回子窗口句柄
*/
func GetWindow(hander uintptr, wCmd uintptr) uintptr {
	ret, _, _ := getWindow.Call(
		hander, wCmd,
	)
	return ret
}

/**
根据进程的标题返回所需要的句柄集合
*/
func EnumWindowsByTitle(title string) []uintptr {
	firstWindow := GetDesktopWindow()
	hander := GetWindow(firstWindow, GW_CHILD)
	maxCount := 256
	var arr = make([]uintptr, 0)
	var str = make([]byte, maxCount)
	for hander != 0 {
		GetWindowText(hander, str, uintptr(maxCount))

		if strings.Contains(string(str), title) {
			fmt.Println(string(str))
			arr = append(arr, hander)
		}
		hander = GetWindow(hander, GW_HWNDNEXT)
	}
	return arr
}

/**
获得桌面的句柄
*/
func GetDesktopWindow() uintptr {
	ret, _, _ := getDesktopWindow.Call()
	return ret
}
