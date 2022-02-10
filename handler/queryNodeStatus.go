package handler

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"pluschainapi/tool"
)

//GET /cosmos/base/tendermint/v1beta1/node_info

func QueryNodeInfo(w http.ResponseWriter, r *http.Request) {
	ret := tool.Result{}
	url := "http://" + apiURL + "/cosmos/base/tendermint/v1beta1/node_info"

	client := &http.Client{}
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err.Error())
		ret.Code = 400
		ret.Msg = err.Error()
		tool.Jsonm(w, ret)
		return
	}
	resp, err1 := client.Do(reqest)
	if err1 != nil {
		log.Fatal(err.Error())
		ret.Code = 400
		ret.Msg = err.Error()
		tool.Jsonm(w, ret)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
		ret.Code = 400
		ret.Msg = err.Error()
		tool.Jsonm(w, ret)
		return
	}
	ret.Code = 200
	ret.Msg = string(body)

	tool.Jsonm(w, ret)
}
