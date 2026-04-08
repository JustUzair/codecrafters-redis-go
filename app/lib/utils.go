package lib

import (
	"strconv"
	"strings"
)

func GetIdTimeSequence(id string) (int, int) {
	time_seq := strings.Split(id, "-")
	time, _ := strconv.Atoi(time_seq[0])

	seq, _ := strconv.Atoi(time_seq[1])

	return time, seq
}

func IsValidTimeSequence(prev_time int, prev_seq int, new_time int, new_seq int) bool {
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
