package common

// Bet represents a bet made by a client. It contains all the necessary information
// about the bet, including the agency, client's name, surname, document, birthdate, and the Number
type Bet struct {
	Agency string `json:"agency"`
	Name string `json:"name"`
	Surname string `json:"surname"`
	Document string `json:"document"`
	Birthdate string `json:"birthdate"`
	Number string `json:"number"`
}