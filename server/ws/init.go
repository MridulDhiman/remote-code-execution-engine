package ws

import "gobackend/utils"

var Ws = NewWsClient("localhost:8080")

func init() {
	if err := Ws.InitWsConnection(); err != nil {
		utils.FailOnError(err, "Could not init ws client connection")
	}
}