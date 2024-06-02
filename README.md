# Torch Sync Server

Sync torch database with the client. Works by listening for PostgreSQL notifications that happen on every INSERT, UPDATE, and DELETE operation. Client connected via websockets gets a message with the updated data whenever the server receives a notification from the database.

## Connecting to the server

Establish a websocket connection to the following address

```
ws://localhost:8000/sync
```

## Messsaging Protocol

There are three supported operations: INSERT, UPDATE, and DELETE. The format of these messages applies for both ways communications - from server to client and from client to server.

```json
{
  "op": "INSERT",
  "itemID": "ds34jhb2134",
  "diffs": {
    "title": {
      "val": "New Task",
      "cl": 123
    },
    "status": {
      "val": "ACTIVE",
      "cl": 12
    }
  }
}
```

```json
{
  "op": "UPDATE",
  "itemID": "ds34jhb2134",
  "diffs": {
    "title": {
      "val": "New Task",
      "cl": 123
    },
    "status": {
      "val": "ACTIVE",
      "cl": 12
    }
  }
}
```

```json
{
  "op": "DELETE",
  "itemID": "ds34jhb2134",
  "cl": 123
}
```

## TODOS

- [ ] Make sure ParentID column cannot accept IDs that are of a different user and that it does not equal to the item ID itself
- [x] Is `custom.disable_trigger` required? (No)
