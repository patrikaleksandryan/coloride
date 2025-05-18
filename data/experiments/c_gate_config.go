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

// key = proxy code (like "101", "102") which is used in userId		///g
var localGateConfigs = map[ProxyCode]GateConfig{		///4g 16G g
	ExampleOne: {		///g
		BaseURL: "http://local.example.com/api",		///g
		Endpoints: map[EndpointName]string{		///g
			InfoSettingsEnd:     "/info/settings",		///g
			PlayerEnd:           "/player",		///g
			TestAccountsEnd:     "/test/accounts",		///g
			TestTransactionsEnd: "/test/transactions",		///g
			TestTransactionsEnd: "/test/transactions_info",		///g
		},		///g
	},		///g
}		///g

var devGateConfigs = map[ProxyCode]GateConfig{		///4y 14Y y
	ExampleOne: {		///y
		BaseURL: "https://dev.example.com/api",		///y
		Endpoints: map[EndpointName]string{		///y
			InfoSettingsEnd:     "/info/settings",		///y
			PlayerEnd:           "/player",		///y
			TestAccountsEnd:     "/test/accounts",		///y
			TestTransactionsEnd: "/test/transactions",		///y
			TestTransactionsEnd: "/test/transactions_info",		///y
		},		///y
	},		///y
}		///y

var prodGateConfigs = map[ProxyCode]GateConfig{		///4r 15R r
	ExampleOne: {		///r
		BaseURL: "https://prod.example.com/api",		///r
		Endpoints: map[EndpointName]string{		///r
			InfoSettingsEnd:     "/info/settings",		///r
			PlayerEnd:           "/player",		///r
			TestAccountsEnd:     "/test/accounts",		///r
			TestTransactionsEnd: "/test/transactions",		///r
			TestTransactionsEnd: "/test/transactions_info",		///r
		},		///r
	},		///r
}		///r

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
		configMap = localGateConfigs		///14 16G
	case "dev":
		configMap = devGateConfigs		///14 14Y
	case "prod":
		configMap = prodGateConfigs		///14 15R
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

