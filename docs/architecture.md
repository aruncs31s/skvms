# SKVMS Architecture Guide

This document describes the current layered architecture of SKVMS, explains how it maps to SOLID principles, and provides concrete guidance on how to improve it using the **Strategy**, **Repository**, and other design patterns.

---

## Table of Contents

1. [Overview](#1-overview)
2. [Directory Structure](#2-directory-structure)
3. [Layer Responsibilities](#3-layer-responsibilities)
   - [Handlers (HTTP layer)](#31-handlers-http-layer)
   - [Services (Business logic layer)](#32-services-business-logic-layer)
   - [Repositories (Data access layer)](#33-repositories-data-access-layer)
   - [Models](#34-models)
   - [DTOs](#35-dtos)
4. [Request Lifecycle](#4-request-lifecycle)
5. [SOLID Principles in SKVMS](#5-solid-principles-in-skvms)
6. [Design Patterns Used & Opportunities](#6-design-patterns-used--opportunities)
7. [Middleware](#7-middleware)
8. [Adding a New Feature: Step-by-Step](#8-adding-a-new-feature-step-by-step)
9. [Improvement Roadmap](#9-improvement-roadmap)

---

## 1. Overview

SKVMS is a **Solar/IoT Device Monitoring & Management System** built in Go with the following stack:

| Component         | Technology              |
|-------------------|-------------------------|
| HTTP framework    | [Gin](https://gin-gonic.com/) |
| ORM               | [GORM](https://gorm.io/) with MySQL |
| Authentication    | JWT (user + device tokens) |
| Authorization     | Role-based (RBAC via Casbin config) |
| Audit logging     | Custom `AuditService`   |
| Structured logging | [Zap](https://github.com/uber-go/zap) |
| Firmware codegen  | `arduino-cli` / PlatformIO |
| Export            | CSV, XLSX, XML, PDF     |

---

## 2. Directory Structure

```
skvms/
├── main.go                    # Application entry point: wiring of dependencies
├── config/                    # Configuration loading (env vars, casbin)
├── internal/
│   ├── handler/
│   │   ├── http/              # HTTP handlers – one file per domain
│   │   └── middleware/        # JWT auth, Audit middleware
│   ├── service/               # Business logic interfaces + implementations
│   ├── repository/            # Database access – one file per model
│   ├── model/                 # GORM database models
│   ├── dto/                   # Request/response data transfer objects
│   ├── router/                # Route registration
│   ├── database/              # DB connection + seeding
│   ├── logger/                # Zap logger initialisation
│   ├── codegen/               # ESP32 firmware generation
│   └── export/                # CSV/XLSX/XML/PDF export
├── utils/                     # Shared utility functions
├── templates/                 # Export templates (PDF, XLSX)
└── docs/                      # This documentation
```

---

## 3. Layer Responsibilities

### 3.1 Handlers (HTTP layer)

**Location:** `internal/handler/http/`

Each handler file corresponds to one domain (e.g. `user_handler.go`, `device_handler.go`).

Responsibilities:
- Parse and validate HTTP request parameters and body.
- Call one or more service methods.
- Map service results to HTTP responses.
- Record audit log entries when state-changing actions succeed.
- **Never** contain business logic or direct database calls.

```go
// Good: handler delegates to service, then audits
func (h *UserHandler) DeleteUser(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }
    if err := h.userService.Delete(c.Request.Context(), uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
        return
    }
    userID, _ := c.Get("user_id")
    username, _ := c.Get("username")
    _ = h.auditService.Log(...)
    c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
```

---

### 3.2 Services (Business logic layer)

**Location:** `internal/service/`

Each service is defined as a **Go interface** with a private struct implementing it. This design enables:
- Easy mocking in tests.
- Multiple implementations (Strategy pattern).
- Inversion of Control via constructor injection.

| Service                    | Responsibilities |
|----------------------------|-----------------|
| `AuthService`              | Register users, login, issue JWTs |
| `UserService`              | User CRUD, profile aggregation |
| `DeviceService`            | Device CRUD, parent-child linking, control commands |
| `DeviceAuthService`        | Issue and validate device-scoped JWTs |
| `DeviceStateService`       | Manage device state definitions |
| `DeviceStateHistoryService`| Record and query device state transitions |
| `DeviceTypesService`       | Device type CRUD |
| `ReadingService`           | Ingest sensor readings, query by date/interval |
| `AuditService`             | Log user actions |
| `AdminService`             | Aggregate platform statistics |
| `LocationService`          | Location CRUD, associate devices, 7-day readings |
| `VersionService`           | Firmware version + feature flag management |
| `SolarService`             | Solar-specific aggregation |

---

### 3.3 Repositories (Data access layer)

**Location:** `internal/repository/`

Repositories are also interface-driven, wrapping GORM calls. No business logic belongs here – only database queries.

| Repository                        | Model(s) accessed |
|-----------------------------------|-------------------|
| `UserRepository`                  | `User` |
| `DeviceRepository`                | `Device`, `DeviceAssignment` |
| `ReadingRepository`               | `Reading` |
| `AuditRepository`                 | `AuditLog` |
| `DeviceTypeRepository`            | `DeviceType` |
| `DeviceStateRepository`           | `DeviceState` |
| `DeviceStateHistoryRepository`    | `DeviceStateHistory` |
| `LocationRepository`              | `Location` |
| `MicrocontrollersRepository`      | `Device` (microcontroller subset) |
| `SolarRepository`                 | `Reading` (solar aggregation) |
| `VersionRepository`               | `Version`, `Feature` |

---

### 3.4 Models

**Location:** `internal/model/`

GORM structs that map directly to database tables. Models should contain:
- Field definitions with GORM tags.
- Relationship declarations (`HasMany`, `BelongsTo`, etc.).
- No business logic.

---

### 3.5 DTOs

**Location:** `internal/dto/`

Data Transfer Objects separate the API contract from the database schema:
- `CreateXxxRequest` / `UpdateXxxRequest` – input validation via `binding` tags.
- `XxxView` – sanitised output (excludes sensitive fields like `password`).

---

## 4. Request Lifecycle

```
HTTP Request
     │
     ▼
┌──────────────┐
│  Gin Router  │  (router/router.go)
└──────┬───────┘
       │  optional middleware
       ▼
┌──────────────┐
│  Middleware  │  JWTAuth, DeviceJWTAuth, AuditMiddleware
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Handler    │  parse request, call service, write response
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Service    │  validate business rules, orchestrate operations
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Repository  │  execute database queries via GORM
└──────┬───────┘
       │
       ▼
   MySQL DB
```

---

## 5. SOLID Principles in SKVMS

### S – Single Responsibility Principle

Each layer has one clear job:
- `AuthHandler.Login` only parses the HTTP request, calls the service, and returns a response.
- `authService.Login` only authenticates the user and issues a token.
- `UserRepository.GetByID` only fetches a user row from the database.

**To improve:** avoid mixing audit log writes inside handlers when the same action is triggered from multiple handlers. Consider centralising audit logic into service methods or middleware.

---

### O – Open/Closed Principle

Services are defined as interfaces. To add a new implementation (e.g. a cached user service) you create a new struct that implements `UserService` without touching the existing code:

```go
type cachedUserService struct {
    inner UserService
    cache Cache
}

func (s *cachedUserService) GetByID(ctx context.Context, id uint) (*dto.UserView, error) {
    if v := s.cache.Get(id); v != nil {
        return v, nil
    }
    return s.inner.GetByID(ctx, id)
}
```

Wire it in `main.go`:
```go
userService := service.NewCachedUserService(
    service.NewUserService(userRepo, deviceService, auditService),
    redisCache,
)
```

---

### L – Liskov Substitution Principle

Any struct implementing `service.UserService` must be substitutable for another. Tests use mock implementations:

```go
type mockUserService struct{}

func (m *mockUserService) List(ctx context.Context) ([]dto.UserView, error) {
    return []dto.UserView{{ID: 1, Username: "test"}}, nil
}
```

---

### I – Interface Segregation Principle

Smaller interfaces enable narrower dependencies. Example already in the codebase:

```go
type UserReader interface {
    GetByID(ctx context.Context, id uint) (*dto.UserView, error)
}

type UserService interface {
    UserReader
    List(ctx context.Context) ([]dto.UserView, error)
    Create(ctx context.Context, req *dto.CreateUserRequest) error
    Update(ctx context.Context, id uint, req *dto.UpdateUserRequest) error
    Delete(ctx context.Context, id uint) error
    GetProfile(ctx context.Context, userID uint) (*dto.UserProfile, error)
}
```

Components that only need to read users (e.g. a notification service) can depend on `UserReader` rather than the full `UserService`.

**To improve:** apply the same pattern to `DeviceService`, which is currently a large interface. Split it into:
- `DeviceReader` – read-only operations
- `DeviceWriter` – create/update/delete
- `DeviceController` – control commands
- `ConnectedDeviceManager` – parent-child linking

---

### D – Dependency Inversion Principle

All dependencies are injected via constructors:

```go
func NewUserService(
    repo repository.UserRepository,   // abstraction, not concrete type
    deviceService DeviceService,       // abstraction
    auditService AuditService,         // abstraction
) UserService {
    return &userService{repo: repo, ...}
}
```

`main.go` is the **composition root** – the only place where concrete types are instantiated and wired together.

---

## 6. Design Patterns Used & Opportunities

### Repository Pattern ✅ (already implemented)

Repositories abstract data access behind interfaces, making it straightforward to swap MySQL for PostgreSQL or add an in-memory store for tests.

---

### Strategy Pattern – Export Formatters ✅ (already implemented)

The `export` package uses a strategy pattern for format selection. Each format (`csv`, `xlsx`, `xml`, `pdf`) is an independent implementation of a common `Exporter` interface, chosen at runtime based on the `format` query parameter.

**How to extend:** add a new format by creating a new struct implementing the exporter interface and registering it in the factory, with zero changes to the handler or service.

```go
// Existing interface (conceptual)
type Exporter interface {
    Export(ctx context.Context, data interface{}, w io.Writer) error
}

// Add a new format:
type JSONExporter struct{}

func (e *JSONExporter) Export(ctx context.Context, data interface{}, w io.Writer) error {
    return json.NewEncoder(w).Encode(data)
}
```

---

### Strategy Pattern – Firmware Build Tools ✅ (already implemented)

`internal/codegen` selects the build tool (`arduino-cli` vs `platformio`) at runtime via a strategy. The `GET /api/codegen/tools` endpoint lists whichever tools are available on the server.

---

### Strategy Pattern – Reading Aggregation (opportunity)

The reading endpoints expose several aggregation strategies (raw, progressive, interval, range). These are currently implemented as separate service methods. Extract them into a `ReadingAggregator` strategy:

```go
// internal/service/reading_aggregator.go

type ReadingAggregator interface {
    Aggregate(ctx context.Context, deviceID uint, params AggregationParams) ([]dto.ReadingView, error)
}

type intervalAggregator struct{ repo repository.ReadingRepository }
type progressiveAggregator struct{ repo repository.ReadingRepository }
type rangeAggregator struct{ repo repository.ReadingRepository }
```

The `ReadingService` selects the appropriate aggregator based on request parameters, keeping each aggregation algorithm in its own testable struct.

---

### Strategy Pattern – Notification Channels (opportunity)

When device state changes occur or thresholds are exceeded, alerts can be sent via different channels (email, SMS, webhook). Implement a `NotificationStrategy`:

```go
type NotificationStrategy interface {
    Notify(ctx context.Context, event DeviceEvent) error
}

type emailNotifier struct{ smtp SmtpClient }
type webhookNotifier struct{ url string }
type smsNotifier struct{ client SmsClient }

// In DeviceStateService:
func (s *deviceStateService) Update(ctx context.Context, id uint, req *dto.UpdateDeviceStateRequest) error {
    // ... update state ...
    for _, notifier := range s.notifiers {
        _ = notifier.Notify(ctx, stateChangeEvent)
    }
    return nil
}
```

---

### Factory Pattern (opportunity)

`main.go` acts as a manual factory (composition root). For larger deployments consider a **wire** or **fx**-based dependency injection container to reduce boilerplate.

---

### Observer / Event Bus (opportunity)

Many state-changing operations (device created, reading posted, state changed) currently trigger inline audit log writes scattered across handlers. Replace this with an **event bus**:

```go
type Event struct {
    Type    string
    Payload interface{}
}

type EventBus interface {
    Publish(ctx context.Context, event Event)
    Subscribe(eventType string, handler func(Event))
}
```

Subscribe audit logging, notifications, and analytics as listeners. Handlers emit events; they don't call subscribers directly.

---

## 7. Middleware

| Middleware         | Location                                | Description |
|--------------------|-----------------------------------------|-------------|
| `JWTAuth`          | `internal/handler/middleware/jwt.go`    | Validates user JWT; sets `user_id` and `username` in Gin context |
| `DeviceJWTAuth`    | `internal/handler/middleware/jwt.go`    | Validates device JWT; sets `device_id` and `user_id` in Gin context |
| `AuditMiddleware`  | `internal/handler/middleware/audit.go`  | Records audit log entries for specific routes (e.g. device updates, location changes) |
| CORS               | `router/router.go` via `gin-contrib/cors` | Allows requests from `localhost:5173` and `localhost:3000` |

---

## 8. Adding a New Feature: Step-by-Step

Example: **Alert threshold configuration** – allow users to set a voltage threshold and receive an alert when a reading exceeds it.

**Step 1 – Model**

```go
// internal/model/alert_threshold.go
type AlertThreshold struct {
    gorm.Model
    DeviceID   uint    `gorm:"not null;index"`
    Metric     string  `gorm:"not null"` // "voltage", "current", etc.
    Threshold  float64 `gorm:"not null"`
    CreatedBy  uint    `gorm:"not null"`
}
```

**Step 2 – Repository**

```go
// internal/repository/alert_threshold_repository.go
type AlertThresholdRepository interface {
    Create(ctx context.Context, t *model.AlertThreshold) error
    GetByDevice(ctx context.Context, deviceID uint) ([]model.AlertThreshold, error)
    Delete(ctx context.Context, id uint) error
}
```

**Step 3 – Service**

```go
// internal/service/alert_threshold_service.go
type AlertThresholdService interface {
    Set(ctx context.Context, req *dto.CreateAlertThresholdRequest) error
    ListByDevice(ctx context.Context, deviceID uint) ([]dto.AlertThresholdView, error)
    Delete(ctx context.Context, id uint) error
    CheckThresholds(ctx context.Context, deviceID uint, reading model.Reading) error
}
```

**Step 4 – DTOs**

```go
// internal/dto/alert_threshold.go
type CreateAlertThresholdRequest struct {
    DeviceID  uint    `json:"device_id" binding:"required"`
    Metric    string  `json:"metric" binding:"required"`
    Threshold float64 `json:"threshold" binding:"required"`
}

type AlertThresholdView struct {
    ID        uint    `json:"id"`
    DeviceID  uint    `json:"device_id"`
    Metric    string  `json:"metric"`
    Threshold float64 `json:"threshold"`
}
```

**Step 5 – Handler**

```go
// internal/handler/http/alert_threshold_handler.go
type AlertThresholdHandler struct {
    service service.AlertThresholdService
}

func (h *AlertThresholdHandler) CreateThreshold(c *gin.Context) { ... }
func (h *AlertThresholdHandler) ListByDevice(c *gin.Context)    { ... }
func (h *AlertThresholdHandler) DeleteThreshold(c *gin.Context) { ... }
```

**Step 6 – Router**

```go
// internal/router/router.go – inside setupAPIRoutes
func (r *Router) setupAlertRoutes(api *gin.RouterGroup) {
    alerts := api.Group("/alerts")
    alerts.POST("/thresholds", middleware.JWTAuth(r.jwtSecret), r.alertHandler.CreateThreshold)
    alerts.GET("/thresholds/device/:id", middleware.JWTAuth(r.jwtSecret), r.alertHandler.ListByDevice)
    alerts.DELETE("/thresholds/:id", middleware.JWTAuth(r.jwtSecret), r.alertHandler.DeleteThreshold)
}
```

**Step 7 – Wire in main.go**

```go
alertRepo := repository.NewAlertThresholdRepository(db)
alertService := service.NewAlertThresholdService(alertRepo, notificationStrategy)
alertHandler := httpHandler.NewAlertThresholdHandler(alertService)
```

---

## 9. Improvement Roadmap

### High Priority

| Area | Recommendation |
|------|----------------|
| **Error handling** | Define a central error type with codes (e.g. `ErrNotFound`, `ErrUnauthorised`) so handlers can map them to HTTP status codes without duplication |
| **Input validation** | Move domain validation (e.g. password length, email format) into service layer rather than relying solely on `binding` tags |
| **Pagination** | Add `page` / `limit` / `cursor` to all list endpoints (`ListUsers`, `ListDevices`, `ListAuditLogs`, etc.) |
| **Tests** | Add unit tests for services using mock repositories, and integration tests for handlers |

### Medium Priority

| Area | Recommendation |
|------|----------------|
| **Caching** | Wrap `DeviceRepository` with a Redis-backed cache for frequently queried devices |
| **Event bus** | Replace inline audit calls with an event-driven system (see §6) |
| **Rate limiting** | Re-enable the commented-out rate-limiting middleware in `main.go` |
| **Graceful shutdown** | Handle `SIGINT`/`SIGTERM` and drain in-flight requests before exit |

### Lower Priority

| Area | Recommendation |
|------|----------------|
| **OpenAPI / Swagger** | Generate an `openapi.yaml` from code annotations using `swaggo/swag` |
| **RBAC enforcement** | Activate Casbin (config files already present in `config/`) for fine-grained permission checks per endpoint |
| **Metrics** | Add Prometheus metrics (request latency, error rates, reading ingest rate) |
| **WebSocket / SSE** | Push real-time reading updates to the dashboard without polling |
| **Multi-tenancy** | Scope all queries to an organisation/tenant ID for SaaS deployments |
