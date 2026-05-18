package tabs

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"webd/internal/gemini"
)

type Tab struct {
	Id string
	Title string
	Url string
	Content any
}

type TabData struct {
	CurrentTabId int
	Data map[string]Tab
	Order []string
}

var tabs = TabData{
	Data: make(map[string]Tab),
};

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

	tabs.Order = append(tabs.Order, tab.Id);
	tabs.Data[tab.Id] = tab;
	tabs.CurrentTabId = len(tabs.Order) -1;

	if err = SaveTabs(); err != nil {
		return tab, err;
	}
	return tab, nil;
}

func PutTab(url string) (Tab, error) {
	tab, err := MakeTab(url);
	if err != nil {
		return tab, err;
	}
	currentTabHashId := tabs.Order[tabs.CurrentTabId];
	tabs.Data[currentTabHashId] = tab;
	tabs.CurrentTabId = len(tabs.Order)-1;

	if err = SaveTabs(); err != nil {
		return tab, err;
	}
	return tab, nil;
}

func All() ([]Tab, error) {
	if err := LoadTabs(); err != nil {
		return nil, err;
	}
	var t []Tab;
	for _,v := range tabs.Order {
		t = append(t, tabs.Data[v]);
	}
	return t, nil;
}

func Delete(id int) error {
	if id < 0 || id >= len(tabs.Order) {
		return fmt.Errorf("index out of range");
	}

	hashId := tabs.Order[id];

	tabs.Order = append(tabs.Order[:id], tabs.Order[id+1:]...);
	delete(tabs.Data, hashId);

	switch {
	case len(tabs.Order) == 0:
		tabs.CurrentTabId = 0;
	case tabs.CurrentTabId == id:
		if id > 0 {
			tabs.CurrentTabId = id-1;
		} else {
			tabs.CurrentTabId = 0;
		}
	case tabs.CurrentTabId > id:
		tabs.CurrentTabId--;
	}

	if err := SaveTabs(); err != nil {
		return err;
	}

	return nil;
}

func Select(id int) (Tab, error) {
	if id >= len(tabs.Order) {
		return Tab{}, errors.New("invalid id, it's out of the length of the tabs list");
	}
	tabs.CurrentTabId = id;

	tabId := tabs.Order[id];
	return tabs.Data[tabId], nil;
}

func Get() (Tab, error) {
	if tabs.CurrentTabId == 0 {
		return Tab{}, errors.New("tab list is empty");
	}	

	currentTabHashId := tabs.Order[tabs.CurrentTabId];
	tab := tabs.Data[currentTabHashId];
	return tab, nil;
}

func SaveTabs() error {
	bt, err := json.MarshalIndent(&tabs, "", "	");
	if err != nil {
		return err;
	}
	err = os.WriteFile("./tabs.json", bt, 0644);
	if err != nil {
		return err;
	}
	return nil;
}

func LoadTabs() error {
	bt, err := os.ReadFile("./tabs.json");
	if err != nil {
		return err;
	}
	err = json.Unmarshal(bt, &tabs);
	if err != nil {
		return err
	}
	return nil;
}

