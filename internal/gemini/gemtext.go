package gemini

import (
	"fmt"
	"strings"
)

type TokenType int;
const (
	TokenText TokenType = iota
	TokenHeading1
	TokenHeading2
	TokenHeading3
	TokenLink
	TokenPreformat
	TokenList
	TokenQuote
)

type Token struct {
	Type TokenType
	Value any
}

type Link struct {
	Text string
	Url string
}

func ParsePage(content string) ([]Token, error) {

	lines := strings.Split(content, "\n");

	var tokens []Token;
	var tok Token;

	preformattedMode := false;

	for i,line := range lines {
		if preformattedMode {
			tok.Type = TokenPreformat;
			tok.Value = line;
		}

		if strings.HasPrefix(line, "###") {
			tok.Type = TokenHeading3
			tok.Value = line[3:];

		} else if strings.HasPrefix(line, "##") {
			tok.Type = TokenHeading2
			tok.Value = line[2:];

		} else if strings.HasPrefix(line, "#") {
			tok.Type = TokenHeading1
			tok.Value = line[1:];

		} else if strings.HasPrefix(line, "* ") {
			tok.Type = TokenList
			tok.Value = line[2:];

		} else if strings.HasPrefix(line, ">") {
			tok.Type = TokenQuote
			tok.Value = line[1:];
		} else if strings.HasPrefix(line, "=> ") {
			tok.Type = TokenLink;
			var l Link;

			parts := strings.Split(line, " ");
			if len(parts) < 2{
				return nil, fmt.Errorf("invalid link line: %v", i);
			}
			l.Url = strings.TrimSpace(parts[1]);

			if len(parts) == 2 {
				l.Text = l.Url;
			} else {
				l.Text = strings.TrimSpace(parts[2]);
			}

			tok.Value = l;

		} else if strings.HasPrefix(line, "```") {
			preformattedMode = !preformattedMode;
		} else {
			tok.Type = TokenText;
			tok.Value = line;
		}
		

		tokens = append(tokens, tok);
	}

	return tokens, nil;
}
