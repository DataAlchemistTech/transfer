package redshift

import (
	"context"
	"fmt"

	"github.com/artie-labs/transfer/clients/utils"

	"github.com/artie-labs/transfer/lib/ptr"

	"github.com/artie-labs/transfer/lib/config"
	_ "github.com/lib/pq"

	"github.com/artie-labs/transfer/lib/config/constants"
	"github.com/artie-labs/transfer/lib/db"
	"github.com/artie-labs/transfer/lib/destination/types"
)

type Store struct {
	credentialsClause string
	bucket            string
	optionalS3Prefix  string
	configMap         *types.DwhToTablesConfigMap
	db.Store
}

func (s *Store) GetConfigMap() *types.DwhToTablesConfigMap {
	if s == nil {
		return nil
	}

	return s.configMap
}

func (s *Store) Label() constants.DestinationKind {
	return constants.Redshift
}

type getTableConfigArgs struct {
	Table              string
	Schema             string
	DropDeletedColumns bool
}

const (
	describeNameCol        = "column_name"
	describeTypeCol        = "data_type"
	describeDescriptionCol = "description"
)

func (s *Store) getTableConfig(ctx context.Context, args getTableConfigArgs) (*types.DwhTableConfig, error) {
	return utils.GetTableConfig(ctx, utils.GetTableCfgArgs{
		Dwh:       s,
		FqName:    fmt.Sprintf("%s.%s", args.Schema, args.Table),
		ConfigMap: s.configMap,
		// This query is a modified fork from: https://gist.github.com/alexanderlz/7302623
		Query: fmt.Sprintf(`select c.column_name,c.data_type,d.description 
from information_schema.columns c 
left join pg_class c1 on c.table_name=c1.relname 
left join pg_catalog.pg_namespace n on c.table_schema=n.nspname and c1.relnamespace=n.oid 
left join pg_catalog.pg_description d on d.objsubid=c.ordinal_position and d.objoid=c1.oid 
where c.table_name='%s' and c.table_schema='%s'`, args.Table, args.Schema),
		ColumnNameLabel:    describeNameCol,
		ColumnTypeLabel:    describeTypeCol,
		ColumnDescLabel:    describeDescriptionCol,
		EmptyCommentValue:  ptr.ToString("<nil>"),
		DropDeletedColumns: args.DropDeletedColumns,
	})
}

func LoadRedshift(ctx context.Context, _store *db.Store) *Store {
	if _store != nil {
		// Used for tests.
		return &Store{
			Store:     *_store,
			configMap: &types.DwhToTablesConfigMap{},
		}
	}

	settings := config.FromContext(ctx)
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		settings.Config.Redshift.Host, settings.Config.Redshift.Port, settings.Config.Redshift.Username,
		settings.Config.Redshift.Password, settings.Config.Redshift.Database)

	return &Store{
		credentialsClause: settings.Config.Redshift.CredentialsClause,
		bucket:            settings.Config.Redshift.Bucket,
		optionalS3Prefix:  settings.Config.Redshift.OptionalS3Prefix,
		Store:             db.Open(ctx, "postgres", connStr),
		configMap:         &types.DwhToTablesConfigMap{},
	}
}
