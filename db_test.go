package main_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	sqlc "TPEspecial.com/ProgramacionWeb/TPEspecial/db/sqlc"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	connStr := "user=u password=p dbname=postgres host=localhost port=5432 sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatalf("failed to ping test database: %v", err)
	}

	t.Cleanup(func() { _ = db.Close() })
	return db
}

func uniqueID() int64 {
	return time.Now().UnixNano()
}

func TestUserCRUD(t *testing.T) {
	ctx := context.Background()
	queries := sqlc.New(openTestDB(t))
	name := fmt.Sprintf("User-%d", uniqueID())
	email := fmt.Sprintf("user-%d@example.com", uniqueID())

	created, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
		Name:     name,
		Email:    email,
		Age:      sql.NullTime{Time: time.Date(1995, 7, 20, 0, 0, 0, 0, time.UTC), Valid: true},
		Industry: sql.NullString{String: "Technology", Valid: true},
		Role:     "Developer",
	})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	got, err := queries.GetUser(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if got.Name != name || got.Email != email {
		t.Fatalf("GetUser returned unexpected data: %+v", got)
	}

	users, err := queries.ListUsers(ctx)
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}
	found := false
	for _, u := range users {
		if u.ID == created.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CreateUser row not found in ListUsers for user ID=%d", created.ID)
	}

	updatedName := fmt.Sprintf("Updated-%s", name)
	updatedEmail := fmt.Sprintf("updated-%d@example.com", uniqueID())
	err = queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:       created.ID,
		Name:     updatedName,
		Email:    updatedEmail,
		Age:      sql.NullTime{Time: time.Date(1996, 3, 10, 0, 0, 0, 0, time.UTC), Valid: true},
		Industry: sql.NullString{String: "Finance", Valid: true},
		Role:     "Manager",
	})
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}

	updated, err := queries.GetUser(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetUser after update failed: %v", err)
	}
	if updated.Name != updatedName || updated.Email != updatedEmail || updated.Role != "Manager" {
		t.Fatalf("UpdateUser did not persist expected values: %+v", updated)
	}

	if err := queries.DeleteUser(ctx, created.ID); err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}

	_, err = queries.GetUser(ctx, created.ID)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("GetUser after DeleteUser should return sql.ErrNoRows, got: %v", err)
	}
}

func TestEventCRUD(t *testing.T) {
	ctx := context.Background()
	queries := sqlc.New(openTestDB(t))

	owner, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
		Name:  fmt.Sprintf("EventOwner-%d", uniqueID()),
		Email: fmt.Sprintf("event-owner-%d@example.com", uniqueID()),
		Role:  "Organizer",
	})
	if err != nil {
		t.Fatalf("CreateUser for event test failed: %v", err)
	}
	defer func() {
		_ = queries.DeleteUser(ctx, owner.ID)
	}()

	created, err := queries.CreateEvent(ctx, sqlc.CreateEventParams{
		Title:       fmt.Sprintf("Event-%d", uniqueID()),
		Description: sql.NullString{String: "Welcome to the event", Valid: true},
		Location:    "Buenos Aires",
		AgeRange:    "18+",
		UserID:      owner.ID,
	})
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	got, err := queries.GetEvent(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetEvent failed: %v", err)
	}
	if got.Title != created.Title || got.Location != created.Location {
		t.Fatalf("GetEvent returned unexpected event: %+v", got)
	}

	events, err := queries.ListEvents(ctx)
	if err != nil {
		t.Fatalf("ListEvents failed: %v", err)
	}
	found := false
	for _, e := range events {
		if e.ID == created.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("created event not found in ListEvents, id=%d", created.ID)
	}

	err = queries.UpdateEvent(ctx, sqlc.UpdateEventParams{
		ID:          created.ID,
		Title:       "Updated Event",
		Description: sql.NullString{String: "Updated description", Valid: true},
		Location:    "Córdoba",
		AgeRange:    "21+",
		UserID:      owner.ID,
	})
	if err != nil {
		t.Fatalf("UpdateEvent failed: %v", err)
	}

	updated, err := queries.GetEvent(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetEvent after update failed: %v", err)
	}
	if updated.Title != "Updated Event" || updated.Location != "Córdoba" || updated.AgeRange != "21+" {
		t.Fatalf("UpdateEvent did not persist expected values: %+v", updated)
	}

	if err := queries.DeleteEvent(ctx, created.ID); err != nil {
		t.Fatalf("DeleteEvent failed: %v", err)
	}

	_, err = queries.GetEvent(ctx, created.ID)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("GetEvent after DeleteEvent should return sql.ErrNoRows, got: %v", err)
	}
}
