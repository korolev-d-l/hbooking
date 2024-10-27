# Hidden Workshop Booking System

## Project Description
This repository contains a reference solution to a junior-level test task. 
The project implements a booking system for a "hidden workshop" on Earth. 
The goal is to allow users to book workshop time slots in various time zones while ensuring 
no overlapping reservations. The implementation uses **Go**, **PostgreSQL**, and libraries 
such as **gin** and **pgx**. 
Data consistency is ensured through the use of transactions.

## Features
1. **Booking System**
    - Users can book workshops for specific dates and times in their local time zones.
    - Time slots can range from 30 minutes to 4 hours.
    - The system checks for overlapping bookings to prevent conflicts before confirming a reservation.

2. **Data**
    - Workshop hours are predefined for at least three different workshops, each with set hours in the local time zone. Note that some workshops may operate past midnight.
    - Workshop schedules are stored in the database.

3. **Technologies**
    - **Go** — backend language for the application’s business logic.
    - **PostgreSQL** — database for storing workshop and booking information.
    - **gin** and **pgx** libraries used for HTTP API handling and database operations, respectively.

## Repository Structure
- `cmd/hbooking/main.go` — application entry point.
- `docs/task.pdf` — original test task description.
- `internal/` — core application logic.
    - `hbooking/repository` — module for database interactions and schema handling.
    - `hbooking/handlers` — app handlers.
    - `hbooking/domain` — ubiquitous domain models.

## Getting Started
### Requirements
- **Go**
- **PostgreSQL**

### Installation and Run
1. Clone the repository:
    ```bash
    git clone https://github.com/spatecon/hbooking.git
    cd hbooking
    ```

2. Download dependencies:
    ```bash
    go mod download
    ```

3. Start the server:
    ```bash
    export DATABASE_URL="postgres://hts-user:hts-pass@localhost:5432/hts?sslmode=disable"
    export HTTP_PORT=8080
    go run cmd/hbooking/main.go
    ```
   
4. Add workshop_schedule data:
```postgresql
INSERT INTO workshop_schedules (workshop_id, workshop_timezone, begin_at, end_at) VALUES
(1, 'Europe/London', '08:00:00', '12:00:00'),
(2, 'Europe/Moscow', '10:00:00', '22:00:00'),
(3, 'Asia/Tokyo', '12:00:00', '15:00:00');
```

## API Endpoints
- `POST /api/v1/bookings/{workshop_id}` — create a new booking.
- `GET /api/v1/bookings/{workshop_id}` — retrieve a list of existing bookings.

### Request Examples

#### Create a Booking
```bash
curl --request POST \
  --url http://localhost:8080/api/v1/bookings/1 \
  --header 'content-type: application/json' \
  --data '{
	"client_id": "i@ela.sh",
  	"begin_at": "27-10-2024 14:00",
  	"end_at": "27-10-2024 14:30",
  	"client_timezone": "Europe/Amsterdam"
}'
```


####
#### Get Bookings
```bash
curl --request GET \
  --url http://localhost:8080/api/v1/bookings/1
```

## Data Consistency
Transactions are used in the project to maintain data integrity 
during booking operations that involve overlap checks.

---

> Note: If bookings overlap too often, the workshops might attract the attention of a certain **"Doctor."** 😉
