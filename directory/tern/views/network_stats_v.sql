/*
 A stats endpoint would give the following stats

- number of open contracts
- total number of contracts
- median contract length of open contracts
- median contract rate of open contracts (pay-as-you-go and subscription)
- number of online providers, and per service
- total number of queries (sum of nonces), and per service
- total number of queries in the last 24hrs, and per service
- total income, and per service
- total income in the last 24hrs and per service
 */

create or replace view network_stats_v as
(
select (select count(1) from contracts)                                 as total_contracts,
       (select count(1) from contracts c where c.closed_height = 0)     as open_contracts,
       (SELECT percentile_cont(0.5) within group (order by duration) -- percentile_disc
        from contracts c
        where c.closed_height = 0)                                      as median_open_contract_length,
       (SELECT percentile_cont(0.5) within group (order by rate)
        from contracts c
        where c.closed_height = 0)                                      as median_open_contract_rate,
       (select count(1) from providers where status = 'Online')         as total_online_providers,
       (select coalesce(sum(nonce), 0) from contract_settlement_events) as total_queries, -- nonce here is serviced request count
       (select coalesce(sum(paid), 0) from contract_settlement_events)  as total_paid );
