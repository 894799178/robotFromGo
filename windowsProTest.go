package main

import (
	"fmt"
	"github.com/axgle/mahonia"
	"syscall"
	"unsafe"
)

var (
	handle1      uintptr
	moduserver32 = syscall.NewLazyDLL("user32.dll")
	//procSendMessage = moduser32.NewProc("SendMessageA")
	findWindowA    = moduserver32.NewProc("FindWindowA")
	getWindowTextA = moduserver32.NewProc("GetWindowTextA")
	getNextWindow  = moduserver32.NewProc("GetWindow")
)

func main() {
	handle1 = FindWindow("", "CocosCreator | sqzz - Google Chrome")
	fmt.Println(handle1)
	var str = make([]byte, 256)
	GetWindowText(handle1, str, 256)
	fmt.Println(string(str))
	window := GetNextWindow(handle1, 5)
	fmt.Println(window)
}

/**
查找顶级窗口
输入空字符串,默认匹配所有
*/
func FindWindow(className string, windowsName string) uintptr {
	ret, _, err := findWindowA.Call(
		stringToUintptr(className), stringToUintptr(windowsName),
	)
	if err != nil {
		fmt.Println("err:", err)
	}
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
func GetNextWindow(hander uintptr, wCmd uintptr) uintptr {
	ret, _, err := getNextWindow.Call(
		hander, wCmd,
	)
	if err != nil {
		fmt.Println("err:", err)
	}
	return ret
}
