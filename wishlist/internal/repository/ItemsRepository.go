package repository

import (
	"context"
	"database/sql"
	"errors"
	"wishlist/internal/domain"

	"github.com/google/uuid"
)

type ItemRepositoryDb struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepositoryDb {
	return &ItemRepositoryDb{db: db}
}

func (r *ItemRepositoryDb) CreateItem(ctx context.Context, item *domain.Item, wishlistId uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)

	defer func(tx *sql.Tx) {
		if err != nil {
			_ = tx.Rollback()
		}
	}(tx)

	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "insert into items(id, name, price, image, link, description) values ($1, $2, $3, $4, $5, $6)",
		item.Id,
		item.Name,
		item.Price,
		item.ImageUrl,
		item.Link,
		item.Description,
	)

	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "insert into wishlists_items(wishlist_id, item_id) VALUES ($1, $2)",
		wishlistId,
		item.Id,
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ItemRepositoryDb) DeleteItem(ctx context.Context, itemId uuid.UUID, wishlistId uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)

	defer func(tx *sql.Tx) {
		if err != nil {
			_ = tx.Rollback()
		}
	}(tx)

	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, "delete from wishlists_items where item_id=$1 and wishlist_id = $2", itemId, wishlistId)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrItemNotFound
	}

	_, err = tx.ExecContext(ctx, "delete from items where id=$1 and not exists(select 1 from wishlists_items where item_id = $1)", itemId)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ItemRepositoryDb) UpdateItem(ctx context.Context, item *domain.Item) error {
	_, err := r.db.ExecContext(ctx, "update items set name = $1, price = $2, image = $3, link = $4, description = $5, updated_at = now() where id = $6",
		item.Name,
		item.Price,
		item.ImageUrl,
		item.Link,
		item.Description,
		item.Id)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrItemNotFound
	}

	return err
}

func (r *ItemRepositoryDb) UpdateItemReserve(ctx context.Context, itemId uuid.UUID, reservedBy uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "update items set reserved_by = $1, updated_at = now() where id = $2", reservedBy, itemId)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrItemNotFound
	}

	return err

}

func (r *ItemRepositoryDb) PutItemInWishlist(ctx context.Context, itemId uuid.UUID, wishlistId uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "insert into wishlists_items(wishlist_id, item_id) values ($1, $2) on conflict do nothing", wishlistId, itemId)

	return err
}

func (r *ItemRepositoryDb) GetItemOwnerId(ctx context.Context, itemId uuid.UUID) (uuid.UUID, error) {
	var ownerId uuid.UUID

	query := `
	select w.owner_id from items i 
    join wishlists_items wi on i.id = wi.item_id
	join wishlists w on w.id = wi.wishlist_id
	where i.id = $1
	LIMIT 1
	`

	err := r.db.QueryRowContext(ctx, query, itemId).Scan(&ownerId)

	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, domain.ErrItemNotFound
	}

	if err != nil {
		return uuid.Nil, err
	}

	return ownerId, nil
}
