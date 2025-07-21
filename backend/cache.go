package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	JWT "github.com/benjaminRoberts01375/Web-Tech-Stack/jwt"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
	"github.com/valkey-io/valkey-go"
)

type CacheSpec interface {
	// Cache Management
	Setup()
	Close()

	// Basic cache functions
	Set(key string, value string, duration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error

	SetHash(key string, values map[string]string, duration time.Duration) error
	GetHash(key string) (map[string]string, error)
	DeleteHash(key string) error

	RenameKey(oldKey string, newKey string) error

	IncrementHashField(key string, field string, amount int, expiration time.Time) error
	IncrementKey(key string, expiration time.Time) error
}

type CacheLayer struct { // Implements main 5 functions
	DB            valkey.Client
	ContainerName string
	Port          int
	Password      string `json:"password"`
}

type CacheClient[client CacheSpec] struct { // Holds some DB that satisfies the CacheSpec interface. Action functions here
	raw client
}

type AnalyticsTimeStep struct {
	key          func(step int) string
	maximumUnits int
}

const (
	cacheUserSignIn = JWT.LoginDuration
)

var (
	cacheAnalyticsMinute = AnalyticsTimeStep{maximumUnits: 60, key: func(step int) string {
		return time.Now().Truncate(time.Minute).Add(time.Duration(step) * time.Minute).Format(time.RFC3339)
	}}
	cacheAnalyticsHour = AnalyticsTimeStep{maximumUnits: 24, key: func(step int) string {
		return time.Now().Truncate(time.Hour).Add(time.Duration(step) * time.Hour).Format(time.RFC3339)
	}}
	cacheAnalyticsDay = AnalyticsTimeStep{maximumUnits: 30, key: func(step int) string {
		now := time.Now()
		year, month, day := now.Date()
		return time.Date(year, month, day, 0, 0, 0, 0, now.Location()).AddDate(0, 0, step).Format(time.RFC3339)
	}}
	cacheAnalyticsMonth = AnalyticsTimeStep{maximumUnits: 12, key: func(step int) string {
		now := time.Now()
		year, month, _ := now.Date()
		return time.Date(year, month, 1, 0, 0, 0, 0, now.Location()).AddDate(0, step, 0).Format(time.RFC3339)
	}}
	cacheAnalyticsTime = []AnalyticsTimeStep{cacheAnalyticsMinute, cacheAnalyticsHour, cacheAnalyticsDay, cacheAnalyticsMonth}
)

const cacheDataValid = "valid"

// Basic cache functions

func (cache *CacheLayer) Setup() {
	cachePort, err := strconv.Atoi(os.Getenv("CACHE_PORT"))
	if err != nil || cachePort <= 0 {
		panic("Failed to parse CACHE_PORT: " + err.Error())
	}
	cacheContainerName := os.Getenv("CACHE_CONTAINER_NAME")
	if cacheContainerName == "" {
		panic("No cache container name specified")
	}
	cacheIDLength, err := strconv.Atoi(os.Getenv("CACHE_ID_LENGTH"))
	if err != nil || cacheIDLength <= 0 {
		panic("Failed to parse CACHE_ID_LENGTH: " + err.Error())
	}

	Coms.ReadExternalConfig("valkey.json", &cache)
	cache.Port = cachePort
	cache.ContainerName = cacheContainerName

	// TODO: Handle username and client name
	url := fmt.Sprintf("%s:%d", cache.ContainerName, cache.Port)
	options := valkey.ClientOption{
		InitAddress: []string{url},
		// Username:    config.CacheUsername,
		Password: cache.Password,
		// ClientName:  config.CacheClientName
	}
	client, err := valkey.NewClient(options)
	if err != nil {
		Coms.PrintErr(err)
		panic("Could not connect to Valkey: " + err.Error())
	}
	cache.DB = client
	Coms.Println("Connected to Valkey")
}

func (cache *CacheLayer) Close() {
	cache.DB.Close()
}

func (cache CacheLayer) Get(key string) (string, error) {
	ctx := context.Background()
	return cache.DB.Do(ctx, cache.DB.B().Get().Key(key).Build()).ToString()
}

func (cache CacheLayer) GetHash(key string) (map[string]string, error) {
	ctx := context.Background()
	return cache.DB.Do(ctx, cache.DB.B().Hgetall().Key(key).Build()).AsStrMap()
}

func (cache CacheLayer) Delete(key string) error {
	return cache.DB.Do(context.Background(), cache.DB.B().Del().Key(key).Build()).Error()
}

func (cache CacheLayer) DeleteHash(key string) error {
	return cache.DB.Do(context.Background(), cache.DB.B().Hdel().Key(key).Field("purpose").Build()).Error()
}

func (cache CacheLayer) Set(key string, value string, duration time.Duration) error {
	ctx := context.Background()
	return cache.DB.Do(ctx, cache.DB.B().Set().Key(key).Value(value).Ex(duration).Build()).Error()
}

func (cache CacheLayer) SetHash(key string, values map[string]string, duration time.Duration) error {
	ctx := context.Background()
	hash := cache.DB.B().Hset().Key(key).FieldValue()
	for field, value := range values {
		hash = hash.FieldValue(field, value)
	}
	err := cache.DB.Do(ctx, hash.Build()).Error()
	if err != nil {
		return err
	}
	return cache.DB.Do(ctx, cache.DB.B().Expire().Key(key).Seconds(int64(duration.Seconds())).Build()).Error()
}

func (cache CacheLayer) RenameKey(oldKey string, newKey string) error {
	ctx := context.Background()
	return cache.DB.Do(ctx, cache.DB.B().Rename().Key(oldKey).Newkey(newKey).Build()).Error()
}

func (cache CacheLayer) IncrementHashField(key string, field string, amount int, expiration time.Time) error {
	ctx := context.Background()
	err := cache.DB.Do(ctx, cache.DB.B().Hincrby().Key(key).Field(field).Increment(int64(amount)).Build()).Error()
	if err != nil {
		return err
	}
	remainingTime := time.Until(expiration)
	return cache.DB.Do(ctx, cache.DB.B().Expire().Key(key).Seconds(int64(remainingTime.Seconds())).Build()).Error()
}

func (cache CacheLayer) IncrementKey(key string, expiration time.Time) error {
	ctx := context.Background()
	err := cache.DB.Do(ctx, cache.DB.B().Incr().Key(key).Build()).Error()
	if err != nil {
		return err
	}
	remainingTime := time.Until(expiration)
	return cache.DB.Do(ctx, cache.DB.B().Expire().Key(key).Seconds(int64(remainingTime.Seconds())).Build()).Error()
}

// Higher-level cache functions

func (cache *CacheClient[client]) setUserSignIn(JWT string) error {
	err := cache.raw.Set("JWT:"+JWT, "valid", cacheUserSignIn)
	if err != nil {
		Coms.PrintErrStr("Valkey Set Error: " + err.Error())
		return err
	}
	return nil
}

func (cache *CacheClient[client]) getUserSignIn(JWT string) (string, error) {
	return cache.raw.Get("JWT:" + JWT)
}

func (cache *CacheClient[client]) deleteUserSignIn(JWT string) error {
	return cache.raw.Delete(JWT)
}

func (cache *CacheClient[client]) incrementAnalytics(serviceID string, resource string, country string, ip string, responseCode int) error {
	for _, timeStep := range cacheAnalyticsTime {
		recordTime := timeStep.key(0)
		expiration, err := time.Parse(time.RFC3339, timeStep.key(timeStep.maximumUnits))
		quantity := strconv.Itoa(timeStep.maximumUnits)
		if err != nil {
			Coms.PrintErrStr("Could not parse expiration time of " + recordTime + ": " + err.Error())
			return err
		}

		err = cache.raw.IncrementKey("Analytics:"+serviceID+":"+quantity+":"+recordTime+":quantity", expiration)
		if err != nil {
			Coms.PrintErrStr("Could not increment minute analytics key: " + err.Error())
			return err
		}
		err = cache.raw.IncrementHashField("Analytics:"+serviceID+":"+quantity+":"+recordTime+":country", country, 1, expiration)
		if err != nil {
			Coms.PrintErrStr("Could not increment minute analytics country: " + err.Error())
			return err
		}
		err = cache.raw.IncrementHashField("Analytics:"+serviceID+":"+quantity+":"+recordTime+":ip", ip, 1, expiration)
		if err != nil {
			Coms.PrintErrStr("Could not increment minute analytics ip: " + err.Error())
			return err
		}
		err = cache.raw.IncrementHashField("Analytics:"+serviceID+":"+quantity+":"+recordTime+":resource", resource, 1, expiration)
		if err != nil {
			Coms.PrintErrStr("Could not increment minute analytics resource: " + err.Error())
			return err
		}
		err = cache.raw.IncrementHashField("Analytics:"+serviceID+":"+quantity+":"+recordTime+":response_code", strconv.Itoa(responseCode), 1, expiration)
		if err != nil {
			Coms.PrintErrStr("Could not increment minute analytics response code: " + err.Error())
			return err
		}
	}
	return nil
}

func (cache *CacheClient[client]) getAnalyticsService(service ServiceData, timeStep AnalyticsTimeStep) map[int]Analytic {
	analytics := map[int]Analytic{}
	quantity := strconv.Itoa(timeStep.maximumUnits)
	for timePeriod := range timeStep.maximumUnits {
		quantityRaw, err := cache.raw.Get("Analytics:" + service.ID + ":" + quantity + ":" + timeStep.key(-timePeriod) + ":quantity")
		if err != nil {
			continue
		}
		countryRaw, err := cache.raw.GetHash("Analytics:" + service.ID + ":" + quantity + ":" + timeStep.key(-timePeriod) + ":country")
		if err != nil {
			continue
		}
		ipRaw, err := cache.raw.GetHash("Analytics:" + service.ID + ":" + quantity + ":" + timeStep.key(-timePeriod) + ":ip")
		if err != nil {
			continue
		}
		resourceRaw, err := cache.raw.GetHash("Analytics:" + service.ID + ":" + quantity + ":" + timeStep.key(-timePeriod) + ":resource")
		if err != nil {
			continue
		}
		quantity, err := strconv.Atoi(quantityRaw)
		if err != nil {
			continue
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

		analytics[timePeriod] = Analytic{
			Quantity: quantity,
			Country:  country,
			IP:       ip,
			Resource: resource,
		}
	}

	return analytics
}
