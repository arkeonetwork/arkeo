package db

var sqlGetNetworkStats = `select * from network_stats_v limit 1`

var sqlGetNetworkStatsByService = `
	 SELECT 
		(
			SELECT count(1) AS count
			FROM open_contracts_v c1
			INNER JOIN providers_v p1 ON c1.provider_id = p1.id 
			WHERE p1.service = $1
		) as open_contracts,
		( 
			SELECT count(1) AS count
			FROM contracts c2
			INNER JOIN providers_v p2 ON c2.provider_id = p2.id 
			WHERE p2.service = $1
		) AS total_contracts,
		(
			SELECT percentile_cont((0.5)::double precision) WITHIN GROUP (ORDER BY ((c3.duration)::double precision)) AS percentile_cont
			FROM contracts c3
			INNER JOIN providers p3 ON c3.provider_id = p3.id
			WHERE ((c3.height + c3.duration + c3.settlement_duration) > c3.settlement_height)
			AND (p3.service = $1)
		) as median_open_contract_length,
		(
			SELECT percentile_cont((0.5)::double precision) WITHIN GROUP (ORDER BY ((c4.rate_amount)::double precision)) AS percentile_cont
			FROM contracts c4
			INNER JOIN providers p4 ON c4.provider_id = p4.id
			WHERE ((c4.height + c4.duration + c4.settlement_duration) > c4.settlement_height)
			AND (p4.service = $1)
		) as median_open_contract_rate,
		(
			SELECT count(1) AS count
			FROM providers p5
			WHERE (p5.status = 'ONLINE'::text)
			AND (p5.service = $1)
		) as total_online_providers,
		(
			SELECT COALESCE(sum(cse1.nonce), (0)::numeric) AS "coalesce"
			FROM contract_settlement_events cse1
			INNER JOIN contracts c5 ON cse1.contract_id = c5.id
			INNER JOIN providers p6 ON c5.provider_id = p6.id
			WHERE (p6.status = 'ONLINE'::text)
			AND (p6.service = $1)
		) as total_queries,
		(
			SELECT COALESCE(sum(cse2.paid), (0)::numeric) AS "coalesce"
			FROM contract_settlement_events cse2
			INNER JOIN contracts c6 ON cse2.contract_id = c6.id
			INNER JOIN providers p7 ON c6.provider_id = p7.id
			WHERE (p7.status = 'ONLINE'::text)
			AND (p7.service = $1)
		) as total_paid
	`
