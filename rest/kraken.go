package rest

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Kraken - object wraps API
type Kraken struct {
	key    string
	secret string
	client *http.Client
}

// New - constructor of Kraken object
func New(key string, secret string) *Kraken {
	if key == "" || secret == "" {
		log.Print("[WARNING] You are not set api key and secret!")
	}
	return &Kraken{
		key:    key,
		secret: secret,
		client: http.DefaultClient,
	}
}

func (api *Kraken) getSign(requestURL string, data url.Values) string {
	sha := sha256.New()
	sha.Write([]byte(data.Get("nonce") + data.Encode()))
	hashData := sha.Sum(nil)
	hmacObj := hmac.New(sha512.New, []byte(api.secret))
	hmacObj.Write(append([]byte(requestURL), hashData...))
	hmacData := hmacObj.Sum(nil)
	return base64.StdEncoding.EncodeToString(hmacData)
}

func (api *Kraken) prepareRequest(method string, isPrivate bool, data url.Values) (*http.Request, error) {
	requestURL := ""
	if isPrivate {
		requestURL = fmt.Sprintf("%s/%s/private/%s", APIUrl, APIVersion, method)
	} else {
		requestURL = fmt.Sprintf("%s/%s/public/%s", APIUrl, APIVersion, method)
	}
	req, err := http.NewRequest("POST", requestURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("Error during request creation: %s", err.Error())
	}

	if isPrivate {
		data.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))
		urlPath := fmt.Sprintf("/%s/private/%s", APIVersion, method)
		req.Header.Add("API-Key", api.key)
		req.Header.Add("API-Sign", api.getSign(urlPath, data))
	}
	return req, nil
}

func (api *Kraken) parseResponse(response *http.Response, retType interface{}) (interface{}, error) {
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Error during response parsing: invalid status code %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error during response parsing: can not read response body (%s)", err.Error())
	}
	var retData KrakenResponse
	if retType != nil {
		retData.Result = retType
	}

	err = json.Unmarshal(body, &retData)
	if err != nil {
		return nil, fmt.Errorf("Error during response parsing: json marshalling (%s)", err.Error())
	}

	if len(retData.Error) > 0 {
		return nil, fmt.Errorf("Kraken return errors: %s", retData.Error)
	}

	return retData.Result, nil
}

func (api *Kraken) request(method string, isPrivate bool, data url.Values, retType interface{}) (interface{}, error) {
	req, err := api.prepareRequest(method, isPrivate, data)
	if err != nil {
		return nil, err
	}
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error during request execution: %s", err.Error())
	}
	defer resp.Body.Close()
	return api.parseResponse(resp, retType)
}
