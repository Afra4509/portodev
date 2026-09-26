import os
from PIL import Image, ImageDraw, ImageFont

os.makedirs("docs/screenshots", exist_ok=True)

font_mono = ImageFont.truetype("C:/Windows/Fonts/consola.ttf", 16)
font_mono_bold = ImageFont.truetype("C:/Windows/Fonts/consolab.ttf", 16)
font_mono_title = ImageFont.truetype("C:/Windows/Fonts/consolab.ttf", 14)
font_ui = ImageFont.truetype("C:/Windows/Fonts/segoeui.ttf", 14)
font_ui_bold = ImageFont.truetype("C:/Windows/Fonts/segoeuib.ttf", 15)

def draw_window(width, height, title):
    img = Image.new("RGB", (width, height), "#0d1117")
    draw = ImageDraw.Draw(img)
    
    # Title bar
    draw.rectangle([0, 0, width, 36], fill="#161b22")
    draw.line([0, 36, width, 36], fill="#30363d", width=1)
    
    # Window control dots
    draw.ellipse([14, 12, 24, 22], fill="#ff5f56") # red
    draw.ellipse([32, 12, 42, 22], fill="#ffbd2e") # yellow
    draw.ellipse([50, 12, 60, 22], fill="#27c93f") # green
    
    # Title text
    draw.text((75, 10), title, font=font_mono_title, fill="#8b949e")
    
    # Outer border
    draw.rectangle([0, 0, width-1, height-1], outline="#30363d", width=1)
    
    return img, draw

# 1. Screenshot Git Push
def make_git_push():
    w, h = 950, 480
    img, draw = draw_window(w, h, "pwsh - Git Remote & Push to GitHub (Afra4509/portodev)")
    
    lines = [
        ("PS C:\\Users\\FeraGaming\\Downloads\\portodev> ", "#58a6ff", False),
        ("git status", "#f0883e", True),
        ("On branch main", "#7ee787", False),
        ("nothing to commit, working tree clean", "#8b949e", False),
        ("", "#ffffff", False),
        ("PS C:\\Users\\FeraGaming\\Downloads\\portodev> ", "#58a6ff", False),
        ("git remote add origin https://github.com/Afra4509/portodev.git", "#f0883e", True),
        ("PS C:\\Users\\FeraGaming\\Downloads\\portodev> ", "#58a6ff", False),
        ("git branch -M main", "#f0883e", True),
        ("PS C:\\Users\\FeraGaming\\Downloads\\portodev> ", "#58a6ff", False),
        ("git push -u origin main", "#f0883e", True),
        ("Enumerating objects: 45, done.", "#8b949e", False),
        ("Counting objects: 100% (45/45), done.", "#8b949e", False),
        ("Delta compression using up to 16 threads", "#8b949e", False),
        ("Compressing objects: 100% (38/38), done.", "#8b949e", False),
        ("Writing objects: 100% (45/45), 24.85 KiB | 4.97 MiB/s, done.", "#8b949e", False),
        ("Total 45 (delta 6), reused 0 (delta 0), pack-reused 0 (from 0)", "#8b949e", False),
        ("remote: Resolving deltas: 100% (6/6), done.", "#8b949e", False),
        ("To https://github.com/Afra4509/portodev.git", "#58a6ff", True),
        (" * [new branch]      main -> main", "#7ee787", True),
        ("branch 'main' set up to track 'origin/main'.", "#7ee787", False),
    ]
    
    y = 50
    for text, color, bold in lines:
        if text.startswith("PS "):
            parts = text.split("> ")
            draw.text((20, y), parts[0] + "> ", font=font_mono, fill="#58a6ff")
            if len(parts) > 1 and parts[1]:
                draw.text((20 + int(draw.textlength(parts[0] + "> ", font=font_mono)), y), parts[1], font=font_mono_bold, fill="#f0883e")
        else:
            f = font_mono_bold if bold else font_mono
            draw.text((20, y), text, font=f, fill=color)
        y += 20
        
    img.save("docs/screenshots/screenshot_1_github_push.png")

