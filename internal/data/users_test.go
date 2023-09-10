package data

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"

	"database/sql"

	"github.com/testcontainers/testcontainers-go/wait"
)

func TestUserModel(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:12",
		ExposedPorts: []string{"5432/tcp"},
		AutoRemove:   true,
		Env: map[string]string{
			"POSTGRES_USER":     "greenlight",
			"POSTGRES_PASSWORD": "pa55word",
			"POSTGRES_DB":       "greenlight",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer postgres.Terminate(ctx)

	db, err := sql.Open("postgres", "postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	duration, err := time.ParseDuration("15m")
	if err != nil {
		t.Fatal(err)
	}
	db.SetConnMaxIdleTime(duration)

	if err = db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}

	var dbUserModel = DBUserModel{DB: db}
	_, err = dbUserModel.DB.ExecContext(ctx, "DELETE FROM users")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Insert", func(t *testing.T) {
		password := "pa55word"
		user := User{
			Name:  "alice smith",
			Email: "alice@foo.com",
			Password: Password{
				Plaintext: &password,
			},
			Activated: true,
		}
		err := user.Password.Set(*user.Password.Plaintext)
		if err != nil {
			t.Fatal(err)
		}

		err = dbUserModel.Insert(&user)
		if err != nil {
			t.Fatal(err)
		}
		_, err = dbUserModel.DB.ExecContext(ctx, "DELETE FROM users")
		if err != nil {
			t.Fatal(err)
		}

	})

	t.Run("Insert Duplicate", func(t *testing.T) {

		password := "pa55word"
		user := User{
			Name:  "alice smith",
			Email: "alice@foo.com",
			Password: Password{
				Plaintext: &password,
			},
			Activated: true,
		}
		err := user.Password.Set(*user.Password.Plaintext)
		if err != nil {
			t.Fatal(err)
		}

		err = dbUserModel.Insert(&user)
		if err != nil {
			t.Fatal(err)
		}

		// second time should return an error
		err = dbUserModel.Insert(&user)
		if err == nil {
			t.Fatal(err)
		}

		_, err = dbUserModel.DB.ExecContext(ctx, "DELETE FROM users")
		if err != nil {
			t.Fatal(err)
		}

	})

	t.Run("Insert & Get By email", func(t *testing.T) {
		password := "pa55word"
		user := User{
			Name:  "alice smith",
			Email: "alice@foo.com",
			Password: Password{
				Plaintext: &password,
			},
			Activated: true,
		}
		err := user.Password.Set(*user.Password.Plaintext)
		if err != nil {
			t.Fatal(err)
		}

		err = dbUserModel.Insert(&user)
		if err != nil {
			t.Fatal(err)
		}

		gotUser, err := dbUserModel.GetByEmail(user.Email)
		if err != nil {
			t.Fatal(err)
		}

		// make sure they match
		if gotUser.Email != user.Email {
			t.Fatal(err)
		}

		_, err = dbUserModel.DB.ExecContext(ctx, "DELETE FROM users")
		if err != nil {
			t.Fatal(err)
		}

	})

	t.Run("Insert & Update", func(t *testing.T) {

		password := "pa55word"
		user := User{
			Name:  "alice smith",
			Email: "alice@foo.com",
			Password: Password{
				Plaintext: &password,
			},
			Activated: true,
		}
		err := user.Password.Set(*user.Password.Plaintext)
		if err != nil {
			t.Fatal(err)
		}

		err = dbUserModel.Insert(&user)
		if err != nil {
			t.Fatal(err)
		}

		user.Activated = false
		err = dbUserModel.Update(&user)
		if err != nil {
			t.Fatal(err)
		}

		_, err = dbUserModel.DB.ExecContext(ctx, "DELETE FROM users")
		if err != nil {
			t.Fatal(err)
		}

	})

}
