package handler

import (
	"fmt"
	"net/http"
)

var (
	nodeURI string = "tcp://127.0.0.1:26657"
	grpcURL string = "39.109.104.67:9090"
	apiURL  string = "39.109.104.67:1317"
	chainID string = "plugchain"
)

// Index serves index page
func Index(w http.ResponseWriter, r *http.Request) {
	// 往w里写入内容，就会在浏览器里输出
	fmt.Fprintf(w, "This index page!")
}
