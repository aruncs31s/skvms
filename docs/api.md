# SKVMS API Reference

All routes are prefixed with `/api`. Authentication uses **Bearer JWT tokens** in the `Authorization` header.

**Two token types exist:**
- **User JWT** – issued by `/api/login` or `/api/register`, used by human users.
- **Device JWT** – issued by `/api/device-auth/token`, used by embedded devices (e.g. ESP32) to post readings.

---

## Table of Contents

1. [Authentication](#1-authentication)
2. [Users](#2-users)
3. [Devices](#3-devices)
4. [Device Types](#4-device-types)
5. [Device States](#5-device-states)
6. [Device Authentication](#6-device-authentication)
7. [Readings](#7-readings)
8. [Sensors](#8-sensors)
9. [Locations](#9-locations)
10. [Audit Logs](#10-audit-logs)
11. [Admin](#11-admin)
12. [Versions & Features](#12-versions--features)
13. [Code Generation (Codegen)](#13-code-generation-codegen)
14. [Export](#14-export)

---

## 1. Authentication

### POST `/api/register`

Register a new user and receive a JWT token immediately.

**Auth:** Public

**Request body:**
```json
{
  "name": "Alice",
  "username": "alice",
  "email": "alice@example.com",
  "password": "secret123",
  "role": "user"
}
```

| Field      | Type   | Required | Description                    |
|------------|--------|----------|--------------------------------|
| `username` | string | ✅       | Unique username                |
| `password` | string | ✅       | Plain-text password (hashed server-side) |
| `email`    | string | ❌       | User email address             |
| `name`     | string | ❌       | Display name (defaults to `User_<username>`) |
| `role`     | string | ❌       | Role (`user`, `admin`). Defaults to `user` |

**Response `200 OK`:**
```json
{
  "token": "<jwt>",
  "user": {
    "id": 1,
    "name": "Alice",
    "username": "alice",
    "email": "alice@example.com"
  }
}
```

---

### POST `/api/login`

Authenticate with username and password.

**Auth:** Public

**Request body:**
```json
{
  "username": "alice",
  "password": "secret123"
}
```

**Response `200 OK`:**
```json
{
  "token": "<jwt>",
  "user": {
    "id": 1,
    "name": "Alice",
    "username": "alice",
    "email": "alice@example.com",
    "role": "user"
  }
}
```

---

## 2. Users

### GET `/api/users`

List all users.

**Auth:** User JWT required

**Response `200 OK`:**
```json
{
  "users": [
    {
      "id": 1,
      "name": "Alice",
      "username": "alice",
      "email": "alice@example.com",
      "role": "user"
    }
  ]
}
```

---

### GET `/api/users/:id`

Get a specific user by ID.

**Auth:** User JWT required

**Path param:** `id` – user ID (integer)

**Response `200 OK`:**
```json
{
  "user": {
    "id": 1,
    "name": "Alice",
    "username": "alice",
    "email": "alice@example.com",
    "role": "user"
  }
}
```

---

### GET `/api/profile`

Get the authenticated user's profile including their devices and recent audit activity.

**Auth:** User JWT required

**Response `200 OK`:**
```json
{
  "profile": {
    "id": 1,
    "name": "Alice",
    "username": "alice",
    "email": "alice@example.com",
    "role": "user",
    "devices": [...],
    "activity": [...]
  }
}
```

---

### POST `/api/users`

Create a new user (alternative to `/api/register`; does not return a token).

**Auth:** Public

**Request body:** same as `/api/register`

**Response `201 Created`:**
```json
{ "message": "user created successfully" }
```

---

### PUT `/api/users/:id`

Update an existing user.

**Auth:** User JWT required

**Path param:** `id` – user ID

**Request body:**
```json
{
  "name": "Alice Updated",
  "username": "alice_new",
  "email": "newemail@example.com",
  "role": "admin",
  "password": "newpassword"
}
```

**Response `200 OK`:**
```json
{ "message": "user updated successfully" }
```

---

### DELETE `/api/users/:id`

Delete a user.

**Auth:** User JWT required

**Path param:** `id` – user ID

**Response `200 OK`:**
```json
{ "message": "user deleted successfully" }
```

---

## 3. Devices

### GET `/api/devices`

List all devices.

**Auth:** Public

**Response `200 OK`:**
```json
{
  "devices": [
    {
      "id": 1,
      "name": "Solar Controller 1",
      "type": "microcontroller",
      ...
    }
  ]
}
```

---

### GET `/api/devices/my`

List only devices belonging to the authenticated user.

**Auth:** User JWT required

**Response `200 OK`:**
```json
{ "devices": [...] }
```

---

### GET `/api/devices/:id`

Get a specific device by ID.

**Auth:** Public

**Path param:** `id` – device ID

**Response `200 OK`:**
```json
{ "device": { "id": 1, "name": "Solar Controller 1", ... } }
```

**Response `404 Not Found`:**
```json
{ "error": "device not found" }
```

---

### POST `/api/devices`

Create a new device.

**Auth:** User JWT required

**Request body:**
```json
{
  "name": "Solar Controller 1",
  "description": "Main solar MPPT controller",
  "type_id": 1,
  "location_id": 2
}
```

**Response `201 Created`:**
```json
{ "device": { "id": 5, "name": "Solar Controller 1", ... } }
```

---

### PUT `/api/devices/:id`

Partially update a device (audit logged as `device_update`).

**Auth:** User JWT required

**Path param:** `id` – device ID

**Request body:** fields to update (same schema as create)

**Response `200 OK`:**
```json
{ "message": "device updated successfully" }
```

---

### PUT `/api/devices/:id/full`

Fully replace a device record (audit logged as `device_full_update`).

**Auth:** User JWT required

**Path param:** `id` – device ID

**Request body:** complete device schema

**Response `200 OK`:**
```json
{ "message": "device fully updated successfully" }
```

---

### DELETE `/api/devices/:id`

Delete a device (audit logged).

**Auth:** User JWT required

**Path param:** `id` – device ID

**Response `200 OK`:**
```json
{ "message": "device deleted successfully" }
```

---

### POST `/api/devices/:id/control`

Send a control command to a device.

**Auth:** User JWT required

**Path param:** `id` – device ID

**Request body:**
```json
{ "action": 1 }
```

**Response `200 OK`:**
```json
{ "message": { "state": "on" } }
```

---

### GET `/api/devices/search`

Search devices by name or description.

**Auth:** Public

**Query params:**

| Param | Description        |
|-------|--------------------|
| `q`   | Search query string |

**Response `200 OK`:** array of matching devices

---

### GET `/api/devices/search/microcontrollers`

Search microcontroller-type devices.

**Auth:** Public

**Query params:** `q` – search string

---

### GET `/api/devices/microcontrollers`

List all microcontroller devices.

**Auth:** Public

**Response `200 OK`:**
```json
{ "devices": [...] }
```

---

### GET `/api/ba`

Get microcontroller statistics.

**Auth:** Public

**Response `200 OK`:**
```json
{ "stats": { ... } }
```

---

### GET `/api/devices/:id/connected`

Get all devices connected to (child of) the given device.

**Auth:** Public

**Path param:** `id` – parent device ID

**Response `200 OK`:**
```json
{ "connected_devices": [...] }
```

---

### POST `/api/devices/:id/connected`

Link an existing device as a child of the given device.

**Auth:** User JWT required

**Request body:**
```json
{ "child_id": 3 }
```

**Response `200 OK`:**
```json
{ "message": "connected device added successfully" }
```

---

### POST `/api/devices/:id/connected/new`

Create a brand-new device and immediately attach it as a child.

**Auth:** User JWT required

**Request body:**
```json
{
  "name": "DHT22 Sensor",
  "type_id": 2
}
```

**Response `200 OK`:**
```json
{ "connected_device": { "id": 10, ... } }
```

---

### DELETE `/api/devices/:id/connected/:cid`

Remove the parent-child link between `:id` and `:cid`.

**Auth:** User JWT required

**Path params:**
- `id` – parent device ID
- `cid` – child device ID

**Response `200 OK`:**
```json
{ "message": "connected device removed successfully" }
```

---

### GET `/api/devices/:id/connected/:cid/readings`

Get readings of a child device within a date range.

**Auth:** Public

**Path params:** `id` (parent), `cid` (child)

**Query params:**

| Param        | Format       | Default      |
|--------------|--------------|--------------|
| `start_date` | `2006-01-02` | Today 00:00  |
| `end_date`   | `2006-01-02` | Today 23:59  |

**Response `200 OK`:**
```json
{
  "readings": [...],
  "latest": [...]
}
```

---

### GET `/api/devices/:id/readings`

Get today's readings for a device (most recent first).

**Auth:** Public

**Query params:**

| Param   | Default | Description            |
|---------|---------|------------------------|
| `limit` | `50`    | Maximum records to return |

**Response `200 OK`:**
```json
{
  "latest": { "voltage": 12.5, "current": 2.1, ... },
  "readings": [...]
}
```

---

### GET `/api/devices/:id/readings/range`

Get readings for a device within a date range, with statistics.

**Auth:** Public

**Query params:**

| Param        | Format       | Required |
|--------------|--------------|----------|
| `start_date` | `2006-01-02` | ✅       |
| `end_date`   | `2006-01-02` | ✅       |

**Response `200 OK`:**
```json
{
  "readings": [...],
  "stats": {
    "avg_voltage": 12.3,
    "avg_current": 2.0,
    ...
  }
}
```

---

### GET `/api/devices/:id/readings/progressive`

Get hourly averaged readings for a device (used for progressive/trend charts).

**Auth:** Public

**Response `200 OK`:**
```json
{ "readings": [{ "voltage": 12.1, "current": 2.0, "hour": "2026-02-14T10:00:00Z" }] }
```

---

### GET `/api/devices/:id/readings/interval`

Get readings sampled at a configurable time interval.

**Auth:** Public

**Query params:**

| Param        | Format              | Default  |
|--------------|---------------------|----------|
| `start_date` | `2006-01-02`        | ✅       |
| `end_date`   | `2006-01-02`        | ✅       |
| `interval`   | Go duration (`1h`)  | `1h`     |
| `count`      | integer             | `24`     |

**Response `200 OK`:**
```json
{ "readings": [...] }
```

---

### GET `/api/devices/:id/versions`

Get all firmware versions associated with a device.

**Auth:** Public

**Response `200 OK`:**
```json
{ "versions": [...] }
```

---

### POST `/api/devices/:id/versions`

Create a new firmware version for a device.

**Auth:** Public (no JWT check in current router setup)

**Request body:**
```json
{
  "version": "1.2.0",
  "previous_version": 3,
  "features": [1, 2, 4]
}
```

**Response `201 Created`:** the created version object

---

### GET `/api/device/:id/features`

Get all features associated with a device (via its versions).

**Auth:** User JWT required

**Response `200 OK`:**
```json
{ "features": null }
```

> ⚠️ This endpoint is currently a stub and always returns `null`. See [architecture.md](architecture.md) for improvement guidance.

---

## 4. Device Types

### GET `/api/device-types`

List all device types.

**Auth:** Public

**Response `200 OK`:**
```json
{ "device_types": [...] }
```

---

### GET `/api/devices/types`

Alias for `/api/device-types`.

**Auth:** Public

---

### POST `/api/devices/types`

Create a new device type.

**Auth:** User JWT required

**Request body:**
```json
{
  "name": "Solar Inverter",
  "category": "hardware"
}
```

---

### GET `/api/devices/types/hardware`

List hardware device types only.

**Auth:** User JWT required

---

### GET `/api/devices/types/sensors`

List sensor device types only.

**Auth:** Public

---

### GET `/api/devices/:id/type`

Get the type of a specific device.

**Auth:** Public

**Path param:** `id` – device ID

---

## 5. Device States

### GET `/api/devices/states`

List all defined device states.

**Auth:** Public

**Response `200 OK`:**
```json
{ "device_states": [...] }
```

---

### GET `/api/devices/states/:id`

Get a specific device state by ID.

**Auth:** Public

**Response `200 OK`:**
```json
{ "device_state": { "id": 1, "name": "online", ... } }
```

---

### POST `/api/devices/states`

Create a new device state definition.

**Auth:** User JWT required

**Request body:**
```json
{
  "name": "maintenance",
  "device_id": 5
}
```

**Response `201 Created`:**
```json
{ "message": "device state created successfully" }
```

---

### PUT `/api/devices/states/:id`

Update an existing device state.

**Auth:** User JWT required

**Response `200 OK`:**
```json
{ "message": "device state updated successfully" }
```

---

### GET `/api/devices/:id/states/history`

Get the state-change history for a device.

**Auth:** User JWT required

**Query params:**

| Param       | Format       | Description          |
|-------------|--------------|----------------------|
| `from_date` | `2006-01-02` | Filter start date    |
| `to_date`   | `2006-01-02` | Filter end date      |
| `states`    | integer[]    | Filter by state IDs  |

**Response `200 OK`:**
```json
{
  "device_id": 5,
  "history": [
    { "state_id": 1, "state_name": "online", "changed_at": "2026-02-14T10:00:00Z" }
  ]
}
```

---

## 6. Device Authentication

### POST `/api/device-auth/token`

Generate a device-specific JWT token that an embedded device (e.g. ESP32) uses to authenticate when posting sensor readings.

**Auth:** User JWT required

**Request body:**
```json
{ "device_id": 5 }
```

**Response `200 OK`:**
```json
{
  "token": "<device-jwt>",
  "user_id": 1,
  "device_id": 5
}
```

> Store this token securely on the embedded device. Pass it as `Authorization: Bearer <token>` when calling `POST /api/readings`.

---

## 7. Readings

### POST `/api/readings`

Submit a new sensor reading from an embedded device.

**Auth:** Device JWT required (from `/api/device-auth/token`)

**Request body:**
```json
{
  "voltage": 12.45,
  "current": 2.13,
  "power": 26.5,
  "temperature": 35.0
}
```

**Response `200 OK`:**
```json
{ "message": "reading saved" }
```

---

## 8. Sensors

### GET `/api/devices/sensors`

List all sensor-type devices.

**Auth:** Public

**Response `200 OK`:** array of sensor devices

---

### POST `/api/devices/sensors`

Create a new sensor device.

**Auth:** User JWT required

**Request body:** same as `POST /api/devices`

**Response `201 Created`:** the created sensor device object

---

### GET `/api/devices/sensors/:id`

Get a specific sensor device.

**Auth:** Public

**Path param:** `id` – sensor device ID

---

### GET `/api/devices/sensors/:id/connected`

Get devices connected to the given sensor.

**Auth:** Public

---

### GET `/api/devices/sensors/search`

Search sensor devices.

**Auth:** Public

**Query params:** `q` – search string

---

## 9. Locations

### GET `/api/locations`

List all locations.

**Auth:** Public

**Response `200 OK`:**
```json
{
  "locations": [
    { "id": 1, "code": "LOC-001", "name": "Rooftop A", ... }
  ]
}
```

---

### GET `/api/locations/:id`

Get a specific location by ID.

**Auth:** Public

**Response `200 OK`:**
```json
{ "location": { "id": 1, "code": "LOC-001", "name": "Rooftop A" } }
```

**Response `404 Not Found`:**
```json
{ "error": "location not found" }
```

---

### GET `/api/locations/search`

Search locations by name or code.

**Auth:** Public

**Query params:** `q` – search string

**Response `200 OK`:**
```json
{ "locations": [...] }
```

---

### POST `/api/locations`

Create a new location (audit logged as `location_create`).

**Auth:** User JWT required

**Request body:**
```json
{
  "code": "LOC-002",
  "name": "Rooftop B",
  "description": "South-facing roof"
}
```

**Response `201 Created`:**
```json
{ "message": "location created successfully" }
```

---

### PUT `/api/locations/:id`

Update a location (audit logged as `location_update`).

**Auth:** User JWT required

**Response `200 OK`:**
```json
{ "message": "location updated successfully" }
```

---

### DELETE `/api/locations/:id`

Delete a location (audit logged as `location_delete`).

**Auth:** User JWT required

**Response `200 OK`:**
```json
{ "message": "location deleted successfully" }
```

---

### GET `/api/locations/:id/devices`

List all devices assigned to a location.

**Auth:** Public

**Response `200 OK`:**
```json
{ "devices": [...] }
```

---

### GET `/api/locations/:id/readings/seven`

Get per-device hourly readings for the last 7 days for all devices in a location.

Each day contains up to 168 readings (7 days × 24 hours × 60 min ÷ 60 = 168 hourly averages).

**Auth:** Requires `Authorization` header (user JWT)

**Response `200 OK`:**
```json
{
  "readings": [
    {
      "device_id": 5,
      "data": [
        { "voltage": 12.3, "current": 2.5, "created_at": "2026-02-14T23:30:00+05:30" }
      ]
    }
  ]
}
```

---

## 10. Audit Logs

### GET `/api/audit`

List audit log entries.

**Auth:** User JWT required

**Query params:**

| Param    | Description                                  |
|----------|----------------------------------------------|
| `action` | Filter by action type (e.g. `device_create`) |
| `limit`  | Maximum number of records (default: `100`)   |

**Response `200 OK`:**
```json
{
  "logs": [
    {
      "id": 1,
      "username": "alice",
      "action": "device_create",
      "details": "Created device: Solar Controller 1",
      "ip_address": "192.168.1.10",
      "device_id": 5,
      "created_at": "2026-02-14 10:00:00"
    }
  ]
}
```

**Common action values:**
- `login`
- `user_create`, `user_update`, `user_delete`
- `device_create`, `device_update`, `device_full_update`, `device_delete`
- `device_control`
- `device_token_generated`
- `device_state_create`, `device_state_update`, `device_state_delete`
- `location_create`, `location_update`, `location_delete`

---

## 11. Admin

### GET `/api/admin/stats`

Get platform-wide statistics (user count, device count, reading count, etc.).

**Auth:** User JWT required

**Response `200 OK`:**
```json
{
  "user_count": 10,
  "device_count": 25,
  "reading_count": 150000,
  "audit_count": 500
}
```

---

## 12. Versions & Features

Versions track firmware/software versions that can be linked to devices.
Features are individual capability flags associated with a version.

### POST `/api/versions`

Create a new version record.

**Auth:** User JWT required

**Request body:**
```json
{ "version": "1.0.0" }
```

**Response `201 Created`:** the created version object

---

### GET `/api/versions`

List all versions.

**Auth:** Public

**Response `200 OK`:**
```json
{ "versions": [...] }
```

---

### GET `/api/versions/:id`

Get a specific version.

**Auth:** User JWT required

---

### PUT `/api/versions/:id`

Update a version.

**Auth:** User JWT required

**Request body:**
```json
{ "version": "1.0.1" }
```

---

### DELETE `/api/versions/:id`

Delete a version.

**Auth:** User JWT required

**Response `200 OK`:**
```json
{ "message": "version deleted" }
```

---

### POST `/api/features`

Create a feature flag linked to a version.

**Auth:** User JWT required

**Request body:**
```json
{
  "version_id": 1,
  "name": "ota_updates",
  "enabled": true
}
```

**Response `201 Created`:** the created feature object

---

### GET `/api/features/version/:verid`

Get all features for a specific version.

**Auth:** User JWT required

**Path param:** `verid` – version ID

**Response `200 OK`:**
```json
{ "features": [...] }
```

---

### PUT `/api/features/:id`

Update a feature flag.

**Auth:** User JWT required

**Request body:**
```json
{
  "feature_name": "ota_updates",
  "enabled": false
}
```

---

### DELETE `/api/features/:id`

Delete a feature flag.

**Auth:** User JWT required

**Response `200 OK`:**
```json
{ "message": "feature deleted" }
```

---

## 13. Code Generation (Codegen)

These endpoints generate and manage ESP32 Arduino/PlatformIO firmware source code and compiled binaries.

### GET `/api/codegen/tools`

List available build tools (e.g. `arduino-cli`, `platformio`).

**Auth:** Public

**Response `200 OK`:**
```json
{ "available_tools": ["arduino-cli", "platformio"] }
```

---

### POST `/api/codegen/generate`

Generate firmware source and compile it. Returns the build ID.

**Auth:** User JWT required

**Request body:**
```json
{
  "ip": "192.168.1.50",
  "host_ip": "192.168.1.1",
  "host_ssid": "MyWifi",
  "host_password": "wifi_pass",
  "device_id": 5,
  "token": "<device-jwt>",
  "port": 8080,
  "build_tool": "arduino-cli"
}
```

**Response `200 OK`:**
```json
{
  "message": "Firmware built successfully",
  "build_tool": "arduino-cli",
  "binary_size": 204800,
  "build_id": "abc123"
}
```

---

### POST `/api/codegen/build`

Same as `/generate` but also returns a `download_url` for the binary.

**Auth:** User JWT required

**Response `200 OK`:**
```json
{
  "message": "Firmware built successfully",
  "build_id": "abc123",
  "download_url": "http://localhost:8080/api/codegen/download/abc123"
}
```

---

### POST `/api/codegen/build-and-download`

Build firmware and immediately stream the binary file to the client.

**Auth:** User JWT required

**Response:** binary file stream (`application/octet-stream`)
- `X-Build-ID` and `X-Build-Tool` headers are set

---

### GET `/api/codegen/download/:build_id`

Download a previously compiled firmware binary.

**Auth:** User JWT required

**Path param:** `build_id` – ID returned from a previous build request

**Response:** binary file stream (`application/octet-stream`)

---

### POST `/api/codegen/upload`

Build firmware and push it to an ESP32 via OTA (Over-the-Air) update.

**Auth:** User JWT required

**Request body:** same as `/generate` plus:
```json
{ "device_ip": "192.168.1.50" }
```

**Response `200 OK`:**
```json
{
  "message": "Firmware uploaded successfully via OTA",
  "device_ip": "192.168.1.50"
}
```

---

### DELETE `/api/codegen/builds/:build_id`

Remove a build's artifacts from the server.

**Auth:** User JWT required

**Path param:** `build_id`

**Response `200 OK`:**
```json
{ "message": "build cleaned up", "build_id": "abc123" }
```

---

## 14. Export

Export readings or device lists in various file formats.

### GET `/api/export/formats`

List supported export formats.

**Auth:** Public

**Response `200 OK`:**
```json
{ "formats": ["csv", "xlsx", "xml", "pdf"] }
```

---

### GET `/api/export/readings`

Export sensor readings for a device to a file.

**Auth:** User JWT required

**Query params:**

| Param        | Required | Description                              |
|--------------|----------|------------------------------------------|
| `format`     | ✅       | `csv`, `xlsx`, `xml`, or `pdf`           |
| `device_id`  | ✅       | ID of the device                         |
| `start_date` | ❌       | Start date (`2006-01-02`). Defaults to today |
| `end_date`   | ❌       | End date (`2006-01-02`). Defaults to today |
| `template`   | ❌       | Custom PDF template path                 |

**Response:** file download (content type varies by format)

**Example:**
```bash
curl -H "Authorization: Bearer <token>" \
  "http://localhost:8080/api/export/readings?format=csv&device_id=5&start_date=2026-01-01&end_date=2026-01-31" \
  -o readings.csv
```

---

### GET `/api/export/devices`

Export the full device list to a file.

**Auth:** User JWT required

**Query params:**

| Param      | Required | Description                    |
|------------|----------|--------------------------------|
| `format`   | ✅       | `csv`, `xlsx`, `xml`, or `pdf` |
| `template` | ❌       | Custom PDF template path        |

**Response:** file download

---

## Error Responses

All endpoints return a consistent error structure:

```json
{
  "error": "human-readable error message",
  "details": "optional additional detail"
}
```

| HTTP Status | Meaning                                     |
|-------------|---------------------------------------------|
| `400`       | Bad request – invalid input or missing param |
| `401`       | Unauthorized – missing or invalid token     |
| `404`       | Not found                                   |
| `500`       | Internal server error                       |
