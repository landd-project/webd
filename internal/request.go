package internal

import (
	"fmt"
	"strings"
)

/*
	FETCH url
	TRUST url
	UNTRUST url

	TAB ALL
	TAB GET 1
	TAB DEL 1
*/

func ParseRequest(req string) error {
	parts := strings.Split(req, " ");

	if len(parts) < 2 {
		return fmt.Errorf("invalid request: `%v`", req);
	}

	url := strings.TrimSpace(parts[1]);
	command := parts[0];

	switch command {
	case "FETCH":
		err := RequestPage(url);
		if err != nil {
			return err;
		}
	case "TRUST":
	case "UNTRUST":
	case "TAB":
	default:
		return fmt.Errorf("invalid command on request: `%v`", command);
	}

	return nil;
}
