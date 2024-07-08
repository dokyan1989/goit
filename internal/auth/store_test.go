package auth

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/dokyan1989/goit/misc/t/seeder"
	"github.com/go-faker/faker/v4"
	"github.com/google/go-cmp/cmp"
)

func TestStore_CreateUser(t *testing.T) {
	fakeStatus := func() string {
		v := []string{"NEW", "ACTIVE", "INACTIVE"}
		p, err := faker.RandomInt(0, len(v)-1, 1)
		if err != nil {
			return v[0]
		}
		return v[p[0]]
	}

	tests := []struct {
		name    string
		seedURL string
		params  CreateUserParams
		want    int64
		wantErr bool
	}{
		{
			name: "create a new user",
			params: CreateUserParams{
				UserName:  faker.Name(),
				Password:  faker.Password(),
				Email:     faker.Email(),
				FirstName: faker.FirstName(),
				LastName:  faker.LastName(),
				Status:    fakeStatus(),
			},
			want:    1,
			wantErr: false,
		},
		{
			name:    "create a duplicated user name",
			seedURL: "./sample",
			params: CreateUserParams{
				UserName:  "bmackerness1", // duplicated with user_id=2 in sample
				Password:  faker.Password(),
				Email:     faker.Email(),
				FirstName: faker.FirstName(),
				LastName:  faker.LastName(),
				Status:    fakeStatus(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			t.Cleanup(func() {
				err := postgresC.TruncateAllTables(ctx)
				if err != nil {
					t.Fatalf("truncate all tables in db: %v", err)
				}
			})

			if tt.seedURL != "" {
				seeder.MustRun(ctx, t, postgresC, fmt.Sprintf("file://%s", filepath.Join(workingDir, tt.seedURL)))
			}

			s := &Store{db: db}

			got, err := s.CreateUser(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Store.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_FindUserByID(t *testing.T) {
	tests := []struct {
		name    string
		seedURL string
		id      int64
		want    User
		wantErr bool
	}{
		{
			name:    "return error not found",
			id:      1,
			want:    User{},
			wantErr: true,
		},
		{
			name:    "return a user",
			seedURL: "./sample",
			id:      1,
			want: User{
				UserID:      1,
				UserName:    "kgreenaway0",
				Password:    "$2a$04$fW6oRPmdyVJRZVg1ex6rwOvo7eLeesyA8OEdC8W/Gbktn5MvPAhQW",
				Email:       "kgreenaway0@csmonitor.com",
				FirstName:   "Ky",
				LastName:    "Greenaway",
				CreatedAt:   time.Date(2023, time.August, 12, 22, 58, 46, 0, time.Local),
				UpdatedAt:   time.Date(2023, time.December, 25, 7, 55, 47, 0, time.Local),
				LastLoginAt: sql.NullTime{Time: time.Date(2024, time.June, 29, 2, 26, 05, 0, time.Local), Valid: true},
				Status:      "LOCKED",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			t.Cleanup(func() {
				err := postgresC.TruncateAllTables(ctx)
				if err != nil {
					t.Fatalf("truncate all tables in db: %v", err)
				}
			})

			if tt.seedURL != "" {
				seeder.MustRun(ctx, t, postgresC, fmt.Sprintf("file://%s", filepath.Join(workingDir, tt.seedURL)))
			}

			s := &Store{db: db}

			got, err := s.FindUserByID(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestStore_UpdateUserStatus(t *testing.T) {
	type args struct {
		id     int64
		status string
	}
	tests := []struct {
		name    string
		seedURL string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "return error not found",
			args:    args{},
			wantErr: true,
		},
		{
			name:    "return a user",
			seedURL: "./sample",
			args:    args{id: 1, status: "ACTIVE"},
			want:    "ACTIVE",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			t.Cleanup(func() {
				err := postgresC.TruncateAllTables(ctx)
				if err != nil {
					t.Fatalf("truncate all tables in db: %v", err)
				}
			})

			if tt.seedURL != "" {
				seeder.MustRun(ctx, t, postgresC, fmt.Sprintf("file://%s", filepath.Join(workingDir, tt.seedURL)))
			}

			s := &Store{db: db}

			err := s.UpdateUserStatus(ctx, tt.args.id, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				user, err := findUserByID(ctx, db, tt.args.id)
				if err != nil {
					t.Fatal(err.Error())
				}

				if tt.want != tt.args.status {
					t.Errorf("result mismatch (-want +got):\n-\t%q\n+\t%q", tt.want, user.Status)
				}
			}
		})
	}
}

func TestStore_FindUserByNameAndPassword(t *testing.T) {
	tests := []struct {
		name    string
		seedURL string
		uname   string
		upass   string
		wantErr bool
	}{
		{
			name:    "return the user",
			uname:   "uname",
			upass:   "upass",
			wantErr: false,
		},
		{
			name:    "return error password is incorrect",
			uname:   "uname",
			upass:   "upass_incorrect",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			t.Cleanup(func() {
				err := postgresC.TruncateAllTables(ctx)
				if err != nil {
					t.Fatalf("truncate all tables in db: %v", err)
				}
			})

			s := &Store{db: db}

			_, err := s.CreateUser(ctx, CreateUserParams{
				UserName:  "uname",
				Password:  "upass",
				Email:     faker.Email(),
				FirstName: faker.FirstName(),
				LastName:  faker.LastName(),
				Status:    "NEW",
			})
			if err != nil {
				t.Error("failed to create user")
				return
			}

			_, err = s.FindUserByNameAndPassword(ctx, tt.uname, tt.upass)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.FindUserByNameAndPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
