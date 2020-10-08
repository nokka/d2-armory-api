# Armory

The armory backend is responsible for storing and parsing d2s characters. 
It also exposes an HTTP API to get characters as JSON.

--- 

## API

### Get a character by name
```http
GET /api/v1/characters?name=nokka
```

### Deprecated handler for consumers who rely on it
```http
GET /retrieving/v1/character?name=nokka
```

### Health check
```http
GET /health
```

---

## Environment variables

| Name           	| Default         	|
|----------------	|-----------------	|
| HTTP_ADDRESS   	| `:80`           	|
| MONGO_HOST     	| `mongodb:27017` 	|
| MONGO_DB       	| `armory`        	|
| MONGO_USERNAME 	|                 	|
| MONGO_PASSWORD 	|                 	|
| D2S_PATH       	|                 	|
| CACHE_DURATION 	| `3m`            	|

--- 

## Data storage
The armory API relies on [mongodb](https://www.mongodb.com/) to store the data.

