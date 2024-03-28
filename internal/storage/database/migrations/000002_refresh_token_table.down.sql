alter table refresh_tokens
drop constraint if exists fk_refresh_tokens_parent_token_id;

alter table refresh_tokens
drop constraint if exists fk_refresh_tokens_user_id;

drop table refresh_tokens;
