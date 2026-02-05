# Sensor API Documentation

This document provides JSON request/response examples for the sensor-related APIs in the SKVMS system.

## Base URL
All API endpoints are prefixed with `/api/devices/sensors`

## Authentication
Most endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

---

## 1. List All Sensors

**Endpoint:** `GET /api/devices/sensors`

**Description:** Retrieves a list of all sensor devices in the system.

**Authentication:** Not required

**Response:**
```json
[
  {
    "id": 1,
    "name": "Temperature Sensor 001",
    "type": "DS18B20 Temperature Sensor",
    "hardware_type": "Sensor",
    "status": "Online",
    "ip_address": "192.168.1.100",
    "mac_address": "AA:BB:CC:DD:EE:FF",
    "firmware_version": "v1.2.3",
    "address": "Building A, Floor 1",
    "city": "Kannur"
  },
  {
    "id": 2,
    "name": "Current Sensor 001",
    "type": "ACS712 Current Sensor",
    "hardware_type": "CurrentSensor",
    "status": "Offline",
    "ip_address": "192.168.1.101",
    "mac_address": "AA:BB:CC:DD:EE:GG",
    "firmware_version": "v1.1.0",
    "address": "Building B, Floor 2",
    "city": "Kannur"
  }
]
```

---

## 2. Create Sensor Device

**Endpoint:** `POST /api/devices/sensors`

**Description:** Creates a new sensor device.

**Authentication:** Required (JWT token)

**Request Body:**
```json
{
  "name": "New Temperature Sensor",
  "type": 3,
  "ip_address": "192.168.1.102",
  "mac_address": "AA:BB:CC:DD:EE:HH",
  "firmware_version_id": 1,
  "address": "Building C, Room 101",
  "city": "Kannur"
}
```

**Response:**
```json
{
  "id": 3,
  "name": "New Temperature Sensor",
  "type": "DS18B20 Temperature Sensor",
  "hardware_type": "Sensor",
  "status": "Offline",
  "ip_address": "192.168.1.102",
  "mac_address": "AA:BB:CC:DD:EE:HH",
  "firmware_version": "v1.2.3",
  "address": "Building C, Room 101",
  "city": "Kannur"
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "name is required"
}
```

---

## 3. Get Sensor Device

**Endpoint:** `GET /api/devices/sensors/{id}`

**Description:** Retrieves details of a specific sensor device by ID.

**Authentication:** Not required

**Parameters:**
- `id` (path): The sensor device ID

**Response:**
```json
{
  "id": 1,
  "name": "Temperature Sensor 001",
  "type": "DS18B20 Temperature Sensor",
  "hardware_type": "Sensor",
  "status": "Online",
  "ip_address": "192.168.1.100",
  "mac_address": "AA:BB:CC:DD:EE:FF",
  "firmware_version": "v1.2.3",
  "address": "Building A, Floor 1",
  "city": "Kannur"
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "sensor device not found"
}
```

---

## 4. Get Connected Devices

**Endpoint:** `GET /api/devices/sensors/{id}/connected`

**Description:** Retrieves devices connected to a specific sensor device.

**Authentication:** Not required

**Parameters:**
- `id` (path): The sensor device ID

**Response:**
```json
[
  {
    "id": 4,
    "name": "ESP32 Microcontroller",
    "type": "ESP32-WROOM-32",
    "hardware_type": "MicroController",
    "status": "Online",
    "ip_address": "192.168.1.103",
    "mac_address": "AA:BB:CC:DD:EE:II",
    "firmware_version": "v2.0.1",
    "address": "Building A, Floor 1",
    "city": "Kannur"
  }
]
```

---

## 5. Search Sensor Devices

**Endpoint:** `GET /api/devices/sensors/search?q={query}`

**Description:** Searches for sensor devices by name.

**Authentication:** Not required

**Parameters:**
- `q` (query): Search query string

**Response:**
```json
[
  {
    "id": 1,
    "name": "Temperature Sensor 001"
  },
  {
    "id": 2,
    "name": "Temperature Sensor 002"
  }
]
```

**Error Response (400 Bad Request):**
```json
{
  "error": "query parameter 'q' is required"
}
```

---

## Data Types

### DeviceView
```json
{
  "id": "number (uint)",
  "name": "string",
  "type": "string",
  "hardware_type": "string",
  "status": "string",
  "ip_address": "string",
  "mac_address": "string",
  "firmware_version": "string",
  "address": "string",
  "city": "string"
}
```

### CreateDeviceRequest
```json
{
  "name": "string (required)",
  "type": "number (uint, required)",
  "ip_address": "string (optional)",
  "mac_address": "string (optional)",
  "firmware_version_id": "number (uint, optional)",
  "address": "string (optional)",
  "city": "string (optional)"
}
```

### GenericDropdown
```json
{
  "id": "number (uint)",
  "name": "string"
}
```

## Hardware Types
Sensor devices can have the following hardware types:
- `Sensor`
- `CurrentSensor`
- `PowerMeter`
- `VoltageMeter`

## Status Values
- `Online`
- `Offline`
- `Maintenance`
- `Error`</content>
<parameter name="filePath">/home/aruncs/Projects/Smart-City/skvms/sensor-api-examples.md