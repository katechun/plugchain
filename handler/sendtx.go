package handler

import (
	"context"
	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/oracleNetworkProtocol/liquidity/app"
	"google.golang.org/grpc"
	"net/http"
	"pluschainapi/tool"
	"strconv"

	cliTx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	log "github.com/sirupsen/logrus"
)

func SendTx(w http.ResponseWriter, r *http.Request) {

	ret := tool.Result{}
	from := tool.ResultVal(r, "from")
	to := tool.ResultVal(r, "to")
	amount := tool.ResultVal(r, "amount")
	privkey := tool.ResultVal(r, "privkey")
	gas := tool.ResultVal(r, "gas")

	amount1, err := strconv.Atoi(amount)
	if err != nil {
		log.Fatal(err.Error())
		ret.Code = 400
		ret.Msg = err.Error()
		tool.Jsonm(w, ret)
		return
	}

	gas1, err := strconv.Atoi(gas)
	if err != nil {
		log.Fatal(err.Error())
		ret.Code = 400
		ret.Msg = err.Error()
		tool.Jsonm(w, ret)
		return
	}
	// 选择您的编解码器：Amino 或 Protobuf
	encCfg := app.MakeEncodingConfig()

	// 创建一个新的 TxBuilder。
	txBuilder := encCfg.TxConfig.NewTxBuilder()

	addr1, _ := types.AccAddressFromBech32(from)
	addr2, _ := types.AccAddressFromBech32(to)
	//发起者私钥
	priv := privkey
	privB, _ := hex.DecodeString(priv)
	priv1 := secp256k1.PrivKey{Key: privB}
	accountSeq := uint64(1)
	acountNumber := uint64(0)

	msg1 := banktypes.NewMsgSend(addr1, addr2, types.NewCoins(types.NewInt64Coin("plug", int64(amount1))))
	err = txBuilder.SetMsgs(msg1)
	if err != nil {
		log.Fatal(err.Error())
		ret.Code = 400
		ret.Msg = err.Error()
		tool.Jsonm(w, ret)
		return
	}

	txBuilder.SetGasLimit(200000)
	txBuilder.SetFeeAmount(types.NewCoins(types.NewInt64Coin("plug", int64(gas1))))
	txBuilder.SetMemo("mengniu")
	// txBuilder.SetTimeoutHeight(...)

	// 生成的 Protobuf 编码字节。
	txBytes, err := encCfg.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		log.Fatal(err.Error())
		ret.Code = 400
		ret.Msg = err.Error()
		tool.Jsonm(w, ret)
		return
	}

	//第一轮：我们收集所有签名者信息。 我们使用“设置空签名”技巧来做到这一点
	sign := signing.SignatureV2{
		PubKey: priv1.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  encCfg.TxConfig.SignModeHandler().DefaultMode(),
			Signature: nil,
		},

		Sequence: accountSeq,
	}

	err = txBuilder.SetSignatures(sign)
	if err != nil {
		log.Fatal(err.Error())
		ret.Code = 400
		ret.Msg = err.Error()
		tool.Jsonm(w, ret)
		return
	}

	//第二轮： 设置所有签名者信息，因此每个签名者都可以签名。
	sign = signing.SignatureV2{}
	signerD := xauthsigning.SignerData{
		ChainID:       chainID,
		AccountNumber: acountNumber,
		Sequence:      accountSeq,
	}
	sign, err = cliTx.SignWithPrivKey(
		encCfg.TxConfig.SignModeHandler().DefaultMode(), signerD,
		txBuilder, cryptotypes.PrivKey(&priv1), encCfg.TxConfig, accountSeq)

	if err != nil {
		log.Fatal(err.Error())
		ret.Code = 400
		ret.Msg = err.Error()
		tool.Jsonm(w, ret)
		return
	}
	err = txBuilder.SetSignatures(sign)
	if err != nil {
		log.Fatal(err.Error())
		ret.Code = 400
		ret.Msg = err.Error()
		tool.Jsonm(w, ret)
		return
	}

	// 创建一个grpc服务
	grpcConn, err := grpc.Dial(
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
		log.Fatal(err.Error())
		ret.Code = 400
		ret.Msg = err.Error()
		tool.Jsonm(w, ret)
		return
	}

	if grpcRes.TxResponse.Code != 0 {
		log.Fatal(err.Error())
		ret.Code = 400
		ret.Msg = "Send coin error!"
		tool.Jsonm(w, ret)
		return
	}

	ret.Code = 200
	ret.Msg = "txhash:" + grpcRes.TxResponse.TxHash
	tool.Jsonm(w, ret)
}
