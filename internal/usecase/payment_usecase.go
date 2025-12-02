package usecase

import (
	"context"
	"fmt"
	"os"

	"github.com/Rakhulsr/foodcourt/config"
	"github.com/Rakhulsr/foodcourt/internal/model"
	"github.com/xendit/xendit-go/v7/invoice"
)

type PaymentUsecase struct{}

func NewPaymentService() *PaymentUsecase {
	return &PaymentUsecase{}
}

func (s *PaymentUsecase) CreateInvoice(order model.Order) (*invoice.Invoice, error) {

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	successURL := fmt.Sprintf("%s/order/success/%s", baseURL, order.OrderCode)
	failURL := fmt.Sprintf("%s/cart", baseURL)

	createInvoiceRequest := *invoice.NewCreateInvoiceRequest(
		order.OrderCode,
		float64(order.TotalAmount),
	)

	createInvoiceRequest.SetDescription(fmt.Sprintf("Pembayaran Order %s - %s", order.OrderCode, order.CustomerName))
	createInvoiceRequest.SetSuccessRedirectUrl(successURL)
	createInvoiceRequest.SetFailureRedirectUrl(failURL)
	createInvoiceRequest.SetCurrency("IDR")
	createInvoiceRequest.SetReminderTime(1)

	resp, _, err := config.XenditClient.InvoiceApi.
		CreateInvoice(context.Background()).
		CreateInvoiceRequest(createInvoiceRequest).
		Execute()

	if err != nil {
		return nil, err
	}

	return resp, nil
}
