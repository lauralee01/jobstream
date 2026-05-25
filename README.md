# JobStream 🚀  
A modern job aggregation platform built with Go, Nuxt 3, and PostgreSQL — designed to collect jobs from multiple providers into a unified, fast, and scalable search experience.

---

## ✨ Features

### **Multi‑Source Job Aggregation**
Fetch jobs from multiple providers through pluggable fetchers:
- Remotive  
- Adzuna  
- Greenhouse  
- Arbeitnow  
- Easily extensible for more providers

### **Advanced Filtering & Search**
Search jobs by:
- Keyword  
- Location  
- Category  
- Platform  
- Pagination  
- Sorting  

### **Smart Category Normalization**
Jobs are automatically categorized using a centralized classification engine:
Engineering • Product • Design • Marketing • Data • Finance • HR • Security

### **Clean Backend Architecture**
Designed for scalability and maintainability:
- Repository Pattern  
- Service Layer  
- DTO Mapping  
- Domain Models  
- Dependency Injection  
- Context Propagation  
- Fetcher Interfaces  

### **Modern Frontend**
Built with:
- Nuxt 3  
- Vue 3  
- Tailwind CSS  
- SSR support  
- Dynamic filtering  
- Responsive UI  
- Dark mode  

---

## 🏗️ Tech Stack

### **Backend**
- Go  
- PostgreSQL  
- pgx  
- Docker  
- REST API  

### **Frontend**
- Nuxt 3  
- Vue 3  
- TailwindCSS  

### **Infrastructure**
- Neon (Postgres)  
- Render (Backend)  
- Vercel (Frontend)  

---

## 🚀 Running Locally

### **Clone the Repository**

git clone https://github.com/yourusername/jobstream.git
cd jobstream


### **Start PostgreSQL**
docker compose up -d

### **Run Migrations**
for file in migrations/*.sql; do
psql "postgres://username:password@localhost:5432/db_name?sslmode=disable" -f "$file"
done


---

## 🖥️ Backend Setup
cd backend
go mod tidy


### **Environment Variables**
Create `.env`:
DATABASE_URL=postgres://username:password@localhost:5432/db_name?sslmode=disable
PORT=8080
POSTGRES_USER=POSTGRES_USER
POSTGRES_PASSWORD=POSTGRES_PASSWORD
POSTGRES_DB=POSTGRES_DB

# job fetcher adzuna app key
ADZUNA_APP_ID=your_app_id
ADZUNA_APP_KEY=your_app_key


### **Start Server**
go run cmd/main.go


---

## 🌐 Frontend Setup
cd frontend
npm install


### **Environment Variables**
Create `.env`:
NUXT_PUBLIC_API_BASE=http://localhost:8080/api/v1


### **Start Dev Server**
npm run dev

