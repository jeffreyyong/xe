package date

import "time"

const (
	layoutISO = "2006-01-02"
)

// GenerateStartAndEnd generates the start and end date
// in ISO string format given the number of days
// Note: end date is always the current date
func GenerateStartAndEnd(days int) (string, string) {
	end := time.Now()
	start := end.AddDate(0, 0, -days)

	endDate := end.Format(layoutISO)
	startDate := start.Format(layoutISO)

	return startDate, endDate
}
