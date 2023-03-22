create view open_contracts_v as
(
with indexed_height as (select height
                        from indexer_status
                        limit 1)
select c.*, c.height as start_height,(select height from indexed_height) as current_height,
       c.duration-((select height from indexed_height)-c.height) as remaining
from contracts c
where c.height + c.duration > (select height from indexed_height)
);