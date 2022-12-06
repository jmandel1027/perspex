package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	if fallback == "" {
		otelzap.L().Ctx(context.TODO()).Info(zap.String(key, "env missing, no fallback").String)
	}

	return fallback
}

// MustGet will return the env or panic if it is not present
func MustGet(k string, fallback string) string {
	v := getEnv(k, fallback)

	return v
}

// MustGetInt will return the env or fallback value if it is not present
func MustGetInt(k string, fallback string) int {
	v := getEnv(k, fallback)

	i, err := strconv.Atoi(v)
	if err != nil {
		log.Print(k, fmt.Sprintf("ENV err: [%s] %s", k, err.Error()))
	}

	return i
}

// MustGetBool will return the env as boolean, or fallback value if not present
func MustGetBool(k string, fallback string) bool {
	v := getEnv(k, fallback)

	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Print(k, fmt.Sprintf("ENV err: [%s] %s", k, err.Error()))
	}

	return b
}
