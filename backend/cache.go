package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	Config "github.com/benjaminRoberts01375/Web-Tech-Stack/config"
	JWT "github.com/benjaminRoberts01375/Web-Tech-Stack/jwt"
	Printing "github.com/benjaminRoberts01375/Web-Tech-Stack/logging"
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
	time         func(step int) time.Time
	maximumUnits int
}

func (analytics AnalyticsTimeStep) timeStr(step int) string {
	return analytics.time(step).Format(time.RFC3339)
}

const (
	cacheUserSignIn = JWT.LoginDuration
)

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

	Config.ReadExternalConfig("valkey.json", &cache)
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
		Printing.PrintErr(err)
		panic("Could not connect to Valkey: " + err.Error())
	}
	cache.DB = client
	Printing.Println("Connected to Valkey")
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
		Printing.PrintErrStr("Valkey Set Error: " + err.Error())
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
		recordTime := timeStep.timeStr(0)
		expiration := timeStep.time(timeStep.maximumUnits)
		quantity := strconv.Itoa(timeStep.maximumUnits)
		err := cache.raw.IncrementKey("Analytics:"+serviceID+":"+quantity+":"+recordTime+":quantity", expiration)
		if err != nil {
			Printing.PrintErrStr("Could not increment minute analytics key: " + err.Error())
			return err
		}
		err = cache.raw.IncrementHashField("Analytics:"+serviceID+":"+quantity+":"+recordTime+":country", country, 1, expiration)
		if err != nil {
			Printing.PrintErrStr("Could not increment minute analytics country: " + err.Error())
			return err
		}
		err = cache.raw.IncrementHashField("Analytics:"+serviceID+":"+quantity+":"+recordTime+":ip", ip, 1, expiration)
		if err != nil {
			Printing.PrintErrStr("Could not increment minute analytics ip: " + err.Error())
			return err
		}
		err = cache.raw.IncrementHashField("Analytics:"+serviceID+":"+quantity+":"+recordTime+":resource", resource, 1, expiration)
		if err != nil {
			Printing.PrintErrStr("Could not increment minute analytics resource: " + err.Error())
			return err
		}
		err = cache.raw.IncrementHashField("Analytics:"+serviceID+":"+quantity+":"+recordTime+":response_code", strconv.Itoa(responseCode), 1, expiration)
		if err != nil {
			Printing.PrintErrStr("Could not increment minute analytics response code: " + err.Error())
			return err
		}
	}
	return nil
}

func (cache *CacheClient[client]) getAnalyticsService(service ServiceData, timeStep AnalyticsTimeStep) map[time.Time]Analytic {
	analytics := map[time.Time]Analytic{}
	quantity := strconv.Itoa(timeStep.maximumUnits)
	for timePeriod := range timeStep.maximumUnits {
		quantityRaw, err := cache.raw.Get("Analytics:" + service.ID + ":" + quantity + ":" + timeStep.timeStr(-timePeriod) + ":quantity")
		if err != nil {
			continue
		}
		countryRaw, err := cache.raw.GetHash("Analytics:" + service.ID + ":" + quantity + ":" + timeStep.timeStr(-timePeriod) + ":country")
		if err != nil {
			continue
		}
		ipRaw, err := cache.raw.GetHash("Analytics:" + service.ID + ":" + quantity + ":" + timeStep.timeStr(-timePeriod) + ":ip")
		if err != nil {
			continue
		}
		resourceRaw, err := cache.raw.GetHash("Analytics:" + service.ID + ":" + quantity + ":" + timeStep.timeStr(-timePeriod) + ":resource")
		if err != nil {
			continue
		}
		responseCodesRaw, err := cache.raw.GetHash("Analytics:" + service.ID + ":" + quantity + ":" + timeStep.timeStr(-timePeriod) + ":response_code")
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
			Quantity:     quantity,
			Country:      country,
			IP:           ip,
			Resource:     resource,
			ResponseCode: responseCodes,
		}
	}

	return analytics
}

func (cache *CacheClient[client]) deleteService(service ServiceLink) error {
	for _, timeStep := range cacheAnalyticsTime {
		recordTime := timeStep.timeStr(0)
		quantity := strconv.Itoa(timeStep.maximumUnits)
		err := cache.raw.Delete("Analytics:" + service.ID + ":" + quantity + ":" + recordTime + ":quantity")
		if err != nil {
			Printing.PrintErrStr("Could not delete minute analytics key: " + err.Error())
			return err
		}
		err = cache.raw.DeleteHash("Analytics:" + service.ID + ":" + quantity + ":" + recordTime + ":country")
		if err != nil {
			Printing.PrintErrStr("Could not delete minute analytics country: " + err.Error())
			return err
		}
		err = cache.raw.DeleteHash("Analytics:" + service.ID + ":" + quantity + ":" + recordTime + ":ip")
		if err != nil {
			Printing.PrintErrStr("Could not delete minute analytics ip: " + err.Error())
			return err
		}
		err = cache.raw.DeleteHash("Analytics:" + service.ID + ":" + quantity + ":" + recordTime + ":resource")
		if err != nil {
			Printing.PrintErrStr("Could not delete minute analytics resource: " + err.Error())
			return err
		}
		err = cache.raw.DeleteHash("Analytics:" + service.ID + ":" + quantity + ":" + recordTime + ":response_code")
		if err != nil {
			Printing.PrintErrStr("Could not delete minute analytics response code: " + err.Error())
			return err
		}
	}
	return nil
}
