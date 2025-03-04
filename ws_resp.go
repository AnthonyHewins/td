package td

import "fmt"

// General response message that usually denotes errors
// in requests
type WSResp struct {
	Code WSRespCode `json:"code"`
	Msg  string     `json:"msg"`
}

func (w *WSResp) Error() string {
	return fmt.Sprintf("%d: %s", w.Code, w.Msg)
}
