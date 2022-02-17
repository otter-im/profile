package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func ServiceEnvironment() string {
	return envString("SERVICE_ENV", "dev")
}

func ServiceHost() string {
	return envString("SERVICE_HOST", "0.0.0.0")
}

func ServicePort() string {
	return fmt.Sprintf("%v", envInt("SERVICE_PORT", 3000))
}

func envString(key string, v string) string {
	result, ok := os.LookupEnv(key)
	if !ok {
		return v
	}
	return result
}

func envInt(key string, v int) int {
	resultString := envString(key, fmt.Sprintf("%v", v))
	if resultString == "" {
		return v
	}

	result, err := strconv.Atoi(resultString)
	if err != nil {
		log.Printf("invalid environment variable \"%s\", using default of %v\n", key, v)
		return v
	}
	return result
}

func envBool(key string, v bool) bool {
	resultString := envString(key, fmt.Sprintf("%v", v))
	if resultString == "true" {
		return true
	}
	return false
}
