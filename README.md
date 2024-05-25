# Torch Sync Server

Sync torch database with the client. Works by listening for PostgreSQL notifications that happen on every INSERT, UPDATE, and DELETE operation. Client connected via websockets gets a message with the updated data whenever the server receives a notification from the database.

## Connecting to the server

Add `userID` as a query parameter to the sync endpoint.

```
ws://localhost:8000/sync?userID=123454567
```

## Messsaging Protocol

From client to server

```json
{
  "cmd": "UPDATE",
  "data": {
    "itemID": "5bax1usfu2uk",
    "title": "New Test"
  }
}
```

From server to client

```json
{
	"op": "UPDATE",
	"itemID": "absjfdnfds35m21",
	"diff": {
		"title": "New Task",
		"updatedAt": "2023-04-05T0503"
	}
}
```

## TODOS

- [] Make sure ParentID column cannot accept IDs that are of a different user and that it does not equal to the item ID itself
