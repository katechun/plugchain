package test

import (
	"context"
	"encoding/hex"
	"fmt"
	cliTx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/oracleNetworkProtocol/liquidity/app"
	"pluschainapi/handler"

	"google.golang.org/grpc"
	"net/http"
	"pluschainapi/tool"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func SendTx(w http.ResponseWriter, r *http.Request) {
	ret := tool.Result{}
	from := tool.ResultVal(r, "from")
	to := tool.ResultVal(r, "to")
	priv := tool.ResultVal(r, "priv")
	amount := tool.ResultVal(r, "amount")
	gas := tool.ResultVal(r, "gas")

	fmt.Println("from:", from, " to:", to)

	// 选择您的编解码器：Amino 或 Protobuf
	encCfg := app.MakeEncodingConfig()

	// 创建一个新的 TxBuilder。
	txBuilder := encCfg.TxConfig.NewTxBuilder()
	fmt.Println("11111")
	addr1, err := sdk.AccAddressFromBech32(from)
	if err != nil {
		fmt.Println("Err:", err)
		ret.Code = 400
		ret.Msg = err.Error()
		tool.Jsonm(w, ret)
		return
	}

	fmt.Println("2222x")
	addr2, err := sdk.AccAddressFromBech32(to)

	if err != nil {
		log.Fatal(err.Error())
		ret.Code = 400
		ret.Msg = err.Error()
		tool.Jsonm(w, ret)
		return
	}

	fmt.Println("addr1:", addr1, " addr2:", addr2)

	//发起者私钥
	privB, err := hex.DecodeString(priv)
	if err != nil {
		log.Fatal(err.Error())
		ret.Code = 400
		ret.Msg = err.Error()
		tool.Jsonm(w, ret)
		return
	}

	priv1 := secp256k1.PrivKey{Key: privB}
	accountSeq := uint64(1)
	acountNumber := uint64(0)

	fmt.Println("22222")
	amount1, err := strconv.Atoi(amount)
	msg1 := banktypes.NewMsgSend(addr1, addr2, sdk.NewCoins(sdk.NewInt64Coin("plug", int64(amount1))))
	err = txBuilder.SetMsgs(msg1)
	if err != nil {
		log.Fatal(err.Error())
		ret.Code = 400
		ret.Msg = err.Error()
		tool.Jsonm(w, ret)
		return
	}

	fmt.Println("33333")
	txBuilder.SetGasLimit(200000)
	gas1, err := strconv.Atoi(gas)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewInt64Coin("plug", int64(gas1))))
	txBuilder.SetMemo("transfer coin")
	// txBuilder.SetTimeoutHeight(...)

	//第一轮：我们收集所有签名者信息。 我们使用“设置空签名”技巧来做到这一点
	sign := signing.SignatureV2{
		PubKey: priv1.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  encCfg.TxConfig.SignModeHandler().DefaultMode(),
			Signature: nil,
		},

		Sequence: accountSeq,
	}

	fmt.Println("44444")
	err = txBuilder.SetSignatures(sign)
	if err != nil {
		if err != nil {
			log.Fatal(err.Error())
			ret.Code = 400
			ret.Msg = err.Error()
			tool.Jsonm(w, ret)
			return
		}
	}

	fmt.Println("5555")
	//第二轮： 设置所有签名者信息，因此每个签名者都可以签名。
	sign = signing.SignatureV2{}
	signerD := xauthsigning.SignerData{
		ChainID:       handler.chainID,
		AccountNumber: acountNumber,
		Sequence:      accountSeq,
	}
	sign, err = cliTx.SignWithPrivKey(
		encCfg.TxConfig.SignModeHandler().DefaultMode(), signerD,
		txBuilder, cryptotypes.PrivKey(&priv1), encCfg.TxConfig, accountSeq)

	if err != nil {
		panic(err)
	}
	err = txBuilder.SetSignatures(sign)
	if err != nil {
		panic(err)
	}

	fmt.Println("66666")
	// 生成的 Protobuf 编码字节。
	txBytes, err := encCfg.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		log.Fatal(err.Error())
		ret.Code = 400
		ret.Msg = err.Error()
		tool.Jsonm(w, ret)
		return
	}

	// 创建一个grpc服务
	grpcConn, err := grpc.Dial(
		handler.grpcURL,     // 你的 gRPC 服务器地址。
		grpc.WithInsecure(), // SDK 不支持任何传输安全机制。
	)
	if err != nil {
		log.Fatal(err.Error())
		ret.Code = 400
		ret.Msg = err.Error()
		tool.Jsonm(w, ret)
		return
	}
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
		log.Fatal(err.Error())
		ret.Code = 400
		ret.Msg = err.Error()
		tool.Jsonm(w, ret)
		return
	}

	//fmt.Println("777777")
	//if grpcRes.TxResponse.Code != 0 {
	//	//log.Fatal(err.Error())
	//	ret.Code = 400
	//	ret.Msg = "Send coin error!  error:" + strconv.Itoa(int(grpcRes.TxResponse.Code))
	//	tool.Jsonm(w, ret)
	//	return
	//}

	ret.Code = 200
	ret.Msg = "txhash:" + grpcRes.TxResponse.TxHash + " resp::" + grpcRes.TxResponse.String()
	tool.Jsonm(w, ret)

}
