package fncommon

import (
	"os"
	"strconv"
)

func GetEnv(key, fallback string) string {
	println("Looking up env var ", key)
	if value, ok := os.LookupEnv(key); ok {
		println("Found value ", value, " for env var ", key)
		return value
	}
	println("No value found for env var ", key, ", falling back to ", fallback)
	return fallback
}

func GetEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		// linter liked this better than if/else
		var err error
		var i int
		if i, err = strconv.Atoi(value); err != nil {
			panic(err) // not sure how to handle this
		}
		return i
	}
	return fallback
}
