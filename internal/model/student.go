package model

type Student struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name" binding:"required"`
	Serie        string `json:"serie" binding:"required"`
	Dateofbirth  string `json:"date_of_birth" binding:"required"`
	NameOfMother string `json:"name_of_mother" binding:"required"`
	NameOfDad    string `json:"name_of_dad" binding:"required"`
	Ra           string `json:"ra" binding:"required"`
	SchoolID     int    `json:"school_id" binding:"required"`
}

type StudentResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Serie        string `json:"serie"`
	Dateofbirth  string `json:"date_of_birth"`
	NameOfMother string `json:"name_of_mother"`
	NameOfDad    string `json:"name_of_dad"`
	Ra           string `json:"ra"`
	School       School `json:"school"`
}

type StudentServiceResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Serie        string `json:"serie"`
	Dateofbirth  string `json:"date_of_birth"`
	NameOfMother string `json:"name_of_mother"`
	NameOfDad    string `json:"name_of_dad"`
	Ra           string `json:"ra"`
}
