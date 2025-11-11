package usecase

import (
	"errors"

	"github.com/Rakhulsr/foodcourt/internal/dto"
	"github.com/Rakhulsr/foodcourt/internal/model"
	"github.com/Rakhulsr/foodcourt/internal/repository"
)

type BoothUseCase interface {
	ListActive() (*dto.BoothListResponse, error)
	ListAll() (dto.BoothListResponse, error)
	GetByID(id uint) (*dto.BoothResponse, error)
	Create(req dto.BoothCreateRequest) (*dto.BoothResponse, error)
	Update(id uint, req dto.BoothUpdateRequest) (*dto.BoothResponse, error)
	Delete(id uint) error
}

type boothUseCase struct {
	repo repository.BoothRepository
}

func NewBoothUseCase(repo repository.BoothRepository) BoothUseCase {
	return &boothUseCase{repo: repo}
}

func (u *boothUseCase) ListActive() (*dto.BoothListResponse, error) {
	booths, err := u.repo.FindActive()
	if err != nil {
		return nil, err
	}

	resp := &dto.BoothListResponse{Total: len(booths)}
	for _, b := range booths {
		resp.Booths = append(resp.Booths, dto.BoothResponse{
			ID:       b.ID,
			Name:     b.Name,
			WhatsApp: b.WhatsApp,
			IsActive: b.IsActive,
		})
	}
	return resp, nil
}

func (u *boothUseCase) ListAll() (dto.BoothListResponse, error) {
	booths, err := u.repo.FindAll()
	if err != nil {
		return dto.BoothListResponse{}, err
	}

	resp := dto.BoothListResponse{Total: len(booths)}
	for _, b := range booths {
		resp.Booths = append(resp.Booths, dto.BoothResponse{
			ID:        b.ID,
			Name:      b.Name,
			WhatsApp:  b.WhatsApp,
			IsActive:  b.IsActive,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		})
	}

	return resp, nil

}

func (u *boothUseCase) GetByID(id uint) (*dto.BoothResponse, error) {
	booth, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.BoothResponse{
		ID:       booth.ID,
		Name:     booth.Name,
		WhatsApp: booth.WhatsApp,
		IsActive: booth.IsActive,
	}, nil
}

func (u *boothUseCase) Create(req dto.BoothCreateRequest) (*dto.BoothResponse, error) {

	existing, err := u.repo.FindByExactName(req.Name)
	if err == nil && existing != nil {
		return nil, errors.New("name already exists")
	}

	booth := &model.Booth{
		Name:     req.Name,
		WhatsApp: req.WhatsApp,
		IsActive: req.IsActive,
	}

	if err := u.repo.Create(booth); err != nil {
		return nil, err
	}

	return &dto.BoothResponse{
		ID:        booth.ID,
		Name:      booth.Name,
		WhatsApp:  booth.WhatsApp,
		IsActive:  booth.IsActive,
		CreatedAt: booth.CreatedAt,
		UpdatedAt: booth.UpdatedAt,
	}, nil
}

func (u *boothUseCase) Update(id uint, req dto.BoothUpdateRequest) (*dto.BoothResponse, error) {
	booth, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		booth.Name = req.Name
	}
	if req.WhatsApp != "" {
		booth.WhatsApp = req.WhatsApp
	}
	booth.IsActive = req.IsActive

	if err := u.repo.Update(booth); err != nil {
		return nil, err
	}

	return &dto.BoothResponse{
		ID:        booth.ID,
		Name:      booth.Name,
		WhatsApp:  booth.WhatsApp,
		IsActive:  booth.IsActive,
		CreatedAt: booth.CreatedAt,
		UpdatedAt: booth.UpdatedAt,
	}, nil
}

func (u *boothUseCase) Delete(id uint) error {
	return u.repo.Delete(id)
}
