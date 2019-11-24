package date

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateStartAndEnd(t *testing.T) {
	days := 7
	startDate, endDate := GenerateStartAndEnd(days)
	start, err := time.Parse(layoutISO, startDate)
	assert.NoError(t, err)

	end, err := time.Parse(layoutISO, endDate)
	assert.NoError(t, err)

	daysDiff := end.Sub(start).Hours() / 24
	assert.Equal(t, days, int(daysDiff), "Days difference is wrong")
}
