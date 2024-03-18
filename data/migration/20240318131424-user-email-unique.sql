
-- +migrate Up
alter table users add constraint user_email_unique unique (email);

-- +migrate Down
alter table users drop constraint user_email_unique;