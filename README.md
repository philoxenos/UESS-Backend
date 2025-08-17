# User Email Authentication Backend

This is a Go backend service that authenticates user emails against a local JSON database and supports the user registration/authentication flow described in the project requirements.

## Setup and Running

1. Make sure you have Go installed on your system.
2. Run the server:

```
go run main.go
```

The server will start on port 8080.

## API Endpoints

### 1. Authenticate User

**Endpoint**: `/authenticate`
**Method**: POST
**Description**: Checks if a user email exists in the database.

**Request Body**:
```json
{
  "email": "user@example.com",
  "name": "John",
  "surname": "Doe",
  "createdAt": "2023-08-17T12:00:00Z"
}
```

**Response** (If user exists):
```json
{
  "status": "success",
  "exists": true,
  "user": {
    "email": "user@example.com",
    "name": "John",
    "surname": "Doe",
    "createdAt": "2023-08-17T12:00:00Z",
    "role": "user"
  }
}
```

**Response** (If user doesn't exist):
```json
{
  "status": "success",
  "exists": false
}
```

### 2. Update User

**Endpoint**: `/update`
**Method**: PUT
**Description**: Updates a user's password and role, or creates a new user if one doesn't exist.

**Request Body**:
```json
{
  "email": "user@example.com",
  "name": "John",
  "surname": "Doe",
  "createdAt": "2023-08-17T12:00:00Z",
  "password": "newpassword",
  "role": "user"
}
```

**Response**:
```json
{
  "status": "success"
}
```

## Database

The server uses a simple JSON file (`db.json`) as its database. You can manually add users to this file in the following format:

```json
{
  "users": [
    {
      "email": "user@example.com",
      "name": "John",
      "surname": "Doe",
      "createdAt": "2023-08-17T12:00:00Z",
      "password": "password123",
      "role": "user"
    }
  ]
}
```

## Integration with Android App

1. Android app serializes user data (email, name, surname, createdAt) to JSON
2. Backend checks if the email exists in the database
3. If the email exists, the Android app prompts for a password
4. After password entry, Android app updates the local SQLite Cipher database
5. This approach ensures users can log in even when offline (after initial authentication)
