create schema if not exists "order";

create table if not exists "order"."order" (
    id bigserial primary key,
    checkout_id text not null,
    user_id bigint not null,
    vendor_id bigint not null,
    status text not null,
    payment_status text not null default 'pending_funds',
    fulfillment_status text not null default 'created',
    delivery_address_id bigint not null default 0,
    product_id bigint not null,
    product_name text not null,
    product_image_url text not null,
    quantity integer not null,
    unit_price text not null,
    total_price text not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    unique (user_id, checkout_id, product_id)
);

create table if not exists "order".payment (
    order_id bigint not null,
    payment_id integer,
    type integer not null,
    status integer not null,

    constraint pk_payment primary key (order_id),
    constraint fk_payment_order
        foreign key (order_id)
        references "order"."order" (id)
        on delete cascade
);

create table if not exists "order".delivery (
    order_id bigint not null,
    entity_id integer not null, -- pick_up_id or client_address_id
    type integer not null,
    status integer not null,

    constraint pk_delivery primary key (order_id),
    constraint fk_delivery_order
        foreign key (order_id)
        references "order"."order" (id)
        on delete cascade
);
