-- +migrate Up
create table table1 (
    id int auto_increment primary key,
    status enum ('apple', 'orange', 'mango') not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp null,

    alive bool generated always as (
        case when `deleted_at` is null then 1 else null end
    ) virtual,

    constraint patient_status_alive unique (status, alive)
);

-- +migrate Down
drop table table1;