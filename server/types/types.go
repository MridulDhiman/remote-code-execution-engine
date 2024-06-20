
package types

type Executor struct {
  ID string `json:"id"`
  Code string `json:"code"`
  Input string `json:"input"`
  Output string `json:"output"`
  Status string `json:"status"`
  Lang string `json:"lang"`
}

type CodeExecutionInputBody struct {
	Code string `json:"code"`
	Input string `json:"input"`
  Lang string `json:"lang"`
}


type AckResponse struct {
  Status string `json:"status"`
}


type ErrorResponse struct {
  Message string `json:"message"`
}

