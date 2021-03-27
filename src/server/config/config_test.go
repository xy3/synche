package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

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
			"empty config path string", args{""}, false,
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
