package main

import (
	"encoding/json"
	"fmt"
	"gobackend/rabbitmq"
	"gobackend/types"
	"gobackend/utils"
	"log"
	"net/http"
	"github.com/gorilla/mux"
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
              if  err:= f(w,r); err!= nil {
				log.Fatal(err)
			  }
	}
}

func (s* APIServer) Run() {
   router:= mux.NewRouter()
   router.HandleFunc("/code", makeHTTPHandlerFunc(s.handleCode))
   fmt.Printf("Listening on port %v", s.listenAddr);
   http.ListenAndServe(s.listenAddr, router)
}

func (s* APIServer) GetAPIServer() *APIServer {
return s
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
	if lang != "javascript" && lang != "python" && lang != "c++" && lang != "c" {
		return utils.WriteJSON(w, http.StatusUnprocessableEntity, types.ErrorResponse{Message: "Language Not Supported"})
	}

   // add to queue => message queue
   go s.mqConn.AddToQueue(*codeExecutionInputBody)
return utils.WriteJSON(w, http.StatusOK, types.AckResponse{ Status: "pending"})
}