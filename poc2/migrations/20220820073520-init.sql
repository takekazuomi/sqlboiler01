-- +migrate Up
create table users
(
    id   int auto_increment primary key,

    name varchar(100) not null,
    memo varchar(100) not null default "",

    constraint unique_name unique (name)
);

create table videos
(
    id       int auto_increment  primary key,
    user_id int                  not null,

    name     varchar(100)         not null,

    deleted boolean default false null,

    constraint fk_user_id foreign key (user_id) references users (id)
);

create table tags
(
    id   int auto_increment primary key,

    name varchar(100) not null,

    deleted  boolean default false not null
);

create table tags_videos
(
    video_id int          not null,
    tag_id   int          not null,

    primary key (video_id, tag_id),
    constraint fk_tag_id foreign key (tag_id) references tags (id),
    constraint fk_video_id foreign key (video_id) references videos (id)
);

-- +migrate Down

drop table tags_videos;
drop table tags;
drop table videos;
drop table users;

