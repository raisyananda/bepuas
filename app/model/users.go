package model

type users struct {
	id: UUID PRIMARY KEY 
	username: VARCHAR(50) UNIQUE NOT NULL 
	email: VARCHAR(100) UNIQUE NOT NULL 
	password_hash: VARCHAR(255) NOT NULL 
	full_name: VARCHAR(100) NOT NULL 
	role_id: UUID FOREIGN KEY -> roles.id 
	is_active: BOOLEAN DEFAULT true 
	created_at: TIMESTAMP DEFAULT NOW() 
	updated_at: TIMESTAMP DEFAULT NOW() 
} 