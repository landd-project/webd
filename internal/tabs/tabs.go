package tabs

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
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
	CurrentTabId int;
	Tabs map[string]Tab;
}

var tabData TabData;

var tabsOrder []string;

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

	err = SetCurrentTab(tabData.CurrentTabId);
	if err != nil {
		return tab, err;
	}
	tabsOrder = append(tabsOrder, tab.Id);
	tabData.Tabs[tab.Id] = tab;

	return tab, nil;
}

func PutTab(url string) (Tab, error) {
	tab, err := MakeTab(url);
	if err != nil {
		return tab, err;
	}
	currentTab := GetCurrentTab();
	tabData.Tabs[currentTab.Id] = tab;

	err = SetCurrentTab(tabData.CurrentTabId);
	if err != nil {
		return tab, err;
	}
	return tab, nil;
}

func GetCurrentTab() Tab {
	id := tabsOrder[tabData.CurrentTabId];
	return tabData.Tabs[id];
}

func SetCurrentTab(id int) error {
	tabData.CurrentTabId = id;
	return nil;
}

func All() []Tab {
	var t []Tab;
	for _,v := range tabsOrder {
		t = append(t, tabData.Tabs[v]);
	}
	return t;
}

func Delete(id int) error {
	last := len(tabsOrder)-1;
	tabsOrder = append(tabsOrder[:last], tabsOrder[last+1:]...);
	delete(tabData.Tabs, tabsOrder[id]);

	return nil;
}

func SaveTabs() error {
	bt, err := json.Marshal(&tabData);
	if err != nil {
		return err;
	}

	// TODO: fix path
	err = os.WriteFile("tabs.json", bt, 0644);
	if err != nil {
		return err;
	}
	return nil;
}

func LoadTabs() {

}

