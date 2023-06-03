package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/huandu/go-sqlbuilder"
	"github.com/pkg/errors"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/directory/types"
	"github.com/arkeonetwork/arkeo/directory/utils"
	"github.com/arkeonetwork/arkeo/sentinel"
	atypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
)

type ArkeoProvider struct {
	Entity  `json:"-"`
	Pubkey  string `json:"pubkey" db:"pubkey"`
	Service string `json:"service" db:"service"`
	// this is a DECIMAL type in the db
	Bond                string       `json:"bond" db:"bond"`
	MetadataURI         string       `json:"metadata_uri" db:"metadata_uri"`
	MetadataNonce       uint64       `json:"metadata_nonce" db:"metadata_nonce"`
	Status              string       `json:"status" db:"status,text"`
	MinContractDuration int64        `json:"min_contract_duration" db:"min_contract_duration"`
	MaxContractDuration int64        `json:"max_contract_duration" db:"max_contract_duration"`
	SettlementDuration  int64        `json:"settlement_duration" db:"settlement_duration"`
	SubscriptionRate    cosmos.Coins `json:"subscription_rates" db:"-"`
	PayAsYouGoRate      cosmos.Coins `json:"paygo_rates" db:"-"`
}

func (d *DirectoryDB) InsertProvider(ctx context.Context, provider *ArkeoProvider) (*Entity, error) {
	if provider == nil {
		return nil, fmt.Errorf("nil provider")
	}
	conn, err := d.getConnection(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	defer conn.Release()

	bond, err := strconv.ParseInt(provider.Bond, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "error converting bond to int64 (%s)", provider.Bond)
	}
	return insert(ctx, conn, sqlInsertProvider, provider.Pubkey, provider.Service, bond)
}

func (d *DirectoryDB) UpdateProvider(ctx context.Context, provider *ArkeoProvider) (*Entity, error) {
	if provider == nil {
		return nil, fmt.Errorf("nil provider")
	}
	conn, err := d.getConnection(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// update provide records
	var providerID int64
	var created, updated time.Time
	err = tx.QueryRow(ctx, sqlUpdateProvider,
		provider.Pubkey,
		provider.Service,
		provider.Bond,
		provider.MetadataURI,
		provider.MetadataNonce,
		provider.Status,
		provider.MinContractDuration,
		provider.MaxContractDuration,
		provider.SettlementDuration,
	).Scan(&providerID, &created, &updated)
	if err != nil {
		return nil, fmt.Errorf("fail to update provider,err: %w", err)
	}
	entity := &Entity{ID: providerID, Created: created, Updated: updated}

	// delete current subscription rate and pay-as-you-go rates before inserting new ones
	_, err = tx.Exec(ctx, sqlDeleteSubscriptionRates, provider.Pubkey, provider.Service)
	if err != nil {
		return entity, fmt.Errorf("fail to delete subscriber rate: %w", err)
	}
	_, err = tx.Exec(ctx, sqlDeletePayAsYouGoRates, provider.Pubkey, provider.Service)
	if err != nil {
		return entity, fmt.Errorf("fail to delete PayAsYouGo rate: %w", err)
	}
	if provider.SubscriptionRate.Len() > 0 {
		// insert new subscription and pay-as-you-go rates
		query, args := d.getRateArgs(providerID, sqlInsertSubscriptionRates, provider.SubscriptionRate)
		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			return entity, fmt.Errorf("fail to insert subscription rate: %w", err)
		}
	}
	if provider.PayAsYouGoRate.Len() > 0 {
		query, args := d.getRateArgs(providerID, sqlInsertPayAsYouGoRates, provider.PayAsYouGoRate)
		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			return entity, fmt.Errorf("fail to insert PayAsYouGo rate: %w", err)
		}
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	return entity, err
}

func (d *DirectoryDB) getRateArgs(providerID int64, query string, coins cosmos.Coins) (string, []interface{}) {
	var args []interface{}
	type insertRate struct {
		ProviderID  int64
		TokenName   string
		TokenAmount int64
	}
	rates := make([]insertRate, len(coins))
	for i, rate := range coins {
		rates[i] = insertRate{providerID, strings.ToLower(rate.Denom), rate.Amount.Int64()}
	}

	for i, row := range rates {
		if i > 0 {
			query += ","
		}
		query += "($1, $2, $3)"
		args = append(args, row.ProviderID, row.TokenName, row.TokenAmount)
	}
	return query, args
}

