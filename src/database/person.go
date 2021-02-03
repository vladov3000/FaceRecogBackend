package database

import "go.mongodb.org/mongo-driver/bson"

type Person struct {
	Title       string    `json:"title" bson:"title"`
	Subtitle    string    `json:"subtitle" bson:"subtitle"`
	Encoding    bson.A    `json:"encoding" bson:"encoding"`
	ExtraFields []string  `json:"extra_field,omitempty" bson:"extra_field"`
	BBox        []float32 `json:"bbox,omitempty" bson:"-"`
	Known       bool      `json:"known,omitempty" bson:"-"`
}

func NewPerson(title string, subtitle string, encoding []float32, extra []string) Person {
	return Person{
		Title:       title,
		Subtitle:    subtitle,
		Encoding:    Float32SliceToBsonA(encoding),
		ExtraFields: extra,
	}
}

func Float32SliceToBsonA(slice []float32) bson.A {
	res := bson.A{}
	for _, i := range slice {
		res = append(res, i)
	}
	return res
}
