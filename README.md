# d2 armory api

The armory backend is responsible for storing and parsing d2s characters. 
It also exposes an HTTP API to get characters as JSON.

--- 

## Environment variables

| Name                	| Default         	|
|---------------------	|-----------------	|
| HTTP_ADDRESS        	| `:80`           	|
| MONGO_HOST          	| `mongodb:27017` 	|
| MONGO_DB            	| `armory`        	|
| MONGO_USERNAME      	|                 	|
| MONGO_PASSWORD      	|                 	|
| D2S_PATH            	|                 	|
| CACHE_DURATION      	| `3m`            	|
| STATISTICS_USER     	|                 	|
| STATISTICS_PASSWORD 	|                 	|
| CORS_ENABLED        	| `false`         	|
| LOG_REQUESTS        	| `false`         	|

--- 

## API

#### Get a character by name
Gets the character by name, either served through the mongoDB
cache or by parsing the d2s binary if the cache duration has expired.
```http
GET /api/v1/characters?name=nokka
```

#### Deprecated handler for consumers who rely on it
Deprecated handler used by < v1.0.0 users.
```http
GET /retrieving/v1/character?name=nokka
```

#### Health check
```http
GET /health
```

---

## Data storage
The armory API relies on [mongodb](https://www.mongodb.com/) to store the data.

