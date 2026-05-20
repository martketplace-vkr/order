alter table "order"."order"
    add column if not exists payment_status text not null default 'pending_funds',
    add column if not exists fulfillment_status text not null default 'created',
    add column if not exists delivery_address_id bigint not null default 0;

update "order"."order"
set
    fulfillment_status = case
        when status in ('created', 'waiting_for_payment') then 'created'
        else status
    end,
    payment_status = case
        when status = 'success' then 'captured'
        when status in ('cancelled_by_client', 'cancelled_by_seller', 'cancelled') then 'cancelled'
        else 'pending_funds'
    end
where payment_status = 'pending_funds'
  and fulfillment_status = 'created';

update "order"."order" o
set delivery_address_id = d.entity_id
from "order".delivery d
where d.order_id = o.id
  and d.type = 2
  and o.delivery_address_id = 0;
