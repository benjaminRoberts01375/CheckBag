package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
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
	cachePasswordSet   CacheType = CacheType{duration: time.Minute * 15, purpose: "Set Password"}
	cacheChangeEmail   CacheType = CacheType{duration: time.Minute * 15, purpose: "Change Email"}
	cacheNewUserSignUp CacheType = CacheType{duration: time.Minute * 15, purpose: "User Sign Up"}
	cacheUserSignIn    CacheType = CacheType{duration: JWT.LoginDuration, purpose: "User Sign In"}
)

const cacheKeyLength = 16

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

// Higher-level cache functions

func (cache *CacheClient[client]) setNewUser(email string) (string, error) {
	id := generateRandomString(cacheKeyLength)
	return id, cache.raw.Set(id, email, cacheNewUserSignUp)
}

func (cache *CacheClient[client]) getAndDeleteNewUser(id string) (string, error) {
	email, cacheType, err := cache.raw.Get(id)
	if err != nil {
		return "", err
	}
	if cacheType != cacheNewUserSignUp {
		return "", errors.New("invalid cache type")
	}
	err = cache.raw.Delete(id)
	if err != nil {
		return "", err
	}
	return email, nil
}

func (cache *CacheClient[client]) setUserSignIn(JWT string, userID string) error {
	err := cache.raw.Set(JWT, userID, cacheUserSignIn)
	if err != nil {
		Coms.PrintErrStr("Valkey Set Error: " + err.Error())
		return err
	}
	return nil
}

func (cache *CacheClient[client]) getUserSignIn(JWT string) (string, error) {
	userID, cacheType, err := cache.raw.Get(JWT)
	if err != nil {
		return "", errors.New("failed to get user ID from JWT: " + err.Error())
	}
	if cacheType != cacheUserSignIn {
		return "", errors.New("invalid cache type")
	}
	return userID, nil
}

func (cache *CacheClient[client]) deleteUserSignIn(JWT string) error {
	return cache.raw.Delete(JWT)
}

func (cache *CacheClient[client]) setChangeEmail(userID string, newEmail string) (string, error) {
	id := generateRandomString(cacheKeyLength)
	return id, cache.raw.Set(id, userID+":"+newEmail, cacheChangeEmail)
}

func (cache *CacheClient[client]) getAndDeleteChangeEmail(id string) (string, string, error) {
	userIDEmail, cacheType, err := cache.raw.Get(id)
	if err != nil {
		return "", "", err
	}
	if cacheType != cacheChangeEmail {
		return "", "", errors.New("invalid cache type")
	}
	err = cache.raw.Delete(id)
	if err != nil {
		return "", "", err
	}
	return userIDEmail, strings.Split(userIDEmail, ":")[1], nil
}

func (cache *CacheClient[client]) setForgotPassword(email string) (string, error) {
	id := generateRandomString(cacheKeyLength)
	return id, cache.raw.Set(id, email, cachePasswordSet)
}

func (cache *CacheClient[client]) getForgotPassword(id string) (string, error) {
	email, cacheType, err := cache.raw.Get(id)
	if err != nil {
		return "", err
	}
	if cacheType != cachePasswordSet {
		return "", errors.New("invalid cache type")
	}
	return email, nil
}

func (cache *CacheClient[client]) getAndDeleteResetPassword(id string) (string, error) {
	email, cacheType, err := cache.raw.Get(id)
	if err != nil {
		return "", err
	}
	if cacheType != cachePasswordSet {
		return "", errors.New("invalid cache type")
	}
	err = cache.raw.Delete(id)
	if err != nil {
		return "", err
	}
	return email, nil
}

// Utilities

func generateRandomString(length int) string {
	// Charset is URL safe and easy to read
	const charset = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ123456789"

	stringBase := make([]byte, length)
	for i := range stringBase {
		stringBase[i] = charset[rand.Intn(len(charset))]
	}
	return string(stringBase)
}

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
