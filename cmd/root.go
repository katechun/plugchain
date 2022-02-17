package cmd

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"pluschainapi/handler"
)

var (
	cfgFile string
)

func init() {
	flags := rootCmd.PersistentFlags()
	flags.StringP("bind-address", "b", "0.0.0.0:7770", "Address to bind to")
	viper.BindPFlags(flags)
}

var rootCmd = &cobra.Command{
	Use:   "pluschainapi",
	Short: "plushchainapi short",
	Long:  `plushchainapi long`,
	Run: func(cmd *cobra.Command, args []string) {
		RunServer()
	},
}

// Execute runs the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func RunServer() {
	addr := viper.GetString("bind-address")

	r := mux.NewRouter()

	// 设置路由，如果访问/，则调用index方法
	r.HandleFunc("/", handler.Index)
	r.HandleFunc("/query/pledge/{validator_address}", handler.QueryPledge)
	r.HandleFunc("/query/commission/{validator_address}", handler.QueryCommission)
	r.HandleFunc("/query/outstanding_rewards/{validator_address}", handler.QueryOutstandingRewards)
	r.HandleFunc("/query/node_info", handler.QueryNodeInfo)
	r.HandleFunc("/query/delegators/{delegator_address}/rewards/{validator_address}", handler.QueryUserRewards)
	r.HandleFunc("/send/from/{from}/to/{to}/amount/{amount}/gas/{gas}/priv/{priv}", handler.SendTx)
	//r.HandleFunc("/send/to/{to}/amount/{amount}/gas/{gas}/user/{user}/passwd/{passwd}", test.Send)
	r.HandleFunc("/watch", handler.Watch)

	log.Printf("Start pluschainapi  at http://%s/", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Printf("Error occur when start server %v", err)
		os.Exit(1)
	}
}
