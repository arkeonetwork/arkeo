create or replace view open_contracts_v as
(
with indexed_height as (select height
                        from indexer_status
                        limit 1)
select c.*,
       c.height as start_height,
       (select height from indexed_height) as current_height,
       case when c.closed_height = 0 then
         c.duration-((select height from indexed_height)-c.height)
       else 0
       end as remaining
from contracts c
where c.closed_height = 0
  and c.height + c.duration > (select height from indexed_height)
);
