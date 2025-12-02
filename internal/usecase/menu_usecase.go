package usecase

import (
	"errors"

	"github.com/Rakhulsr/foodcourt/internal/dto"
	"github.com/Rakhulsr/foodcourt/internal/model"
	"github.com/Rakhulsr/foodcourt/internal/repository"
)

type MenuUseCase interface {
	ListActive() (*dto.MenuListResponse, error)
	ListAll() (*dto.MenuListResponse, error)
	GetByID(id uint) (*dto.MenuResponse, error)
	ListActiveByBoothID(id uint) (*dto.MenuListResponse, error)

	FindByName(keyword string) (*dto.MenuListResponse, error)
	FindByCategory(category string) (*dto.MenuListResponse, error)

	Create(req dto.MenuCreateRequest, imagePath string) (*dto.MenuResponse, error)
	Update(id uint, req dto.MenuUpdateRequest, imagePath string) (*dto.MenuResponse, error)
	Delete(id uint) error
}

type menuUseCase struct {
	repo      repository.MenuRepository
	boothRepo repository.BoothRepository
}

func NewMenuUseCase(repo repository.MenuRepository, bRepo repository.BoothRepository) *menuUseCase {
	return &menuUseCase{repo: repo, boothRepo: bRepo}
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
			ImagePath:   m.ImagePath,
			Category:    m.Category,
			Description: m.Description,
			Booth: struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
			}{ID: m.Booth.ID, Name: m.Booth.Name},
		})
	}
	return resp, nil
}

func (u *menuUseCase) ListAll() (*dto.MenuListResponse, error) {

	menus, err := u.repo.FindAll()
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
			ImagePath:   m.ImagePath,
			Description: m.Description,
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
		Category:    menu.Category,
		ImagePath:   menu.ImagePath,
		Description: menu.Description,
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
			Description: m.Description,
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
			ImagePath:   m.ImagePath,
			Category:    m.Category,
			Description: m.Description,
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
			Description: m.Description,
			Booth: struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
			}{ID: m.Booth.ID, Name: m.Booth.Name},
		})
	}
	return resp, nil
}

func (u *menuUseCase) Create(req dto.MenuCreateRequest, imgPath string) (*dto.MenuResponse, error) {

	booth, err := u.boothRepo.FindByID(req.BoothID)
	if err != nil {
		return nil, errors.New("booth Not Found")
	}

	menu := &model.Menu{
		BoothID:     req.BoothID,
		Name:        req.Name,
		Price:       req.Price,
		IsAvailable: true,
		ImagePath:   imgPath,
		Category:    req.Category,
		Description: req.Description,
	}

	if err := u.repo.Create(menu); err != nil {
		return nil, err
	}

	return &dto.MenuResponse{
		ID:          menu.ID,
		Name:        menu.Name,
		Price:       menu.Price,
		IsAvailable: menu.IsAvailable,
		Category:    menu.Category,
		ImagePath:   menu.ImagePath,
		Description: menu.Description,
		Booth: struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}{ID: booth.ID, Name: booth.Name},
	}, nil

}

func (u *menuUseCase) Update(id uint, req dto.MenuUpdateRequest, imagePath string) (*dto.MenuResponse, error) {
	menu, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.Price != 0 {
		menu.Price = req.Price
	}
	if req.Category != "" {
		menu.Category = req.Category
	}
	menu.IsAvailable = req.IsAvailable

	if imagePath != "" {
		menu.ImagePath = imagePath
	}

	if err := u.repo.Update(menu); err != nil {
		return nil, err
	}

	return &dto.MenuResponse{
		ID:          menu.ID,
		Name:        menu.Name,
		Price:       menu.Price,
		IsAvailable: menu.IsAvailable,
		Category:    menu.Category,
		ImagePath:   menu.ImagePath,
		Description: menu.Description,
		Booth: struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}{ID: menu.Booth.ID, Name: menu.Booth.Name},
	}, nil
}

func (u *menuUseCase) Delete(id uint) error {
	return u.repo.Delete(id)
}
