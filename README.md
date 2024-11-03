# Match Maker

Match Maker is a service designed to manage player lobbies and match players based on certain criteria. This service provides APIs to join lobbies and create matches. Once the expected number of players have joined (or 30 seconds of waiting has passed), the match creates.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Configuration](#configuration)
- [Running Tests](#running-tests)
- [Contributing](#contributing)
- [License](#license)

## Installation
To install the Match Maker service, clone the repository and build the project:

```bash
git clone git@github.com:TanyEm/match-maker.git
cd match-maker
make build
## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Configuration](#configuration)
- [Running Tests](#running-tests)
- [Contributing](#contributing)
- [License](#license)

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

Check a match for a player in the lobby.

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

## Configuration

The service can be configured using environment variables:

 - PORT: The port on which the service will run (default: 8080).
 - SHUTDOWN_DURATION: The duration to wait before shutting down the service (default: 3s).

## Running Tests

To run the tests for the Match Maker service, use the following command:

```bash
make test
```

