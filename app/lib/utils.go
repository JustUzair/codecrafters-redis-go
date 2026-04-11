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
