package model

type students struct { 
  id: UUID PRIMARY KEY 
  user_id: UUID FOREIGN KEY -> users.id 
  student_id: VARCHAR(20) UNIQUE NOT NULL 
  program_study: VARCHAR(100) 
  academic_year: VARCHAR(10) 
  advisor_id: UUID FOREIGN KEY -> lecturers.id 
  created_at: TIMESTAMP DEFAULT NOW() 
} 