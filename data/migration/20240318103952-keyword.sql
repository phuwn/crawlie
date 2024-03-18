
-- +migrate Up
create table keywords (
	name text not null unique primary key,
	created_at timestamp(6) with time zone default now(),
	updated_at timestamp(6) with time zone,
	deleted_at timestamp(6) with time zone,
	ad_words_count integer default 0,
	links_count integer default 0,
	search_results_count text,
	status integer default 0,
	html_cache text,
	last_crawled_at timestamp(6)
);

create table user_keywords (
	user_id uuid references users (id) on delete cascade not null,
	keyword text references keywords (name) on delete cascade not null,
	file_name varchar(200),
	created_at timestamp(6) with time zone default now(),
	updated_at timestamp(6) with time zone,
	deleted_at timestamp(6) with time zone,
	primary key(user_id, keyword)
);

-- +migrate Down
drop table user_keywords;
drop table keywords;