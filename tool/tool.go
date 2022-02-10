package tool

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Result struct {
	Code int
	Msg  string
}

func Jsonm(w http.ResponseWriter, ret Result) {
	ret_json, _ := json.Marshal(ret)
	if _, err := io.WriteString(w, string(ret_json)); err != nil {
		log.Fatal(err)
	}
}

func ResultVal(r *http.Request, name string) string {
	val := mux.Vars(r)
	vv := val[name]
	return vv
}
