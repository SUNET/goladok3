package ladokmocks

import (
	"testing"

	"github.com/SUNET/goladok3/ladoktypes"
	"github.com/stretchr/testify/assert"
)

func TestMockCertificateAndKey(t *testing.T) {
	envs := []string{ladoktypes.EnvIntTestAPI, ladoktypes.EnvTestAPI, ladoktypes.EnvProdAPI}

	for _, env := range envs {
		t.Run(env, func(t *testing.T) {
			certPEM, cert, keyPEM, key := MockCertificateAndKey(t, env, 0, 100)
			assert.NotNil(t, certPEM)
			assert.NotNil(t, cert)
			assert.NotNil(t, keyPEM)
			assert.NotNil(t, key)
			assert.Equal(t, env, cert.Subject.OrganizationalUnit[1])
		})
	}
}
