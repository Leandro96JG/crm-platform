package domain

import "github.com/google/uuid"

type OrderItem struct {
	ItemID            string
	OrderID           string
	PlanchaID         string
	MaterialID        string
	SheetQuantity     int
	UnitPrice         float64
	Subtotal          float64
	CustomDesignFile  string
	CustomDesignNotes string
	SortOrder         int
}

func NewOrderItem(
	orderID string,
	planchaID string,
	materialID string,
	sheetQuantity int,
	unitPrice float64,
	customDesignFile string,
	customDesignNotes string,
	sortOrder int,
) (OrderItem, error) {
	itemID, err := uuid.NewUUID()
	if err != nil {
		return OrderItem{}, err
	}

	return OrderItem{
		ItemID:            itemID.String(),
		OrderID:           orderID,
		PlanchaID:         planchaID,
		MaterialID:        materialID,
		SheetQuantity:     sheetQuantity,
		UnitPrice:         unitPrice,
		Subtotal:          unitPrice * float64(sheetQuantity),
		CustomDesignFile:  customDesignFile,
		CustomDesignNotes: customDesignNotes,
		SortOrder:         sortOrder,
	}, nil
}
