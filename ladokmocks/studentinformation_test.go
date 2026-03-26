package ladokmocks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockStudentinformationStudent(t *testing.T) {
	s := MockStudentinformationStudent()
	assert.NotNil(t, s)
	assert.Equal(t, "TestFornamn", s.Fornamn)
	assert.Equal(t, "TestEfternamn", s.Efternamn)
}

func TestStudentJSON(t *testing.T) {
	b := StudentJSON(Students[0])
	assert.NotNil(t, b)
	assert.Contains(t, string(b), Students[0].Personnummer)
	assert.Contains(t, string(b), Students[0].StudentUID)
}

func TestStudentJSONAllStudents(t *testing.T) {
	for i, s := range Students {
		b := StudentJSON(s)
		assert.NotNilf(t, b, "Student %d should produce valid JSON", i)
	}
}
