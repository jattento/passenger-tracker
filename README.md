# passenger-tracker
Volume Senior Software Engineer Take-Home Programming Assignment for Golang

# API Documentation
## /calculate Endpoint

### Description
The `calculate` endpoint accepts POST requests. The request must include one or more flights inside the body, where each flight is represented by a pair of airports that are connected by a direct flight.

### Request Body
The request body should be a JSON array object, containing JSON array objects, contaning two values, the first being the start airport and the second the end airport:

```json
[["EZE","MIA"],["MIA","BRA"],["BRA","EZE"],["EZE","NYC"]]
```

### Response
The response in case of success is a JSON array object contaning two values, the first being the start airport and the second the end airport:
```json
["EZE","NYC"]
```

### Errors

| Status code | 	Description                                                                                                             |
| --- |--------------------------------------------------------------------------------------------------------------------------|
| 400 | 	Invalid input data. All airport names should contain 3 letters.                                                         |
| 400 | 	Not all flights are connected.                                                                                          |
| 400 | 	Start and end airport are the same, it is not possible to determine with the current information which is the starting. |
| 405 | 	Only POST requests are supported on this endpoint.                                                                      |
| 500 | 	Internal error.                                                                                                         |

Example error response:
```json
{
    "code": 400,
    "message": "not all flights are connected"
}
```

### Usage

Start the web service with the following command: `go run ./cmd/webservice` it will start listining at the port `8080`

Execute this simple example and try it out!
```bash
curl -sS -X POST \
-H "Content-Type: application/json" \
-d '[["BRA","EZE"],["EZE","MIA"],["MIA","BRA"],["EZE","NYC"]]' \
http://localhost:8080/calculate \
| jq .
``` 
