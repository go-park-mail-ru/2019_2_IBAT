package auth

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

// const TimeFormat = time.RFC3339
// const CookieLength = 10

// var Loc *time.Location

func init() {
	Loc, _ = time.LoadLocation("Europe/Moscow")
}

type SessionManager struct {
	redisConn redis.Conn
}

func NewSessionManager(conn redis.Conn) *SessionManager {
	return &SessionManager{
		redisConn: conn,
	}
}

func (st SessionManager) Get(cookie string) (AuthStorageValue, bool) {
	data, err := redis.Bytes(st.redisConn.Do("GET", cookie))
	if err != nil {
		log.Println("cant get data:", err)
		return AuthStorageValue{}, false
	}

	record := AuthStorageValue{}
	err = json.Unmarshal(data, &record)

	if err != nil {
		fmt.Println("Unmarshalling error") //
		return AuthStorageValue{}, false
	} //cannot be error

	expiresAt, err := time.Parse(TimeFormat, record.Expires)

	if err != nil {
		fmt.Println("Parse error")
		return AuthStorageValue{}, false
	} //cannot be error

	now := time.Now().In(Loc)
	diff := expiresAt.Sub(now)

	if diff < 0 {
		// delete(st.Storage, cookie)
		_, _ = redis.String(st.redisConn.Do("DEL", cookie))

		return AuthStorageValue{}, false
	}

	return record, true
}

func (st SessionManager) Set(id uuid.UUID, class string) (AuthStorageValue, string, error) {
	expires := time.Now().In(Loc).Add(24 * time.Hour)

	record := AuthStorageValue{
		ID:      id,
		Expires: expires.Format(TimeFormat),
		Role:    class,
	}

	cookie := generateCookie()
	dataSerialized, _ := json.Marshal(record)

	result, err := redis.String(st.redisConn.Do("SET", cookie, dataSerialized))

	if err != nil {
		return AuthStorageValue{}, "", err
	}

	if result != "OK" {
		return AuthStorageValue{}, "", fmt.Errorf("result not OK")
	}

	return record, cookie, nil
}

func (st SessionManager) Delete(cookie string) bool {
	_, err := redis.Int(st.redisConn.Do("DEL", cookie))

	if err != nil {
		return false
	}

	fmt.Println("Cookie was deleted")
	return true
}

// func generateCookie() string {
// 	rand.Seed(time.Now().UnixNano())
// 	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
// 		"abcdefghijklmnopqrstuvwxyz" + "0123456789")

// 	var b strings.Builder
// 	for i := 0; i < CookieLength; i++ {
// 		b.WriteRune(chars[rand.Intn(len(chars))])
// 	}

// 	return b.String()
// }
