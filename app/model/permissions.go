package model

type permissions struct { 
  id: UUID PRIMARY KEY 
  name: VARCHAR(100) UNIQUE NOT NULL 
  resource: VARCHAR(50) NOT NULL 
  action: VARCHAR(50) NOT NULL 
  description: TEXT 
} 