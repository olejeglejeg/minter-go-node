package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/rpc/core/types"
	"net/http"
	"strconv"
)

func Block(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	height, _ := strconv.ParseInt(vars["height"], 10, 64)

	result := new(core_types.ResultBlock)
	_, err := client.Call("block", map[string]interface{}{
		"height": height,
	}, result)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Code:   0,
			Result: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(Response{
		Code:   0,
		Result: result,
	})

	if err != nil {
		panic(err)
	}
}
