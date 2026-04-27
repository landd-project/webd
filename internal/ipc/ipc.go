package ipc

import (
	"fmt"
	"strings"

	"webd/internal/gemini"
)

/*
	FETCH url

	TAB ALL
	TAB GET 1
	TAB DEL 1
*/

type Response struct {
	Ok bool `json:"ok"`
	Error string `json:"error"`
	Data any `json:"data"`
}

func ParseRequest(req string) Response {

	var r = Response{
		Ok: true,
		Error: "",
		Data: nil,
	};

	parts := strings.Split(req, " ");
	if len(parts) < 2 {
		r.Ok = false;
		r.Error = fmt.Sprintf("invalid number of parameters in request, expected: 2 or more but founds: %v", len(parts));
		return r;
	}

	url := strings.TrimSpace(parts[1]);
	command := parts[0];

	switch command {
	case "FETCH":
		response, err := gemini.RequestPage(url);
		if err != nil {
			r.Ok = false;
			r.Error = err.Error();
			return r;
		}

		g, err := gemini.ParseResponse(response);
		if err != nil {
			r.Ok = false;
			r.Error = err.Error();
			return r;
		}
		r.Data = g;
		
	case "TAB":
		r.Ok = false;
		r.Error = "TODO: not implemented";
		return r;
	default:
		r.Ok = false;
		r.Error = fmt.Sprintf("invalid command on request: `%v`", command);
		return r;
	}

	return r;
}

