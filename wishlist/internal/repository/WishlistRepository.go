package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
	"wishlist/internal/domain"

	"github.com/google/uuid"
)

type WishlistRepositoryDb struct {
	db *sql.DB
}

func NewWishlistRepositoryDb(db *sql.DB) *WishlistRepositoryDb {
	return &WishlistRepositoryDb{
		db: db,
	}
}

func (r *WishlistRepositoryDb) GetWishlists(ctx context.Context, ownerId uuid.UUID) ([]domain.Wishlist, error) {
	var wishlists = make([]domain.Wishlist, 0)

	rows, err := r.db.QueryContext(ctx, "select id, owner_id, name, expires_at from wishlists where owner_id = $1", ownerId)
	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()

	//if errors.Is(err, sql.ErrNoRows) {
	//	return wishlists, nil
	//}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var wishlist domain.Wishlist

		err = rows.Scan(
			&wishlist.Id,
			&wishlist.OwnerId,
			&wishlist.Name,
			&wishlist.ExpiresAt,
		)

		if err != nil {
			log.Default().Printf("error scanning row: %v", err)
			return nil, err
		}

		wishlists = append(wishlists, wishlist)
	}
	if err := rows.Err(); err != nil {
		log.Default().Printf("error scanning rows: %v", err)
		return nil, err
	}

	return wishlists, nil

}

func (r *WishlistRepositoryDb) Create(ctx context.Context, wishlist *domain.Wishlist) error {
	_, err := r.db.ExecContext(ctx, "insert into wishlists(id, owner_id, name, expires_at) values ($1, $2, $3, $4)",
		wishlist.Id,
		wishlist.OwnerId,
		wishlist.Name,
		wishlist.ExpiresAt,
	)

	return err
}

func (r *WishlistRepositoryDb) Delete(ctx context.Context, wishlistId uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, "delete from wishlists where id = $1", wishlistId)

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

	return err
}

func (r *WishlistRepositoryDb) Update(ctx context.Context, wishlist *domain.Wishlist) error {
	_, err := r.db.ExecContext(ctx, "update wishlists set name = $1, expires_at = $2, updated_at = now() where id = $3",
		wishlist.Name,
		wishlist.ExpiresAt,
		wishlist.Id,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrItemNotFound
	}

	return err
}

func (r *WishlistRepositoryDb) GetWishlist(ctx context.Context, id uuid.UUID) (*domain.Wishlist, error) {
	var wishlist *domain.Wishlist

	query := `
	select w.id, w.owner_id, w.name, w.expires_at, i.id, i.name, i.link, i.image, i.price, i.description, i.reserved_by from wishlists w 
	left join wishlists_items wi on w.id = wi.wishlist_id
	left join items i on wi.item_id = i.id
	where w.id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, id)

	defer func() {
		_ = rows.Close()
	}()

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrWishlistNotFound
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {

		var (
			wishlistId uuid.UUID
			ownerId    uuid.UUID
			name       string
			expiresAt  time.Time

			itemId          *uuid.UUID
			itemName        *string
			itemLink        *string
			itemImage       *string
			itemPrice       *float64
			itemDescription *string
			itemReservedBy  *uuid.UUID
		)

		err := rows.Scan(
			&wishlistId,
			&ownerId,
			&name,
			&expiresAt,
			&itemId,
			&itemName,
			&itemLink,
			&itemImage,
			&itemPrice,
			&itemDescription,
			&itemReservedBy,
		)

		if err != nil {
			log.Default().Printf("error scanning row: %v", err)
			return nil, err
		}

		if wishlist == nil {
			wishlist = &domain.Wishlist{
				Id:        wishlistId,
				OwnerId:   ownerId,
				Name:      name,
				ExpiresAt: expiresAt,
				Items:     []domain.Item{},
			}
		}

		if itemId != nil {
			wishlist.Items = append(wishlist.Items, domain.Item{
				Id:          *itemId,
				Name:        *itemName,
				Link:        *itemLink,
				Price:       *itemPrice,
				Description: *itemDescription,
				ReservedBy:  itemReservedBy,
				ImageUrl:    *itemImage,
			})
		}
	}

	if err := rows.Err(); err != nil {
		log.Default().Printf("error scanning rows: %v", err)
		return nil, err
	}

	return wishlist, err
}

func (r *WishlistRepositoryDb) GetWishlistOwner(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	var ownerId uuid.UUID

	err := r.db.QueryRowContext(ctx, "select owner_id from wishlists where id = $1", id).Scan(&ownerId)

	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, domain.ErrWishlistNotFound
	}

	return ownerId, err
}
