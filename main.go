package main

import (
	"fmt"
	"time"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

type Post struct {
	Title string
	Permalink string
	Url string
	Self bool
	Stickied bool
	Score int
	Flair string
	Author string
	Date time.Time
}

func main() {
	session, err := mgo.Dial("mongodb://localhost:27017/myDatabase");
	if err != nil {
		panic(err)
	}
	defer session.Close();

	c := session.DB("myDatabase").C("posts");

	var results []Post

	err = c.Find(nil).All(&results)

	if err != nil {
		fmt.Println(err);
	}

	for _, element := range results {
		fmt.Println(element.Title)
	}
}
