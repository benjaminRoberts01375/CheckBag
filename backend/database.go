package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"

	Printing "github.com/benjaminRoberts01375/CheckBag/backend/logging"
	"github.com/valkey-io/valkey-go"
)

type BasicDB interface {
	// Basic cache functions
	Set(ctx context.Context, key string, value string, duration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error

	SetHash(ctx context.Context, key string, values map[string]string) error
	GetHash(ctx context.Context, key string) (map[string]string, error)
	DeleteHash(ctx context.Context, key string) error

	RenameKey(ctx context.Context, oldKey string, newKey string) error

	IncrementHashField(ctx context.Context, key string, field string, amount int, expiration time.Time) error
	IncrementKey(ctx context.Context, key string, amount int, expiration time.Time) error

	AddToList(ctx context.Context, key string, value string) error
	RemoveFromList(ctx context.Context, key string, value string) error
	GetList(ctx context.Context, key string) ([]string, error)
}

type AdvancedDB interface {
	incrementAnalytics(ctx context.Context, serviceID string, resource string, country string, ip string, responseCode int, receivedBytes int, sentBytes int) error
	getAnalyticsService(ctx context.Context, service ServiceData, timeStep AnalyticsTimeStep) map[time.Time]Analytic
	deleteService(ctx context.Context, service ServiceLink) error
	addAPIKey(ctx context.Context, APIKey string, keyID string, name string) error
	removeAPIKey(ctx context.Context, APIKeyID string) error
	getAPIKeyInfo(ctx context.Context) ([]APIKeyInfo, error)
	apiKeyExists(ctx context.Context, APIKey string) bool
	getVersion(ctx context.Context) (string, error)
	setVersion(ctx context.Context, version string) error
}

type ValkeyDB struct {
	db     valkey.Client
	prefix string
}

type DB struct {
	basicDB BasicDB
}

type AnalyticsTimeStep struct {
	time         func(step int) time.Time
	maximumUnits int
}

func (analytics AnalyticsTimeStep) timeStr(step int) string {
	return analytics.time(step).Format(time.RFC3339)
}

var (
	cacheAnalyticsMinute = AnalyticsTimeStep{maximumUnits: 60, time: func(step int) time.Time {
		return time.Now().Truncate(time.Minute).Add(time.Duration(step) * time.Minute)
	}}
	cacheAnalyticsHour = AnalyticsTimeStep{maximumUnits: 24, time: func(step int) time.Time {
		return time.Now().Truncate(time.Hour).Add(time.Duration(step) * time.Hour)
	}}
	cacheAnalyticsDay = AnalyticsTimeStep{maximumUnits: 30, time: func(step int) time.Time {
		now := time.Now()
		year, month, day := now.Date()
		return time.Date(year, month, day, 0, 0, 0, 0, now.Location()).AddDate(0, 0, step)
	}}
	cacheAnalyticsMonth = AnalyticsTimeStep{maximumUnits: 12, time: func(step int) time.Time {
		now := time.Now()
		year, month, _ := now.Date()
		return time.Date(year, month, 1, 0, 0, 0, 0, now.Location()).AddDate(0, step, 0)
	}}
	cacheAnalyticsTime = []AnalyticsTimeStep{cacheAnalyticsMinute, cacheAnalyticsHour, cacheAnalyticsDay, cacheAnalyticsMonth}
)

// Basic cache functions

func SetupDB() AdvancedDB {
	dbPort, err := strconv.Atoi(os.Getenv("CACHE_PORT"))
	if err != nil || dbPort <= 0 {
		panic("Failed to parse CACHE_PORT: " + err.Error())
	}
	dbAddress := os.Getenv("CACHE_ADDRESS")
	if dbAddress == "" {
		panic("No cache container name specified")
	}
	cacheIDLength, err := strconv.Atoi(os.Getenv("CACHE_ID_LENGTH"))
	if err != nil || cacheIDLength <= 0 {
		panic("Failed to parse CACHE_ID_LENGTH: " + err.Error())
	}
	dbPassword := os.Getenv("CACHE_PASSWORD")
	if dbPassword == "" {
		Printing.Println("No cache password specified")
	}

	// Connect to Valkey
	dbURL := fmt.Sprintf("%s:%d", dbAddress, dbPort)
	dbConnectionOptions := valkey.ClientOption{
		InitAddress: []string{dbURL},
		Password:    dbPassword,
	}
	dbClient, err := valkey.NewClient(dbConnectionOptions)
	if err != nil {
		panic("Failed to connect to Valkey: " + err.Error())
	}

	// Save DB
	db := DB{
		basicDB: &ValkeyDB{
			db:     dbClient,
			prefix: "CheckBag:",
		},
	}
	db.versioning()
	return db
}

func (db DB) versioning() {
	expectedDBVersion := "2"
	ctx := context.Background()
	actualDBVersion, err := db.basicDB.Get(ctx, "version")
	if err != nil {
		Printing.PrintErrStr("Could not get version from DB, setting to " + expectedDBVersion)
		db.setVersion(ctx, expectedDBVersion)
	} else if actualDBVersion != "2" {
		panic("Expected database version " + expectedDBVersion + " but got " + actualDBVersion)
	}
}

func (db *ValkeyDB) Get(ctx context.Context, key string) (string, error) {
	return db.db.Do(ctx, db.db.B().Get().Key(db.prefix+key).Build()).ToString()
}

func (db *ValkeyDB) GetHash(ctx context.Context, key string) (map[string]string, error) {
	return db.db.Do(ctx, db.db.B().Hgetall().Key(db.prefix+key).Build()).AsStrMap()
}

func (db *ValkeyDB) Delete(ctx context.Context, key string) error {
	return db.db.Do(ctx, db.db.B().Del().Key(db.prefix+key).Build()).Error()
}

func (db *ValkeyDB) DeleteHash(ctx context.Context, key string) error {
	return db.db.Do(ctx, db.db.B().Hdel().Key(db.prefix+key).Field("purpose").Build()).Error()
}

// If duration is 0, the value is set without an expiration time
func (db *ValkeyDB) Set(ctx context.Context, key string, value string, duration time.Duration) error {
	if duration == 0 {
		return db.db.Do(ctx, db.db.B().Set().Key(db.prefix+key).Value(value).Build()).Error()
	}
	return db.db.Do(ctx, db.db.B().Set().Key(db.prefix+key).Value(value).Ex(duration).Build()).Error()
}

func (db *ValkeyDB) SetHash(ctx context.Context, key string, values map[string]string) error {
	hash := db.db.B().Hset().Key(db.prefix + key).FieldValue()
	for field, value := range values {
		hash = hash.FieldValue(field, value)
	}
	return db.db.Do(ctx, hash.Build()).Error()
}

func (db *ValkeyDB) RenameKey(ctx context.Context, oldKey string, newKey string) error {
	return db.db.Do(ctx, db.db.B().Rename().Key(db.prefix+oldKey).Newkey(db.prefix+newKey).Build()).Error()
}

func (db *ValkeyDB) IncrementHashField(ctx context.Context, key string, field string, amount int, expiration time.Time) error {
	err := db.db.Do(ctx, db.db.B().Hincrby().Key(db.prefix+key).Field(field).Increment(int64(amount)).Build()).Error()
	if err != nil {
		return err
	}
	remainingTime := time.Until(expiration)
	return db.db.Do(ctx, db.db.B().Expire().Key(db.prefix+key).Seconds(int64(remainingTime.Seconds())).Build()).Error()
}

func (db *ValkeyDB) IncrementKey(ctx context.Context, key string, amount int, expiration time.Time) error {
	err := db.db.Do(ctx, db.db.B().Incrby().Key(db.prefix+key).Increment(int64(amount)).Build()).Error()
	if err != nil {
		return err
	}
	remainingTime := time.Until(expiration)
	return db.db.Do(ctx, db.db.B().Expire().Key(db.prefix+key).Seconds(int64(remainingTime.Seconds())).Build()).Error()
}

func (db *ValkeyDB) AddToList(ctx context.Context, key string, value string) error {
	return db.db.Do(ctx, db.db.B().Lpush().Key(db.prefix+key).Element(value).Build()).Error()
}

func (db *ValkeyDB) RemoveFromList(ctx context.Context, key string, value string) error {
	return db.db.Do(ctx, db.db.B().Lrem().Key(db.prefix+key).Count(1).Element(value).Build()).Error()
}

func (db *ValkeyDB) GetList(ctx context.Context, key string) ([]string, error) {
	return db.db.Do(ctx, db.db.B().Lrange().Key(db.prefix+key).Start(0).Stop(-1).Build()).AsStrSlice()
}

// Higher-level DB functions

func (db DB) incrementAnalytics(ctx context.Context, serviceID string, resource string, country string, ip string, responseCode int, receivedBytes int, sentBytes int) error {
	for _, timeStep := range cacheAnalyticsTime {
		recordTime := timeStep.timeStr(0)
		expiration := timeStep.time(timeStep.maximumUnits)
		quantity := strconv.Itoa(timeStep.maximumUnits)
		baseKey := "Analytics:" + serviceID + ":" + quantity + ":" + recordTime + ":"
		err := db.basicDB.IncrementKey(ctx, baseKey+"quantity", 1, expiration)
		if err != nil {
			Printing.PrintErrStr("Could not increment quantity analytics key: " + err.Error())
			return err
		}
		err = db.basicDB.IncrementKey(ctx, baseKey+"received_bytes", receivedBytes, expiration)
		if err != nil {
			Printing.PrintErrStr("Could not increment received bytes analytics:" + err.Error())
			return err
		}
		err = db.basicDB.IncrementKey(ctx, baseKey+"sent_bytes", sentBytes, expiration)
		if err != nil {
			Printing.PrintErrStr("Could not increment sent bytes analytics")
			return err
		}
		err = db.basicDB.IncrementHashField(ctx, baseKey+"country", country, 1, expiration)
		if err != nil {
			Printing.PrintErrStr("Could not increment analytics country: " + err.Error())
			return err
		}
		err = db.basicDB.IncrementHashField(ctx, baseKey+"ip", ip, 1, expiration)
		if err != nil {
			Printing.PrintErrStr("Could not increment analytics ip: " + err.Error())
			return err
		}
		err = db.basicDB.IncrementHashField(ctx, baseKey+"resource", resource, 1, expiration)
		if err != nil {
			Printing.PrintErrStr("Could not increment analytics resource: " + err.Error())
			return err
		}
		err = db.basicDB.IncrementHashField(ctx, baseKey+":response_code", strconv.Itoa(responseCode), 1, expiration)
		if err != nil {
			Printing.PrintErrStr("Could not increment analytics response code: " + err.Error())
			return err
		}
	}
	return nil
}

func (db DB) getAnalyticsService(ctx context.Context, service ServiceData, timeStep AnalyticsTimeStep) map[time.Time]Analytic {
	analytics := map[time.Time]Analytic{}
	quantity := strconv.Itoa(timeStep.maximumUnits)
	for timePeriod := range timeStep.maximumUnits {
		baseKey := "Analytics:" + service.ID + ":" + quantity + ":" + timeStep.timeStr(-timePeriod) + ":"
		quantityRaw, err := db.basicDB.Get(ctx, baseKey+"quantity")
		if err != nil {
			continue
		}
		sentBytesRaw, _ := db.basicDB.Get(ctx, baseKey+"sent_bytes")
		receivedBytesRaw, _ := db.basicDB.Get(ctx, baseKey+"received_bytes")
		countryRaw, err := db.basicDB.GetHash(ctx, baseKey+"country")
		if err != nil {
			continue
		}
		ipRaw, err := db.basicDB.GetHash(ctx, baseKey+"ip")
		if err != nil {
			continue
		}
		resourceRaw, err := db.basicDB.GetHash(ctx, baseKey+"resource")
		if err != nil {
			continue
		}
		responseCodesRaw, err := db.basicDB.GetHash(ctx, baseKey+"response_code")
		if err != nil {
			continue
		}

		quantity, err := strconv.Atoi(quantityRaw)
		if err != nil {
			quantity = 0
		}
		sentBytes, err := strconv.Atoi(sentBytesRaw)
		if err != nil {
			sentBytes = 0
		}
		receivedBytes, err := strconv.Atoi(receivedBytesRaw)
		if err != nil {
			receivedBytes = 0
		}
		country := make(map[string]int)
		for countryName, countryCount := range countryRaw {
			country[countryName], err = strconv.Atoi(countryCount)
			if err != nil {
				continue
			}
		}
		ip := make(map[string]int)
		for ipAddress, ipCount := range ipRaw {
			ip[ipAddress], err = strconv.Atoi(ipCount)
			if err != nil {
				continue
			}
		}
		resource := make(map[string]int)
		for resourceName, resourceCount := range resourceRaw {
			resource[resourceName], err = strconv.Atoi(resourceCount)
			if err != nil {
				continue
			}
		}
		responseCodes := make(map[int]int)
		for responseCode, responseCodeCount := range responseCodesRaw {
			responseCodeInt, err := strconv.Atoi(responseCode)
			if err != nil {
				continue
			}
			responseCodes[responseCodeInt], err = strconv.Atoi(responseCodeCount)
			if err != nil {
				continue
			}
		}

		analytics[timeStep.time(-timePeriod)] = Analytic{
			Quantity:      quantity,
			Country:       country,
			IP:            ip,
			Resource:      resource,
			ResponseCode:  responseCodes,
			SentBytes:     sentBytes,
			ReceivedBytes: receivedBytes,
		}
	}

	return analytics
}

func (db DB) getVersion(ctx context.Context) (string, error) {
	version, err := db.basicDB.Get(ctx, "version")
	if err != nil {
		return "", errors.New("Failed to get version: " + err.Error())
	}
	return version, nil
}

func (db DB) setVersion(ctx context.Context, version string) error {
	err := db.basicDB.Set(ctx, "version", version, 0)
	if err != nil {
		return errors.New("Failed to set version: " + err.Error())
	}
	return nil
}

func (db DB) deleteService(ctx context.Context, service ServiceLink) error {
	for _, timeStep := range cacheAnalyticsTime {
		recordTime := timeStep.timeStr(0)
		quantity := strconv.Itoa(timeStep.maximumUnits)
		baseKey := "Analytics:" + service.ID + ":" + quantity + ":" + recordTime + ":"
		err := db.basicDB.Delete(ctx, baseKey+"quantity")
		if err != nil {
			Printing.PrintErrStr("Could not delete minute analytics key: " + err.Error())
			return err
		}
		err = db.basicDB.DeleteHash(ctx, baseKey+"country")
		if err != nil {
			Printing.PrintErrStr("Could not delete minute analytics country: " + err.Error())
			return err
		}
		err = db.basicDB.DeleteHash(ctx, baseKey+"ip")
		if err != nil {
			Printing.PrintErrStr("Could not delete minute analytics ip: " + err.Error())
			return err
		}
		err = db.basicDB.DeleteHash(ctx, baseKey+"resource")
		if err != nil {
			Printing.PrintErrStr("Could not delete minute analytics resource: " + err.Error())
			return err
		}
		err = db.basicDB.DeleteHash(ctx, baseKey+"response_code")
		if err != nil {
			Printing.PrintErrStr("Could not delete minute analytics response code: " + err.Error())
			return err
		}
	}
	return nil
}

func (db DB) addAPIKey(ctx context.Context, APIKey string, keyID string, name string) error {
	if name == "" {
		name = "Unnamed API"
	}
	hash := map[string]string{
		"name": name,
		"id":   keyID,
	}

	err := db.basicDB.SetHash(ctx, "APIKey:"+APIKey, hash)
	if err != nil {
		return err
	}
	err = db.basicDB.AddToList(ctx, "APIKeys", APIKey)
	if err != nil {
		return err
	}
	Printing.Println("Added API key: " + APIKey)
	return nil
}

func (db DB) removeAPIKey(ctx context.Context, APIKeyID string) error {
	keys, err := db.basicDB.GetList(ctx, "APIKeys")
	if err != nil {
		return err
	}
	for _, key := range keys {
		keyInfo, err := db.basicDB.GetHash(ctx, "APIKey:"+key)
		if err != nil {
			return err
		}
		data, _ := json.Marshal(keyInfo)
		Printing.Println("Key info: " + string(data))

		if keyInfo["id"] == APIKeyID {
			err := db.basicDB.DeleteHash(ctx, "APIKey:"+key)
			if err != nil {
				return err
			}
			err = db.basicDB.RemoveFromList(ctx, "APIKeys", key)
			if err != nil {
				return err
			}
			Printing.Println("Removed API key: " + key)
			return nil
		}
	}
	return errors.New("API key not found")
}

func (db DB) getAPIKeyInfo(ctx context.Context) ([]APIKeyInfo, error) {
	keys, err := db.basicDB.GetList(ctx, "APIKeys")
	if err != nil {
		return nil, err
	}
	keysInfo := make([]APIKeyInfo, len(keys))
	for i, key := range keys {
		keyInfo, err := db.basicDB.GetHash(ctx, "APIKey:"+key)
		if err != nil {
			return nil, err
		}
		keysInfo[i].Name = keyInfo["name"]
		keysInfo[i].ID = keyInfo["id"]
	}

	return keysInfo, nil
}

func (db DB) apiKeyExists(ctx context.Context, APIKey string) bool {
	keys, err := db.basicDB.GetList(ctx, "APIKeys")
	if err != nil {
		return false
	}
	return slices.Contains(keys, APIKey)
}
