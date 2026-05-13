package helper

import "strconv"

func ParseUint(
	s string,
) uint64 {

	v, _ := strconv.ParseUint(
		s,
		10,
		64,
	)

	return v
}

func ParseIntDefault(
	s string,
	def int,
) int {

	if s == "" {
		return def
	}

	v, err := strconv.Atoi(s)

	if err != nil {
		return def
	}

	return v
}
