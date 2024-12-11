Here’s the `README.md` text you can copy and paste directly:

```markdown
# Job Application Automation

A robust application designed to automate job applications through web-based form submission using Playwright and integrated database management. The system ensures that candidates' information is processed efficiently while handling errors and retries effectively.

---

## **Table of Contents**
- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Setup and Installation](#setup-and-installation)
- [Usage](#usage)
- [Technologies Used](#technologies-used)
- [Routes](#routes)
- [Scripts](#scripts)
- [Database Schema](#database-schema)
- [Future Enhancements](#future-enhancements)

---

## **Features**
- **Job Validation:** Ensures the job entry exists in the database before processing.
- **Web Automation:** Automates job application submission using Playwright.
- **Retry Mechanism:** Handles failed processes with a retry mechanism and configurable limits.
- **Database Integration:** Tracks application statuses in a MySQL database.
- **Scheduler:** Periodically processes pending tasks and ensures completion.
- **RESTful API:** Exposes a REST API to trigger job applications.

---

## **Project Structure**
```
job-application-automation/
├── controllers/
│   └── controllers.go
├── database/
│   └── database.go
├── models/
│   └── job.go
│   └── tracker.go
├── routes/
│   └── routes.go
├── services/
│   └── services.go
├── scripts/
│   └── job_script.js
├── main.go
└── README.md
```

---

## **Prerequisites**
- **Go:** Version 1.18 or higher
- **Node.js:** Version 14 or higher (for Playwright)
- **MySQL:** Version 5.7 or higher
- **Playwright:** Installed via `npm`
- **Gorm:** Go ORM for database management

---

## **Setup and Installation**

### **1. Clone the Repository**
```bash
git clone https://github.com/your-username/job-application-automation.git
cd job-application-automation
```

### **2. Install Dependencies**
- **Go Modules:**
  ```bash
  go mod tidy
  ```
- **Node.js Dependencies:**
  ```bash
  cd scripts
  npm install playwright mysql2
  cd ..
  ```

### **3. Configure Database**
- Update the `dsn` in `database/database.go` with your MySQL credentials:
  ```
  dsn := "root:password@tcp(localhost:3306)/job_db?charset=utf8mb4&parseTime=True&loc=Local"
  ```

### **4. Run Database Migrations**
```bash
go run main.go
```

### **5. Start the Application**
```bash
go run main.go
```

---

## **Usage**

### **Trigger Job Application**
Make a `POST` request to the `/apply` endpoint:
```json
{
  "job_id": 101,
  "candidate_id": 1001
}
```

### **Example cURL Command**
```bash
curl -X POST http://localhost:8082/apply \
-H "Content-Type: application/json" \
-d '{"job_id": 101, "candidate_id": 1001}'
```

---

## **Technologies Used**
- **Go:** For backend development
- **Echo:** Lightweight web framework
- **MySQL:** Relational database
- **Playwright:** Web automation and form submission
- **Gorm:** ORM for Go
- **Node.js:** To execute automation scripts

---

## **Routes**
| Endpoint       | Method | Description                      |
|----------------|--------|----------------------------------|
| `/apply`       | POST   | Triggers job application process |

---

## **Scripts**
- **Script Path:** `scripts/job_script.js`
- **Functionality:**
  - Automates job application using Playwright.
  - Updates the database with application status.

### **Run Script Manually**
```bash
node scripts/job_script.js <jobID> <candidateID>
```

---

## **Database Schema**

### **Jobs Table**
| Field          | Type         | Description                  |
|----------------|--------------|------------------------------|
| job_id         | INT          | Primary key                 |
| job_site       | VARCHAR(255) | Job site URL                |
| script_details | VARCHAR(255) | Script file path            |

### **Trackers Table**
| Field          | Type         | Description                  |
|----------------|--------------|------------------------------|
| tracker_id     | INT          | Primary key                 |
| job_id         | INT          | Foreign key to Jobs         |
| candidate_id   | INT          | Candidate identifier        |
| status         | VARCHAR(50)  | Current status              |
| output         | TEXT         | Script output               |
| error          | TEXT         | Error details (if any)      |
| retry_count    | INT          | Number of retry attempts    |
| max_retries    | INT          | Max retries allowed         |
| timestamp      | DATETIME     | Last update timestamp       |

---

## **Future Enhancements**
- Add authentication and authorization for API endpoints.
- Extend support for multiple job application platforms.
- Implement a UI for tracking application progress.
- Enhance error reporting and logging.

---

## **Contributors**
- [Your Name](https://github.com/your-username)

---

## **License**
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
```

