package main

import (
	"fmt"
	"time"
	"sort"
	"strings"
	"gopkg.in/mgo.v2"
	"os"
	"gopkg.in/mgo.v2/bson"
)

type Post struct {
	Title string
	Permalink string
	Url string
	Self bool
	Stickied bool
	Score int `bson:"score"`
	Flair string
	Author string
	Date time.Time `bson:"date"`
	Included bool
}

const topCount = 5

type ScoreSorter []Post

func (s ScoreSorter) Len() int           { return len(s) }
func (s ScoreSorter) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ScoreSorter) Less(i, j int) bool { return s[i].Score > s[j].Score }

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func main() {
	session, err := mgo.Dial("mongodb://vatyx:vatyx@ds041613.mlab.com:41613/dota2logger");
	if err != nil {
		panic(err)
	}
	defer session.Close();

	c := session.DB("dota2logger").C("posts");

	var results []Post

	query := bson.M {
		"date": bson.M {
			"$gt": time.Unix(1466353688, 0),
		},
	}

	err = c.Find(query).All(&results)

	if err != nil {
		fmt.Println(err);
	}

	sort.Sort(ScoreSorter(results))

	fmt.Println("number of posts: ", len(results))

	for _, element := range results {
		fmt.Println(element.Score)
	}

	topPosts := topPosts(results)
	topPics := topPics(results)
	topSelf := topSelf(results)
	topVideos := topVideos(results)
	topFluff := topFluff(results)
	topDiscussion := topDiscussion(results)
	topEsports := topEsports(results)

	stickies := stickies(results)

	f, err := os.Create("output.txt")
	defer f.Close()

	f.WriteString("##**Top of /r/dota2 for this week**\n")

	fmt.Println("top posts")
	f.WriteString("\n##Top Posts\n")
	for index, element := range topPosts {
		str := formatedString(index, element)
		f.WriteString(str)
		fmt.Println(element.Title)
	}

	fmt.Println("top self")
	f.WriteString("\n##Top Self\n")
	for index, element := range topSelf {
		str := fmt.Sprintf("%d. [%s](%s)\n", index+1, element.Title, element.Url)
		f.WriteString(str)
		fmt.Println(element.Title)
	}

	fmt.Println("top pics")
	f.WriteString("\n##Top Images\n")
	for index, element := range topPics {
		str := formatedString(index, element)
		f.WriteString(str)
		fmt.Println(element.Title)
	}

	fmt.Println("top videos")
	f.WriteString("\n##Top Videos\n")
	for index, element := range topVideos {
		str := formatedString(index, element)
		f.WriteString(str)
		fmt.Println(element.Url)
	}

	fmt.Println("top discussion")
	f.WriteString("\n##Top Discussion\n")
	for index, element := range topDiscussion {
		str := formatedString(index, element)
		f.WriteString(str)
		fmt.Println(element.Url)
	}

	fmt.Println("top esports")
	f.WriteString("\n##Top eSports\n")
	for index, element := range topEsports {
		str := formatedString(index, element)
		f.WriteString(str)
		fmt.Println(element.Url)
	}

	fmt.Println("top fluff")
	f.WriteString("\n##Top Fluff\n")
	for index, element := range topFluff {
		str := formatedString(index, element)
		f.WriteString(str)
		fmt.Println(element.Url)
	}

	fmt.Println("top stickies")
	f.WriteString("\n##Stickies of this week\n")
	for index, element := range stickies {
		str := fmt.Sprintf("%d. [%s](%s)\n", index+1, element.Title, element.Url)
		f.WriteString(str)
		fmt.Println(element.Url)
	}
}

func formatedString(index int, element Post) string {
	return fmt.Sprintf("%d. [%s](%s) ([Comments](http://reddit.com%s)) by /u/%s\n", index+1, element.Title, element.Url, element.Permalink, element.Author)

}

func topPosts(results []Post) []Post {
	var ret []Post
	for i := 0; len(ret) < topCount; i++ {
		if !results[i].Stickied {
			results[i].Included = true
			ret = append(ret, results[i])
		}
	}

	return ret
}

func topPics(results []Post) []Post {
	var ret []Post
	for i := 0; len(ret) < topCount + 2; i++ {
		if !results[i].Included && (strings.Contains(results[i].Url, "imgur") || strings.Contains(results[i].Url, "gyazo") || strings.Contains(results[i].Url, "puu.sh")) {
			results[i].Included = true
			ret = append(ret, results[i])
		}
	}

	return ret
}

func topSelf(results []Post) []Post {
	var ret []Post
	for i := 0; len(ret) < topCount; i++ {
		if !results[i].Included && results[i].Self && !results[i].Stickied {
			results[i].Included = true
			ret = append(ret, results[i])
		}
	}

	return ret
}

func topVideos(results []Post) []Post {
	var ret []Post
	for i := 0; len(ret) < topCount + 2; i++ {
		if !results[i].Included && (strings.Contains(results[i].Url, "youtube") || strings.Contains(results[i].Url, "gfycat") || strings.Contains(results[i].Url, "oddshot") || strings.Contains(results[i].Url, "youtu.be") || strings.Contains(results[i].Url, "livecap")) {
			results[i].Included = true
			ret = append(ret, results[i])
		}
	}

	return ret
}

func stickies(results []Post) []Post {
	var ret []Post
	for i := 0; i < len(results); i++ {
		if results[i].Stickied {
			results[i].Included = true
			ret = append(ret, results[i])
		}
	}

	return ret
}

func topDiscussion(results []Post) []Post {
	var ret []Post
	for i := 0; len(ret) < topCount; i++ {
		if !results[i].Included && strings.Contains(results[i].Flair, "Discussion") && !results[i].Stickied {
			results[i].Included = true
			ret = append(ret, results[i])
		}
	}

	return ret
}

func topEsports(results []Post) []Post {
	var ret []Post
	for i := 0; len(ret) < topCount; i++ {
		if !results[i].Included && strings.Contains(results[i].Flair, "eSports") && !results[i].Stickied {
			results[i].Included = true
			ret = append(ret, results[i])
		}
	}

	return ret
}

func topFluff(results []Post) []Post {
	var ret []Post
	for i := 0; len(ret) < topCount; i++ {
		if !results[i].Included && strings.Contains(results[i].Flair, "Fluff") && !results[i].Stickied {
			results[i].Included = true
			ret = append(ret, results[i])
		}
	}

	return ret
}
