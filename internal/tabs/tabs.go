package tabs

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"webd/internal/gemini"
)

type Tab struct {
	Id string
	Title string
	Url string
	Content any
}

var CurrentTab Tab = Tab{};

var tabs = make(map[string]Tab);

func generateTabId() (string, error){
	bytes := make([]byte, 16);
	if _, err := rand.Read(bytes); err != nil {
		return "", err;
	}
	return hex.EncodeToString(bytes), nil;
}

func MakeTab(url string) (Tab, error) {
	var t Tab;
	t.Title = url;
	t.Url = url;
	id, err := generateTabId();
	if err != nil {
		return t, err;
	}
	t.Id = id;

	resp, err := gemini.RequestPage(url, 0);
	if err != nil {
		return t, err;
	}
	t.Content = resp.Body;

	if resp.Meta == "text/gemini" {
		tokens, err := gemini.ParseGemtext(resp.Body);
		if err != nil {
			return t, err;
		}
		t.Content = tokens;
		if tokens[0].Type == gemini.TokenHeading1 {
			t.Title = tokens[0].Value.(string);
		}
	}

	return t, nil;
}
func NewTab(url string) (Tab, error) {
	tab, err := MakeTab(url);
	if err != nil {
		return tab, err;
	}

	err = SetCurrentTab(tab);
	if err != nil {
		return tab, err;
	}
	tabs[tab.Id] = tab;

	return tab, nil;
}

func PutTab(url string) (Tab, error) {
	currentTab := GetCurrentTab();

	tab, err := MakeTab(url);
	if err != nil {
		return tab, err;
	}

	tabs[currentTab.Id] = tab;

	err = SetCurrentTab(tab);
	if err != nil {
		return tab, err;
	}
	return tab, nil;
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

func All() map[string]Tab {
	return tabs;
}

func Delete(id string) error {
	delete(tabs, id);

	return nil;
}


