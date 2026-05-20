alter table "order"."order"
    drop column if exists delivery_address_id,
    drop column if exists fulfillment_status,
    drop column if exists payment_status;
