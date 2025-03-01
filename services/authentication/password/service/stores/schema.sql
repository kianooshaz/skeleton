create table passwords (
        id         uuid primary key,
        user_id uuid not null,
        password_hash text not null,
        created_at timestamp with time zone not null default now(),
        deleted_at timestamp with time zone null
    );