# SKVMS – Solar/IoT Device Monitoring & Management System

SKVMS is a Go REST API server for managing solar panels, IoT sensors, and microcontroller devices (e.g. ESP32). It handles device registration, sensor reading ingestion, firmware code generation, audit logging, data export, and more.

## Documentation

| Document | Description |
|----------|-------------|
| [docs/api.md](docs/api.md) | Complete API reference – every endpoint with request/response examples |
| [docs/architecture.md](docs/architecture.md) | Architecture guide, SOLID principles, Strategy pattern, and improvement roadmap |

---

## Quick Start

### Prerequisites

- Go 1.21+
- MySQL 8.0+
- (Optional) `arduino-cli` or PlatformIO for firmware code generation

### Environment variables

Create a `.env` file (or set environment variables):

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=skvms
DB_PASSWORD=secret
DB_NAME=skvms
JWT_SECRET=change-me-in-production
SERVER_PORT=8080
LOG_DIR=logs
LOG_LEVEL=info
```

### Run

```bash
go run main.go
```

The server starts on `http://localhost:8080`.

---

## API Overview

All routes are prefixed with `/api`. Two authentication token types are used:

- **User JWT** – obtained from `POST /api/login` or `POST /api/register`
- **Device JWT** – obtained from `POST /api/device-auth/token` (used by embedded devices)

### Route groups

| Prefix | Description |
|--------|-------------|
| `/api/login`, `/api/register` | Authentication |
| `/api/users` | User management |
| `/api/devices` | Device CRUD, connected devices, readings |
| `/api/devices/sensors` | Sensor-type devices |
| `/api/devices/states` | Device state definitions and history |
| `/api/device-types`, `/api/devices/types` | Device type catalogue |
| `/api/device-auth/token` | Generate device JWT tokens |
| `/api/readings` | Ingest sensor readings (device JWT) |
| `/api/locations` | Location management and 7-day readings |
| `/api/audit` | Audit log viewer |
| `/api/admin/stats` | Platform statistics |
| `/api/versions`, `/api/features` | Firmware version and feature flag management |
| `/api/codegen` | ESP32 firmware generation and OTA upload |
| `/api/export` | Export data as CSV, XLSX, XML, or PDF |

See [docs/api.md](docs/api.md) for the full reference.

---

## Architecture

SKVMS follows a clean layered architecture:

```
Handler → Service → Repository → MySQL
```

- **Handlers** parse HTTP input and write responses.
- **Services** contain all business logic; defined as interfaces for easy testing and extension.
- **Repositories** handle database access via GORM; defined as interfaces.
- **DTOs** separate the API contract from the database schema.

The codebase uses the **Strategy pattern** for export formats (CSV/XLSX/XML/PDF) and firmware build tools (`arduino-cli`/PlatformIO). See [docs/architecture.md](docs/architecture.md) for detailed guidance on SOLID principles and how to add new features.

---

## Example Usage

### Register and log in

```bash
# Register
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"secret","email":"alice@example.com"}'

# Login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"secret"}'
# → { "token": "<jwt>", "user": { ... } }
```

### Create a device and post a reading

```bash
TOKEN="<jwt from login>"

# Create a device
curl -X POST http://localhost:8080/api/devices \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Solar Controller 1","type_id":1}'

# Generate a device token
curl -X POST http://localhost:8080/api/device-auth/token \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"device_id":1}'
# → { "token": "<device-jwt>", ... }

DEVICE_TOKEN="<device-jwt>"

# Post a sensor reading from the device
curl -X POST http://localhost:8080/api/readings \
  -H "Authorization: Bearer $DEVICE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"voltage":12.45,"current":2.13}'
```

### Get 7-day readings for a location

```bash
curl http://localhost:8080/api/locations/1/readings/seven \
  -H "Authorization: Bearer $TOKEN"
```

Each location returns per-device hourly averages over the last 7 days (up to 168 data points per device: 7 × 24 hourly averages).

```json
{
  "readings": [
    {
      "voltage": 11.93,
      "current": 2.52,
      "created_at": "2026-02-14T23:30:00+05:30"
    }
  ]
}
```

### Export readings to CSV

```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/export/readings?format=csv&device_id=1&start_date=2026-01-01&end_date=2026-01-31" \
  -o readings.csv
```