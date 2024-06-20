package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"gobackend/types"
	"log"
	"net/http"
	"sync"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func EncodeInput(input types.CodeExecutionInputBody) ([]byte, error) {
	var inputBuffer bytes.Buffer
	if err := gob.NewEncoder(&inputBuffer).Encode(input); err != nil {
		return nil, err
	}

	return inputBuffer.Bytes(), nil
}

func DecodeInput(buffer []byte) (*types.CodeExecutionInputBody, error) {
	var inputBuffer bytes.Buffer = *bytes.NewBuffer(buffer)
	codeExecutionInputBody := new(types.CodeExecutionInputBody)
	if err := gob.NewDecoder(&inputBuffer).Decode(codeExecutionInputBody); err != nil {
		return nil, err
	}

	return codeExecutionInputBody, nil
}

func WriteJSON(w http.ResponseWriter, status int, input any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(input); err != nil {
		return err
	}

	return nil
}

// make a global map
var m map[string][2]string = make(map[string][2]string)

var mutex sync.RWMutex

func mapSetter(lang string, runtime string, fileName string) {
	mutex.Lock()
	m[lang] = [2]string{runtime, fileName}
	mutex.Unlock()
}

func mapGetter(lang string) [2]string {
	mutex.RLock()
	result := m[lang]
	mutex.RUnlock()

	return result
}

func initializeMap() {
mapSetter("javascript", "node", "code.js")
mapSetter("python", "python", "code.py")
mapSetter("cpp", "gcc", "code.cpp")
mapSetter("c", "gcc", "code.c")
mapSetter("rust", "rustc", "code.rs")

}
func GetRuntimeFromLang(lang string) string {
initializeMap()
	
result := mapGetter(lang)
return result[0]
}


func GetFilenameFromLang(lang string) string {
	initializeMap()

	result := mapGetter(lang)
	return result[1]
}