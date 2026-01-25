# CampusCommons

CampusCommons is a full-stack web forum application built to emulate a reddit style website. It supports a classic forum hierarchy:

**Topics -> Posts -> Comments**

Users can log in with a username, create and manage topics and posts, and participate via comments. The frontend communicates with the backend through a RESTful API.

---

## Features

### Authentication
- **Username-only login**: Users enter a username; an account is created automatically if it’s new.
- **Session management**: User session is stored in browser `localStorage`.

### Forum Structure
- **Topics**
  - Title and description
  - Created by a user
- **Posts**
  - Belong to a topic
  - Contain a title and body
- **Comments**
  - Belong to a post
  - Contain body text only

### CRUD Support
- **Topics & Posts**: Full CRUD
  - Create via forms
  - Read via list + detail views
  - Update via editing with pre-filled fields
  - Delete with cascade deletion
- **Comments**: Create + Read only  
  - Update/Delete for comments were intentionally omitted to prioritize core requirements within scope.

---

## Tech Stack

### Backend
- **Language**: Go (1.22+)
- **HTTP**: `net/http` (standard library)
- **Database**: SQLite (embedded, file-based) via `mattn/go-sqlite3`
- **CORS**: `rs/cors`

### Frontend
- **Framework**: React 18 + TypeScript
- **Routing**: `react-router-dom`
- **State**: Local component state + `localStorage`
- **Styling**: Minimal inline styles (kept simple and readable)

### Dev / Tooling
- **Version Control**: Git + GitHub
- **Local Dev Ports**
  - Backend: `http://localhost:8080`
  - Frontend: `http://localhost:3000`

---

## Project Structure

CampusCommons/
├── backend/                # Go backend
│   ├── database/           # DB connection & schema
│   ├── handlers/           # API route handlers
│   ├── data/               # SQLite database file (ignored in Git)
│   └── main.go             # Server entry point
│
├── frontend/               # React + TypeScript frontend
│   ├── src/
│   │   ├── pages/          # Page components (Login, TopicList, etc.)
│   │   ├── types/          # TypeScript interfaces
│   │   └── services/       # API service functions
│   └── ...
│
├── README.md
└── FINAL_REPORT.md         # Final project documentation


## AI declaration
AI tools were used in the development of this project to support frontend design decisions. In particular, AI assistance was used to refine the visual styling of the React pages, including colour selection, layout/positioning, and general UI presentation. All application logic, feature implementation, and integration between the React frontend and Go backend were implemented by me.


## How to Run Locally

### 1. Start the backend

cd backend
go run main.go
# Server runs on http://localhost:8080

### 2. Start the frontend
cd ../frontend
npm install
npm start
# App opens at http://localhost:3000

