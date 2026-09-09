package env

import (
	"os"
	"strconv"
)

func GetString(key, payload string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return payload
}

func GetInt(key string, payload int) int {
	val, ok := os.LookupEnv(key)
	if ok {
		i, err := strconv.Atoi(val)
		if err == nil {
			return i
		}
	}
	return payload
}

func GetBool(key string, payload bool) bool {
	val, ok := os.LookupEnv(key)
	if ok {
		b, err := strconv.ParseBool(val)
		if err == nil {
			return b
		}
	}
	return payload
}
