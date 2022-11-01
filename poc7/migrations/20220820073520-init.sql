-- +migrate Up
create table users
(
    id   int auto_increment primary key,

    name varchar(100) not null,
    memo varchar(100) not null default "",

    created_at datetime not null,
    updated_at datetime not null,

    constraint unique_name unique (name)
);

-- +migrate Down

drop table users;

