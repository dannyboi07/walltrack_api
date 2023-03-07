package controller

import (
	"encoding/json"
	"net/http"
	"walltrack/schema"
	"walltrack/util"
)

func AddTransaction(w http.ResponseWriter, r *http.Request) {

	var jDec = json.NewDecoder(r.Body)
	jDec.DisallowUnknownFields()

	var (
		AddTransactionReq schema.AddTransaction
		statusCode        int
		err               error
	)
	statusCode, err = util.JsonParseErr(jDec.Decode(&AddTransactionReq))
	if err != nil {
		util.WriteApiErrMessage(w, statusCode, err.Error())
		util.Log.Println("Failed to parse request, err:", err)
		return
	}

	statusCode, err = AddTransactionReq.Validate()
	if err != nil {
		util.WriteApiErrMessage(w, statusCode, err.Error())
		util.Log.Println("Invalid add transaction request, err:", err)
		return
	}
}
