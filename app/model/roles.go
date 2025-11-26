package model

type roles struct {
  id: UUID PRIMARY KEY 
  name: VARCHAR(50) UNIQUE NOT NULL 
  description: TEXT 
  created_at: TIMESTAMP DEFAULT NOW() 
} 