create table usernames (
        -- Max username length is 32, chosen arbitrarily. 
        -- Open to better suggestions!
        id varchar(32) primary key,
        user_id uuid not null,
        organization_id uuid not null,
        status bigint not null,
        created_at timestamp with time zone not null default now(),
        updated_at timestamp with time zone not null default now(),
        deleted_at timestamp with time zone null
    );