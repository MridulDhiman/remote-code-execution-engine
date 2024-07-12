package main

import (
	"encoding/json"
	"fmt"
	"gobackend/rabbitmq"
	"gobackend/types"
	"gobackend/utils"
	"gobackend/ws"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{} 
)


type APIServer struct {
	listenAddr string
	mqConn *rabbitmq.MQConnection
}

// initialize server
func NewAPIServer (listenAddr string, mqConn *rabbitmq.MQConnection) *APIServer {
	return &APIServer{listenAddr: listenAddr, mqConn: mqConn}
}

type apiFunc func (http.ResponseWriter, *http.Request)  error

func makeHTTPHandlerFunc (f apiFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r* http.Request)  {
		// idhar websocket setup kr denge

              if  err:= f(w,r); err!= nil {
				log.Fatal(err)
			  }
	}
}

func (s* APIServer) Run() {
   router:= mux.NewRouter()
   router.HandleFunc("/code", utils.RateLimiter(makeHTTPHandlerFunc(s.handleCode)))
   router.HandleFunc("/echo", utils.RateLimiter(makeHTTPHandlerFunc(s.WsServer)))
   fmt.Printf("Listening on port %v", s.listenAddr);
   http.ListenAndServe(s.listenAddr, router)
}

func (s* APIServer) GetAPIServer() *APIServer {
return s
}

func (s* APIServer) WsServer(w http.ResponseWriter, r* http.Request) error {
   wsConn, err := upgrader.Upgrade(w,r, nil)
   if err != nil {
	return err
   }

   defer wsConn.Close()

   for {
	mt, message, err:= wsConn.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		break
	}
	log.Printf("recv: %s", message)
		err = wsConn.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}

   }
	return nil
}

func (s* APIServer) handleCode (w http.ResponseWriter, r* http.Request) error {
	if(r.Method == "POST") {
   return s.handleCreateCode(w,r)
	}
return nil
}

func (s* APIServer) handleCreateCode (w http.ResponseWriter, r* http.Request) error {
	codeExecutionInputBody:= new(types.CodeExecutionInputBody)
	if err:= json.NewDecoder(r.Body).Decode(codeExecutionInputBody); err != nil {
		return err
	}

	lang:= codeExecutionInputBody.Lang
	code:= codeExecutionInputBody.Code
	if utils.IsLanguageSupported(lang) {
		return utils.WriteJSON(w, http.StatusUnprocessableEntity, types.ErrorResponse{Message: "Language Not Supported"})
	}

	if utils.HasSuspiciousPatterns(code) {
		return utils.WriteJSON(w, http.StatusUnprocessableEntity, types.ErrorResponse{Message: "Potentially malicious code detected..."})
	}

   // add to queue => message queue
   go s.mqConn.AddToQueue(*codeExecutionInputBody)
   go ws.Ws.SendMsg("Initiated");
return utils.WriteJSON(w, http.StatusOK, types.AckResponse{ Status: "pending"})
}