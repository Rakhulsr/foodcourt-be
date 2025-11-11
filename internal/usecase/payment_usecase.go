package usecase

import (
	"context"
	"fmt"

	"github.com/Rakhulsr/foodcourt/config"
	"github.com/Rakhulsr/foodcourt/internal/model"
	"github.com/xendit/xendit-go/v7/invoice"
)

type PaymentUsecase struct{}

func NewPaymentService() *PaymentUsecase {
	return &PaymentUsecase{}
}

func (s *PaymentUsecase) CreateInvoice(order model.Order) (*invoice.Invoice, error) {
	desc := fmt.Sprintf("Pembayaran pesanan: %s ", order.OrderCode)
	feUrlSucess := fmt.Sprintf("https://your-frontend/thank-you/%s", order.OrderCode)
	feUrlFailed := fmt.Sprintf("https://your-frontend/payment-failed/%s", order.OrderCode)

	req := invoice.CreateInvoiceRequest{
		ExternalId: order.OrderCode,

		Description: &desc,
		Amount:      float64(order.TotalAmount),
		// PayerEmail:         "",
		Customer: &invoice.CustomerObject{
			GivenNames: *invoice.NewNullableString(&order.CustomerName),
			// Jika nanti ingin tambah phone/email, taruh disini.
			// Email: xendit.PtrString(order.Email),
		},
		SuccessRedirectUrl: &feUrlSucess,
		FailureRedirectUrl: &feUrlFailed,
		ReminderTime:       func() *float32 { v := float32(15); return &v }(),
	}

	resp, _, err := config.XenditClient.InvoiceApi.
		CreateInvoice(context.Background()).
		CreateInvoiceRequest(req).
		Execute()

	return resp, err
}
