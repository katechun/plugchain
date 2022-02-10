package handler

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"pluschainapi/tool"
)

//GET /cosmos/distribution/v1beta1/validators/{validator_address}/outstanding_rewards

func QueryOutstandingRewards(w http.ResponseWriter, r *http.Request) {
	ret := tool.Result{}
	addr := tool.ResultVal(r, "address")
	url := "http://" + apiURL + "/cosmos/distribution/v1beta1/validators/" + addr + "/outstanding_rewards"

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
