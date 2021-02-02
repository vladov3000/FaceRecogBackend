package database

import "go.mongodb.org/mongo-driver/bson"

type Person struct {
	Title       string    `json:"title" bson:"title"`
	Subtitle    string    `json:"subtitle" bson:"subtitle"`
	Encoding    bson.A    `json:"-" bson:"encoding"`
	ExtraFields bson.M    `json:"extra_field" bson:"extra_field"`
	BBox        []float32 `json:"bbox" bson:"-"`
}

func NewPerson(title string, subtitle string, encoding []float32, extra bson.M) Person {
	res := Person{
		Title:       title,
		Subtitle:    subtitle,
		ExtraFields: extra,
	}

	res.Encoding = bson.A{}
	for _, i := range encoding {
		res.Encoding = append(res.Encoding, i)
	}

	return res
}
