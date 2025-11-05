package httputil



type Response struct{
	Error string `json:"error,omitempty"`	
	
}


func Error(msg string) Response {
	return Response{
		Error: msg,
	}
}
 