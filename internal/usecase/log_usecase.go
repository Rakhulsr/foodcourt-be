package usecase

import (
	"github.com/Rakhulsr/foodcourt/internal/model"
	"github.com/Rakhulsr/foodcourt/internal/repository"
)

type LogUseCase interface {
	RecordLog(orderID uint, boothID uint, targetPhone string) error
	GetLogs(page int, limit int) ([]model.WhatsAppLog, int64, error)
}

type logUseCase struct {
	repo repository.WhatsAppLogRepository
}

func NewLogUseCase(repo repository.WhatsAppLogRepository) LogUseCase {
	return &logUseCase{repo: repo}
}

func (u *logUseCase) RecordLog(orderID uint, boothID uint, targetPhone string) error {
	log := &model.WhatsAppLog{
		OrderID:     &orderID,
		BoothID:     &boothID,
		MessageType: "manual_click",
		Status:      "clicked",
		Response:    "Redirected to " + targetPhone,
	}
	return u.repo.Create(log)
}

func (u *logUseCase) GetLogs(page int, limit int) ([]model.WhatsAppLog, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return u.repo.FindAll(page, limit)
}
