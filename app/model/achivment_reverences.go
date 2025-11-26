package model

type achievement_references struct { 
  id: UUID PRIMARY KEY 
  student_id: UUID FOREIGN KEY -> students.id 
  mongo_achievement_id: VARCHAR(24) NOT NULL 
  status: ENUM('draft', 'submitted', 'verified', 'rejected') 
  submitted_at: TIMESTAMP 
  verified_at: TIMESTAMP 
  verified_by: UUID FOREIGN KEY -> users.id 
  rejection_note: TEXT 
  created_at: TIMESTAMP DEFAULT NOW() 
  updated_at: TIMESTAMP DEFAULT NOW() 
} 