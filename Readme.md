# Online Booking System Web

# Customizable Online Booking System  

## Overview  
A web-based booking system designed for service-based businesses like salons, restaurants, and consultants. Currently, the following features are being implemented:  

- **User Authentication**: Local login system.  
- **CRUD Operations**: Manage organizations, locations, booking resources, and reservations.  
- **Authentication Middleware**: Secure API endpoints and pages.  
- **Page Rendering**: Handlebars.js for templating, styled with SCSS.  

## Tech Stack  

### Backend  
- **Golang (Fiber)** for API and server-side rendering.  
- **MongoDB** for data storage.  

### Frontend  
- **Templating**: Handlebars.js.  
- **Styling**: SCSS.  
- **Interactivity**: jQuery.  

## Next Steps  
- Implement user roles and permissions.  
- Add real-time availability for booking resources.  
- Improve UI/UX with dynamic updates. 

## How to Install Locally

### 1. Prerequisites  
Ensure you have the following installed on your system:  
- [Go](https://go.dev/dl/) (latest version)  
- [MongoDB](https://www.mongodb.com/try/download/community) (or use MongoDB Atlas)  
- [Node.js](https://nodejs.org/) (for SCSS compilation)  
- [Docker](https://www.docker.com/) (optional for containerized setup)  

---

### 2. Clone the Repository  
```sh
# git clone https://github.com/TheAttentionSeeker5050/golang-fiber-web-server-4-online-booking-system.git
git clone git@github.com:TheAttentionSeeker5050/golang-fiber-web-server-4-online-booking-system.git
cd customizable-booking-system
```

---

### 3. Set Up Environment Variables  
Create a `.env` file in the root directory and add the following:  

```env
MONGODB_URI=

# Cookie settings
COOKIE_EXPIRES_IN_DAYS=7

# JWT Secret
JWT_SECRET=

# Google OAuth2.0 Credentials (Optional)
OAUTH2_GOOGLE_CLIENT_ID=
OAUTH2_GOOGLE_CLIENT_SECRET=
OAUTH2_GOOGLE_REDIRECT_URL=
OAUTH2_GOOGLE_ORIGIN=
```

> **Note:** Replace sensitive values with your own credentials. If you are not using Google OAuth, you can comment out those lines.

---

### 4. Install Dependencies  
**Backend (Go + Fiber)**  
```sh
go mod tidy
```

### Frontend (SCSS Compilation - Optional if using plain CSS)  
```sh
npm run sass:build # For listening to changes
npm run sass # for compiling into production css code without listening to changes
```

---

### 5. Run the Application  
```sh
go run server.go
```
The server should start at **`http://localhost:8080`**.

---

### 6. Access the Application  
- **Dashboard**: `http://localhost:8080/dashboard`  
- **Login**: `http://localhost:8080/login`  
- **Reservations**: `http://localhost:8080/reservations`  

---

### 7. Docker Setup (Optional)  
To run the system using Docker, first build the image and then start the container.  

```sh
docker build -t booking-system .
docker run -p 8080:8080 --env-file .env booking-system
```

---

### 8. Troubleshooting  
- **MongoDB connection issues?** Check if your `MONGODB_URI` is correct and that MongoDB is running.  
- **Port already in use?** Change the default port in `main.go` or stop conflicting processes.  
- **OAuth login not working?** Ensure Google OAuth credentials are set up correctly in the Google Developer Console.  

## Monitoring
To monitor resource usage run the following command: 
```sh
top -p <PID>
```

The PID stands for process ID which you can find when executing the `go run server.go` command