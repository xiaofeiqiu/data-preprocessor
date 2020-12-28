package dbservice

import (
	"sort"
	"time"
)

type timeSlice []time.Time

func (s timeSlice) Less(i, j int) bool { return s[i].Before(s[j]) }
func (s timeSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s timeSlice) Len() int           { return len(s) }

// return max date and min-20 day
// min-20 is used to calculate diff
func GetMinAndMaxDate(data []DataInputEntity) (string, string) {
	ts := timeSlice{}
	for _, v := range data {
		ts = append(ts, v.Date)
	}
	sort.Sort(ts)

	ts[0] = ts[0].AddDate(0, 0, -20)

	return ts[0].Format("2006-01-02"), ts[len(ts)-1].Format("2006-01-02")
}


