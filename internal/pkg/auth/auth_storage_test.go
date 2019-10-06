package auth

import (
	"sync"
	"testing"
	"time"

	. "2019_2_IBAT/internal/pkg/interfaces"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMapAuthStorage_Get(t *testing.T) {

	type fields struct {
		Storage map[string]AuthStorageValue
		Mu      *sync.Mutex
	}
	type args struct {
		cookie string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   AuthStorageValue
		want1  bool
	}{
		{
			name: "GetTest1",
			fields: fields{
				Storage: map[string]AuthStorageValue{
					"werwe": {
						ID:      uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
						Expires: time.Now().In(Loc).Add(-1 * time.Hour).Format(TimeFormat),
					},
					"werre": {
						ID:      uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					},
				},
				Mu: &sync.Mutex{},
			},
			args: args{
				cookie: "werre",
			},
			want: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
			want1: true,
		},
		{
			name: "GetTest2",
			fields: fields{
				Storage: map[string]AuthStorageValue{
					"werwe": {
						ID:      uuid.MustParse("6aa7b810-9dad-11d1-80b1-00c04fd430c8"),
						Expires: time.Now().In(Loc).Format(TimeFormat),
					},
					"werre": {
						ID:      uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					},
				},
				Mu: &sync.Mutex{},
			},
			args: args{
				cookie: "werwe",
			},
			want:  AuthStorageValue{},
			want1: false,
		},
		{
			name: "GetTest3",
			fields: fields{
				Storage: map[string]AuthStorageValue{
					"werwe": {
						ID:      uuid.MustParse("6aa7b810-9dad-11d1-80b1-00c04fd430c8"),
						Expires: time.Now().In(Loc).Format(TimeFormat),
					},
					"value": {
						ID:      uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
						Expires: time.Now().In(Loc).Add(48 * time.Hour).Format(TimeFormat),
					},
				},
				Mu: &sync.Mutex{},
			},
			args: args{
				cookie: "value",
			},
			want: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
				Expires: time.Now().In(Loc).Add(48 * time.Hour).Format(TimeFormat),
			},
			want1: true,
		},
		{
			name: "GetTest4",
			fields: fields{
				Storage: map[string]AuthStorageValue{
					"werwe": {
						ID:      uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					},
					"werre": {
						ID:      uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
						Expires: time.Now().In(Loc).Add(-24 * time.Hour).Format(TimeFormat),
					},
				},
				Mu: &sync.Mutex{},
			},
			args: args{
				cookie: "werre",
			},
			want:  AuthStorageValue{},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &MapAuthStorage{
				Storage: tt.fields.Storage,
				Mu:      &sync.Mutex{},
			}
			got, got1 := st.Get(tt.args.cookie)

			require.Equal(t, tt.want, got, "The two values should be the same.")

			if got1 != tt.want1 {
				t.Errorf("MapAuthStorage.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMapAuthStorage_Set(t *testing.T) {
	type fields struct {
		Storage map[string]AuthStorageValue
		Mu      *sync.Mutex
	}
	type args struct {
		id    uuid.UUID
		class string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "SetTest1",
			fields: fields{
				Storage: map[string]AuthStorageValue{
					"werwe": {
						ID:      uuid.New(),
						Role:    SeekerStr,
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					},
					"werre": {
						ID:      uuid.New(),
						Role:    EmployerStr,
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					},
				},
				Mu: &sync.Mutex{},
			},
			args: args{
				id:    uuid.New(),
				class: SeekerStr,
			},
			want: "any string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &MapAuthStorage{
				Storage: tt.fields.Storage,
				Mu:      &sync.Mutex{},
			}
			authData, _ := st.Set(tt.args.id, tt.args.class)

			require.Equal(t, tt.args.id, authData.ID, "The two values should be the same.")
		})
	}
}

func TestMapAuthStorage_Delete(t *testing.T) {
	type fields struct {
		Storage map[string]AuthStorageValue
		Mu      *sync.Mutex
	}
	type args struct {
		cookie string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "DeleteTest1",
			fields: fields{
				Storage: map[string]AuthStorageValue{
					"werwe": {
						ID:      uuid.New(),
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					},
					"werre": {
						ID:      uuid.New(),
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					},
				},
				Mu: &sync.Mutex{},
			},
			args: args{"werre"},
			want: "any string",
		},
		{
			name: "DeleteTest2",
			fields: fields{
				Storage: map[string]AuthStorageValue{
					"werwe": {
						ID:      uuid.New(),
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					},
					"werre": {
						ID:      uuid.New(),
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					},
				},
				Mu: &sync.Mutex{},
			},
			args: args{"werwe"},
			want: "any string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := MapAuthStorage{
				Storage: tt.fields.Storage,
				Mu:      &sync.Mutex{},
			}
			st.Delete(tt.args.cookie)
			authData, ok := st.Get(tt.args.cookie)
			if (ok && authData != AuthStorageValue{}) {
				t.Errorf("Error during deleting a session")
				t.Log(ok)
				t.Log(authData)
			}
		})
	}
}
