package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"

	. "2019_2_IBAT/pkg/pkg/models"
)

// const TimeFormat = time.RFC3339
const CookieLength = 32

var Loc *time.Location

func init() {
	Loc, _ = time.LoadLocation("Europe/Moscow")
}

type SessionManager struct {
	// redisConn redis.Conn
	redisPool *redis.Pool
}

func NewSessionManager(pool *redis.Pool) *SessionManager {
	return &SessionManager{
		redisPool: pool,
	}
}

func (st *SessionManager) Get(cookie string) (AuthStorageValue, bool) {
	// log.Println("AuthStorage: Get started")
	redisConn := st.redisPool.Get()
	defer redisConn.Close()

	data, err := redis.Bytes(redisConn.Do("GET", cookie))
	// defer redisConn.Close()
	if err != nil {
		log.Println("AuthStorage: Can not get auth info:", err)
		return AuthStorageValue{}, false
	}

	record := AuthStorageValue{}
	err = json.Unmarshal(data, &record)

	if err != nil {
		fmt.Println("AuthStorage: Unmarshalling error") //
		return AuthStorageValue{}, false
	} //cannot be error

	expiresAt, err := time.Parse(TimeFormat, record.Expires)

	if err != nil {
		fmt.Println("AuthStorage: Time parse error")
		return AuthStorageValue{}, false
	} //cannot be error

	now := time.Now().In(Loc)
	diff := expiresAt.Sub(now)

	if diff < 0 {
		// delete(st.Storage, cookie)
		_, _ = redis.String(redisConn.Do("DEL", cookie))

		return AuthStorageValue{}, false
	}

	return record, true
}

func (st *SessionManager) Set(id uuid.UUID, class string) (AuthStorageValue, string, error) {
	expires := time.Now().In(Loc).Add(24 * time.Hour)

	record := AuthStorageValue{
		ID:      id,
		Expires: expires.Format(TimeFormat),
		Role:    class,
	}

	cookie := generateCookie()
	dataSerialized, _ := json.Marshal(record)

	redisConn := st.redisPool.Get()
	defer redisConn.Close()
	result, err := redis.String(redisConn.Do("SET", cookie, dataSerialized))

	if err != nil {
		return AuthStorageValue{}, "", err
	}

	if result != "OK" {
		return AuthStorageValue{}, "", fmt.Errorf("result not OK")
	}

	return record, cookie, nil
}

func (st *SessionManager) Delete(cookie string) bool {
	redisConn := st.redisPool.Get()
	defer redisConn.Close()

	_, err := redis.Int(redisConn.Do("DEL", cookie))

	if err != nil {
		return false
	}

	fmt.Println("Cookie was deleted")
	return true
}

func generateCookie() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" + "0123456789")

	var b strings.Builder
	for i := 0; i < CookieLength; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	return b.String()
}

func RedNewPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}
