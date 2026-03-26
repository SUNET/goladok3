package ladoktypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenderString(t *testing.T) {
	tts := []struct {
		name string
		have *Student
		want string
	}{
		{name: "female", have: &Student{KonID: 1}, want: "female"},
		{name: "male", have: &Student{KonID: 2}, want: "male"},
		{name: "n/a", have: &Student{KonID: 10}, want: "n/a"},
		{name: "zero", have: &Student{KonID: 0}, want: "n/a"},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.have.GenderString())
		})
	}
}
