package main

import (
	"context"
	"errors"
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
	Set(key string, value string, duration CacheType) error
	Get(key string) (string, CacheType, error)
	Delete(key string) error

	SetHash(key string, values map[string]string, duration CacheType) error
	GetHash(key string) (map[string]string, CacheType, error)
	DeleteHash(key string) error

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

type CacheType struct {
	duration time.Duration
	purpose  string
}

var (
	cachePasswordSet   CacheType           = CacheType{duration: time.Minute * 15, purpose: "Set Password"}
	cacheChangeEmail   CacheType           = CacheType{duration: time.Minute * 15, purpose: "Change Email"}
	cacheNewUserSignUp CacheType           = CacheType{duration: time.Minute * 15, purpose: "User Sign Up"}
	cacheUserSignIn    CacheType           = CacheType{duration: JWT.LoginDuration, purpose: "User Sign In"}
	cacheAnalyticsTime []AnalyticsTimeStep = []AnalyticsTimeStep{cacheAnalyticsMinute, cacheAnalyticsHour, cacheAnalyticsDay, cacheAnalyticsMonth}
)

type AnalyticsTimeStep string

const (
	cacheAnalyticsMinute AnalyticsTimeStep = "Minute"
	cacheAnalyticsHour   AnalyticsTimeStep = "Hour"
	cacheAnalyticsDay    AnalyticsTimeStep = "Day"
	cacheAnalyticsMonth  AnalyticsTimeStep = "Month"
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

func (cache CacheLayer) Get(key string) (string, CacheType, error) {
	ctx := context.Background()
	rawResult, err := cache.DB.Do(ctx, cache.DB.B().Hgetall().Key(key).Build()).AsStrMap()
	if err != nil {
		return "", CacheType{}, errors.New("failed to get value from cache: " + err.Error())
	}
	value := rawResult["value"]
	cacheType, err := cache.getCacheType(rawResult["purpose"])

	return value, cacheType, err
}

func (cache CacheLayer) GetHash(key string) (map[string]string, CacheType, error) {
	ctx := context.Background()
	rawResult, err := cache.DB.Do(ctx, cache.DB.B().Hgetall().Key(key).Build()).AsStrMap()
	if err != nil {
		return nil, CacheType{}, err
	}
	cacheType, err := cache.getCacheType(rawResult["purpose"])

	return rawResult, cacheType, err
}

func (cache CacheLayer) Delete(key string) error {
	return cache.DB.Do(context.Background(), cache.DB.B().Hdel().Key(key).Field("value").Field("purpose").Build()).Error()
}

func (cache CacheLayer) DeleteHash(key string) error {
	return cache.DB.Do(context.Background(), cache.DB.B().Hdel().Key(key).Field("purpose").Build()).Error()
}

func (cache CacheLayer) Set(key string, value string, cacheType CacheType) error {
	ctx := context.Background()
	Coms.Println("Valkey Set with type: " + cacheType.purpose)
	err := cache.DB.Do(ctx, cache.DB.B().Hset().Key(key).FieldValue().FieldValue("value", value).FieldValue("purpose", cacheType.purpose).Build()).Error()
	if err != nil {
		Coms.PrintErrStr("Valkey Set Error: " + err.Error())
		return err
	}
	return cache.DB.Do(ctx, cache.DB.B().Expire().Key(key).Seconds(int64(cacheType.duration.Seconds())).Build()).Error()
}

func (cache CacheLayer) SetHash(key string, values map[string]string, duration CacheType) error {
	ctx := context.Background()
	hash := cache.DB.B().Hset().Key(key).FieldValue().FieldValue("purpose", duration.purpose)
	for field, value := range values {
		hash = hash.FieldValue(field, value)
	}
	err := cache.DB.Do(ctx, hash.Build()).Error()
	if err != nil {
		return err
	}
	return cache.DB.Do(ctx, cache.DB.B().Expire().Key(key).Seconds(int64(duration.duration.Seconds())).Build()).Error()
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
	err := cache.raw.Set(JWT, "valid", cacheUserSignIn)
	if err != nil {
		Coms.PrintErrStr("Valkey Set Error: " + err.Error())
		return err
	}
	return nil
}

func (cache *CacheClient[client]) getUserSignIn(JWT string) (string, error) {
	userData, cacheType, err := cache.raw.Get(JWT)
	if err != nil {
		return "", errors.New("failed to get user ID from JWT: " + err.Error())
	}
	if cacheType != cacheUserSignIn {
		return "", errors.New("invalid cache type")
	}
	return userData, nil
}

func (cache *CacheClient[client]) deleteUserSignIn(JWT string) error {
	return cache.raw.Delete(JWT)
}

func (cache *CacheClient[client]) incrementAnalytics(serviceID string, resource string, country string, ip string, responseCode int) error {
	for _, timeStep := range cacheAnalyticsTime {
		err := cache.raw.IncrementKey("Analytics:" + serviceID + ":" + string(timeStep) + ":" + "quantity")
		if err != nil {
			Coms.PrintErrStr("Could not increment minute analytics key: " + err.Error())
			return err
		}
		err = cache.raw.IncrementHashField("Analytics:"+serviceID+":"+string(timeStep)+":"+"country", country, 1)
		if err != nil {
			Coms.PrintErrStr("Could not increment minute analytics country: " + err.Error())
			return err
		}
		err = cache.raw.IncrementHashField("Analytics:"+serviceID+":"+string(timeStep)+":"+"ip", ip, 1)
		if err != nil {
			Coms.PrintErrStr("Could not increment minute analytics ip: " + err.Error())
			return err
		}
		err = cache.raw.IncrementHashField("Analytics:"+serviceID+":"+string(timeStep)+":"+"resource", resource, 1)
		if err != nil {
			Coms.PrintErrStr("Could not increment minute analytics resource: " + err.Error())
			return err
		}
		err = cache.raw.IncrementHashField("Analytics:"+serviceID+":"+string(timeStep)+":"+"response_code", strconv.Itoa(responseCode), 1)
		if err != nil {
			Coms.PrintErrStr("Could not increment minute analytics response code: " + err.Error())
			return err
		}
	}
	return nil
}

// Utilities

func (cache CacheLayer) getCacheType(purpose string) (CacheType, error) {
	switch purpose {
	case cachePasswordSet.purpose:
		return cachePasswordSet, nil
	case cacheChangeEmail.purpose:
		return cacheChangeEmail, nil
	case cacheNewUserSignUp.purpose:
		return cacheNewUserSignUp, nil
	case cacheUserSignIn.purpose:
		return cacheUserSignIn, nil
	default:
		return CacheType{}, errors.New("invalid cache type")
	}
}
