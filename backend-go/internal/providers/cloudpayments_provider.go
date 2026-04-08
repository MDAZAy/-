package providers

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"vpn-bot/backend-go/internal/config"
)

const cloudPaymentsAPIBaseURL = "https://api.cloudpayments.ru"

type CloudPaymentsProvider struct {
	publicID      string
	apiSecret     string
	publicBaseURL string
	client        *http.Client
}

type cloudPaymentsCreateOrderRequest struct {
	Amount              float64                `json:"Amount"`
	Currency            string                 `json:"Currency"`
	InvoiceID           string                 `json:"InvoiceId"`
	Description         string                 `json:"Description"`
	AccountID           string                 `json:"AccountId"`
	Email               string                 `json:"Email,omitempty"`
	SendEmail           bool                   `json:"SendEmail"`
	RequireConfirmation bool                   `json:"RequireConfirmation,omitempty"`
	SuccessRedirectURL  string                 `json:"SuccessRedirectUrl,omitempty"`
	FailRedirectURL     string                 `json:"FailRedirectUrl,omitempty"`
	JsonData            map[string]interface{} `json:"JsonData,omitempty"`
}

type cloudPaymentsCreateOrderResponse struct {
	Success bool                          `json:"Success"`
	Message string                        `json:"Message"`
	Model   cloudPaymentsCreateOrderModel `json:"Model"`
}

type cloudPaymentsCreateOrderModel struct {
	ID  string `json:"Id"`
	URL string `json:"Url"`
}

type cloudPaymentsWebhookPayload struct {
	TransactionID string
	InvoiceID     string
	Status        string
	Amount        string
	Reason        string
	ReasonCode    string
}

