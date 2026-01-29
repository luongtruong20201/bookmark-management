// Package fixture provides helpers to create and manage database-backed test fixtures.
// It defines a small interface and utilities to migrate schemas and seed data
// against an in-memory test database.
package fixture

import (
	"testing"

	sqldb "github.com/luongtruong20201/bookmark-management/pkg/sql"
	"gorm.io/gorm"
)

type Fixture interface {
	// SetupDB injects the database connection used by the fixture.
	SetupDB(*gorm.DB)
	// Migrate applies database schema migrations required for the fixture.
	Migrate() error
	// GenerateData seeds initial test data into the database.
	GenerateData() error
	// DB returns the underlying database connection after setup.
	DB() *gorm.DB
}

type base struct {
	db *gorm.DB
}

// SetupDB assigns the given GORM database instance to the base fixture.
func (b *base) SetupDB(db *gorm.DB) {
	b.db = db
}

// DB returns the GORM database instance associated with the base fixture.
func (b *base) DB() *gorm.DB {
	return b.db
}

// NewFixture initializes a new test fixture backed by an in-memory database.
// It creates a mock DB, runs the fixture's migrations and data generation,
// and returns the ready-to-use *gorm.DB for tests.
func NewFixture(t *testing.T, fix Fixture) *gorm.DB {
	fix.SetupDB(sqldb.InitMockDB(t))
	if err := fix.Migrate(); err != nil {
		t.Fatal("failed to migrate db for testing", err)
	}

	if err := fix.GenerateData(); err != nil {
		t.Fatal("failed to generate data for testing", err)
	}

	return fix.DB()
}
