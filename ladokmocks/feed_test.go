package ladokmocks

import (
	"encoding/json"
	"testing"

	"github.com/SUNET/goladok3/ladoktypes"
	"github.com/stretchr/testify/assert"
)

func TestMockSuperFeed(t *testing.T) {
	assert.Equal(t, 10, MockSuperFeed(10).ID)
	assert.Len(t, MockSuperFeed(1).SuperEvents, 7)
}

func TestJSONSuperFeed(t *testing.T) {
	b := JSONSuperFeed(42)
	assert.NotNil(t, b)

	var sf ladoktypes.SuperFeed
	err := json.Unmarshal(b, &sf)
	assert.NoError(t, err)
	assert.Equal(t, 42, sf.ID)
}

func TestFeedXML(t *testing.T) {
	b := FeedXML(100)
	assert.NotNil(t, b)
	assert.Contains(t, string(b), "100")
}
