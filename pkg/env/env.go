package env

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// String returns the env string value if the key exists.
// Otherwise, returns initial value.
func String(key, initial string) string {
	v, exists := os.LookupEnv(key)
	if !exists {
		return initial
	}

	v = strings.TrimSpace(v)
	if len(v) == 0 {
		return initial
	}

	return v
}

// Int returns the env integer value if the key exists.
// Otherwise, returns initial value.
func Int(key string, initial int) int {
	v := String(key, "")
	if v == "" {
		return initial
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}

	return n
}

// Duration returns the env duration value if the key exists.
// Otherwise, returns initial value.
func Duration(key string, initial time.Duration) time.Duration {
	v := String(key, "")
	if v == "" {
		return initial
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		panic(err)
	}

	return d
}

// Bool returns the env boolean value if the key exists.
// Otherwise, returns initial value.
func Bool(key string, initial bool) bool {
	v := String(key, "")
	if v == "" {
		return initial
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		panic(err)
	}

	return b
}
