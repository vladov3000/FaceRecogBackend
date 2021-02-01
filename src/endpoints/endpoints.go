package endpoints

import "net/http"

type ReqHandler = func(w http.ResponseWriter, r *http.Request)