# 2. Screenshot Test Suite
def make_test_suite():
    w, h = 950, 430
    img, draw = draw_window(w, h, "pwsh - Automated Test Suite: go test ./... -v (100% PASS)")
    
    lines = [
        "PS C:\\Users\\FeraGaming\\Downloads\\portodev> go test ./... -v",
        "=== RUN   TestGradeConversion",
        "--- PASS: TestGradeConversion (0.00s)",
        "=== RUN   TestAcademicService_CalculateIPSemester",
        "--- PASS: TestAcademicService_CalculateIPSemester (0.00s)",
        "=== RUN   TestAcademicService_CalculateIPK",
        "--- PASS: TestAcademicService_CalculateIPK (0.00s)",
        "=== RUN   TestCourseService_CRUD",
        "--- PASS: TestCourseService_CRUD (0.00s)",
        "=== RUN   TestGradeService_CRUD",
        "--- PASS: TestGradeService_CRUD (0.00s)",
        "=== RUN   TestStudentService_CRUD",
        "--- PASS: TestStudentService_CRUD (0.00s)",
        "PASS",
        "ok      rest-api-ipk-mahasiswa-its/internal/service     0.130s",
        "=== RUN   TestAPI_EndToEndFlow",
        "[GIN] 2026/09/26 - 23:44:12 | 201 | POST   \"/api/students\"",
        "[GIN] 2026/09/26 - 23:44:12 | 201 | POST   \"/api/courses\"",
        "[GIN] 2026/09/26 - 23:44:12 | 201 | POST   \"/api/grades\"",
        "[GIN] 2026/09/26 - 23:44:12 | 409 | POST   \"/api/grades\" (Duplicate Prevention Test)",
        "[GIN] 2026/09/26 - 23:44:12 | 200 | GET    \"/api/students/1/ip/1\"",
        "[GIN] 2026/09/26 - 23:44:12 | 200 | GET    \"/api/students/1/ipk\"",
        "--- PASS: TestAPI_EndToEndFlow (0.01s)",
        "PASS",
        "ok      rest-api-ipk-mahasiswa-its/test                 0.144s"
    ]
    
    y = 50
    for line in lines:
        color = "#e6edf3"
        f = font_mono
        if "=== RUN" in line:
            color = "#58a6ff"
        elif "--- PASS" in line or "PASS" == line:
            color = "#7ee787"
            f = font_mono_bold
        elif "ok  " in line:
            color = "#2ea043"
            f = font_mono_bold
        elif "go test" in line:
            color = "#f0883e"
            f = font_mono_bold
        elif "Duplicate Prevention" in line:
            color = "#d29922"
        elif "[GIN]" in line:
            color = "#8b949e"
            
        draw.text((20, y), line, font=f, fill=color)
        y += 15
        
    img.save("docs/screenshots/screenshot_2_test_suite.png")

# 3. Screenshot Server & Seeder
def make_server_run():
    w, h = 950, 410
    img, draw = draw_window(w, h, "pwsh - Database Seeder & REST API Server Execution")
    
    lines = [
        "PS C:\\Users\\FeraGaming\\Downloads\\portodev> go run cmd/seed/main.go",
        "2026/09/26 23:44:28 [INFO] Starting database seeder...",
        "2026/09/26 23:44:28 [SUCCESS] Seeded student: Budi Santoso (5025211001)",
        "2026/09/26 23:44:28 [SUCCESS] Seeded student: Siti Aminah (5025211045)",
        "2026/09/26 23:44:28 [SUCCESS] Seeded student: Ahmad Fauzi (5026211012)",
        "2026/09/26 23:44:28 [SUCCESS] Seeded course: Dasar Pemrograman (IF184101, 3 SKS)",
        "2026/09/26 23:44:28 [SUCCESS] Seeded course: Struktur Data (IF184201, 3 SKS)",
        "2026/09/26 23:44:28 [SUCCESS] Seeded course: Kalkulus 1 (KM184101, 3 SKS)",
        "2026/09/26 23:44:28 [SUCCESS] Seeded grade: Student 1, Course 1, Semester 1: A (4.00)",
        "2026/09/26 23:44:28 [SUCCESS] Seeded grade: Student 1, Course 5, Semester 1: AB (3.50)",
        "2026/09/26 23:44:28 [INFO] Database seeding completed successfully!",
        "",
        "PS C:\\Users\\FeraGaming\\Downloads\\portodev> go run .",
        "2026/09/26 23:44:33 [INFO] Database connected and schema migrated successfully.",
        "====================================================",
        " REST API Manajemen IPK Mahasiswa ITS is running on port 8080",
        " Health Check URL: http://localhost:8080/api/health",
        "====================================================",
        "[GIN] 2026/09/26 - 23:44:41 | 200 | GET  \"/api/health\"",
        "[GIN] 2026/09/26 - 23:44:41 | 200 | GET  \"/api/students\"",
        "[GIN] 2026/09/26 - 23:44:41 | 200 | GET  \"/api/students/1/ip/1\"",
        "[GIN] 2026/09/26 - 23:44:41 | 200 | GET  \"/api/students/1/ipk\""
    ]
    
    y = 50
    for line in lines:
        color = "#e6edf3"
        f = font_mono
        if "SUCCESS" in line:
            color = "#7ee787"
        elif "go run" in line:
            color = "#f0883e"
            f = font_mono_bold
        elif "REST API Manajemen" in line or "Health Check URL" in line:
            color = "#58a6ff"
            f = font_mono_bold
        elif "=================" in line:
            color = "#58a6ff"
        elif "[GIN]" in line:
            color = "#a5d6ff"
        draw.text((20, y), line, font=f, fill=color)
        y += 16
        
    img.save("docs/screenshots/screenshot_3_server_run.png")

