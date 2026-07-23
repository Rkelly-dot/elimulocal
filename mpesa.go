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
