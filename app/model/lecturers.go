package model

type lecturers struct { 
  id: UUID PRIMARY KEY 
  user_id: UUID FOREIGN KEY -> users.id 
  lecturer_id: VARCHAR(20) UNIQUE NOT NULL 
  department: VARCHAR(100) 
  created_at: TIMESTAMP DEFAULT NOW() 
} 