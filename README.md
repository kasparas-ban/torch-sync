# Torch Sync Server

Sync torch database with the client. Works by listening for PostgreSQL notifications that happen on every INSERT, UPDATE, and DELETE operation. Client connected via websockets gets a message with the updated data whenever the server receives a notification from the database.
