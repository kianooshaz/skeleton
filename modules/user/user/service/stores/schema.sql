create table users
(
    id         uuid primary key,
    created_at timestamp with time zone not null default now()
);