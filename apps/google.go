// Copyright 2014 beego authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
//
// Maintain by https://github.com/slene

package apps

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/httplib"

	"github.com/alexandru-porcescu/social-auth"
)

type Google struct {
	BaseProvider
}

func (p *Google) GetType() social.SocialType {
	return social.SocialGoogle
}

func (p *Google) GetName() string {
	return "Google"
}

func (p *Google) GetPath() string {
	return "google"
}

func (p *Google) GetSocialData(tok *social.Token) (*social.SocialData, error) {

	uri := "https://www.googleapis.com/userinfo/v2/me"
	req := httplib.Get(uri)
	req.SetTransport(social.DefaultTransport)
	req.Header("Authorization", "Bearer "+tok.AccessToken)

	resp, err := req.Response()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()

	type email struct {
		Value   string     `json:"value"`
		Type   string      `json:"type"`
	}

	type response struct {
		Id   string        `json:"id"`
		Name   string      `json:"name"`
		NickName   string  `json:"nickname"`
		Emails   []email   `json:"emails"`
		Error   interface{} `json:"error"`
	}

	var result response

	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}

	if result.Error != nil {
		return nil, fmt.Errorf("%v", result.Error)
	}

	sData := &social.SocialData{
		Id: result.Id,
		Name: result.Name,
		NickName: result.NickName,
		Email: result.Emails[0].Value,
	}

	return sData, nil
}

var _ social.Provider = new(Google)

func NewGoogle(clientId, secret string) *Google {
	p := new(Google)
	p.App = p
	p.ClientId = clientId
	p.ClientSecret = secret
	p.Scope = "email profile https://www.googleapis.com/auth/plus.login"
	p.AuthURL = "https://accounts.google.com/o/oauth2/auth"
	p.TokenURL = "https://accounts.google.com/o/oauth2/token"
	p.RedirectURL = social.DefaultAppUrl + "login/google/access"
	p.AccessType = "offline"
	p.ApprovalPrompt = "auto"
	return p
}