func NewCloudPaymentsProvider(cfg config.Config) PaymentProvider {
	return &CloudPaymentsProvider{
		publicID:      cfg.CloudPaymentsPublicID,
		apiSecret:     cfg.CloudPaymentsAPIToken,
		publicBaseURL: strings.TrimRight(cfg.PublicBaseURL, "/"),
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (p *CloudPaymentsProvider) Name() string {
	return "cloudpayments"
}

func (p *CloudPaymentsProvider) CreatePayment(userID uint, planID uint, amount float64, returnURL string) (PaymentIntent, error) {
	if p.publicID == "" || p.apiSecret == "" {
		return PaymentIntent{}, fmt.Errorf("cloudpayments credentials are not configured")
	}

	invoiceID := fmt.Sprintf("cp_%d_%d_%s", userID, planID, randomHex(8))
	if strings.TrimSpace(returnURL) == "" {
		returnURL = p.publicBaseURL + "/payments/return"
	}

	requestBody := cloudPaymentsCreateOrderRequest{
		Amount:             amount,
		Currency:           "RUB",
		InvoiceID:          invoiceID,
		Description:        fmt.Sprintf("VPN subscription plan %d for user %d", planID, userID),
		AccountID:          strconv.FormatUint(uint64(userID), 10),
		SendEmail:          false,
		SuccessRedirectURL: returnURL + "?status=success",
		FailRedirectURL:    returnURL + "?status=failed",
		JsonData: map[string]interface{}{
			"user_id": userID,
			"plan_id": planID,
		},
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return PaymentIntent{}, fmt.Errorf("marshal cloudpayments create order: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, cloudPaymentsAPIBaseURL+"/orders/create", bytes.NewReader(payload))
	if err != nil {
		return PaymentIntent{}, fmt.Errorf("create cloudpayments request: %w", err)
	}
	req.SetBasicAuth(p.publicID, p.apiSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return PaymentIntent{}, fmt.Errorf("send cloudpayments create order: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PaymentIntent{}, fmt.Errorf("read cloudpayments create order response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return PaymentIntent{}, fmt.Errorf("cloudpayments create order failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var result cloudPaymentsCreateOrderResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return PaymentIntent{}, fmt.Errorf("decode cloudpayments create order response: %w", err)
	}
	if !result.Success {
		return PaymentIntent{}, fmt.Errorf("cloudpayments create order failed: %s", strings.TrimSpace(result.Message))
	}
	if strings.TrimSpace(result.Model.URL) == "" {
		return PaymentIntent{}, fmt.Errorf("cloudpayments create order returned empty payment url")
	}

	return PaymentIntent{
		ExternalID: invoiceID,
		PaymentURL: result.Model.URL,
		RawPayload: string(body),
		Status:     "pending",
	}, nil
}

func (p *CloudPaymentsProvider) ParseWebhook(payload []byte, headers http.Header) (PaymentWebhookEvent, error) {
	if p.publicID == "" || p.apiSecret == "" {
		return PaymentWebhookEvent{}, fmt.Errorf("cloudpayments credentials are not configured")
	}

	if err := p.validateSignature(payload, headers); err != nil {
		return PaymentWebhookEvent{}, err
	}

	data, err := parseCloudPaymentsWebhook(payload, headers.Get("Content-Type"))
	if err != nil {
		return PaymentWebhookEvent{}, err
	}
	if data.InvoiceID == "" {
		return PaymentWebhookEvent{}, fmt.Errorf("%w: empty InvoiceId", ErrInvalidPaymentWebhook)
	}

	status := normalizeCloudPaymentsStatus(data.Status)
	if status == "pending" && (data.Reason != "" || data.ReasonCode != "") {
		status = "failed"
	}
	return PaymentWebhookEvent{
		ExternalID: data.InvoiceID,
		Status:     status,
		Event:      status,
		RawPayload: string(payload),
		IsSuccess:  status == "succeeded",
	}, nil
}

func (p *CloudPaymentsProvider) WebhookSuccessResponse() map[string]interface{} {
	return map[string]interface{}{
		"code": 0,
	}
}

func (p *CloudPaymentsProvider) validateSignature(payload []byte, headers http.Header) error {
	signature := strings.TrimSpace(headers.Get("X-Content-HMAC"))
	if signature == "" {
		signature = strings.TrimSpace(headers.Get("Content-HMAC"))
	}
	if signature == "" {
		return fmt.Errorf("%w: missing CloudPayments HMAC header", ErrInvalidPaymentWebhook)
	}

	mac := hmac.New(sha256.New, []byte(p.apiSecret))
	mac.Write(payload)
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(signature), []byte(expected)) {
		return fmt.Errorf("%w: invalid CloudPayments HMAC", ErrInvalidPaymentWebhook)
	}

	return nil
}

func parseCloudPaymentsWebhook(payload []byte, contentType string) (cloudPaymentsWebhookPayload, error) {
	if strings.Contains(strings.ToLower(contentType), "application/json") {
		var input map[string]interface{}
		if err := json.Unmarshal(payload, &input); err != nil {
			return cloudPaymentsWebhookPayload{}, fmt.Errorf("%w: decode CloudPayments json webhook: %v", ErrInvalidPaymentWebhook, err)
		}

		return cloudPaymentsWebhookPayload{
			TransactionID: fmt.Sprint(input["TransactionId"]),
			InvoiceID:     fmt.Sprint(input["InvoiceId"]),
			Status:        fmt.Sprint(input["Status"]),
			Amount:        fmt.Sprint(input["Amount"]),
			Reason:        fmt.Sprint(input["Reason"]),
			ReasonCode:    fmt.Sprint(input["ReasonCode"]),
		}, nil
	}

	values, err := url.ParseQuery(string(payload))
	if err != nil {
		return cloudPaymentsWebhookPayload{}, fmt.Errorf("%w: decode CloudPayments form webhook: %v", ErrInvalidPaymentWebhook, err)
	}

	return cloudPaymentsWebhookPayload{
		TransactionID: values.Get("TransactionId"),
		InvoiceID:     values.Get("InvoiceId"),
		Status:        values.Get("Status"),
		Amount:        values.Get("Amount"),
		Reason:        values.Get("Reason"),
		ReasonCode:    values.Get("ReasonCode"),
	}, nil
}

func normalizeCloudPaymentsStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "completed", "authorized", "success", "paid":
		return "succeeded"
	case "declined", "cancelled", "canceled", "rejected", "fail", "failed":
		return "failed"
	default:
		if strings.TrimSpace(status) == "" {
			return "pending"
		}
		return strings.ToLower(strings.TrimSpace(status))
	}
}

func randomHex(size int) string {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}
