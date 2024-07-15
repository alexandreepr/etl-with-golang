# ETL-with-golang

## Description

This repository is a simple ETL app written in golang (Gin + Air + GORM + PostgresSQL).

## Requirements

- docker
- docker-compose

## Steps to Run the Application

1. Clone the repository.
2. Run `docker-compose up --build` in the root folder of the project.
3. The application will run on http://0.0.0.0:8000/.

## API Endpoints

### GET /health

- **Response**:
  - 200 OK: Server is running
  ```json
  {
    "status": "All good, captain!"
  }
  ```

### POST /api/v1/file-import

- **Request**:
  - Content-Type: multipart/form-data
  - Body: form field "file" containing a .txt file
- **Responses**:

  - 200 OK: File uploaded and processed successfully

  ```json
  {
      "message": "File uploaded and data saved successfully.",
      "importacao_id": "<uuid>"
  }
  - 400 Bad Request: Invalid file or request
  - 500 Internal Server Error: Server-side processing error

  ```

### GET /api/v1/import-report

- **Query Parameters**:
  - importacaoId: UUID of the import operation (required)
- **Responses**:
  - 200 OK: Report retrieved successfully
  ```json
  {
      "totalRows": <int>,
      "invalidCPFCount": <int>,
      "invalidLojaMaisFrequenteCNPJCount": <int>,
      "invalidLojaUltimaCompraCNPJCount": <int>
  }
  - 400 Bad Request: Invalid importacaoId
  - 404 Not Found: Import operation records not found
  - 500 Internal Server Error: Server-side processing error
  ```

## Run tests

1. Run `docker ps`
2. Get the ID of the Golang container
3. Run `docker exec -it <container-id> sh`
4. Run `go test ./...`

## Access PgAdmin

1. Open [http://localhost:5050/browser](http://localhost:5050/browser) in your browser
2. Login with:
   - Username: `admin@admin`
   - Password: `root`
3. Create a new server with:
   - Host: `postgres_db`
   - User: `mamun`
   - Password: `123`
