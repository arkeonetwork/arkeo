package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/arkeonetwork/arkeo/directory/types"
	"github.com/arkeonetwork/arkeo/directory/utils"
	"github.com/arkeonetwork/arkeo/sentinel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/huandu/go-sqlbuilder"
	"github.com/pkg/errors"
)

type ArkeoProvider struct {
	Entity  `json:"-"`
	Pubkey  string `json:"pubkey" db:"pubkey"`
	Service string `json:"service" db:"service"`
	// this is a DECIMAL type in the db
	Bond                string               `json:"bond" db:"bond"`
	MetadataURI         string               `json:"metadata_uri" db:"metadata_uri"`
	MetadataNonce       uint64               `json:"metadata_nonce" db:"metadata_nonce"`
	Status              types.ProviderStatus `json:"status" db:"status,text"`
	MinContractDuration int64                `json:"min_contract_duration" db:"min_contract_duration"`
	MaxContractDuration int64                `json:"max_contract_duration" db:"max_contract_duration"`
	SubscriptionRate    int64                `json:"subscription_rate" db:"subscription_rate"`
	PayAsYouGoRate      int64                `json:"paygo_rate" db:"paygo_rate"`
}

func (d *DirectoryDB) InsertProvider(provider *ArkeoProvider) (*Entity, error) {
	if provider == nil {
		return nil, fmt.Errorf("nil provider")
	}
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	bond, err := strconv.ParseInt(provider.Bond, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "error converting bond to int64 (%s)", provider.Bond)
	}
	return insert(conn, sqlInsertProvider, provider.Pubkey, provider.Service, bond)
}

func (d *DirectoryDB) UpdateProvider(provider *ArkeoProvider) (*Entity, error) {
	if provider == nil {
		return nil, fmt.Errorf("nil provider")
	}
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return update(conn,
		sqlUpdateProvider,
		provider.Pubkey,
		provider.Service,
		provider.Bond,
		provider.MetadataURI,
		provider.MetadataNonce,
		provider.Status,
		provider.MinContractDuration,
		provider.MaxContractDuration,
		provider.SubscriptionRate,
		provider.PayAsYouGoRate,
	)
}

func (d *DirectoryDB) FindProvider(pubkey, service string) (*ArkeoProvider, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	provider := ArkeoProvider{}
	if err = selectOne(conn, sqlFindProvider, &provider, pubkey, service); err != nil {
		return nil, errors.Wrapf(err, "error selecting")
	}
	// not found
	if provider.Pubkey == "" {
		return nil, nil
	}
	return &provider, nil
}

const provSearchCols = `
	p.id,
	p.created,
	p.pubkey,
	p.service, 
	coalesce(p.status,'Offline') as status,
	coalesce(p.metadata_uri,'') as metadata_uri,
	coalesce(p.metadata_nonce,0) as metadata_nonce,
	coalesce(p.subscription_rate,0) as subscription_rate,
	coalesce(p.paygo_rate,0) as paygo_rate,
	coalesce(p.min_contract_duration,0) as min_contract_duration,
	coalesce(p.max_contract_duration,0) as max_contract_duration,
	coalesce(p.bond,0) as bond
`

func (d *DirectoryDB) SearchProviders(criteria types.ProviderSearchParams) ([]*ArkeoProvider, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

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

	sql, params := sb.BuildWithFlavor(getFlavor())
	log.Debugf("sql: %s\n%v", sql, params)

	providers := make([]*ArkeoProvider, 0, 512)
	if err := pgxscan.Select(context.Background(), conn, &providers, sql, params...); err != nil {
		return nil, errors.Wrapf(err, "error selecting many")
	}

	return providers, nil
}

func (d *DirectoryDB) UpsertValidatorPayoutEvent(evt types.ValidatorPayoutEvent) (*Entity, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return upsert(conn, sqlUpsertValidatorPayoutEvent, evt.Validator, evt.Height, evt.Paid)
}

func (d *DirectoryDB) InsertBondProviderEvent(providerID int64, evt types.BondProviderEvent) (*Entity, error) {
	if evt.BondAbsolute == "" {
		return nil, fmt.Errorf("nil BondAbsolute")
	}
	if evt.BondRelative == "" {
		return nil, fmt.Errorf("nil BondRelative")
	}
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return insert(conn, sqlInsertBondProviderEvent, providerID, evt.Height, evt.TxID, evt.BondRelative, evt.BondAbsolute)
}

func (d *DirectoryDB) InsertModProviderEvent(providerID int64, evt types.ModProviderEvent) (*Entity, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return insert(conn, sqlInsertModProviderEvent, providerID, evt.Height, evt.TxID, evt.MetadataURI, evt.MetadataNonce, evt.Status,
		evt.MinContractDuration, evt.MaxContractDuration, evt.SubscriptionRate, evt.PayAsYouGoRate)
}

func (d *DirectoryDB) UpsertProviderMetadata(providerID, nonce int64, data sentinel.Metadata) (*Entity, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

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
	return insert(conn, sqlUpsertProviderMetadata, providerID, nonce, c.Moniker, c.Website, c.Description, location, c.FreeTierRateLimit)
}
