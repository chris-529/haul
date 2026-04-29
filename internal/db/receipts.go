package db

import (
	"context"
	"errors"

	"github.com/chris-529/haul/internal/models"
)

var ErrNotFound = errors.New("not found")

// Save a given receipt for userID

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
		RETURNING id, created_at`,
		userID,
		receipt.Store,
		receipt.Status,
	).Scan(&receiptID, &receipt.CreatedAt)
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

// Returns all receipts for userID

func GetReceipts(ctx context.Context, userID string) ([]models.Receipt, error) {

	// Query all receipts belonging to userID
	// TODO: Replace with JOIN
	rows, err := Pool.Query(ctx,
		`SELECT id, user_id, store, status, created_at
		FROM receipts
		WHERE user_id = $1
		ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	receipts := []models.Receipt{}

	// For each receipt, build a receipt by also querying the items that belong to it
	for rows.Next() {
		var receipt models.Receipt

		err := rows.Scan(
			&receipt.ID,
			&receipt.UserID,
			&receipt.Store,
			&receipt.Status,
			&receipt.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Find the items for this receipt
		itemRows, err := Pool.Query(ctx,
			`SELECT id, receipt_id, name, price, quantity, unit
			 FROM items
			 WHERE receipt_id = $1
			 ORDER BY id ASC`,
			receipt.ID,
		)
		if err != nil {
			return nil, err
		}

		// Build each item model
		items := []models.Item{}
		for itemRows.Next() {
			var item models.Item

			err := itemRows.Scan(
				&item.ID,
				&item.ReceiptID,
				&item.Name,
				&item.Price,
				&item.Quantity,
				&item.Unit,
			)
			if err != nil {
				itemRows.Close()
				return nil, err
			}

			items = append(items, item)
		}

		if err := itemRows.Err(); err != nil {
			itemRows.Close()
			return nil, err
		}

		itemRows.Close()

		// Assign items to receipt
		receipt.Items = items
		receipts = append(receipts, receipt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return receipts, nil
}

func DeleteReceipt(ctx context.Context, userID string, receiptID string) error {
	tag, err := Pool.Exec(ctx,
		`DELETE FROM receipts
		 WHERE id = $1 AND user_id = $2`,
		receiptID,
		userID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
