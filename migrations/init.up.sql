create table if not exists orders (
    id bigserial primary key,
    checkout_id text not null,
    user_id bigint not null,
    vendor_id bigint not null,
    status text not null,
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
