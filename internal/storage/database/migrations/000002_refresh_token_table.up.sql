create table refresh_tokens (
    id serial primary key,
    parent_token_id integer,
    token varchar(255) not null,
    user_id integer not null,
    expires_at timestamp with time zone not null,
    used boolean default false,
    created_at timestamp with time zone not null default (current_timestamp at time zone 'utc')
);

alter table refresh_tokens 
add constraint fk_refresh_tokens_parent_token_id
foreign key (parent_token_id)
references refresh_tokens(id);

alter table refresh_tokens
add constraint fk_refresh_tokens_user_id
foreign key (user_id)
references users(id);
