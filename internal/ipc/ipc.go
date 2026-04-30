package ipc

import (
	"fmt"
	"strings"

	"webd/internal/gemini"
)

type Response struct {
	Ok bool `json:"ok"`
	Error string `json:"error"`
	Data any `json:"data"`
}

func ParseRequest(req string) Response {
	var r = Response{
		Ok: false,
		Error: "",
		Data: nil,
	};

	parts := strings.Split(req, " ");
	if len(parts) < 2 {
		r.Error = fmt.Sprintf("invalid number of parameters in request, expected: 2 or more but founds: %v", len(parts));
		return r;
	}

	command := parts[0];

	switch command {
	case "FETCH":
		url := strings.TrimSpace(parts[1]);
		response, err := gemini.RequestPage(url, 0);
		if err != nil {
			r.Error = err.Error();
			return r;
		}
		r.Data = response;

	case "PARSE":
		url := strings.TrimSpace(parts[1]);
		response, err := gemini.RequestPage(url, 0);
		if err != nil {
			r.Error = err.Error();
			return r;
		}
		tokens, err := gemini.ParseGemtext(response.Body);
		if err != nil {
			r.Error = err.Error();
			return r;
		}

		r.Data = gemini.GeminiParsedResponse{
			StatusCode: response.StatusCode,
			Meta: response.Meta,
			Body: response.Body,
			Tokens: tokens,
			RedirectCount: response.RedirectCount,
		};
		
	default:
		r.Error = fmt.Sprintf("invalid command on request: `%v`", command);
		return r;
	}

	r.Ok = true;
	return r;
}

