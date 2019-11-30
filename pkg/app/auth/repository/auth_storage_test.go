package repository

import (
	"fmt"
	"testing"
	"time"

	. "2019_2_IBAT/pkg/pkg/models"

	"github.com/alicebob/miniredis"
	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDBUserStorage_1(t *testing.T) {
	conn := redigomock.NewConn()
	pool := &redis.Pool{
		Dial:    func() (redis.Conn, error) { return conn, nil },
		MaxIdle: 10,
	}
	sessManager := NewSessionManager(pool)

	record := AuthStorageValue{
		ID:      uuid.New(),
		Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
		Role:    SeekerStr,
	}

	conn.Command("SET").ExpectError(fmt.Errorf("Low level error!"))

	_, _, err := sessManager.Set(record.ID, record.Role)
	if err == nil {
		t.Errorf("Expected error")
		return
	}
}

func TestDBUserStorage_2(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	sessManager := NewSessionManager(RedNewPool(s.Addr()))

	id := uuid.New()
	role := SeekerStr
	_, cookie, err := sessManager.Set(id, role)

	if err != nil {
		t.Errorf("Unexpected err %s\n", err.Error())
		return
	}

	authRec, ok := sessManager.Get(cookie)

	if !ok {
		t.Error("Expected ok\n")
		return
	}

	require.Equal(t, id, authRec.ID, "The two values should be the same.")
	require.Equal(t, role, authRec.Role, "The two values should be the same.")

	ok = sessManager.Delete(cookie)

	if !ok {
		t.Error("Expected ok\n")
		return
	}
}

func TestDBUserStorage_3(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	sessManager := NewSessionManager(RedNewPool(s.Addr()))

	cookie := "cookie"

	ok := sessManager.Delete(cookie)
	if ok {
		t.Error("Expected ok == false\n")
		return
	}

	_, ok = sessManager.Get(cookie)

	if ok {
		t.Error("Expected ok == false\n")
		return
	}
}
