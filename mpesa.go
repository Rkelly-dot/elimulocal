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
func (d *DarajaClient) initiateSTKPush(phone string, amount int, reference string) (*STKPushResponse, error) {
	token, err := d.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("could not get access token: %w", err)
	}

	timestamp := time.Now().Format("20060102150405")

	rawPassword := d.ShortCode + d.Passkey + timestamp
	password := base64.StdEncoding.EncodeToString([]byte(rawPassword))

	payload := STKPushRequest{
		BusinessShortCode: d.ShortCode,
		Password:          password,
		Timestamp:         timestamp,
		TransactionType:   "CustomerPayBillOnline",
		Amount:            amount,
		PartyA:            phone,
		PartyB:            d.ShortCode,
		PhoneNumber:       phone,
		CallBackURL:       d.CallbackURL,
		AccountReference:  reference,
		TransactionDesc:   "ElimuLocal Premium - " + reference,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := d.BaseURL + "/mpesa/stkpush/v1/processrequest"
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("STK push request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var stkResp STKPushResponse
	if err := json.Unmarshal(respBody, &stkResp); err != nil {
		return nil, fmt.Errorf("could not parse STK response: %w", err)
	}

	log.Printf("STK Push initiated: CheckoutRequestID=%s ResponseCode=%s",
		stkResp.CheckoutRequestID, stkResp.ResponseCode)

	return &stkResp, nil
}
func mpesaPayHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, loggedIn := getSessionUser(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		data := newPageData(r, "Subscribe - ElimuLocal")
		renderTemplate(w, "subscribe.html", data)
		return
	}

	if r.Method == "POST" {
		phone := r.FormValue("phone")
		if phone == "" {
			data := newPageData(r, "Subscribe - ElimuLocal")
			data.Message = "Please enter your M-Pesa phone number."
			renderTemplate(w, "subscribe.html", data)
			return
		}

		phone = normalisePhone(phone)

		amount := 2000 

		var paymentID int64
		result, err := db.Exec(
			`INSERT INTO payments
			 (user_id, phone_number, amount, account_reference, status, created_at, updated_at)
			 VALUES (?, ?, ?, ?, 'pending', ?, ?)`,
			currentUser.ID, phone, amount, "ELIMULOCAL-PREMIUM",
			time.Now().Format("2006-01-02 15:04:05"),
			time.Now().Format("2006-01-02 15:04:05"),
		)
		if err != nil {
			log.Printf("Could not save payment record: %v", err)
		} else {
			paymentID, _ = result.LastInsertId()
		}

		stkResp, err := daraja.initiateSTKPush(phone, amount, "ELIMULOCAL-PREMIUM")
		if err != nil {
			log.Printf("STK push failed: %v", err)
			data := newPageData(r, "Subscribe - ElimuLocal")
			data.Message = "Could not initiate M-Pesa payment. Please try again."
			renderTemplate(w, "subscribe.html", data)
			return
		}

		if paymentID > 0 {
			db.Exec(
				`UPDATE payments
				 SET checkout_request_id=?, merchant_request_id=?, updated_at=?
				 WHERE id=?`,
				stkResp.CheckoutRequestID,
				stkResp.MerchantRequestID,
				time.Now().Format("2006-01-02 15:04:05"),
				paymentID,
			)
		}

		data := newPageData(r, "Payment Initiated - ElimuLocal")
		data.Message = "✅ Check your phone! Enter your M-Pesa PIN to complete the payment. " +
			"Reference: " + stkResp.CheckoutRequestID
		renderTemplate(w, "subscribe.html", data)
		return
	}
}
func mpesaCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Callback: could not read body: %v", err)
		w.WriteHeader(http.StatusOK) // always return 200 to Daraja
		return
	}
	defer r.Body.Close()

	log.Printf("M-Pesa callback received: %s", string(body))

	var callback STKCallback
	if err := json.Unmarshal(body, &callback); err != nil {
		log.Printf("Callback: could not parse body: %v", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	stk := callback.Body.STKCallback
	resultCode := stk.ResultCode
	checkoutRequestID := stk.CheckoutRequestID

	// ResultCode 0 means success — any other code means failure
	status := "failed"
	var receiptNumber string

	if resultCode == 0 {
		status = "completed"
		// Extract the M-Pesa receipt number from callback metadata
		for _, item := range stk.CallbackMetadata.Item {
			if item.Name == "MpesaReceiptNumber" {
				if v, ok := item.Value.(string); ok {
					receiptNumber = v
				}
			}
		}
	}

	// Update the payment record matching this checkout request
	db.Exec(
		`UPDATE payments
		 SET status=?, mpesa_receipt_number=?, result_code=?, result_desc=?, updated_at=?
		 WHERE checkout_request_id=?`,
		status,
		receiptNumber,
		fmt.Sprintf("%d", resultCode),
		stk.ResultDesc,
		time.Now().Format("2006-01-02 15:04:05"),
		checkoutRequestID,
	)

	if status == "completed" {
		log.Printf("Payment completed: receipt=%s checkout=%s",
			receiptNumber, checkoutRequestID)
		// TODO Month 4: activate campus premium tier here
	} else {
		log.Printf("Payment failed: code=%d desc=%s checkout=%s",
			resultCode, stk.ResultDesc, checkoutRequestID)
	}

	// Always respond 200 to Daraja — if you return any other
	// status code, Daraja will retry the callback repeatedly
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ResultCode":0,"ResultDesc":"Success"}`))
}