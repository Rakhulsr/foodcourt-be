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
	var paymentURL string

	for _, itemReq := range req.Items {
		menu, err := u.menuRepo.FindByID(itemReq.MenuID)
		if err != nil {
			return nil, fmt.Errorf("menu ID %d tidak ditemukan", itemReq.MenuID)
		}

		if !menu.IsAvailable {
			return nil, fmt.Errorf("menu '%s' sedang tidak tersedia", menu.Name)
		}
		if !menu.Booth.IsActive {
			return nil, fmt.Errorf("booth '%s' sedang tutup", menu.Booth.Name)
		}

		total += menu.Price * itemReq.Quantity

		orderItems = append(orderItems, model.OrderItem{

			MenuID:          itemReq.MenuID,
			BoothID:         menu.BoothID,
			Quantity:        itemReq.Quantity,
			PriceAtPurchase: menu.Price,
			Notes:           itemReq.Notes,
		})
	}

	order := model.Order{
		OrderCode:     "ORD-" + utils.RandomString(8),
		CustomerName:  req.CustomerName,
		TableNumber:   req.TableNumber,
		TotalAmount:   total,
		PaymentMethod: req.PaymentMethod,
		OrderStatus:   "pending",
		PaymentStatus: "pending",
		Items:         orderItems,
	}

	if err := u.orderRepo.Create(&order); err != nil {
		return nil, err
	}

	if order.PaymentMethod == "qris" {
		inv, err := u.paymentUc.CreateInvoice(order)
		if err != nil {

			return nil, fmt.Errorf("gagal membuat invoice xendit: %w", err)
		}

		paymentURL = inv.InvoiceUrl

		err = u.orderRepo.UpdateInvoice(order.OrderCode, *inv.Id, inv.InvoiceUrl)
		if err != nil {
			return nil, err
		}
	} else {
		paymentURL = ""
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
		PaymentURL: paymentURL,
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
		"pending": true, "confirmed": true, "preparing": true,
		"ready": true, "completed": true, "cancelled": true,
	}
	if !validStatuses[newStatus] {
		return errors.New("invalid status")
	}

	order, err := u.orderRepo.FindByCode(orderCode)
	if err != nil {
		return err
	}

	if newStatus == "cancelled" {
		u.orderRepo.UpdatePaymentStatus(orderCode, "expired")
	}

	if order.PaymentMethod == "cash" {

		if newStatus == "confirmed" || newStatus == "preparing" || newStatus == "ready" || newStatus == "completed" {
			if order.PaymentStatus != "paid" {
				u.orderRepo.UpdatePaymentStatus(orderCode, "paid")
			}
		}

		if newStatus == "pending" {
			u.orderRepo.UpdatePaymentStatus(orderCode, "pending")
		}
	}

	if order.OrderStatus == "completed" || order.OrderStatus == "cancelled" {
		return errors.New("pesanan sudah final (selesai/batal) dan tidak dapat diubah lagi")
	}

	return u.orderRepo.UpdateOrderStatus(orderCode, newStatus)
}

func (u *orderUsecase) ProcessXenditCallback(payload dto.XenditCallbackRequest) error {

	order, err := u.orderRepo.FindByCode(payload.ExternalID)
	if err != nil {
		return err
	}

	if payload.Status == "PAID" || payload.Status == "SETTLED" {

		if err := u.orderRepo.UpdatePaymentStatus(order.OrderCode, "paid"); err != nil {
			return err
		}

		if err := u.orderRepo.UpdateOrderStatus(order.OrderCode, "confirmed"); err != nil {
			return err
		}
	} else if payload.Status == "EXPIRED" {
		u.orderRepo.UpdatePaymentStatus(order.OrderCode, "expired")
		u.orderRepo.UpdateOrderStatus(order.OrderCode, "cancelled")
	}

	return nil
}
