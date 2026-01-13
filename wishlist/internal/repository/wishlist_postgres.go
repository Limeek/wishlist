package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
	"wishlist/internal/entity"

	"github.com/google/uuid"
)

type WishlistRepositoryPostgres struct {
	db *sql.DB
}

func NewWishlistRepositoryPostgres(db *sql.DB) *WishlistRepositoryPostgres {
	return &WishlistRepositoryPostgres{
		db: db,
	}
}

func (r *WishlistRepositoryPostgres) GetWishlists(ctx context.Context, ownerId uuid.UUID) ([]entity.Wishlist, error) {
	var wishlists = make([]entity.Wishlist, 0)

	rows, err := r.db.QueryContext(ctx, "select id, owner_id, name, expires_at from wishlists where owner_id = $1", ownerId)
	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()

	if errors.Is(err, sql.ErrNoRows) {
		return wishlists, nil
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var wishlist entity.Wishlist

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

func (r *WishlistRepositoryPostgres) Create(ctx context.Context, wishlist *entity.Wishlist) error {
	_, err := r.db.ExecContext(ctx, "insert into wishlists(id, owner_id, name, expires_at) values ($1, $2, $3, $4)",
		wishlist.Id,
		wishlist.OwnerId,
		wishlist.Name,
		wishlist.ExpiresAt,
	)

	return err
}

func (r *WishlistRepositoryPostgres) Delete(ctx context.Context, wishlistId uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, "delete from wishlists where id = $1", wishlistId)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return entity.ErrItemNotFound
	}

	return err
}

func (r *WishlistRepositoryPostgres) Update(ctx context.Context, wishlist *entity.Wishlist) error {
	_, err := r.db.ExecContext(ctx, "update wishlists set name = $1, expires_at = $2, updated_at = now() where id = $3",
		wishlist.Name,
		wishlist.ExpiresAt,
		wishlist.Id,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return entity.ErrItemNotFound
	}

	return err
}

func (r *WishlistRepositoryPostgres) GetWishlist(ctx context.Context, id uuid.UUID) (*entity.Wishlist, error) {
	var wishlist *entity.Wishlist

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
		return nil, entity.ErrWishlistNotFound
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
			wishlist = &entity.Wishlist{
				Id:        wishlistId,
				OwnerId:   ownerId,
				Name:      name,
				ExpiresAt: expiresAt,
				Items:     []entity.Item{},
			}
		}

		if itemId != nil {
			wishlist.Items = append(wishlist.Items, entity.Item{
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

func (r *WishlistRepositoryPostgres) GetWishlistOwner(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	var ownerId uuid.UUID

	err := r.db.QueryRowContext(ctx, "select owner_id from wishlists where id = $1", id).Scan(&ownerId)

	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, entity.ErrWishlistNotFound
	}

	return ownerId, err
}
