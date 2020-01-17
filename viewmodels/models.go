package viewmodels

import (
	"encoding/json"
)

//ListItem is used to allow get the value form common query or procedure which extract id , title etc.
type ListItem struct {
	ID       int    `json:",omitempty"`
	BigID    int64  `json:",omitempty"`
	IDStr    string `json:",omitempty"`
	DataInt  int    `json:",omitempty"`
	DataStr  string `json:",omitempty"`
	DataStr2 string `json:",omitempty"`
	Checked  bool   `json:",omitempty"`
}

//AuthInfo is used to take input from user
type AuthInfo struct {
	ClientID      string `json:",omitempty"`
	ClientSecret  string `json:",omitempty"`
	Scopes        string `json:",omitempty"`
	AppUserID     string `json:",omitempty"`
	RefereshToken string `json:",omitempty"`
}

//ToJSON marsals instance of AuthInfo to JSON
func (a *AuthInfo) ToJSON() []byte {
	data, err := json.Marshal(a)
	if err != nil {
		return nil
	}
	return data
}
