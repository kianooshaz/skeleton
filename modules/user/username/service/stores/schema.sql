create table usernames (
        id uuid primary key,
        username_value varchar(50) not null,
        user_id uuid not null,
        status bigint not null,
        is_primary boolean not null,
        created_at timestamp with time zone not null default now(),
        updated_at timestamp with time zone not null default now(),
        deleted_at timestamp with time zone null
    );