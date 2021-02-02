package database

import "go.mongodb.org/mongo-driver/bson"

type Person struct {
	Title    string `json: "title" bson: "title"`
	Subtitle string `json: "subtitle" bson: "subtitle"`
	Encoding bson.A `json: "-" bson: "encoding`
}

func NewPerson(title string, subtitle string, encoding []float32) Person {
	var res Person

	res.Title = title
	res.Subtitle = subtitle

	res.Encoding = bson.A{}
	for _, i := range encoding {
		res.Encoding = append(res.Encoding, i)
	}

	return res
}
