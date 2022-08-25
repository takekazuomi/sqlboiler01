
-- +migrate Up
insert into users (name) values ('john');
insert into users (name) values ('wayne');

insert into videos (user_id, name) values (1, 'batman begins');
insert into videos (user_id, name) values (1, 'dark knight');
insert into videos (user_id, name) values (1, 'dark knight rises');
insert into videos (user_id, name) values (1, 'xman');
insert into videos (user_id, name) values (1, 'matrix');

insert into videos (user_id, name) values (2, 'starship rroopers');
insert into videos (user_id, name) values (2, 'other');

insert into tags (name, deleted) values ('action', false);
insert into tags (name, deleted) values ('drana', false);

insert into tags_videos (video_id, tag_id) values (1,1);
insert into tags_videos (video_id, tag_id) values (1,2);
insert into tags_videos (video_id, tag_id) values (2,2);

-- +migrate Down

delete from tags_videos;
delete from videos;
delete from tags;
delete from users;

alter table tags_videos auto_increment = 0;
alter table videos auto_increment = 0;
alter table tags auto_increment = 0;
alter table users auto_increment = 0;
