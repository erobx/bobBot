package types

type Flight struct {
	UserID 			string 		`dynamodbav:"user_id"`
	FlightKey		string		`dynamodbav:"flight_key"`
	Number 			string		`dynamodbav:"number"`
	AirportCode 	string		`dynamodbav:"airport_code"`
	ArrivalDate		string		`dynamodbav:"arrival_date"`
	ArrivalTime 	string	`dynamodbav:"arrival_time"`
	DepartureTime 	string	`dynamodbav:"departure_time"`
}
