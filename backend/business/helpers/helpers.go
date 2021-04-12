package helpers

import "time"

// GetEarliestTime gives the earliest time Golang has to offer.
// Reference: https://stackoverflow.com/questions/23051973/what-is-the-zero-value-for-time-time-in-go
func GetEarliestTime() *time.Time {
	t := new(time.Time)
	return t
}

// GetMaximalTime gives the maximal time Golang has to offer.
// Reference: https://stackoverflow.com/questions/25065055/what-is-the-maximum-time-time-in-go
func GetMaximalTime() *time.Time {
	t := time.Unix(1<<63-62135596801, 999999999)
	return &t
}

func GetEarliestEpoch() int64 {
	return 0
}

// Reference: https://medium.com/@nate510/the-2286-bug-65697bb1b908#:~:text=A%20better%2Dknown%20and%20p,signed%2C%2032%2Dbit%20integer.
func GetMaximalEpoch() int64 {
	return 2147483647
}
