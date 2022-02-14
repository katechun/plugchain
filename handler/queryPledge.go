package handler

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"pluschainapi/tool"
)

//http://39.109.104.67:1317/cosmos/staking/v1beta1/validators/gxvaloper1rkecwy9pjsfd058w0pwa2perquc8xe638h4u5p/delegations

func QueryPledge(w http.ResponseWriter, r *http.Request) {
	ret := tool.Result{}
	addr := tool.ResultVal(r, "validator_address")
	url := "http://" + apiURL + "/cosmos/staking/v1beta1/validators/" + addr + "/delegations"

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
	fmt.Println(ret.Msg)

	tool.Jsonm(w, ret)
}
