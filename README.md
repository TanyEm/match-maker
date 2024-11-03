# Match Maker

Match Maker is a service designed to manage player lobbies and match players based on certain criteria. This service provides APIs to join lobbies and create matches. Once the expected number of players have joined (or 30 seconds of waiting has passed), the match creates.

## Table of Contents

- [Structure](#structure)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Configuration](#configuration)
- [Running Tests](#running-tests)

## Structure

The project structure is shown in the diagram below:

![match-maker structure](pics/match-maker.png)

## Installation
To install the Match Maker service, clone the repository and build the project:

```bash
git clone git@github.com:TanyEm/match-maker.git
cd match-maker
make build
```

## Usage

To run the Match Maker service, use the following command:

```bash
./bin/match-maker
```

## API Endpoints

### View Full API Documentation in Swagger Editor
1. Open [Swagger Editor](https://editor.swagger.io/).
2. Upload the `swagger.json` file by selecting the "File" menu, then "Import file."
3. The Swagger Editor will display the API documentation for you to interact with.

`GET /ping`

Health check endpoint to verify the service is running.

Response:

```json
{
  "message": "pong"
}
```

`POST /lobby`

Join a lobby with player details.

Request:

```json
{
  "player_id": "player1",
  "level": 1,
  "country": "USA"
}
```

Response:

```json
{
  "join_id": "00000000-0000-0000-0000-000000000000"
}
```

`GET /match`

Check a match for a player in the lobby. **Note**, this is supposed to be a polling-request from client side, meaning wait at most 30 seconds (default MATCH_MAKING_TIME configured on the server side). The client may also have lower timeouts and retry the request until the result is provided.

Request:

```bash
GET /match?join_id=00000000-0000-0000-0000-000000000000
```

Response:

```json
{
  "match_id": "00000000-0000-0000-0000-000000000000"
}
```

In case the match for a player, e.g. there was only 1 player joining the match, a **404** Response is returned:
```json
{
	"error": "no match for the player, try to join the lobby again"
}
```

`GET /leaderboard`

Get leaderboard by match_id.

Request:

```bash
GET /leaderboard?match_id=00000000-0000-0000-0000-000000000000
```

Response:

```json
{
	"match_id": "00000000-0000-0000-0000-000000000000",
	"players": [
		{
			"player_id": "123",
			"level": 4,
			"country": "FIN",
			"score": 0
		},
		{
			"player_id": "1234",
			"level": 4,
			"country": "FIN",
			"score": 0
		}
	]
}
```

## Configuration

The service can be configured using environment variables:

 - PORT: The port on which the service will run (default: 8080).
 - SHUTDOWN_DURATION: The duration to wait before shutting down the service (default: 3s).
 - MATCH_MAKING_TIME: The duration time for match making players in lobby (default: 30s)

## Running Tests

To run the tests for the Match Maker service, use the following command:

```bash
make test
```

If there are issues with mocks, ensure that [gomock](https://github.com/uber-go/mock) `mockgen` is installed:

```bash
go install go.uber.org/mock/mockgen@latest
export PATH=$PATH:$(go env GOPATH)/bin
mockgen -version
```
