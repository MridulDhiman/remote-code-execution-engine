package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"gobackend/types"
	"log"
	"net/http"
)

func FailOnError(err error, msg string) {
	if err != nil {
	  log.Panicf("%s: %s", msg, err)
	}
  }


  func EncodeInput (input types.CodeExecutionInputBody) ([]byte, error) {
var inputBuffer bytes.Buffer
if err:= gob.NewEncoder(&inputBuffer).Encode(input); err != nil {
	return nil, err
}


return inputBuffer.Bytes(), nil
  }

  func DecodeInput (buffer []byte) (*types.CodeExecutionInputBody, error) {
    var inputBuffer bytes.Buffer = *bytes.NewBuffer(buffer)
	   codeExecutionInputBody := new(types.CodeExecutionInputBody)
	   if err := gob.NewDecoder(&inputBuffer).Decode(codeExecutionInputBody); err!= nil {
         return nil, err
	   }

return codeExecutionInputBody, nil
  }


  func WriteJSON (w http.ResponseWriter, status int, input any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	if err:= json.NewEncoder(w).Encode(input); err != nil {
		return err
	}

	return nil
  }