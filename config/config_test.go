package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ConfigValidate(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		cfg *Config
		err error
	}{
		"valid deployment network config (mainnet) should return no errors": {
			cfg: &Config{
				Deployment: &Deployment{
					Network: NetworkMainet,
				},
			},
			err: nil,
		}, "valid deployment network config (testnet) should return no errors": {
			cfg: &Config{
				Deployment: &Deployment{
					Network: NetworkTestnet,
				},
			},
			err: nil,
		}, "valid deployment network config (stn) should return no errors": {
			cfg: &Config{
				Deployment: &Deployment{
					Network: NetworkSTN,
				},
			},
			err: nil,
		}, "valid deployment network config (regtest) should return no errors": {
			cfg: &Config{
				Deployment: &Deployment{
					Network: NetworkRegtest,
				},
			},
			err: nil,
		}, "deployment network type within other word should fail": {
			cfg: &Config{
				Deployment: &Deployment{
					Network: "btestneth",
				},
			},
			err: errors.New("[deployment.network: value btestneth failed to meet requirements]"),
		}, "invalid deployment network config should error": {
			cfg: &Config{
				Deployment: &Deployment{
					Network: "blah",
				},
			},
			err: errors.New("[deployment.network: value blah failed to meet requirements]"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.cfg.Validate()
			if test.err == nil {
				assert.NoError(t, err)
				return
			}
			assert.EqualError(t, err, test.err.Error())

		})
	}
}
