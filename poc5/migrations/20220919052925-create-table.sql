-- +migrate Up
create table table1
(
    id varbinary(16) not null primary key,
    num int not null,
    f1 varchar(100) default '' not null,
    constraint table1_num_uindex unique (num)
);

create table table2
(
    id varbinary(16) not null primary key,
    table1_id binary(16) not null,
    f2 varchar(100) default '' not null,

    constraint table2_table1_id_fk foreign key (table1_id) references table1 (id)
);

-- +migrate Down
drop table table2;
drop table table1;
