package handler

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/types/tx"
)

func sendTx() error {
	// --剪断--

	// 创建一个grpc服务
	grpcConn, _ := grpc.Dial(
		grpcURL,             // 你的 gRPC 服务器地址。
		grpc.WithInsecure(), // SDK 不支持任何传输安全机制。
	)
	defer grpcConn.Close()

	// 通过 gRPC 广播 tx。 我们为 Protobuf Tx 服务创建了一个新客户端。
	txClient := tx.NewServiceClient(grpcConn)
	//然后我们在这个客户端上调用 BroadcastTx 方法。
	grpcRes, err := txClient.BroadcastTx(
		context.Background(),
		&tx.BroadcastTxRequest{
			Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: txBytes,
		},
	)
	if err != nil {
		return err
	}

	fmt.Println(grpcRes.TxResponse.Code) // 如果 tx 成功，则应为 0

	return nil
}
