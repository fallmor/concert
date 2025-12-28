# Concert Booking System

A little concert booking just to practice golang .

## Features

- User registration and authentication
- Show and artist management
- Fan registration for concerts
- Role-based access control (User, Moderator, Admin)
- Ratelimit
- Password reset workflow woith temporal

## Todo
- Change legacy html to typescript - react step by step
- When show expires deletes it (one day after using tempral)
-


## Usage

1. Clone the repository
   go mod download
   Create `.env` file
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_username
   DB_PASSWORD=your_password
   DB_NAME=concert_db
   DB_SSLMODE=disable

2. Start a temporal server : https://docs.temporal.io/cli/server
3. start the worker:
  1. go to temporal folder in the root directory
  2. go run worker.go

### Roles

- **User**: View and register for shows
- **Moderator**: Edit shows, artists, and fan registrations
- **Admin**: Full access including delete operations
