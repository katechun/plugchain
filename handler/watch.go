package handler

import (
	"context"
	"fmt"
	"github.com/tendermint/tendermint/libs/pubsub/query"
	"github.com/tendermint/tendermint/types"
)

func Watch() {

	client := client.NewHTTP(nodeURI, "/websocket")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	query := query.MustParse("tm.event = 'Tx' AND account.owner CONTAINS 'xx' AND  message.action='send'")
	txs := make(chan interface{})
	err := client.Subscribe(ctx, "test-client", query, txs)

	go func() {

		for e := range txs {
			fmt.Println("got ", e.(types.EventDataTx))
		}
	}()

}
