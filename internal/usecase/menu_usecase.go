package usecase

import (
	"github.com/Rakhulsr/foodcourt/internal/dto"
	"github.com/Rakhulsr/foodcourt/internal/repository"
)

type MenuUseCase interface {
	ListActive() (*dto.MenuListResponse, error)
	GetByID(id uint) (*dto.MenuResponse, error)
	ListActiveByBoothID(id uint) (*dto.MenuListResponse, error)

	FindByName(keyword string) (*dto.MenuListResponse, error)
	FindByCategory(category string) (*dto.MenuListResponse, error)
}

type menuUseCase struct {
	repo repository.MenuRepository
}

func NewMenuUseCase(repo repository.MenuRepository) MenuUseCase {
	return &menuUseCase{repo: repo}
}

func (u *menuUseCase) ListActive() (*dto.MenuListResponse, error) {
	menus, err := u.repo.FindActive()
	if err != nil {
		return nil, err
	}

	resp := &dto.MenuListResponse{Total: len(menus)}
	for _, m := range menus {
		resp.Menus = append(resp.Menus, dto.MenuResponse{
			ID:          m.ID,
			Name:        m.Name,
			Price:       m.Price,
			IsAvailable: m.IsAvailable,
			Booth: struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
			}{ID: m.Booth.ID, Name: m.Booth.Name},
		})
	}
	return resp, nil
}

func (u *menuUseCase) GetByID(id uint) (*dto.MenuResponse, error) {
	menu, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.MenuResponse{
		ID:          menu.ID,
		Name:        menu.Name,
		Price:       menu.Price,
		IsAvailable: menu.IsAvailable,
		Booth: struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}{ID: menu.Booth.ID, Name: menu.Booth.Name},
	}, nil
}

func (u *menuUseCase) ListActiveByBoothID(id uint) (*dto.MenuListResponse, error) {
	menuList, err := u.repo.FindActiveByBoothID(id)
	if err != nil {
		return nil, err
	}

	resp := &dto.MenuListResponse{Total: len(menuList)}
	for _, m := range menuList {
		resp.Menus = append(resp.Menus, dto.MenuResponse{
			ID:          m.ID,
			Name:        m.Name,
			Price:       m.Price,
			IsAvailable: m.IsAvailable,
			Booth: struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
			}{ID: m.Booth.ID, Name: m.Booth.Name},
		})
	}
	return resp, nil

}

func (u *menuUseCase) FindByName(keyword string) (*dto.MenuListResponse, error) {
	menus, err := u.repo.FindByName(keyword)
	if err != nil {
		return nil, err
	}

	resp := &dto.MenuListResponse{Total: len(menus)}
	for _, m := range menus {
		resp.Menus = append(resp.Menus, dto.MenuResponse{
			ID:          m.ID,
			Name:        m.Name,
			Price:       m.Price,
			IsAvailable: m.IsAvailable,
			Booth: struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
			}{ID: m.Booth.ID, Name: m.Booth.Name},
		})
	}
	return resp, nil
}

func (u *menuUseCase) FindByCategory(category string) (*dto.MenuListResponse, error) {
	menus, err := u.repo.FindByCategory(category)
	if err != nil {
		return nil, err
	}

	resp := &dto.MenuListResponse{Total: len(menus)}
	for _, m := range menus {
		resp.Menus = append(resp.Menus, dto.MenuResponse{
			ID:          m.ID,
			Name:        m.Name,
			Price:       m.Price,
			IsAvailable: m.IsAvailable,
			Booth: struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
			}{ID: m.Booth.ID, Name: m.Booth.Name},
		})
	}
	return resp, nil
}