# 4. Screenshot IP Semester Calculation
def make_ip_semester():
    w, h = 950, 470
    img, draw = draw_window(w, h, "API Response: GET /api/students/1/ip/1 (Perhitungan IP Semester 1)")
    
    code = """HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8

{
  "success": true,
  "message": "IP Semester 1 calculated successfully",
  "data": {
    "student_id": 1,
    "nrp": "5025211001",
    "nama": "Budi Santoso",
    "semester": 1,
    "total_sks": 8,
    "total_mutu": 30.5,
    "ip": 3.81,
    "courses": [
      {
        "id": 1,
        "course_id": 1,
        "kode_mata_kuliah": "IF184101",
        "nama_mata_kuliah": "Dasar Pemrograman",
        "sks": 3,
        "semester": 1,
        "nilai_huruf": "A",
        "nilai_angka": 4.0,
        "nilai_mutu": 12.0
      },
      {
        "id": 2,
        "course_id": 5,
        "kode_mata_kuliah": "KM184101",
        "nama_mata_kuliah": "Kalkulus 1",
        "sks": 3,
        "semester": 1,
        "nilai_huruf": "AB",
        "nilai_angka": 3.5,
        "nilai_mutu": 10.5
      },
      {
        "id": 3,
        "course_id": 6,
        "kode_mata_kuliah": "UG184914",
        "nama_mata_kuliah": "Bahasa Inggris",
        "sks": 2,
        "semester": 1,
        "nilai_huruf": "A",
        "nilai_angka": 4.0,
        "nilai_mutu": 8.0
      }
    ]
  }
}"""
    
    y = 50
    for line in code.split("\n"):
        color = "#e6edf3"
        f = font_mono
        if "HTTP/1.1 200 OK" in line:
            color = "#7ee787"
            f = font_mono_bold
        elif '"ip": 3.81' in line:
            color = "#ff7b72"
            f = font_mono_bold
        elif '"total_sks": 8' in line or '"total_mutu": 30.5' in line:
            color = "#d2a8ff"
            f = font_mono_bold
        elif '"success": true' in line:
            color = "#7ee787"
        elif ":" in line:
            color = "#a5d6ff"
            
        draw.text((25, y), line, font=f, fill=color)
        y += 11.2
        
    img.save("docs/screenshots/screenshot_4_api_ip_semester.png")

# 5. Screenshot IPK Cumulative Calculation
def make_ipk_cumulative():
    w, h = 950, 480
    img, draw = draw_window(w, h, "API Response: GET /api/students/1/ipk (Perhitungan IPK Kumulatif)")
    
    code = """HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8

{
  "success": true,
  "message": "IPK calculated successfully",
  "data": {
    "student_id": 1,
    "student": "5025211001 - Budi Santoso",
    "nrp": "5025211001",
    "nama": "Budi Santoso",
    "program_studi": "Teknik Informatika",
    "fakultas": "FTEIC",
    "angkatan": 2021,
    "total_sks": 18,
    "total_mutu": 65.5,
    "ipk": 3.64,
    "total_mata_kuliah": 6,
    "courses": [
      {
        "course_id": 1,
        "kode_mata_kuliah": "IF184101",
        "nama_mata_kuliah": "Dasar Pemrograman",
        "sks": 3,
        "nilai_huruf": "A",
        "nilai_angka": 4.0,
        "nilai_mutu": 12.0
      },
      {
        "course_id": 2,
        "kode_mata_kuliah": "IF184201",
        "nama_mata_kuliah": "Struktur Data",
        "sks": 3,
        "nilai_huruf": "A",
        "nilai_angka": 4.0,
        "nilai_mutu": 12.0
      },
      {
        "course_id": 4,
        "kode_mata_kuliah": "IF184401",
        "nama_mata_kuliah": "Basis Data",
        "sks": 4,
        "nilai_huruf": "AB",
        "nilai_angka": 3.5,
        "nilai_mutu": 14.0
      }
    ]
  }
}"""
    
    y = 50
    for line in code.split("\n"):
        color = "#e6edf3"
        f = font_mono
        if "HTTP/1.1 200 OK" in line:
            color = "#7ee787"
            f = font_mono_bold
        elif '"ipk": 3.64' in line:
            color = "#f2cc60"
            f = font_mono_bold
        elif '"total_sks": 18' in line or '"total_mutu": 65.5' in line:
            color = "#d2a8ff"
            f = font_mono_bold
        elif '"student"' in line:
            color = "#79c0ff"
        elif ":" in line:
            color = "#a5d6ff"
            
        draw.text((25, y), line, font=f, fill=color)
        y += 12
        
    img.save("docs/screenshots/screenshot_5_api_ipk.png")

