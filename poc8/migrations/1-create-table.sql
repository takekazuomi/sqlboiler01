-- +migrate Up
create table users
(
    id         int auto_increment       primary key,
    name       varchar(100) default '' not null,

    created_at datetime(6)             not null,
    updated_at datetime(6)             not null,
    deleted_at datetime(6)                null
);

-- +migrate Down
drop table users;