func (d *DirectoryDB) FindProvider(ctx context.Context, pubkey, service string) (*ArkeoProvider, error) {
	conn, err := d.getConnection(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	defer conn.Release()
	provider := ArkeoProvider{}
	if err = selectOne(ctx, conn, sqlFindProvider, &provider, pubkey, service); err != nil {
		return nil, errors.Wrapf(err, "error selecting")
	}

	// fetch subscription and pay-as-you-go rates
	provider.SubscriptionRate, err = d.findRates(conn, provider.ID, sqlFindProviderSubscriptionRates)
	if err != nil {
		return nil, errors.Wrapf(err, "error finding subscription rates")
	}
	provider.PayAsYouGoRate, err = d.findRates(conn, provider.ID, sqlFindProviderPayAsYouGoRates)
	if err != nil {
		return nil, errors.Wrapf(err, "error finding pay-as-you-go rates")
	}

	return &provider, nil
}

func (d *DirectoryDB) findRates(conn IConnection, providerID int64, query string) (cosmos.Coins, error) {
	// Execute the query
	ctx := context.Background()
	rows, err := conn.Query(ctx, query, providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query rates: %v", err)
	}
	defer rows.Close()

	type loadRate struct {
		ID         int64
		ProviderID int64
		Denom      string
		Amount     int64
	}

	// Iterate over the rows and store the results in a slice of slices
	results := make(cosmos.Coins, 0)
	for rows.Next() {
		// You should replace 'YourStruct' with the appropriate struct type
		// and the number of fields in the struct with the number of columns in the table
		var r loadRate
		if err := rows.Scan(&r.ID, &r.ProviderID, &r.Denom, &r.Amount); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		results = append(results, cosmos.NewInt64Coin(r.Denom, r.Amount))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to process rows: %v", err)
	}

	return results, nil
}

const provSearchCols = `
	p.id,
	p.created,
	p.pubkey,
	p.service, 
	coalesce(p.status,'OFFLINE') as status,
	coalesce(p.metadata_uri,'') as metadata_uri,
	coalesce(p.metadata_nonce,0) as metadata_nonce,
	coalesce(p.subscription_rate,0) as subscription_rate,
	coalesce(p.paygo_rate,0) as paygo_rate,
	coalesce(p.min_contract_duration,0) as min_contract_duration,
	coalesce(p.max_contract_duration,0) as max_contract_duration,
	coalesce(p.bond,0) as bond
`

func (d *DirectoryDB) SearchProviders(ctx context.Context, criteria types.ProviderSearchParams) ([]*ArkeoProvider, error) {
	conn, err := d.getConnection(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	defer conn.Release()

	sb := sqlbuilder.NewSelectBuilder()

	sb.Select(provSearchCols).
		From("providers_v p")

	// Filter
	if criteria.Pubkey != "" {
		sb = sb.Where(sb.Equal("p.pubkey", criteria.Pubkey))
	}
	if criteria.Service != "" {
		sb = sb.Where(sb.Equal("p.service", criteria.Service))
	}
	if criteria.IsMaxDistanceSet || criteria.IsMinFreeRateLimitSet || criteria.IsMinPaygoRateLimitSet || criteria.IsMinSubscribeRateLimitSet {
		sb = sb.JoinWithOption(sqlbuilder.LeftJoin, "provider_metadata", "p.id = provider_metadata.provider_id and p.metadata_nonce = provider_metadata.nonce")
	}
	if criteria.IsMaxDistanceSet {
		// note psql using long,lat instead of the normal lat,long per https://www.postgresql.org/docs/current/earthdistance.html
		sb = sb.Where(sb.LessEqualThan(fmt.Sprintf("provider_metadata.location<@>point(%.5f,%.5f)", criteria.Coordinates.Longitude, criteria.Coordinates.Latitude), criteria.MaxDistance))
	}
	if criteria.IsMinFreeRateLimitSet {
		sb = sb.Where(sb.GE("provider_metadata.free_rate_limit", criteria.MinFreeRateLimit))
	}
	if criteria.IsMinPaygoRateLimitSet {
		sb = sb.Where(sb.GE("provider_metadata.paygo_rate_limit", criteria.MinPaygoRateLimit))
	}
	if criteria.IsMinPaygoRateLimitSet {
		sb = sb.Where(sb.GE("provider_metadata.subscribe_rate_limit", criteria.MinSubscribeRateLimit))
	}
	if criteria.IsMinProviderAgeSet {
		sb = sb.Where(sb.GE("p.age", criteria.MinProviderAge))
	}
	if criteria.IsMinOpenContractsSet {
		// p.open_contract_count
		sb = sb.Where(sb.GE("p.contract_count", criteria.MinOpenContracts))
	}
	if criteria.IsMinValidatorPaymentsSet {
		sb = sb.Where(sb.GE("p.total_paid", criteria.MinValidatorPayments))
	}

	// Sort
	switch criteria.SortKey {
	case types.ProviderSortKeyNone:
		// NOP
	case types.ProviderSortKeyAge:
		sb = sb.OrderBy("p.created").Asc()
	case types.ProviderSortKeyContractCount:
		sb = sb.OrderBy("p.contract_count").Desc()
	case types.ProviderSortKeyAmountPaid:
		sb = sb.OrderBy("p.total_paid").Desc()
	default:
		return nil, fmt.Errorf("not a valid sortKey %s", criteria.SortKey)
	}

	q, params := sb.BuildWithFlavor(getFlavor())
	log.Debugf("sql: %s\n%v", q, params)

	providers := make([]*ArkeoProvider, 0, 512)
	if err := pgxscan.Select(ctx, conn, &providers, q, params...); err != nil {
		return nil, errors.Wrapf(err, "error selecting many")
	}

	return providers, nil
}

func (d *DirectoryDB) UpsertValidatorPayoutEvent(ctx context.Context, evt atypes.EventValidatorPayout, height int64) (*Entity, error) {
	conn, err := d.getConnection(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	defer conn.Release()

	return upsert(ctx, conn, sqlUpsertValidatorPayoutEvent, evt.Validator.String(), height, evt.Reward.Int64())
}

func (d *DirectoryDB) InsertBondProviderEvent(ctx context.Context, providerID int64, evt atypes.EventBondProvider, height int64, txID string) (*Entity, error) {
	if evt.BondAbs.IsNil() {
		return nil, fmt.Errorf("nil BondAbsolute")
	}
	if evt.BondRel.IsNil() {
		return nil, fmt.Errorf("nil BondRelative")
	}
	conn, err := d.getConnection(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	defer conn.Release()

	return insert(ctx, conn, sqlInsertBondProviderEvent, providerID, height, txID, evt.BondRel.String(), evt.BondAbs.String())
}

func (d *DirectoryDB) InsertModProviderEvent(ctx context.Context, providerID int64, evt types.ModProviderEvent) (*Entity, error) {
	conn, err := d.getConnection(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	defer conn.Release()

	return insert(ctx, conn, sqlInsertModProviderEvent, providerID, evt.Height, evt.TxID, evt.MetadataURI, evt.MetadataNonce, evt.Status,
		evt.MinContractDuration, evt.MaxContractDuration)
}

func (d *DirectoryDB) UpsertProviderMetadata(ctx context.Context, providerID, nonce int64, data sentinel.Metadata) (*Entity, error) {
	conn, err := d.getConnection(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	defer conn.Release()

	c := data.Configuration

	coordinates, err := utils.ParseCoordinates(c.Location)
	var location sql.NullString // using "" doesn't work here with casting to a point, only a null string ('') works with the SQL
	if err != nil {
		location = sql.NullString{Valid: false}
	} else {
		// note psql using long,lat instead of the normal lat,long per https://www.postgresql.org/docs/current/earthdistance.html
		location = sql.NullString{String: fmt.Sprintf("%.5f,%.5f", coordinates.Longitude, coordinates.Latitude), Valid: true}
	}

	// TODO - always insert instead of upsert, fail on dupe (or read and fail on exists). are there any restrictions on version string?
	return insert(ctx, conn, sqlUpsertProviderMetadata, providerID, nonce, c.Moniker, c.Website, c.Description, location, c.FreeTierRateLimit)
}
