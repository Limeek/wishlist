create table wishlists
(
    id         uuid primary key,
    owner_id   uuid,
    name       text,
    expires_at timestamptz default now(),
    created_at timestamptz default now(),
    updated_at timestamptz default now()
);

create table items
(
    id          uuid primary key,
    name        text,
    price       numeric(12, 2),
    image       text,
    link        text,
    description text,
    reserved_by uuid,
    owner_id    uuid,
    created_at  timestamptz default now(),
    updated_at  timestamptz default now()
);

create table wishlists_items
(
    wishlist_id uuid references wishlists (id),
    item_id     uuid references items (id),
    primary key (wishlist_id, item_id)
);

CREATE INDEX idx_wishlists_items_item_id
    ON wishlists_items (item_id);