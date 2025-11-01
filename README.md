# Hospital Management System (Backend)

A robust and scalable backend system for hospital management, built with Go (Golang). This system provides a comprehensive set of APIs to manage patients, doctors, appointments, and other hospital-related entities.

## Features

- **Patient Management**: Create, read, update, and delete patient records.
- **Doctor Management**: Manage doctor profiles, specialties, and availability.
- **Appointment Scheduling**: Schedule, reschedule, and cancel appointments with support for various constraints.
- **Medical Records**: Securely store and manage patient medical histories.
- **Authentication & Authorization**: Role-based access control (RBAC) to secure API endpoints.
- **RESTful API**: Clean and well-documented API endpoints following REST principles.
- **Data Validation**: Comprehensive input validation for data integrity.
- **Error Handling**: Consistent and informative error handling across the application.

## Installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/developersajadur/Hospital-Management-System-Backend-By-Golang
   cd Hospital-Management-System-Backend-By-Golang
   ```

2. **Install Dependencies**:
   Ensure you have Go (v1.18 or higher) installed on your system.
   
   Then, install the project dependencies:
   ```bash
   go mod download
   ```

3. **Environment Configuration**:
   Copy the example environment file and update it with your configuration.
   ```bash
   cp .env.example .env
   ```
   Edit the `.env` file with your database credentials and other settings.

4. **Database Setup**:
   Ensure your database is running. Then, run the migrations:
   ```bash
   go run cmd/migrate/main.go
   ```

5. **Run the Application**:
   ```bash
   go run cmd/main.go
   ```
   The server will start on the port specified in your environment configuration (default: `8080`).

## Usage

After starting the server, you can interact with the API using tools like `curl` or Postman.

**Example: Create a Patient**

```bash
curl -X POST http://localhost:8080/api/patients \
  -H "Content-Type: application/json" \
  -d '{
        "name": "John Doe",
        "age": 30,
        "gender": "male",
        "contact": "1234567890"
      }'
```

**Example: Retrieve a Patient**

```bash
curl http://localhost:8080/api/patients/1
```

For detailed API specifications, please refer to the [API Reference](#api-reference) section.

## API Reference

### Patients

- `GET /api/patients` - Get all patients
- `POST /api/patients` - Create a new patient
- `GET /api/patients/{id}` - Get a patient by ID
- `PUT /api/patients/{id}` - Update a patient
- `DELETE /api/patients/{id}` - Delete a patient

### Doctors

- `GET /api/doctors` - Get all doctors
- `POST /api/doctors` - Create a new doctor
- `GET /api/doctors/{id}` - Get a doctor by ID
- `PUT /api/doctors/{id}` - Update a doctor
- `DELETE /api/doctors/{id}` - Delete a doctor

### Appointments

- `GET /api/appointments` - Get all appointments
- `POST /api/appointments` - Create a new appointment
- `GET /api/appointments/{id}` - Get an appointment by ID
- `PUT /api/appointments/{id}` - Update an appointment
- `DELETE /api/appointments/{id}` - Delete an appointment

_Note: Replace `{id}` with the actual ID of the resource._

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.