package application

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type planchaService struct {
	planchaRepo        domain.PlanchaRepository
	materialRepo       domain.StickerMaterialRepository
	planchaPriceRepo   domain.PlanchaPriceRepository
}

type PlanchaService interface {
	CreatePlancha(ctx context.Context, plancha domain.Plancha) (string, error)
	GetPlanchaByID(ctx context.Context, planchaID string) (*domain.Plancha, error)
	SearchPlanchas(ctx context.Context, filters domain.PlanchaFilters) (domain.PagingResult[domain.Plancha], error)
	UpdatePlancha(ctx context.Context, planchaID string, update domain.UpdatePlancha) error
	DeletePlancha(ctx context.Context, planchaID string) error

	CreateMaterial(ctx context.Context, material domain.StickerMaterial) (string, error)
	GetMaterialByID(ctx context.Context, materialID string) (*domain.StickerMaterial, error)
	SearchMaterials(ctx context.Context, filters domain.StickerMaterialFilters) (domain.PagingResult[domain.StickerMaterial], error)

	CreatePrice(ctx context.Context, price domain.PlanchaPrice) (string, error)
	GetPlanchaPrices(ctx context.Context, planchaID string) ([]domain.PlanchaPrice, error)
	CalculatePrice(ctx context.Context, planchaID string, materialID string, quantity int) (float64, error)
}

func NewPlanchaService(
	planchaRepo domain.PlanchaRepository,
	materialRepo domain.StickerMaterialRepository,
	planchaPriceRepo domain.PlanchaPriceRepository,
) PlanchaService {
	return &planchaService{
		planchaRepo:      planchaRepo,
		materialRepo:     materialRepo,
		planchaPriceRepo: planchaPriceRepo,
	}
}

func (s *planchaService) CreatePlancha(ctx context.Context, plancha domain.Plancha) (string, error) {
	return s.planchaRepo.Create(ctx, plancha)
}

func (s *planchaService) GetPlanchaByID(ctx context.Context, planchaID string) (*domain.Plancha, error) {
	return s.planchaRepo.GetByID(ctx, planchaID)
}

func (s *planchaService) SearchPlanchas(ctx context.Context, filters domain.PlanchaFilters) (domain.PagingResult[domain.Plancha], error) {
	return s.planchaRepo.Search(ctx, filters)
}

func (s *planchaService) UpdatePlancha(ctx context.Context, planchaID string, update domain.UpdatePlancha) error {
	plancha, err := s.planchaRepo.GetByID(ctx, planchaID)
	if err != nil {
		return err
	}

	plancha.MergeUpdate(update)
	return s.planchaRepo.Update(ctx, *plancha)
}

func (s *planchaService) DeletePlancha(ctx context.Context, planchaID string) error {
	return s.planchaRepo.Delete(ctx, planchaID)
}

func (s *planchaService) CreateMaterial(ctx context.Context, material domain.StickerMaterial) (string, error) {
	return s.materialRepo.Create(ctx, material)
}

func (s *planchaService) GetMaterialByID(ctx context.Context, materialID string) (*domain.StickerMaterial, error) {
	return s.materialRepo.GetByID(ctx, materialID)
}

func (s *planchaService) SearchMaterials(ctx context.Context, filters domain.StickerMaterialFilters) (domain.PagingResult[domain.StickerMaterial], error) {
	return s.materialRepo.Search(ctx, filters)
}

func (s *planchaService) CreatePrice(ctx context.Context, price domain.PlanchaPrice) (string, error) {
	return s.planchaPriceRepo.Create(ctx, price)
}

func (s *planchaService) GetPlanchaPrices(ctx context.Context, planchaID string) ([]domain.PlanchaPrice, error) {
	return s.planchaPriceRepo.SearchByPlancha(ctx, planchaID)
}

func (s *planchaService) CalculatePrice(ctx context.Context, planchaID string, materialID string, quantity int) (float64, error) {
	price, err := s.planchaPriceRepo.GetByPlanchaAndMaterial(ctx, planchaID, materialID)
	if err != nil {
		return 0, err
	}

	return price.CalculatePrice(quantity), nil
}
