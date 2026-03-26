package goladok3

import (
	"context"
	"testing"

	"github.com/SUNET/goladok3/ladoktypes"
	"github.com/stretchr/testify/assert"
)

func TestGetTillfallesdeltagandePagaendeStudent(t *testing.T) {
	client := mockNewClient(t, ladoktypes.EnvIntTestAPI, "https://example.com")
	result, resp, err := client.Studentdeltagande.GetTillfallesdeltagandePagaendeStudent(context.Background(), GetAktivPaLarosateReq{UID: "test"})
	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}
