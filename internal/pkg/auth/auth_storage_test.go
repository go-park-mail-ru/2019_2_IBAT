package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMapAuthStorage_Get(t *testing.T) {

	type fields struct {
		Storage map[string]StorageValue
	}
	type args struct {
		cookie string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   StorageValue
		want1  bool
	}{
		{
			name: "GetTest1",
			fields: fields{
				map[string]StorageValue{
					"werwe": {
						ID:      1,
						Expires: time.Now().In(Loc).Add(-1 * time.Hour).Format(TimeFormat),
					},
					"werre": {
						ID:      2,
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					},
				},
			},
			args: args{
				cookie: "werre",
			},
			want: StorageValue{
				ID:      2,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
			want1: true,
		},
		{
			name: "GetTest2",
			fields: fields{
				map[string]StorageValue{
					"werwe": {
						ID:      1,
						Expires: time.Now().In(Loc).Format(TimeFormat),
					},
					"werre": {
						ID:      2,
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					},
				},
			},
			args: args{
				cookie: "werwe",
			},
			want:  StorageValue{},
			want1: false,
		},
		{
			name: "GetTest3",
			fields: fields{
				map[string]StorageValue{
					"werwe": {
						ID:      1,
						Expires: time.Now().In(Loc).Format(TimeFormat),
					},
					"value": {
						ID:      3,
						Expires: time.Now().In(Loc).Add(48 * time.Hour).Format(TimeFormat),
					},
				},
			},
			args: args{
				cookie: "value",
			},
			want: StorageValue{
				ID:      3,
				Expires: time.Now().In(Loc).Add(48 * time.Hour).Format(TimeFormat),
			},
			want1: true,
		},
		{
			name: "GetTest4",
			fields: fields{
				map[string]StorageValue{
					"werwe": {
						ID:      1,
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					},
					"werre": {
						ID:      2,
						Expires: time.Now().In(Loc).Add(-24 * time.Hour).Format(TimeFormat),
					},
				},
			},
			args: args{
				cookie: "werre",
			},
			want:  StorageValue{},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &MapAuthStorage{
				Storage: tt.fields.Storage,
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
		Storage map[string]StorageValue
	}
	type args struct {
		id uint64
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
				map[string]StorageValue{
					"werwe": {
						ID:      1,
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					},
					"werre": {
						ID:      2,
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					},
				},
			},
			args: args{11},
			want: "any string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &MapAuthStorage{
				Storage: tt.fields.Storage,
			}
			got := st.Set(tt.args.id)

			authData, ok := st.Get(got)
			if !ok {
				t.Errorf("Error during setting a cookie")
			}
			require.Equal(t, tt.args.id, authData.ID, "The two values should be the same.") //switch ids
		})
	}
}

func TestMapAuthStorage_Delete(t *testing.T) {
	type fields struct {
		Storage map[string]StorageValue
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
				map[string]StorageValue{
					"werwe": {
						ID:      1,
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					},
					"werre": {
						ID:      2,
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					},
				},
			},
			args: args{"werre"},
			want: "any string",
		},
		{
			name: "DeleteTest2",
			fields: fields{
				map[string]StorageValue{
					"werwe": {
						ID:      1,
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					},
					"werre": {
						ID:      2,
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					},
				},
			},
			args: args{"werwe"},
			want: "any string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := MapAuthStorage{
				Storage: tt.fields.Storage,
			}
			st.Delete(tt.args.cookie)
			authData, ok := st.Get(tt.args.cookie)
			if (ok && authData != StorageValue{}) {
				t.Errorf("Error during deleting a session")
				t.Log(ok)
				t.Log(authData)
			}
		})
	}
}