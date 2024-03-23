
-- +migrate Up
alter table keywords add column id uuid default gen_random_uuid() not null unique;
alter table user_keywords add column keyword_id uuid references keywords(id) on delete cascade;
update user_keywords set keyword_id = (select id from keywords where user_keywords.keyword = keywords.name);
alter table user_keywords alter column keyword_id set not null;
alter table user_keywords drop constraint user_keywords_pkey;
alter table user_keywords add primary key (user_id, keyword_id);
alter table user_keywords drop column keyword;
alter table keywords drop constraint keywords_pkey;
alter table keywords drop constraint if exists keywords_name_key;
alter table keywords add primary key (id);
create extension citext;
alter table users alter column email type citext;

-- +migrate Down
alter table users alter column email type varchar(200);
drop extension citext;
alter table keywords drop constraint keywords_pkey;
delete from keywords a using keywords b where a.id <> b.id and a.name = b.name;
alter table keywords add primary key (name);
alter table user_keywords add column keyword text references keywords (name) on delete cascade;
alter table user_keywords drop constraint user_keywords_pkey;
update user_keywords set keyword = (select name from keywords where user_keywords.keyword_id = keywords.id);
alter table user_keywords alter column keyword set not null;
alter table user_keywords add constraint user_keywords_pkey primary key (user_id, keyword);
alter table user_keywords drop column keyword_id;
alter table keywords drop column id;
