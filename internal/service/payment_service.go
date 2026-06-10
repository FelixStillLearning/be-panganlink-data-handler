package service

import (
	"crypto/sha512"
	"encoding/hex"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type PaymentService interface {
	InitiatePayment(orderID string, totalHarga float64, buyerName, buyerEmail string) (string, error)
	VerifyPaymentSignature(orderID, statusCode, grossAmount, signature string) bool
}

type paymentService struct {
	snapClient snap.Client
	serverKey  string
}

func NewPaymentService(serverKey string, isProduction bool) PaymentService {
	var sClient snap.Client
	env := midtrans.Sandbox
	if isProduction {
		env = midtrans.Production
	}
	sClient.New(serverKey, env)
	
	// CoreAPI Setup if needed
	var cClient coreapi.Client
	cClient.New(serverKey, env)

	return &paymentService{
		snapClient: sClient,
		serverKey:  serverKey,
	}
}

func (s *paymentService) InitiatePayment(orderID string, totalHarga float64, buyerName, buyerEmail string) (string, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(totalHarga),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: buyerName,
			Email: buyerEmail,
		},
	}
	
	snapResp, err := s.snapClient.CreateTransaction(req)
	if err != nil {
		return "", err
	}
	return snapResp.Token, nil
}

func (s *paymentService) VerifyPaymentSignature(orderID, statusCode, grossAmount, signature string) bool {
	// signature_key = hash(order_id + status_code + gross_amount + server_key)
	importCrypto := orderID + statusCode + grossAmount + s.serverKey
	
	hasher := sha512.New()
	hasher.Write([]byte(importCrypto))
	expectedSignature := hex.EncodeToString(hasher.Sum(nil))

	return signature == expectedSignature
}
