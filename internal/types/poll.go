package types

type Poll struct {
	Question 	string 	`dynamodbav:"question"`
	Answer 		string	`dynamodbav:"answer"`
}
