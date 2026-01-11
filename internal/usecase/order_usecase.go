package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	SendOrderNotificationToSeller(orderCode string) error
}

type orderUsecase struct {
	orderRepo repository.OrderRepository
	menuRepo  repository.MenuRepository
	paymentUc *PaymentUsecase
	waUc      *WhatsAppUsecase
	logUC     LogUseCase
}

func NewOrderUsecase(or repository.OrderRepository, mr repository.MenuRepository, ps *PaymentUsecase, wa WhatsAppUsecase, log LogUseCase) *orderUsecase {
	return &orderUsecase{
		orderRepo: or,
		menuRepo:  mr,
		paymentUc: ps,
		waUc:      &wa,
		logUC:     log,
	}
}
func (u *orderUsecase) CreateOrder(req dto.CreateOrderRequest) (*dto.CreateOrderResponse, error) {
	var orderItems []model.OrderItem
	total := 0
	var paymentURL string

	type NotificationData struct {
		BoothName string
		BoothWA   string
		Items     []string
	}

	notifMap := make(map[uint]*NotificationData)

	for _, itemReq := range req.Items {
		menu, err := u.menuRepo.FindByID(itemReq.MenuID)
		if err != nil {
			return nil, fmt.Errorf("menu ID %d tidak ditemukan", itemReq.MenuID)
		}
		if !menu.IsAvailable {
			return nil, fmt.Errorf("menu '%s' tidak tersedia", menu.Name)
		}
		if !menu.Booth.IsActive {
			return nil, fmt.Errorf("booth '%s' tutup", menu.Booth.Name)
		}

		total += menu.Price * itemReq.Quantity

		orderItems = append(orderItems, model.OrderItem{
			MenuID:          itemReq.MenuID,
			BoothID:         menu.BoothID,
			Quantity:        itemReq.Quantity,
			PriceAtPurchase: menu.Price,
			Notes:           itemReq.Notes,
		})

		if _, exists := notifMap[menu.BoothID]; !exists {
			notifMap[menu.BoothID] = &NotificationData{
				BoothName: menu.Booth.Name,
				BoothWA:   menu.Booth.WhatsApp,
				Items:     []string{},
			}
		}

		noteText := ""
		if itemReq.Notes != "" {
			noteText = fmt.Sprintf(" _(%s)_", itemReq.Notes)
		}
		itemString := fmt.Sprintf("▪️ %dx %s%s", itemReq.Quantity, menu.Name, noteText)
		notifMap[menu.BoothID].Items = append(notifMap[menu.BoothID].Items, itemString)
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

		if err := u.orderRepo.UpdateOrderStatus(order.OrderCode, "preparing"); err != nil {
			return err
		}

	} else if payload.Status == "EXPIRED" {
		u.orderRepo.UpdatePaymentStatus(order.OrderCode, "expired")
		u.orderRepo.UpdateOrderStatus(order.OrderCode, "cancelled")
	}

	return nil
}

func (u *orderUsecase) SendOrderNotificationToSeller(orderCode string) error {
	order, err := u.orderRepo.FindByCode(orderCode)
	if err != nil {
		return err
	}

	type ItemDetail struct {
		MenuName string
		Qty      int
		Notes    string
	}

	type GroupData struct {
		BoothID   uint
		BoothName string
		BoothWA   string
		Items     []ItemDetail
	}

	groups := make(map[uint]*GroupData)

	for _, item := range order.Items {

		if _, exists := groups[item.BoothID]; !exists {
			groups[item.BoothID] = &GroupData{
				BoothID:   item.BoothID,
				BoothName: item.Booth.Name,
				BoothWA:   item.Booth.WhatsApp,
				Items:     []ItemDetail{},
			}
		}

		detail := ItemDetail{
			MenuName: item.Menu.Name,
			Qty:      item.Quantity,
			Notes:    item.Notes,
		}
		groups[item.BoothID].Items = append(groups[item.BoothID].Items, detail)
	}

	for _, group := range groups {
		targetWA := group.BoothWA

		paymentStatus := "BELUM LUNAS ❌"
		if order.PaymentStatus == "paid" {
			paymentStatus = "SUDAH LUNAS ✅"
		}

		msg := fmt.Sprintf("*PESANAN MASUK!* 🔔\nKepada: *%s*\n\nOrder: *%s*\nMeja: *%s*\nPemesan: *%s*\nStatus: *%s*\n\n🍽️ *MENU:*\n",
			group.BoothName, order.OrderCode, order.TableNumber, order.CustomerName, paymentStatus)

		for _, item := range group.Items {
			noteText := ""
			if item.Notes != "" {
				noteText = fmt.Sprintf(" _(%s)_", item.Notes)
			}
			msg += fmt.Sprintf("▪️ %dx %s%s\n", item.Qty, item.MenuName, noteText)
		}
		msg += "\nMohon segera diproses. Terima kasih! 🙏"

		bgCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

		err := u.waUc.SendMessage(bgCtx, targetWA, msg)
		cancel()

		if err != nil {
			return fmt.Errorf("gagal kirim ke %s: %v", group.BoothName, err)
		}

		err = u.logUC.RecordLog(order.ID, group.BoothID, targetWA)
		if err != nil {
			fmt.Printf("⚠️ Gagal menyimpan log WA: %v\n", err)

		}
	}

	return nil
}
