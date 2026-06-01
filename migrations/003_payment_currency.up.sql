alter table "order".payment
    add column if not exists currency_id bigint not null default 1000;
