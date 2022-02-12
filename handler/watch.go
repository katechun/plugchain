package handler

import (
	"context"
	"fmt"
	"github.com/tendermint/tendermint/libs/pubsub/query"
	"github.com/tendermint/tendermint/rpc/client"
	"time"
)

func Watch() {
	timeout := 8 * time.Second

	c := client.NewHTTP(nodeURI, "/websocket")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	q := query.MustParse("tm.event = 'Tx' AND account.owner CONTAINS 'xx' AND  message.action='send'")
	outs, err := c.Subscribe(ctx, "test-client", q.String(), 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {

		for e := range outs {
			fmt.Println("got ", e)
		}
	}()

}
