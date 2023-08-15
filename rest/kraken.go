package rest

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
		return nil, errors.Wrap(err, "error during request creation")
	}

	if isPrivate {
		urlPath := fmt.Sprintf("/%s/private/%s", APIVersion, method)
		req.Header.Add("API-Key", api.key)
		signature, err := api.getSign(urlPath, data)
		if err != nil {
			return nil, errors.Wrap(err, "invalid secret key")
		}
		req.Header.Add("API-Sign", signature)
	}
	return req, nil
}

func (api *Kraken) parseResponse(response *http.Response, retType interface{}) error {
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

	// log.Println(string(body))
	var retData KrakenResponse
	if retType != nil {
		retData.Result = retType
	}

	if err = json.Unmarshal(body, &retData); err != nil {
		return errors.Wrap(err, "error during response parsing: json marshalling")
	}

	if len(retData.Error) > 0 {
		return errors.Errorf("kraken return errors: %s", retData.Error)
	}

	return nil
}

func (api *Kraken) request(method string, isPrivate bool, data url.Values, retType interface{}) error {
	req, err := api.prepareRequest(method, isPrivate, data)
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
