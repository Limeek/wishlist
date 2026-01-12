create table user_profiles(
    user_id uuid primary key,
    first_name text,
    last_name text,
    email text,
    phone text,
    avatar text,
    birth_date date,
    created_at timestamptz default now(),
    updated_at timestamptz default now()
);