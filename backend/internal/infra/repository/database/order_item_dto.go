package database

import "github.com/icrxz/crm-api-core/internal/domain"

type OrderItemDTO struct {
	ItemID            string  `db:"item_id"`
	OrderID           string  `db:"order_id"`
	PlanchaID         string  `db:"plancha_id"`
	MaterialID        string  `db:"material_id"`
	SheetQuantity     int     `db:"sheet_quantity"`
	UnitPrice         float64 `db:"unit_price"`
	Subtotal          float64 `db:"subtotal"`
	CustomDesignFile  string  `db:"custom_design_file"`
	CustomDesignNotes string  `db:"custom_design_notes"`
	SortOrder         int     `db:"sort_order"`
}

func mapOrderItemToDTO(item domain.OrderItem) OrderItemDTO {
	return OrderItemDTO{
		ItemID:            item.ItemID,
		OrderID:           item.OrderID,
		PlanchaID:         item.PlanchaID,
		MaterialID:        item.MaterialID,
		SheetQuantity:     item.SheetQuantity,
		UnitPrice:         item.UnitPrice,
		Subtotal:          item.Subtotal,
		CustomDesignFile:  item.CustomDesignFile,
		CustomDesignNotes: item.CustomDesignNotes,
		SortOrder:         item.SortOrder,
	}
}

func mapDTOToOrderItem(dto OrderItemDTO) domain.OrderItem {
	return domain.OrderItem{
		ItemID:            dto.ItemID,
		OrderID:           dto.OrderID,
		PlanchaID:         dto.PlanchaID,
		MaterialID:        dto.MaterialID,
		SheetQuantity:     dto.SheetQuantity,
		UnitPrice:         dto.UnitPrice,
		Subtotal:          dto.Subtotal,
		CustomDesignFile:  dto.CustomDesignFile,
		CustomDesignNotes: dto.CustomDesignNotes,
		SortOrder:         dto.SortOrder,
	}
}

func mapDTOsToOrderItems(dtos []OrderItemDTO) []domain.OrderItem {
	items := make([]domain.OrderItem, 0, len(dtos))
	for _, dto := range dtos {
		items = append(items, mapDTOToOrderItem(dto))
	}
	return items
}
