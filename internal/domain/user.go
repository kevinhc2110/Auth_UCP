package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        			uuid.UUID `json:"id"`         		
	Identification	string		`json:"identification"`	
	Name      			string    `json:"name"`       		
	Lastname  			string    `json:"lastname"`   		
	Email     			string    `json:"email"`      		
	Password  			string    `json:"-"`          			
	CreatedAt 			time.Time `json:"created_at"` 		
	UpdatedAt 			time.Time `json:"updated_at"` 		
	LastLoginAt 		time.Time `json:"lastlogin_at"`		
	Active					bool			`json:"active"`					
}


