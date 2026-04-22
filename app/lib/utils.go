package lib

import (
	"strconv"
	"strings"
)

func GetIdTimeSequence(id string) (int64, int64) {
	time_seq := strings.Split(id, "-")
	time, _ := strconv.Atoi(time_seq[0])

	seq, _ := strconv.Atoi(time_seq[1])

	return int64(time), int64(seq)
}

func IsValidTimeSequence(prev_time int64, prev_seq int64, new_time int64, new_seq int64) bool {
	if new_time == prev_time {
		return new_seq > prev_seq
	} else {
		return new_time > prev_time
	}
}

func IsAutoSequence(id string) bool {
	time_seq := strings.Split(id, "-")
	seq := time_seq[1]

	return seq == "*"
}

func IsFullAutoId(id string) bool {
	return id == "*"
}

func IsIDInRange(current, start, stop string) bool {
	// Handle special Redis Range symbols
	if start == "-" {
		start = "0-0"
	}
	if stop == "+" {
		stop = "9999999999999-99999"
	} // maximum upper bound

	// If start is just "0", it should be treated as "0-0" for comparison
	if !strings.Contains(start, "-") {
		start += "-0"
	}
	if !strings.Contains(stop, "-") {
		stop += "-999999999"
	} // Cover all sequences in that ms

	return CompareIDs(current, start) >= 0 && CompareIDs(current, stop) <= 0
}

func CompareIDs(id1, id2 string) int {
	t1, s1 := GetIdTimeSequence(id1)
	t2, s2 := GetIdTimeSequence(id2)

	if t1 != t2 {
		if t1 > t2 {
			return 1
		}
		return -1
	}
	if s1 != s2 {
		if s1 > s2 {
			return 1
		}
		return -1
	}
	return 0
}
