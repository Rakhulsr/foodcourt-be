package usecase

import (
	"errors"
	"fmt"

	"github.com/Rakhulsr/foodcourt/internal/dto"
	"github.com/Rakhulsr/foodcourt/internal/model"
	"github.com/Rakhulsr/foodcourt/internal/repository"
	"github.com/Rakhulsr/foodcourt/utils"
)

type OrderUsecase interface {
	CreateOrder(req dto.CreateOrderRequest) (*dto.CreateOrderResponse, error)
	GetOrderByCode(code string) (*model.Order, error)

	ListOrders(page int, limit int, status string) (*dto.OrderListResponse, error)
	UpdateOrderStatus(orderCode string, newStatus string) error

	ProcessXenditCallback(payload dto.XenditCallbackRequest) error
}

type orderUsecase struct {
	orderRepo repository.OrderRepository
	menuRepo  repository.MenuRepository
	paymentUc *PaymentUsecase
}

func NewOrderUsecase(or repository.OrderRepository, mr repository.MenuRepository, ps *PaymentUsecase) *orderUsecase {
	return &orderUsecase{
		orderRepo: or,
		menuRepo:  mr,
		paymentUc: ps,
	}
}

func (u *orderUsecase) CreateOrder(req dto.CreateOrderRequest) (*dto.CreateOrderResponse, error) {
	var orderItems []model.OrderItem
	total := 0

	for _, itemReq := range req.Items {
		menu, err := u.menuRepo.FindByID(itemReq.ID)
		if err != nil {
			return nil, fmt.Errorf("menu item within ID %d is not found", itemReq.ID)
		}

		if !menu.IsAvailable {
			return nil, fmt.Errorf("menu '%s' sedang tidak tersedia (stok habis)", menu.Name)
		}
		if !menu.Booth.IsActive {
			return nil, fmt.Errorf("booth '%s' sedang tutup, tidak bisa memesan", menu.Booth.Name)
		}

		total += menu.Price * itemReq.Quantity

		orderItems = append(orderItems, model.OrderItem{
			MenuID:          itemReq.ID,
			BoothID:         itemReq.BoothID,
			Quantity:        itemReq.Quantity,
			PriceAtPurchase: menu.Price,
			Notes:           itemReq.Notes,
		})
	}

	if len(orderItems) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	order := model.Order{
		OrderCode:     "ORD-" + utils.RandomString(10),
		CustomerName:  req.CustomerName,
		TableNumber:   req.TableNumber,
		PaymentMethod: req.PaymentMethod,
		OrderStatus:   "pending",
		PaymentStatus: "pending",
		TotalAmount:   total,
		Items:         orderItems,
	}

	if err := u.orderRepo.Create(&order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	if order.PaymentMethod == "cash" {
		return &dto.CreateOrderResponse{
			OrderCode: order.OrderCode,
			Message:   "Order created successfully. Please pay at cashier.",
		}, nil
	}

	inv, err := u.paymentUc.CreateInvoice(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create Xendit invoice: %w", err)
	}

	err = u.orderRepo.UpdateInvoice(order.OrderCode, *inv.Id, inv.InvoiceUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to update order with invoice: %w", err)
	}

	return &dto.CreateOrderResponse{
		OrderCode:  order.OrderCode,
		PaymentURL: inv.InvoiceUrl,
	}, nil
}

func (u *orderUsecase) GetOrderByCode(code string) (*model.Order, error) {
	return u.orderRepo.FindByCode(code)
}

func (u *orderUsecase) ListOrders(page int, limit int, status string) (*dto.OrderListResponse, error) {

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	orders, total, err := u.orderRepo.FindAll(page, limit, status)
	if err != nil {
		return nil, err
	}

	return &dto.OrderListResponse{
		Total: total,
		Page:  page,
		Data:  orders,
	}, nil
}

func (u *orderUsecase) UpdateOrderStatus(orderCode string, newStatus string) error {

	validStatuses := map[string]bool{
		"pending": true, "confirmed": true, "preparing": true, "ready": true, "completed": true, "cancelled": true,
	}
	if !validStatuses[newStatus] {
		return errors.New("invalid status")
	}

	return u.orderRepo.UpdateOrderStatus(orderCode, newStatus)
}

func (u *orderUsecase) ProcessXenditCallback(payload dto.XenditCallbackRequest) error {

	if payload.ExternalID == "" {
		return errors.New("external_id is missing")
	}

	switch payload.Status {
	case "PAID", "SETTLED":
		err := u.orderRepo.UpdatePaymentStatus(payload.ExternalID, "paid")
		if err != nil {
			return err
		}

	case "EXPIRED":
		u.orderRepo.UpdatePaymentStatus(payload.ExternalID, "expired")
		u.orderRepo.UpdateOrderStatus(payload.ExternalID, "cancelled")
	}

	return nil
}
