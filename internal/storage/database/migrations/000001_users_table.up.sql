create table users (
    id uuid primary key default gen_random_uuid(),
    username varchar(50) unique not null,
    email varchar(100) unique not null,
    password varchar(100) not null,
    created_at timestamp  with time zone not null default (current_timestamp at time zone 'utc'), 
    updated_at timestamp with time zone not null default (current_timestamp at time zone 'utc'),
    deleted_at timestamp with time zone 
)
