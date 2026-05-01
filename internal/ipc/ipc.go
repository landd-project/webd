package ipc

import (
	"fmt"
	"strconv"
	"strings"

	"webd/internal/gemini"
	"webd/internal/tabs"
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
	case "fetch":
		url := strings.TrimSpace(parts[1]);
		response, err := gemini.RequestPage(url, 0);
		if err != nil {
			r.Error = err.Error();
			return r;
		}
		r.Data = response;

	case "parse":
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
	case "tab":
		subcommand := strings.TrimSpace(parts[1]);
		switch subcommand {
		case "sel":
			if len(parts) < 3 {
				r.Error = fmt.Sprintf("invalid number of parameters in request for select a tab, expected: 3 or more but founds: %v", len(parts));
			}
			id, err := strconv.Atoi(parts[2]);
			if err != nil {
				r.Error = err.Error();
				return r;
			}
			list := tabs.All();

			if id >= len(list) {
				r.Error = "this is not a valid id, it is major than the len of the tab list";
				return r;
			}
			tab := list[id];
			err = tabs.SetCurrentTab(tab);
			if err != nil {
				r.Error = err.Error();
				return r;
			}
			r.Data = id

		case "del":
		case "all":
			list := tabs.All();
			r.Data = list;
		case "new":
			if len(parts) < 3 {
				r.Error = fmt.Sprintf("invalid number of parameters in request for a new tab, expected: 3 or more but founds: %v", len(parts));
			}
			url := strings.TrimSpace(parts[2]);
			err := tabs.NewTab(url);
			if err != nil {
				r.Error = err.Error();
				return r;
			}
			r.Data = url;
		case "get":
			current := tabs.GetCurrentTab();
			r.Data = current;
		}
		
	default:
		r.Error = fmt.Sprintf("invalid command on request: `%v`", command);
		return r;
	}

	r.Ok = true;
	return r;
}

