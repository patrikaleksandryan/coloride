package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type TestController interface {
	InfoSettings(c *gin.Context)
	Player(c *gin.Context)
	Accounts(c *gin.Context)
	Transactions(c *gin.Context)
	TransactionsInfo(c *gin.Context)
}

type testController struct {
	testService service.TestService
	redisClient *redis.Client
	logger      *log.Logger
}

func NewTestController(testService service.TestService, redisClient *redis.Client) TestController {
	newLogger, _, err := logger.CreateLogger()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	return &testController{
		testService: testService,
		redisClient: redisClient,
		logger:      newLogger,
	}
}

func (sc *testController) InfoSettings(ctx *gin.Context) {
	resp, statusCode, err := sc.testService.InfoSettings()
	// Handles only internal errors, because all gate errors should be in resp and just passed further.
	if err != nil {
		sc.logger.Printf("Error \"%d\" at /info/settings endpoint with message: %s", statusCode, err.Error())
		ctx.JSON(statusCode, gin.H{"code": "other", "message": "Internal Error."})
		return
	}

	ctx.Data(statusCode, "application/json", resp)
}

func (sc *testController) Player(ctx *gin.Context) {
	authToken := ctx.Query("session")
	if authToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": "bad_session", "message": "Session is required."})
		return
	}

	claims, err := security.DecodeAuthToken(authToken)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": "bad_session", "message": "Token is expired."})
		} else if errors.Is(err, jwt.ErrTokenMalformed) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": "bad_session", "message": "Malformed token."})
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": "bad_session", "message": "Invalid token."})
		}

		sc.logger.Printf("Error at /play endpoint with message: %s", err.Error())
		return
	}

	// Check if UserID exists in claims
	if claims.UserID == 0 {
		sc.logger.Printf("Error at /play endpoint with message: UserID is missing in claims.")
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": "bad_session", "message": "Malformed token."})
		return
	}

	response, statusCode, err := sc.testService.Player(claims.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"code": "other", "message": "Failed to fetch data."})
		return
	}

	ctx.Data(statusCode, "application/json", response)
}

func (sc *testController) Accounts(ctx *gin.Context) {
	// Reading the original request body
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": "other", "message": "Failed to read request body."})
		return
	}

	// Restore body to pass it further
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Extracting userID
	userID, err := security.ExtractUserIDFromPayload(bodyBytes)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": "other", "message": "Invalid or missing user_id."})
		return
	}

	// Pass request further
	responder, statusCode, _, err := sc.testService.Wallet(userID, config.AccountsEnd, bodyBytes)
	if err != nil {
		sc.logger.Printf("Error at /wallet/accounts endpoint: %s", err.Error())
		ctx.JSON(http.StatusBadGateway, gin.H{"code": "other", "message": "Failed to fetch data."})
		return
	}

	ctx.Data(statusCode, "application/json", responder)
}

// Unfinished endpoint
func (sc *testController) Transactions(ctx *gin.Context) {
	// Reading the original request body
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": "other", "message": "Failed to read request body."})
		return
	}

	// Restore the body to pass it further
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Extracting userID
	userID, err := security.ExtractUserIDFromPayload(bodyBytes)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": "other", "message": "Invalid or missing user_id."})
		return
	}

	// Check if X-Process-Until exists in the header
	processUntilRaw := ctx.GetHeader("X-Process-Until")
	if processUntilRaw == "" {
		sc.logger.Printf("Error at /wallet/transactions endpoint: X-Process-Until header is missing.")
		ctx.JSON(http.StatusBadRequest, gin.H{"code": "other", "message": "X-Process-Until header is required."})
		return
	}

	//TODO should be rafactored to use RFC3339Nano formatting time
	// Time parsing from RFC3339 format
	processUntil, err := time.Parse(time.RFC3339, processUntilRaw)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": "other", "message": "Invalid X-Process-Until format."})
		return
	}

	// Time validation
	if time.Now().UTC().After(processUntil) {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": "other", "message": "Request expired (X-Process-Until in the past)."})
		return
	}

	type TransactionRequest struct {
		UserID       int `json:"user_id"`
		Transactions []struct {
			TxID string `json:"tx_id"`
		} `json:"transactions"`
	}

	// Refactor to use postgres in future
	var req TransactionRequest
	if err := json.Unmarshal(bodyBytes, &req); err == nil {
		ctxRedis := context.Background()
		ttl := 24 * time.Hour

		for _, tx := range req.Transactions {
			if tx.TxID == "" {
				continue
			}
			err := sc.redisClient.Set(ctxRedis, tx.TxID, req.UserID, ttl).Err()
			if err != nil {
				sc.logger.Printf("Failed to store tx_id in Redis: %v", err)
			}
		}
	} else {
		sc.logger.Printf("Failed to parse transaction body: %v", err)
	}

	// Pass the request further
	responder, statusCode, headers, err := sc.testService.Wallet(userID, config.TransactionsInfoEnd, bodyBytes)
	if err != nil {
		sc.logger.Printf("Error at /wallet/transactions endpoint: %s", err.Error())
		ctx.JSON(http.StatusBadGateway, gin.H{"code": "other", "message": "Failed to process request."})
		return
	}

	// Forwarding all headers
	for key, values := range headers {
		for _, value := range values {
			ctx.Header(key, value)
		}
	}

	ctx.Data(statusCode, headers.Get("Content-Type"), responder)
}

func (sc *testController) TransactionsInfo(ctx *gin.Context) {
	// Reading the original request body
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": "other", "message": "Failed to read request body."})
		return
	}

	// Restore the body to pass it further
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var payload struct {
		Transactions []struct {
			TxID string `json:"tx_id"`
		} `json:"transactions"`
	}

	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": "other", "message": "Invalid JSON format."})
		return
	}

	if len(payload.Transactions) == 0 || payload.Transactions[0].TxID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": "other", "message": "Missing tx_id."})
		return
	}

	txID := payload.Transactions[0].TxID

	// Refactor to use postgres in future
	rawUserID, err := sc.getUserIDByTransactionID(txID)
	if err != nil {
		sc.logger.Printf("user_id not found in redis db: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"tx_id": txID, "status": "not_found"})
		return
	}

	// Pass the request further
	responder, statusCode, headers, err := sc.testService.Wallet(rawUserID, config.TransactionsInfoEnd, bodyBytes)
	if err != nil {
		sc.logger.Printf("Error at /wallet/transactions_info endpoint: %s", err.Error())
		ctx.JSON(http.StatusBadGateway, gin.H{"code": "other", "message": "Failed to process request."})
		return
	}

	// Forwarding all headers
	for key, values := range headers {
		for _, value := range values {
			ctx.Header(key, value)
		}
	}

	ctx.Data(statusCode, headers.Get("Content-Type"), responder)
}

// AI Generated code (Refactor to use postgres in future)
func (sc *testController) getUserIDByTransactionID(txID string) (int64, error) {
	ctx := context.Background()
	val, err := sc.redisClient.Get(ctx, txID).Result()
	if err == redis.Nil {
		return 0, fmt.Errorf("tx_id not found in Redis: %s", txID)
	} else if err != nil {
		return 0, fmt.Errorf("failed to fetch from Redis: %w", err)
	}

	userID, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid user_id format in Redis: %w", err)
	}

	return int64(userID), nil
}
