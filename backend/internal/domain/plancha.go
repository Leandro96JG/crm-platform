package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PlanchaRepository interface {
	Create(ctx context.Context, plancha Plancha) (string, error)
	GetByID(ctx context.Context, planchaID string) (*Plancha, error)
	Search(ctx context.Context, filters PlanchaFilters) (PagingResult[Plancha], error)
	Update(ctx context.Context, plancha Plancha) error
	Delete(ctx context.Context, planchaID string) error
}

type Plancha struct {
	PlanchaID       string
	Name            string
	Description     string
	Category        string
	Subcategory     string
	LayoutFileURL   string
	PreviewImageURL string
	Notes           string
	IsActive        bool
	CreatedBy       string
	CreatedAt       time.Time
	UpdatedBy       string
	UpdatedAt       time.Time
}

type PlanchaFilters struct {
	PlanchaID  []string
	Category   []string
	Subcategory []string
	IsActive   *bool
	PagingFilter
}

func NewPlancha(
	name string,
	description string,
	category string,
	subcategory string,
	layoutFileURL string,
	notes string,
	createdBy string,
) (Plancha, error) {
	now := time.Now().UTC()
	planchaID, err := uuid.NewUUID()
	if err != nil {
		return Plancha{}, err
	}

	return Plancha{
		PlanchaID:       planchaID.String(),
		Name:            name,
		Description:     description,
		Category:        category,
		Subcategory:     subcategory,
		LayoutFileURL:   layoutFileURL,
		Notes:           notes,
		IsActive:        true,
		CreatedBy:       createdBy,
		CreatedAt:       now,
		UpdatedBy:       createdBy,
		UpdatedAt:       now,
	}, nil
}

type UpdatePlancha struct {
	Name            *string
	Description     *string
	Category        *string
	Subcategory     *string
	LayoutFileURL   *string
	PreviewImageURL *string
	Notes           *string
	IsActive        *bool
	UpdatedBy       string
}

func (p *Plancha) MergeUpdate(update UpdatePlancha) {
	p.UpdatedAt = time.Now().UTC()
	p.UpdatedBy = update.UpdatedBy

	if update.Name != nil {
		p.Name = *update.Name
	}
	if update.Description != nil {
		p.Description = *update.Description
	}
	if update.Category != nil {
		p.Category = *update.Category
	}
	if update.Subcategory != nil {
		p.Subcategory = *update.Subcategory
	}
	if update.LayoutFileURL != nil {
		p.LayoutFileURL = *update.LayoutFileURL
	}
	if update.PreviewImageURL != nil {
		p.PreviewImageURL = *update.PreviewImageURL
	}
	if update.Notes != nil {
		p.Notes = *update.Notes
	}
	if update.IsActive != nil {
		p.IsActive = *update.IsActive
	}
}
