package tabs

import (
	"fmt"
	"webd/internal/gemini"
)

type Tab struct {
	Title string
	Url string
	Content any
}

var CurrentTab Tab = Tab{};
var tabs []Tab;

func NewTab(url string) error {
	var t Tab;
	t.Url = url;
	t.Title = url;

	resp, err := gemini.RequestPage(url, 0);
	if err != nil {
		return err;
	}
	t.Content = resp.Body;

	if resp.Meta == "text/gemini" {
		tokens, err := gemini.ParseGemtext(resp.Body);
		if err != nil {
			return err;
		}
		t.Content = tokens;
		if tokens[0].Type == gemini.TokenHeading1 {
			t.Title = tokens[0].Value.(string);
		}
	}

	err = SetCurrentTab(t);
	if err != nil {
		return err;
	}

	tabs = append(tabs, t);

	return nil;
}

func GetCurrentTab() Tab {
	return CurrentTab;
}

func SetCurrentTab(tab Tab) error {
	empty := Tab{};
	if tab == empty {
		return fmt.Errorf("could set current tab to an empty tab");
	}

	CurrentTab = tab;	
	return nil;
}

func All() []Tab {
	return tabs;
}
