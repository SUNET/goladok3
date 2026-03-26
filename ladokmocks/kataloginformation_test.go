package ladokmocks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockKataloginformationProfil(t *testing.T) {
	s := MockKataloginformationProfil()
	assert.NotNil(t, s)
	assert.Equal(t, 96, s.LarosateID)
}

func TestMockKataloginformationBehorighetsprofil(t *testing.T) {
	s := MockKataloginformationBehorighetsprofil()
	assert.NotNil(t, s)
	assert.Equal(t, 27, s.LarosateID)
	assert.NotEmpty(t, s.Systemaktiviteter)
}

func TestMockKataloginformationAutentiserad(t *testing.T) {
	s := MockKataloginformationAutentiserad()
	assert.NotNil(t, s)
	assert.Equal(t, "mail@school.se", s.Anvandarnamn)
	assert.Equal(t, 96, s.LarosateID)
}

func TestMockKataloginformationEgna(t *testing.T) {
	s := MockKataloginformationEgna()
	assert.NotNil(t, s)
	assert.Equal(t, 96, s.LarosateID)
}
