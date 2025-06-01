# EngageSync - Running the Project Using Docker

## This project uses Go backend and PostgreSQL database, and can be run using Docker in any of the following three ways.

## Method 1: Clone GitHub Repo and Run with Docker Compose ---

### 1. Clone the repository:
git clone https://github.com/Rishi855/EngageSync.git
cd EngageSync

### 2. Run Docker Compose to pull images and start containers:
docker-compose up -d

### 3. To stop and remove containers:
docker-compose down

## Method 2: Use the docker-compose.yml File Directly ---

### 1. Place the docker-compose.yml file in your working directory.

### 2. Pull the required Docker images:
docker-compose pull

### 3. Start the containers:
docker-compose up -d

### 4. To stop the containers:
docker-compose down

## Method 3: Use the Shell Script Without Docker Compose ---

### You can run the project using a shell script (start-without-docker-compose.sh) without Docker Compose.

./start-engagesync.sh

### Notes:
- Make sure ports 5432 and 8080 are free before running containers.
- Data is persisted using Docker volume engagesync_postgres-data.
- Backend Docker image is hosted on Docker Hub: rushikesh855/engagesync-backend.

# Engagesync DB Setup

This project is built using **Go** and **PostgreSQL**, with database migrations handled using the **`golang-migrate`** package. Below are the installation steps and instructions for setting up the environment and running the project.

---

## 🚀 Prerequisites

Before starting, make sure you have the following installed:

- **Go** (Version 1.18 or later)  
  Download from: [Go Downloads](https://go.dev/dl/)

- **PostgreSQL**  
  Download from: [PostgreSQL Downloads](https://www.postgresql.org/download/)

- **Postman** (for API testing)  
  Download from: [Postman Downloads](https://www.postman.com/downloads/)

- **Swagger UI** (for API documentation)  
  Swagger UI will be available via a local HTML file.

---

## 🛠️ Installation Steps

### Step 1: Install Go and Set Up the Project

1. **Install Go** (if you haven't already):

    - Download and install Go from [https://go.dev/dl/](https://go.dev/dl/).
    - Verify installation:

    ```bash
    go version
    ```

2. **Initialize your Go project**:

    Create a directory for your project, initialize it with Go modules:

    ```bash
    mkdir engagesync
    cd engagesync
    go mod init engagesync
    ```

---

### Step 2: Install PostgreSQL and Set Up Database

1. **Install PostgreSQL** (if not already installed):
   - Download and install PostgreSQL from [PostgreSQL Downloads](https://www.postgresql.org/download/).
   - During installation, set your username and password (e.g., `postgres` and `root`).

2. **Create Database**:

    - Open `psql` or any PostgreSQL client and run the following command to create your database:

    ```sql
    CREATE DATABASE engagesyncdb;
    ```

    - Collect the database connection details:
      - **Host:** `localhost`
      - **Port:** `5432`
      - **User:** `postgres`
      - **Password:** `root`
      - **Database Name:** `engagesyncdb`

---

### Step 3: Install the `migrate` Package

1. **Install the migration package** to handle database migrations:

    ```bash
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    ```

2. Verify installation:

    ```bash
    migrate -version
    ```

---

### Step 4: Update `.env` File with Database Details

1. Create a `.env` file in the root of your project directory.

2. Update the `.env` file with your PostgreSQL details:

    ```bash
    DB_USER=postgres
    DB_PASSWORD=root
    DB_NAME=engagesyncdb
    DB_SSLMODE=disable
    ```

    This `.env` file will be used to load your PostgreSQL connection details dynamically into your application.

---

### Step 5: Run Migrations

1. **Create migration files**: Ensure that the migration SQL files are present under the `migrations/` directory. If not, create the necessary migration files.

2. **Run migrations** using the following command. This will read the details from the `.env` file for the connection string:

    ```bash
    migrate -database "postgres://postgres:root@localhost:5432/engagesyncdb?sslmode=disable" -path migrations up
    ```
    This command will load the `.env` values dynamically for the migration command.

---

### Step 6: Build and Run the Project

1. **Build the project**:

    ```bash
    go build -o 'engagesyncdb'
    ```

2. **Run the built executable**:

    On Windows:

    ```bash
    .\engagesyncdb.exe
    ```

    On Linux/macOS:

    ```bash
    ./engagesyncdb
    ```

---

## 📝 Postman Collection and API Testing

Once the project is running, you can use the **Postman collection** to test the API endpoints.

1. **Download the Postman Collection** file (provided in the repository or shared separately).
2. **Import the collection** into Postman.
3. **Run API requests** to test the endpoints.

---

## 📖 API Documentation (Swagger)

You can also view the API documentation through **Swagger UI**. To get access to the API documentation:

1. Open the `swagger.html` file located in the `docs/` folder (or a similar directory).
2. Open this file in your browser to get a **single view of the API's** available endpoints.

---

## 💡 Notes

- **Database Connection String**: Ensure your PostgreSQL credentials (username, password, and database name) are correct when setting the connection string.
- **Migration Path**: Make sure the migration files are in the correct directory as specified in the `-path` parameter (`migrations/`).
- **Running the Project**: Make sure the PostgreSQL database is running before starting the application.

---

## 🔧 Troubleshooting

- **"Connection refused" errors**: Ensure that PostgreSQL is running on `localhost:5432` and that your credentials are correct.
- **Migrations not applying**: Double-check your migration files for errors or missing SQL commands.
- **Postman collection issues**: Ensure that the Postman collection file is properly imported and that the server is running when making requests.

---
