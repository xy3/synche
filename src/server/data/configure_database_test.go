package data_test

import (
	"github.com/stretchr/testify/require"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"testing"
)

func TestNewDSN(t *testing.T) {
	tests := []struct {
		name string
		dbConfig config.DatabaseConfig
		wantDsn string
	}{
		{
			name: "dsn with empty config",
			dbConfig: config.DatabaseConfig{},
			wantDsn: ":@()/?charset=utf8mb4&parseTime=True&loc=Local",
		},
		{
			name: "dsn with default config",
			dbConfig: config.Config.Database,
			wantDsn: "root:@tcp(127.0.0.1:3306)/synche?charset=utf8mb4&parseTime=True&loc=Local",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotDsn := data.NewDSN(tc.dbConfig)
			require.Equal(t, tc.wantDsn, gotDsn)
		})
	}
}