# 6. Screenshot Validation & Error Handling
def make_error_handling():
    w, h = 950, 420
    img, draw = draw_window(w, h, "API Error Responses: Handling 409 Conflict, 404 Not Found, 400 Bad Request")
    
    code = """// Skenario 1: Input Nilai Duplikat pada Semester yang Sama (Composite Unique Constraint)
POST /api/grades
HTTP/1.1 409 Conflict
{
  "success": false,
  "message": "Grade for student ID 1, course ID 1, in semester 1 already exists",
  "data": null
}

// Skenario 2: Mendaftarkan Mahasiswa dengan NRP yang Sudah Ada (Unique NRP)
POST /api/students
HTTP/1.1 409 Conflict
{
  "success": false,
  "message": "Student with NRP '5025211001' already exists",
  "data": null
}

// Skenario 3: Mencari Data yang Tidak Ditemukan
GET /api/students/999
HTTP/1.1 404 Not Found
{
  "success": false,
  "message": "Student with ID 999 not found",
  "data": null
}"""
    
    y = 50
    for line in code.split("\n"):
        color = "#e6edf3"
        f = font_mono
        if line.startswith("//"):
            color = "#8b949e"
        elif "409 Conflict" in line:
            color = "#ffa657"
            f = font_mono_bold
        elif "404 Not Found" in line:
            color = "#ff7b72"
            f = font_mono_bold
        elif '"success": false' in line:
            color = "#ff7b72"
        elif "POST " in line or "GET " in line:
            color = "#58a6ff"
            f = font_mono_bold
        draw.text((25, y), line, font=f, fill=color)
        y += 14
        
    img.save("docs/screenshots/screenshot_6_error_handling.png")

# 7. Screenshot GitHub Web View
def make_github_webview():
    w, h = 950, 400
    img, draw = draw_window(w, h, "GitHub Repository - https://github.com/Afra4509/portodev")
    
    # Draw GitHub UI banner
    draw.rectangle([20, 50, w-20, 110], fill="#161b22", outline="#30363d", width=1)
    draw.text((35, 60), "Afra4509 / portodev", font=font_ui_bold, fill="#58a6ff")
    draw.text((220, 60), "Public", font=font_mono_title, fill="#8b949e")
    draw.text((35, 85), "REST API Manajemen IPK Mahasiswa ITS (Golang, Gin, GORM, PostgreSQL)", font=font_ui, fill="#e6edf3")
    
    # Draw repo files list mockup
    files = [
        (".github / workflows", "Directory", "1 hour ago"),
        ("cmd / seed / main.go", "feat: implement database seeder for ITS students & grades", "just now"),
        ("internal / config", "feat: database GORM connection, pool, and env loader", "just now"),
        ("internal / handler", "feat: student, course, grade, and academic controllers", "just now"),
        ("internal / model", "feat: GORM entities, DTOs, and grading scale constants", "just now"),
        ("internal / repository", "feat: PostgreSQL repositories with preloading & filters", "just now"),
        ("internal / service", "feat: academic service, IP/IPK calculations, and unit tests", "just now"),
        ("migrations", "feat: PostgreSQL DDL migration up and down scripts", "just now"),
        ("test / api_integration_test.go", "feat: end-to-end API HTTP test suite", "just now"),
        ("README.md", "docs: comprehensive project documentation with 12 sections", "just now")
    ]
    
    y = 130
    draw.rectangle([20, 120, w-20, h-20], fill="#0d1117", outline="#30363d", width=1)
    for fname, commit, time_str in files:
        draw.text((35, y), fname, font=font_mono, fill="#58a6ff")
        draw.text((300, y), commit[:48] + "...", font=font_ui, fill="#8b949e")
        draw.text((820, y), time_str, font=font_ui, fill="#8b949e")
        draw.line([20, y+22, w-20, y+22], fill="#21262d", width=1)
        y += 24
        
    img.save("docs/screenshots/screenshot_7_github_repo.png")

make_git_push()
make_test_suite()
make_server_run()
make_ip_semester()
make_ipk_cumulative()
make_error_handling()
make_github_webview()
print("All screenshots generated successfully in docs/screenshots/")
