package auth

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const TimeFormat = time.RFC3339

var Loc *time.Location

func init() {
	Loc, _ = time.LoadLocation("Europe/Moscow") //should remove
}

type MapAuthStorage struct {
	Storage map[string]StorageValue
}

func (st MapAuthStorage) Get(cookie string) (StorageValue, bool) { //pointer receiver?
	record, ok := st.Storage[cookie]
	if !ok {
		fmt.Println("No such session error")
		return StorageValue{}, false
	}

	expiresAt, err := time.Parse(TimeFormat, record.Expires)

	if err != nil {
		fmt.Println("Parse error")
		return StorageValue{}, false
	} //cannot be error

	now := time.Now().In(Loc)
	diff := expiresAt.Sub(now)

	if diff < 0 {
		delete(st.Storage, cookie)
		return StorageValue{}, false
	}

	return record, true
}

func (st MapAuthStorage) Set(id uint64) string {
	//id collision should be solved
	expires := time.Now().In(Loc).Add(24 * time.Hour)

	record := StorageValue{
		ID:      id,
		Expires: expires.Format(TimeFormat),
	}

	cookie := generateCookie()
	st.Storage[cookie] = record
	fmt.Println(st.Storage)
	return cookie
}

func (st MapAuthStorage) Delete(cookie string) string {
	expires := time.Now().In(Loc).Format(TimeFormat)
	delete(st.Storage, cookie)
	return expires
}

func generateCookie() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 5
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
