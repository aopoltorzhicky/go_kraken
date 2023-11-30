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
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
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

func NewFromEnv() *KrakenFutures {
	rootDir, _ := filepath.Abs("../")
	err := godotenv.Load(filepath.Join(rootDir, ".env"))
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	key := os.Getenv("KRAKEN_KEY")
	pk := os.Getenv("KRAKEN_PRIVATE_KEY")

	return New(key, pk)
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

// GenerateAuthentSignature generates a signature based on the provided parameters.
func (api *KrakenFutures) getSign(endpointPath string, data url.Values, nonce string) (string, error) {
	// Step 1: Concatenate postData + Nonce + endpointPath
	dataToHash := data.Encode() + nonce + endpointPath

	// Step 2: Hash the result of step 1 with SHA-256
	sha256Hasher := sha256.New()
	sha256Hasher.Write([]byte(dataToHash))
	hashedData := sha256Hasher.Sum(nil)

	// Step 3: Base64-decode the api_secret
	decodedSecret, err := base64.StdEncoding.DecodeString(api.secret)
	if err != nil {
		return "", err
	}

	// Step 4: Use the result of step 3 to hash the result of step 2 with HMAC-SHA-512
	hmac512 := hmac.New(sha512.New, decodedSecret)
	hmac512.Write(hashedData)
	hmacHash := hmac512.Sum(nil)

	// Step 5: Base64-encode the result of step 4
	signature := base64.StdEncoding.EncodeToString(hmacHash)

	return signature, nil
}

func (api *KrakenFutures) prepareRequest(reqType string, method string, isPrivate bool, data url.Values) (*http.Request, error) {
	if data == nil {
		data = url.Values{}
	}
	requestURL := ""
	if isPrivate {
		requestURL = fmt.Sprintf("%s/%s", FuturesAPIUrl, method)
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
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	if err != nil {
		return nil, errors.Wrap(err, "error during request creation")
	}

	if isPrivate {
		nonce := fmt.Sprintf("%d", time.Now().UnixMilli())
		urlPath := fmt.Sprintf("/api/v3/%s", method)
		signature, err := api.getSign(urlPath, data, nonce)
		if err != nil {
			return nil, errors.Wrap(err, "invalid secret key")
		}
		req.Header.Add("APIKey", api.key)
		req.Header.Add("Authent", signature)
		req.Header.Add("Nonce", nonce)
	}
	return req, nil
}

func (api *KrakenFutures) parseResponse(response *http.Response, retType interface{}) error {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "error during response parsing: can not read response body")
	}

	if response.StatusCode != 200 {
		fmt.Println(string(body))
		return errors.Errorf("error during response parsing: invalid status code %d", response.StatusCode)
	}

	if response.Body == nil {
		return errors.New("error during response parsing: can not read response body")
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
