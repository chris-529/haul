package db

import (
	"context"

	"github.com/chris-529/haul/internal/models"
)

func SaveReceipt(ctx context.Context, userID string, receipt *models.Receipt) error {
	tx, err := Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var receiptID string

	// Save receipt into DB
	err = tx.QueryRow(ctx,
		`INSERT INTO receipts (user_id, store, status)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		userID,
		receipt.Store,
		receipt.Status,
	).Scan(&receiptID)
	if err != nil {
		return err
	}

	receipt.ID = receiptID
	receipt.UserID = userID

	// Also save the items from the receipt and pass the receipt ID in as a fk
	for i := range receipt.Items {
		item := &receipt.Items[i]

		err = tx.QueryRow(ctx,
			`INSERT INTO items (receipt_id, name, price, quantity, unit)
			 VALUES ($1, $2, $3, $4, $5)
			 RETURNING id`,
			receiptID,
			item.Name,
			item.Price,
			item.Quantity,
			item.Unit,
		).Scan(&item.ID)
		if err != nil {
			return err
		}

		item.ReceiptID = receiptID
	}

	return tx.Commit(ctx)
}
