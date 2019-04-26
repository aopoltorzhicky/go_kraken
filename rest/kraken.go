package rest

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ugorji/go/codec"
)

// clientInterface - for testing purpose
type clientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// Kraken - object wraps API
type Kraken struct {
	key    string
	secret string
	client clientInterface
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

func (api *Kraken) getSign(requestURL string, data url.Values) (string, error) {
	sha := sha256.New()
	sha.Write([]byte(data.Get("nonce") + data.Encode()))
	hashData := sha.Sum(nil)
	s, err := base64.StdEncoding.DecodeString(api.secret)
	if err != nil {
		return "", err
	}
	hmacObj := hmac.New(sha512.New, s)
	hmacObj.Write(append([]byte(requestURL), hashData...))
	hmacData := hmacObj.Sum(nil)
	return base64.StdEncoding.EncodeToString(hmacData), nil
}

func (api *Kraken) prepareRequest(method string, isPrivate bool, data url.Values) (*http.Request, error) {
	if data == nil {
		data = url.Values{}
	}
	requestURL := ""
	if isPrivate {
		requestURL = fmt.Sprintf("%s/%s/private/%s", APIUrl, APIVersion, method)
		data.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))
	} else {
		requestURL = fmt.Sprintf("%s/%s/public/%s", APIUrl, APIVersion, method)
	}
	req, err := http.NewRequest("POST", requestURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("Error during request creation: %s", err.Error())
	}

	if isPrivate {
		urlPath := fmt.Sprintf("/%s/private/%s", APIVersion, method)
		req.Header.Add("API-Key", api.key)
		signature, err := api.getSign(urlPath, data)
		if err != nil {
			return nil, fmt.Errorf("Invalid secret key: %s", err.Error())
		}
		req.Header.Add("API-Sign", signature)
	}
	return req, nil
}

func (api *Kraken) parseResponse(response *http.Response, retType interface{}) (interface{}, error) {
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Error during response parsing: invalid status code %d", response.StatusCode)
	}

	if response.Body == nil {
		return nil, fmt.Errorf("Error during response parsing: can not read response body")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error during response parsing: can not read response body (%s)", err.Error())
	}
	var retData KrakenResponse
	if retType != nil {
		retData.Result = retType
	}

	log.Println(string(body))
	err = codec.NewDecoderBytes(body, new(codec.JsonHandle)).Decode(&retData)
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
