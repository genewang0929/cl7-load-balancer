package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// 1. 從指令列參數讀取 Port
	// os.Args[0] 是程式本身的名稱，os.Args[1] 才是第一個參數
	if len(os.Args) < 2 {
		log.Fatal("請提供 Port 號碼，例如: go run backend.go 8081")
	}
	port := os.Args[1]

	// 2. 定義 Handler
	// 我們讓它回傳它自己的 Port，這樣我們在瀏覽器才看得到是「誰」在服務我們
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request on port %s\n", port)                           // Server端看到的 Log
		fmt.Fprintf(w, "Hello! I am the Backend Server running on Port %s\n", port) // 回傳給 User 的內容
	})

	// 3. 啟動 Server
	log.Printf("Backend Server starting on port :%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
