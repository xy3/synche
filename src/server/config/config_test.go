package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/setup"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	config.TestMode = true
	err := setup.Dirs(files.AppFS, c.RequiredDirs())
	if err != nil {
		log.WithError(err).Fatal("Could not set up the required directories")
	}

	os.Exit(m.Run())
}

func TestInitConfig(t *testing.T) {
	type args struct {
		cfgFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty config path string", args: args{""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitConfig(tt.args.cfgFile); (err != nil) != tt.wantErr {
				require.NotNil(t, Config.Synche.Dir)
				t.Errorf("InitConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewDSN(t *testing.T) {
	tests := []struct {
		name     string
		dbConfig DatabaseConfig
		wantDsn  string
	}{
		{
			name:     "dsn with empty config",
			dbConfig: DatabaseConfig{},
			wantDsn:  ":@()/?charset=utf8mb4&parseTime=True&loc=Local",
		},
		{
			name:     "dsn with default config",
			dbConfig: DatabaseConfig{
				Name:     "synche",
				Username: "root",
				Password: "test123",
				Protocol: "tcp",
				Address:  "127.0.0.1:3306",
			},
			wantDsn:  "root:test123@tcp(127.0.0.1:3306)/synche?charset=utf8mb4&parseTime=True&loc=Local",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotDsn := tc.dbConfig.DSN()
			require.Equal(t, tc.wantDsn, gotDsn)
		})
	}
}
