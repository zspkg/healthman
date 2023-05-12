# Healthman Service
[![codecov](https://codecov.io/github/zspkg/healthman/branch/main/graph/badge.svg?token=IAWZEHOV58)](https://codecov.io/github/zspkg/healthman)
[![Go Report Card](https://goreportcard.com/badge/github.com/zspkg/healthman)](https://goreportcard.com/report/github.com/zspkg/healthman)
[![Go Reference](https://pkg.go.dev/badge/github.com/zspkg/healthman#section-readme.svg)](https://pkg.go.dev/github.com/zspkg/healthman#section-readme)

Service for easier health managing in services. API response are made according to [Redhat standard](https://developers.redhat.com/sites/default/files/blog/2017/11/microprofile-health-spec.pdf).

## Install

```bash
go get github.com/zspkg/healthman
```

## How to use?

Let us take notifications service as an example. Suppose that it depends on `postgres` database since we need to store identifiers, uses `telegram` client to send messages and `redis` for storing some simple data.

To make **Healthman** work one firstly needs to create a structure (or use existing one) that implements `Healthable` interface, that is, contains `Name() string` and `CheckHealth() bool` methods. For example, for the `postgres` connection we already have implementation, but if you were to implement it, it could approximately look like:
```go
package health

import (
	"database/sql"
	"github.com/zspkg/healthman"
)

// Create a struct that will implement Healthable interface.
// Or, you might take an existing struct and 
// add Name() and CheckHealth() methods
type dbHealthier struct {
	DB *sql.DB
}

// Create a constructor
func NewDBHealthier(db *sql.DB) healthman.Healthable {
	return dbHealthier{DB: db}
}

// Simply return some fixed constant string
func (h dbHealthier) Name() string {
	return "db-postgres"
}

// Try pinging a database and if no error occurred, return true.
// Otherwise, return false.
func (h dbHealthier) CheckHealth() bool {
	return h.DB.Ping() == nil
}
```

Then, we create an instance of `HealthChecker` in which we put an array of `Healthable` objects. Then we `Run(ctx)` this checker that launches a goroutine that with specified amount of time pings all services and updates their health statuses (using `sync.Map` as a storage). So, for example, we might use the following lines to run checker:

```go
package somePackage

import (
	"context"
	"github.com/zspkg/healthman"
)

func RunNotificationService(ctx context.Context) {
	// Initializing a HealthyChecker. First, we specify a period
	// of checking. Then, we specify logger (must not be nil) and
	// finally listing all services (that is, Healthables) to be further checked.
	healthyChecker := healthman.NewHealthyChecker(
		&healthman.HealthmanSettings{Period: time.Second},
		yourLogger,
		healthman.NewDBHealthier(yourDBConnection),
		NewRedisHealthier(yourRedisConnection),
		NewTelegramHealthier(yourTgClient))
	// That will run a routine that will update check
	// statuses of all included services
	healthChecker.Run(ctx)
}
```

Then, you can use `healthyChecker.Info()` method to get information about healthy statuses of your services. 

This package also supports integration with API, so you can easily add `/health` endpoint on your service by writing:

```go
package someAPIPackage

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/zspkg/healthman"
	"gitlab.com/distributed_lab/ape"
)

func RunRouter(ctx context.Context) {
	// Initialize healthyChecker before and run it
	// ...
	// Creating a router
	router := chi.NewRouter()
	
	// Adding HealthChecker to a context
	router.Use(ape.CtxMiddleware(healthman.HealthCheckerSetter(healthyChecker)))
	
	// Adding /health endpoint with a handler already implemented in Healthman
	router.Get("/health", healthman.CheckHealth)
}
```

Also, package contains a structure for config that can retrieve configuration parameters from `.yaml` file. So you can use:
```go
configer := NewHealthmanConfiger(kv.NewViperFile(".config.yaml"))
cfgParams := configer.HealthmanConfig()
```
