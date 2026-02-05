# Device Types API Documentation

This document describes the API endpoints for managing device types in the SKVMS system.

## Endpoints

### GET /api/devices/types

Retrieves a list of all available device types.

**Authentication:** Not required

**Query Parameters:**
- `limit` (optional): Number of device types to return (integer)
- `offset` (optional): Number of device types to skip (integer)

**Response:**
```json
{
  "device_types": [
    {
      "id": 1,
      "name": "ESP8266 (NODEMCU)"
    },
    {
      "id": 2,
      "name": "ESP32"
    }
  ]
}
```

**Status Codes:**
- 200: Success
- 500: Internal server error

---

### POST /api/devices/types

Creates a new device type.

**Authentication:** Required (JWT token in Authorization header)

**Request Body:**
```json
{
  "name": "New Device Type",
  "hardware_type": 1
}
```

**Fields:**
- `name` (string, required): The name of the device type
- `hardware_type` (integer, required): The hardware type ID (see Hardware Types section below)

**Response:**
```json
{
  "message": "device type created successfully"
}
```

**Status Codes:**
- 201: Created successfully
- 400: Invalid request payload
- 500: Internal server error

---

### GET /api/devices/types/hardware

Retrieves a list of all available hardware types.

**Authentication:** Required (JWT token in Authorization header)

**Response:**
```json
{
  "device_types": [
    {
      "id": 1,
      "name": "MicroController"
    },
    {
      "id": 2,
      "name": "SingleBoardComputer"
    },
    {
      "id": 3,
      "name": "Sensor"
    },
    {
      "id": 4,
      "name": "Solar Charger"
    },
    {
      "id": 5,
      "name": "VoltageMeter"
    },
    {
      "id": 6,
      "name": "CurrentSensor"
    },
    {
      "id": 7,
      "name": "PowerMeter"
    },
    {
      "id": 8,
      "name": "Actuator"
    }
  ]
}
```

**Status Codes:**
- 200: Success
- 500: Internal server error

## Hardware Types

The following hardware types are supported:

1. **MicroController** - Microcontrollers like ESP8266, ESP32
2. **SingleBoardComputer** - Single board computers like Raspberry Pi
3. **Sensor** - Various sensors (temperature, humidity, motion, etc.)
4. **Solar Charger** - Solar charge controllers/MPPT controllers
5. **VoltageMeter** - Voltage measurement devices
6. **CurrentSensor** - Current measurement devices
7. **PowerMeter** - Devices that measure both voltage and current
8. **Actuator** - Control devices like relays, switches

## Notes

- Device types are used to categorize devices in the system
- Each device must be associated with a device type
- Hardware types determine the capabilities and control options for devices