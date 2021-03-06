package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/MinterTeam/minter-go-node/cmd/utils"
	"github.com/MinterTeam/minter-go-node/core/minter"
	"github.com/MinterTeam/minter-go-node/core/state"
	"github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/rpc/lib/client"
	"strconv"
	"time"
)

var (
	blockchain *minter.Blockchain
	client     *rpcclient.JSONRPCClient
)

func RunApi(b *minter.Blockchain) {
	client = rpcclient.NewJSONRPCClient(*utils.TendermintRpcAddrFlag)
	core_types.RegisterAmino(client.Codec())

	blockchain = b

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/candidate/{pubkey}", GetCandidate).Methods("GET")
	router.HandleFunc("/api/balance/{address}", GetBalance).Methods("GET")
	router.HandleFunc("/api/balanceWS", GetBalanceWatcher)
	router.HandleFunc("/api/transactionCount/{address}", GetTransactionCount).Methods("GET")
	router.HandleFunc("/api/sendTransaction", SendTransaction).Methods("POST")
	router.HandleFunc("/api/sendTransactionSync", SendTransactionSync).Methods("POST")
	router.HandleFunc("/api/sendTransactionAsync", SendTransactionAsync).Methods("POST")
	router.HandleFunc("/api/transaction/{hash}", Transaction).Methods("GET")
	router.HandleFunc("/api/block/{height}", Block).Methods("GET")
	router.HandleFunc("/api/transactions", Transactions).Methods("GET")
	router.HandleFunc("/api/status", Status).Methods("GET")
	router.HandleFunc("/api/coinInfo/{symbol}", GetCoinInfo).Methods("GET")
	router.HandleFunc("/api/estimateCoinExchangeReturn", EstimateCoinExchangeReturn).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	// wait for tendermint to start
	for true {
		result := new(core_types.ResultHealth)
		_, err := client.Call("health", map[string]interface{}{}, result)
		if err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}

	log.Fatal(http.ListenAndServe(":8841", handler))
}

type Response struct {
	Code   uint32      `json:"code"`
	Result interface{} `json:"result,omitempty"`
	Log    string      `json:"log,omitempty"`
}

func GetStateForRequest(r *http.Request) *state.StateDB {
	height, _ := strconv.Atoi(r.URL.Query().Get("height"))

	cState := blockchain.CurrentState()

	if height > 0 {
		cState, _ = blockchain.GetStateForHeight(height)
	}

	return cState
}
