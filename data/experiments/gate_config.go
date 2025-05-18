package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type EndpointName string
type ProxyCode int

const (
	InfoSettingsEnd     EndpointName = "info_settings"
	PlayerEnd           EndpointName = "player"
	TestAccountsEnd     EndpointName = "test_accounts"
	TestTransactionsEnd EndpointName = "test_transactions"
	TestTransactionsEnd EndpointName = "test_transactions_info"

	ExampleOne ProxyCode = 101
	ExampleTwo ProxyCode = 102
)

type GateConfig struct {
	BaseURL   string
	Endpoints map[EndpointName]string
}

// key = proxy code (like "101", "102") which is used in userId
var localGateConfigs = map[ProxyCode]GateConfig{
	ExampleOne: {
		BaseURL: "http://local.example.com/api",
		Endpoints: map[EndpointName]string{
			InfoSettingsEnd:     "/info/settings",
			PlayerEnd:           "/player",
			TestAccountsEnd:     "/test/accounts",
			TestTransactionsEnd: "/test/transactions",
			TestTransactionsEnd: "/test/transactions_info",
		},
	},
}

var devGateConfigs = map[ProxyCode]GateConfig{
	ExampleOne: {
		BaseURL: "https://dev.example.com/api",
		Endpoints: map[EndpointName]string{
			InfoSettingsEnd:     "/info/settings",
			PlayerEnd:           "/player",
			TestAccountsEnd:     "/test/accounts",
			TestTransactionsEnd: "/test/transactions",
			TestTransactionsEnd: "/test/transactions_info",
		},
	},
}

var prodGateConfigs = map[ProxyCode]GateConfig{
	ExampleOne: {
		BaseURL: "https://prod.example.com/api",
		Endpoints: map[EndpointName]string{
			InfoSettingsEnd:     "/info/settings",
			PlayerEnd:           "/player",
			TestAccountsEnd:     "/test/accounts",
			TestTransactionsEnd: "/test/transactions",
			TestTransactionsEnd: "/test/transactions_info",
		},
	},
}

// The first three digits represent the proxy code, and the remaining digits are the user ID.
func ParseRawUserID(rawUserID int64) (proxyCode ProxyCode, userID int64, err error) {
	rawStr := strconv.FormatInt(rawUserID, 10)

	if len(rawStr) < 4 {
		return 0, 0, errors.New("user_id too short: must be at least 4 digits")
	}

	proxyPart := rawStr[:3]
	userPart := rawStr[3:]

	var proxyCodeInt int
	proxyCodeInt, err = strconv.Atoi(proxyPart)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid proxy part: %v", err)
	}
	proxyCode = ProxyCode(proxyCodeInt)

	userID, err = strconv.ParseInt(userPart, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid user part: %v", err)
	}

	return
}

// BuildGateUrl builds the full URL to an endpoint for a given proxy
func BuildGateUrl(proxy ProxyCode, endpointKey EndpointName) (string, error) {
	env := os.Getenv("ENV")
	if env == "" {
		return "", errors.New("missing ENV in environment variables")
	}

	var configMap map[ProxyCode]GateConfig
	switch env {
	case "local":
		configMap = localGateConfigs
	case "dev":
		configMap = devGateConfigs
	case "prod":
		configMap = prodGateConfigs
	default:
		return "", fmt.Errorf("unknown environment: %s", env)
	}

	config, ok := configMap[proxy]
	if !ok {
		return "", fmt.Errorf("proxy '%s' not configured", proxy)
	}

	relPath, ok := config.Endpoints[endpointKey]
	if !ok {
		return "", fmt.Errorf("endpoint '%s' not found for proxy '%s'", endpointKey, proxy)
	}

	return config.BaseURL + relPath, nil
}
