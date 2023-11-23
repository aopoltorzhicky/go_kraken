package futures

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type clientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// KrakenFutures - object wraps Kraken Futures API
type KrakenFutures struct {
	key    string
	secret string
	client clientInterface
}

// New - constructor of KrakenFutures object
func New(key string, secret string) *KrakenFutures {
	if key == "" || secret == "" {
		log.Print("[WARNING] You are not set api key and secret for Kraken Futures!")
	}
	return &KrakenFutures{
		key:    key,
		secret: secret,
		client: http.DefaultClient,
	}
}

func (api *KrakenFutures) getSign(requestURL string, data url.Values) (string, error) {
	sha := sha256.New()

	if _, err := sha.Write([]byte(data.Get("nonce") + data.Encode())); err != nil {
		return "", err
	}
	hashData := sha.Sum(nil)
	s, err := base64.StdEncoding.DecodeString(api.secret)
	if err != nil {
		return "", err
	}
	hmacObj := hmac.New(sha512.New, s)

	if _, err := hmacObj.Write(append([]byte(requestURL), hashData...)); err != nil {
		return "", err
	}
	hmacData := hmacObj.Sum(nil)
	return base64.StdEncoding.EncodeToString(hmacData), nil
}

func (api *KrakenFutures) prepareRequest(reqType string, method string, isPrivate bool, data url.Values) (*http.Request, error) {
	if data == nil {
		data = url.Values{}
	}
	requestURL := ""
	if isPrivate {
		requestURL = fmt.Sprintf("%s/%s", FuturesAPIUrl, method)
		data.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))
	} else {
		requestURL = fmt.Sprintf("%s/%s", FuturesAPIUrl, method)
	}

	var req *http.Request
	var err error
	if reqType == "GET" {
		// For GET requests, add data as query parameters
		requestURL = fmt.Sprintf("%s?%s", requestURL, data.Encode())
		req, err = http.NewRequest(reqType, requestURL, nil)
	} else {
		// For POST and other types, add data in the body
		req, err = http.NewRequest(reqType, requestURL, strings.NewReader(data.Encode()))
	}
	if err != nil {
		return nil, errors.Wrap(err, "error during request creation")
	}

	if isPrivate {
		urlPath := fmt.Sprintf("/%s", method)
		req.Header.Add("API-Key", api.key)
		signature, err := api.getSign(urlPath, data)
		if err != nil {
			return nil, errors.Wrap(err, "invalid secret key")
		}
		req.Header.Add("API-Sign", signature)
	}
	return req, nil
}

func (api *KrakenFutures) parseResponse(response *http.Response, retType interface{}) error {
	if response.StatusCode != 200 {
		return errors.Errorf("error during response parsing: invalid status code %d", response.StatusCode)
	}

	if response.Body == nil {
		return errors.New("error during response parsing: can not read response body")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "error during response parsing: can not read response body")
	}

	if err = json.Unmarshal(body, &retType); err != nil {
		return errors.Wrap(err, "error during response parsing: json marshalling")
	}

	return nil
}

func (api *KrakenFutures) request(reqType string, method string, isPrivate bool, data url.Values, retType interface{}) error {
	req, err := api.prepareRequest(reqType, method, isPrivate, data)
	if err != nil {
		return err
	}
	resp, err := api.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "error during request execution")
	}
	defer resp.Body.Close()
	return api.parseResponse(resp, retType)
}
