package adapters

import (
	"better-admin-backend-service/config"
	"better-admin-backend-service/dtos"
	"encoding/json"
	"fmt"
	"github.com/bettercode-oss/rest"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type GoogleOAuthAdapter struct {
}

func (adapter GoogleOAuthAdapter) Authenticate(code string, setting dtos.GoogleWorkspaceLoginSetting) (dtos.GoogleMember, error) {
	accessToken, err := adapter.getAccessToken(code, setting)
	if err != nil {
		return dtos.GoogleMember{}, err
	}

	client := rest.Client{}
	googleMember := dtos.GoogleMember{}
	err = client.
		Request().
		SetResult(&googleMember).
		Get(fmt.Sprintf("%v?access_token=%v", config.Config.GoogleOAuth.AuthUri, accessToken))

	if err != nil {
		return googleMember, errors.Wrap(err, "google authenticate error")
	}

	return googleMember, nil
}

func (GoogleOAuthAdapter) getAccessToken(code string, setting dtos.GoogleWorkspaceLoginSetting) (string, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", setting.ClientId)
	data.Set("client_secret", setting.ClientSecret)
	data.Set("redirect_uri", setting.RedirectUri)
	data.Set("grant_type", "authorization_code")

	client := &http.Client{}
	r, err := http.NewRequest("POST", config.Config.GoogleOAuth.TokenUri, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		return "", errors.Wrap(err, "google oauth error")
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil {
		return "", errors.Wrap(err, "google oauth error")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "google oauth error")
	}

	responseBody := map[string]interface{}{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return "", errors.Wrap(err, "google oauth error")
	}

	return responseBody["access_token"].(string), nil
}
