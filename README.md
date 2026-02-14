## SKVMS 



###  Locations

##### Seven Days Readings

```bash
curl --location 'localhost:8080/api/locations/1/readings/seven' \
--header 'Authorization: ••••••'
```

Response be like . Remember , this is not neccesory if you already fetching current days reading.

1. Day will have 168 readings =  7 * ( 24 * 60)  / 60 

```json
 "readings": [
        {
            "voltage": 11.933858945125898,
            "current": 2.5190971118295886,
            "created_at": "2026-02-14T23:30:00+05:30"
        },
        {
            "voltage": 12.316988862944225,
            "current": 3.722999762532789,
            "created_at": "2026-02-14T22:30:00+05:30"
        },
        {
            "voltage": 12.597399154300117,
            "current": 3.415593740142388,
            "created_at": "2026-02-14T21:30:00+05:30"
        },
        {
            "voltage": 11.829738502375223,
            "current": 2.1565101725429026,
            "created_at": "2026-02-14T20:30:00+05:30"
        },
 ]
```