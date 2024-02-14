package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Mpesa is an application that will be making a transaction
type Mpesa struct {
    consumerKey    string
    consumerSecret string
    baseURL        string
    client         *http.Client
}

// MpesaOpts stores all the configuration keys we need to set up a Mpesa app,
type MpesaOpts struct {
    ConsumerKey    string
    ConsumerSecret string
    BaseURL        string
}

// MpesaAccessTokenResponse is the response sent back by Safaricom when we make a request to generate a token
type MpesaAccessTokenResponse struct {
    AccessToken  string `json:"access_token"`
    ExpiresIn    string `json:"expires_in"`
    RequestID    string `json:"requestId"`
    ErrorCode    string `json:"errorCode"`
    ErrorMessage string `json:"errorMessage"`
}

// NewMpesa sets up and returns an instance of Mpesa
func NewMpesa(m *MpesaOpts) *Mpesa {
    client := &http.Client{
        Timeout: 30 * time.Second,
    }

  return &Mpesa{
      consumerKey:    m.ConsumerKey,
      consumerSecret: m.ConsumerSecret,
      baseURL:        m.BaseURL,
      client:         client,
  }
}

// makeRequest performs all the http requests for the specific app
func (m *Mpesa) makeRequest(req *http.Request) ([]byte, error) {
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// generateAccessToken sends a http request to generate new access token
func (m *Mpesa) generateAccessToken() (*MpesaAccessTokenResponse, error) {
	url := fmt.Sprintf("%s/oauth/v1/generate?grant_type=client_credentials", m.baseURL)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(m.consumerKey, m.consumerSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.makeRequest(req)
	if err != nil {
		return nil, err
	}

	accessTokenResponse := new(MpesaAccessTokenResponse)
	if err := json.Unmarshal(resp, &accessTokenResponse); err != nil {
		return nil, err
	}

	return accessTokenResponse, nil
}

func main() {
	mpesa := NewMpesa(&MpesaOpts{
		ConsumerKey:    "kXrA3da7AMp5YAP7ngtcyjkZEpbw4gskTwZBjoj1k0YR6wVo",
		ConsumerSecret: "GlGOHdtNlLFPp7rti02fvprZ1zQ5g16mADG86G6BA7yISLsS5NwAj4xQJAxNPtZx",
		BaseURL:        "https://sandbox.safaricom.co.ke",
	})

	accessTokenResponse, err := mpesa.generateAccessToken()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%+v\n", accessTokenResponse)
}
