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

	IncrementHashField(key string, field string, amount int) error
	IncrementKey(key string) error
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
	key          string
	maximumUnits int
}

// timeToNextStep returns the time until the next step of the analytics in UTC
func (timeStep AnalyticsTimeStep) timeToNextStep() time.Duration {
	now := time.Now()
	switch timeStep.key {
	case "Minute":
		nextMinute := now.Truncate(time.Minute).Add(time.Minute)
		return nextMinute.Sub(now)
	case "Hour":
		nextHour := now.Add(time.Hour).Truncate(time.Hour)
		return nextHour.Sub(now)
	case "Day":
		midnightTonight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		return midnightTonight.Sub(now)
	case "Month":
		nextMonth := now.AddDate(0, 1, 0)
		beginningOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, now.Location())
		return beginningOfNextMonth.Sub(now)
	default:
		return time.Duration(0)
	}
}

const (
	cacheUserSignIn = JWT.LoginDuration
)

var (
	cacheAnalyticsMinute = AnalyticsTimeStep{key: "Minute", maximumUnits: 60}
	cacheAnalyticsHour   = AnalyticsTimeStep{key: "Hour", maximumUnits: 24}
	cacheAnalyticsDay    = AnalyticsTimeStep{key: "Day", maximumUnits: 30}
	cacheAnalyticsMonth  = AnalyticsTimeStep{key: "Month", maximumUnits: 12}
	cacheAnalyticsTime   = []AnalyticsTimeStep{cacheAnalyticsMinute, cacheAnalyticsHour, cacheAnalyticsDay, cacheAnalyticsMonth}
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

func (cache CacheLayer) IncrementHashField(key string, field string, amount int) error {
	ctx := context.Background()
	return cache.DB.Do(ctx, cache.DB.B().Hincrby().Key(key).Field(field).Increment(int64(amount)).Build()).Error()
}

func (cache CacheLayer) IncrementKey(key string) error {
	ctx := context.Background()
	return cache.DB.Do(ctx, cache.DB.B().Incr().Key(key).Build()).Error()
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
		err := cache.raw.IncrementKey("Analytics:" + serviceID + ":" + timeStep.key + "1:quantity")
		if err != nil {
			Coms.PrintErrStr("Could not increment minute analytics key: " + err.Error())
			return err
		}
		err = cache.raw.IncrementHashField("Analytics:"+serviceID+":"+timeStep.key+"1:country", country, 1)
		if err != nil {
			Coms.PrintErrStr("Could not increment minute analytics country: " + err.Error())
			return err
		}
		err = cache.raw.IncrementHashField("Analytics:"+serviceID+":"+timeStep.key+"1:ip", ip, 1)
		if err != nil {
			Coms.PrintErrStr("Could not increment minute analytics ip: " + err.Error())
			return err
		}
		err = cache.raw.IncrementHashField("Analytics:"+serviceID+":"+timeStep.key+"1:resource", resource, 1)
		if err != nil {
			Coms.PrintErrStr("Could not increment minute analytics resource: " + err.Error())
			return err
		}
		err = cache.raw.IncrementHashField("Analytics:"+serviceID+":"+timeStep.key+"1:response_code", strconv.Itoa(responseCode), 1)
		if err != nil {
			Coms.PrintErrStr("Could not increment minute analytics response code: " + err.Error())
			return err
		}
	}
	return nil
}

func (cache *CacheClient[client]) advanceAnalytics(timeStep AnalyticsTimeStep, services []ServiceLink) error {
	for _, service := range services {
		// Remove last unit
		err := cache.raw.Delete("Analytics:" + service.ID + ":" + timeStep.key + strconv.Itoa(timeStep.maximumUnits) + ":quantity")
		if err != nil {
			Coms.PrintErrStr("Could not remove last " + timeStep.key + " analytics quantity for service " + service.Title + ": " + err.Error())
			return err
		}
		err = cache.raw.DeleteHash("Analytics:" + service.ID + ":" + timeStep.key + strconv.Itoa(timeStep.maximumUnits) + ":country")
		if err != nil {
			Coms.PrintErrStr("Could not remove last " + timeStep.key + " analytics country for service " + service.Title + ": " + err.Error())
			return err
		}
		err = cache.raw.DeleteHash("Analytics:" + service.ID + ":" + timeStep.key + strconv.Itoa(timeStep.maximumUnits) + ":ip")
		if err != nil {
			Coms.PrintErrStr("Could not remove last " + timeStep.key + " analytics ip for service " + service.Title + ": " + err.Error())
			return err
		}
		err = cache.raw.DeleteHash("Analytics:" + service.ID + ":" + timeStep.key + strconv.Itoa(timeStep.maximumUnits) + ":resource")
		if err != nil {
			Coms.PrintErrStr("Could not remove last " + timeStep.key + " analytics resource for service " + service.Title + ": " + err.Error())
			return err
		}
		err = cache.raw.DeleteHash("Analytics:" + service.ID + ":" + timeStep.key + strconv.Itoa(timeStep.maximumUnits) + ":response_code")
		if err != nil {
			Coms.PrintErrStr("Could not remove last " + timeStep.key + " analytics response code for service " + service.Title + ": " + err.Error())
			return err
		}
		// Advance current analytics to next time step unit
		for i := timeStep.maximumUnits - 1; i > 0; i-- {
			cache.raw.RenameKey("Analytics:"+service.ID+":"+timeStep.key+strconv.Itoa(i)+":quantity", "Analytics:"+service.ID+":"+timeStep.key+strconv.Itoa(i+1)+":quantity")
			cache.raw.RenameKey("Analytics:"+service.ID+":"+timeStep.key+strconv.Itoa(i)+":country", "Analytics:"+service.ID+":"+timeStep.key+strconv.Itoa(i+1)+":country")
			cache.raw.RenameKey("Analytics:"+service.ID+":"+timeStep.key+strconv.Itoa(i)+":ip", "Analytics:"+service.ID+":"+timeStep.key+strconv.Itoa(i+1)+":ip")
			cache.raw.RenameKey("Analytics:"+service.ID+":"+timeStep.key+strconv.Itoa(i)+":resource", "Analytics:"+service.ID+":"+timeStep.key+strconv.Itoa(i+1)+":resource")
			cache.raw.RenameKey("Analytics:"+service.ID+":"+timeStep.key+strconv.Itoa(i)+":response_code", "Analytics:"+service.ID+":"+timeStep.key+strconv.Itoa(i+1)+":response_code")

		}
	}
	return nil
}
