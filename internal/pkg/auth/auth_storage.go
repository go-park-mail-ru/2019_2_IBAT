package auth

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const TimeFormat = time.RFC3339
const CookieLength = 10

var Loc *time.Location

func init() {
	Loc, _ = time.LoadLocation("Europe/Moscow")
}

type MapAuthStorage struct {
	Storage map[string]AuthStorageValue
	Mu      *sync.Mutex
}

func (st MapAuthStorage) Get(cookie string) (AuthStorageValue, bool) {
	st.Mu.Lock()
	record, ok := st.Storage[cookie]
	st.Mu.Unlock()

	if !ok {
		fmt.Println("No such session error")
		return AuthStorageValue{}, false
	}

	expiresAt, err := time.Parse(TimeFormat, record.Expires)

	if err != nil {
		fmt.Println("Parse error")
		return AuthStorageValue{}, false
	} //cannot be error

	now := time.Now().In(Loc)
	diff := expiresAt.Sub(now)

	if diff < 0 {
		st.Mu.Lock()
		delete(st.Storage, cookie)
		st.Mu.Unlock()

		return AuthStorageValue{}, false
	}

	return record, true
}

func (st MapAuthStorage) Set(id uuid.UUID, class string) (AuthStorageValue, string) {
	expires := time.Now().In(Loc).Add(24 * time.Hour)

	record := AuthStorageValue{
		ID:      id,
		Expires: expires.Format(TimeFormat),
		Class:   class,
	}

	cookie := generateCookie()
	st.Mu.Lock()
	st.Storage[cookie] = record
	st.Mu.Unlock()

	// fmt.Println(st.Storage)
	return record, cookie
}

func (st MapAuthStorage) Delete(cookie string) bool {
	st.Mu.Lock()
	fmt.Println("Cookie was deleted")
	delete(st.Storage, cookie)
	st.Mu.Unlock()

	_, ok := st.Storage[cookie]
	if ok {
		return false
	}

	return true
}

func generateCookie() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	var b strings.Builder
	for i := 0; i < CookieLength; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
