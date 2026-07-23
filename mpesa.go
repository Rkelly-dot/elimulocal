package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)
type OAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn	string `json:"expires_in"`
}
type STKPushRequest struct {
	BusinessShortCode	string `json:"BusinessShortCode"`
	Password		  	string `json:"Password"`
	Timestamp			string `json:"Timestamp"`
	TransactionType		string `json:"TransactionType"`
	Amount				int	   `json:"Amount"`
	PartyA            	string `json:"PartyA"`
	PartyB            	string `json:"PartyB"`
	PhoneNumber			string `json:"PhoneNumber"`
	CallBackURL			string `json:"CallBackURL"`
	AccountReference	string `json:"AccountReference"`
	TransactionDesc		string `json:"TransactionDesc"`
}
type STKPushResponse struct {
	MerchantRequestID	string `json:"MerchantRequestID"`
	CheckoutRequestID	string `json:"CheckoutRequestID"`
	ResponseCode		string `json:"ResponseCode"`
	ResponseDescription	string `json:"ResponseDescription"`
	CustomerMessage		string `json:"CustomerMessage"`
}
type STKCallback struct {
	Body struct {
		STKCallback struct {
			MerchantRequestID string `json:"MerchantRequestID"`
			CheckoutRequestID string `json:"CheckoutRequestID"`
			ResultCode        int    `json:"ResultCode"`
			ResultDesc        string `json:"ResultDesc"`
			CallbackMetadata  struct {
				Item []struct {
					Name  string      `json:"Name"`
					Value interface{} `json:"Value"`
				} `json:"Item"`
			} `json:"CallbackMetadata"`
		} `json:"stkCallback"`
	} `json:"Body"`
}
type DarajaClient struct {
	ConsumerKey		string
	ConsumerSecret	string
	ShortCode		string
	Passkey			string
	CallbackURL 	string
	BaseURL			string
}
var daraja *DarajaClient

func initMpesa() {
	env := os.Getenv("MPESA_ENV")
	baseURL := "https://sandbox.safaricom.co.ke"
	if env == "production" {
		baseURL = "https://api.safaricom.co.ke"
	}

	daraja = &DarajaClient{
		ConsumerKey:    os.Getenv("MPESA_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("MPESA_CONSUMER_SECRET"),
		ShortCode:      os.Getenv("MPESA_SHORTCODE"),
		Passkey:        os.Getenv("MPESA_PASSKEY"),
		CallbackURL:    os.Getenv("MPESA_CALLBACK_URL"),
		BaseURL:        baseURL,
	}

	if daraja.ConsumerKey == "" || daraja.ConsumerSecret == "" {
		log.Println("Warning: M-Pesa credentials not configured")
		return
	}

	fmt.Println("M-Pesa Daraja client initialised (" + env + ")")
}
func (d *DarajaClient) getAccessToken() (string, error) {
	credentials := base64.stdEncoding.EncodeToString(
		[]byte(d.ConsumerKey + ":" + d.ConsumerSecret),
	)
	url := d.BaseURL + "/oauth/v1/generate?grant_type=client_credentials"

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Basic "+credentials)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResp OAuthResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("could not parse token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("empty access token received: %s", string(body))
	}

	return tokenResp.AccessToken, nil
}