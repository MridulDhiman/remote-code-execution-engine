package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"gobackend/types"
	"log"
	"net/http"
	"regexp"
	"golang.org/x/time/rate"
)

var (
	allowedLanguages = map[string]bool{
		"javascript" : true,
		"python": true,
		"cpp": true,
		"c" : true,
		"rust" : true,
	}

	JsConfig = types.NewLangConfig("javascript", false, "node", "code.js")
	PyConfig = types.NewLangConfig("python", false, "python", "code.py")
	CppConfig = types.NewLangConfig("cpp", true, "g++", "code.cpp")
	CConfig = types.NewLangConfig("c", true, "gcc", "code.c")
	RustConfig = types.NewLangConfig("rust", true, "rustc", "code.rs")

	langConfig = []*types.LangConfig{JsConfig, PyConfig, CppConfig, CConfig, RustConfig}
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


func IsLanguageSupported (lang string)  bool {
	return !allowedLanguages[lang]
}

func GetLangConfig (lang string) *types.LangConfig {
	for _, config := range langConfig {
		if config.Lang == lang {
			return config
		}
	}

	return &types.LangConfig{}
}

func HasSuspiciousPatterns (code string) bool {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)SELECT|UPDATE|DELETE|CREATE`), // database operations
		regexp.MustCompile(`(?i)fork|exec|system|eval\(`), // system calls and code evaluation
		regexp.MustCompile(`(?i)socket|curl\(`), // network operations
	}

	for _, pattern := range patterns {
		if pattern.MatchString(code) {
			return true
		}
	}

	return false
}

func RateLimiter (next func(w http.ResponseWriter, r* http.Request)) http.HandlerFunc {
	limiter:= rate.NewLimiter(2, 5) // avg. 2 req. per sec, and allowed burst of 5 requests

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
			 WriteJSON(w, http.StatusTooManyRequests, types.ErrorResponse{Message: "The API is at capacity. Try again later"})
			 return 
        } else {
            next(w, r)
        }
    })
}