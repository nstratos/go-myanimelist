package mal_test

import (
	"context"
	"fmt"
	"net/url"

	"github.com/nstratos/go-myanimelist/mal"
)

func ExampleForumService_Boards() {
	ctx := context.Background()

	c := mal.NewClient(nil)

	// Ignore the 3 following lines. A stub server is used instead of the real
	// API to produce testable examples. See: https://go.dev/blog/examples
	server := newStubServer()
	defer server.Close()
	c.BaseURL, _ = url.Parse(server.URL)

	forum, _, err := c.Forum.Boards(ctx)
	if err != nil {
		fmt.Printf("Forum.Boards error: %v", err)
		return
	}
	for _, category := range forum.Categories {
		fmt.Printf("%s\n", category.Title)
		for _, b := range category.Boards {
			fmt.Printf("|-> %s\n", b.Title)
			for _, b := range b.Subboards {
				fmt.Printf("    |-> %s\n", b.Title)
			}
		}
		fmt.Println("---")
	}
	// Output:
	// MyAnimeList
	// |-> Updates & Announcements
	// |-> MAL Guidelines & FAQ
	// |-> DB Modification Requests
	//     |-> Anime DB
	//     |-> Character & People DB
	//     |-> Manga DB
	// |-> Support
	// |-> Suggestions
	// |-> MAL Contests
	// ---
	// Anime & Manga
	// |-> News Discussion
	// |-> Anime & Manga Recommendations
	// |-> Series Discussion
	//     |-> Anime Series
	//     |-> Manga Series
	// |-> Anime Discussion
	// |-> Manga Discussion
	// ---
	// General
	// |-> Introductions
	// |-> Games, Computers & Tech Support
	// |-> Music & Entertainment
	// |-> Current Events
	// |-> Casual Discussion
	// |-> Creative Corner
	// |-> Forum Games
	// ---
}

func ExampleForumService_Topics() {
	ctx := context.Background()

	c := mal.NewClient(nil)

	// Ignore the 3 following lines. A stub server is used instead of the real
	// API to produce testable examples. See: https://go.dev/blog/examples
	server := newStubServer()
	defer server.Close()
	c.BaseURL, _ = url.Parse(server.URL)

	topics, _, err := c.Forum.Topics(ctx,
		mal.Query("JoJo opening"),
		mal.SortTopicsRecent,
		mal.Limit(2),
	)
	if err != nil {
		fmt.Printf("Forum.Topics error: %v", err)
		return
	}
	for _, t := range topics {
		fmt.Printf("ID: %5d, Title: %5q created by %q\n", t.ID, t.Title, t.CreatedBy.Name)
	}
	// Output:
	// ID: 1877721, Title: "What is the best JoJo opening?" created by "Ringtomb"
	// ID: 1851738, Title: "JoJo's Bizarre Adventures but its Yu Yu Hakusho Opening" created by "TinTin_29"
}

func ExampleForumService_TopicDetails() {
	ctx := context.Background()

	c := mal.NewClient(nil)

	// Ignore the 3 following lines. A stub server is used instead of the real
	// API to produce testable examples. See: https://go.dev/blog/examples
	server := newStubServer()
	defer server.Close()
	c.BaseURL, _ = url.Parse(server.URL)

	topicDetails, _, err := c.Forum.TopicDetails(ctx, 1877721, mal.Limit(3), mal.Offset(0))
	if err != nil {
		fmt.Printf("Forum.TopicDetails error: %v", err)
		return
	}
	fmt.Printf("Topic title: %q\n", topicDetails.Title)
	if topicDetails.Poll != nil {
		fmt.Printf("Poll: %q\n", topicDetails.Poll.Question)
		for _, o := range topicDetails.Poll.Options {
			fmt.Printf("- %-25s %2d\n", o.Text, o.Votes)
		}
	}
	for _, p := range topicDetails.Posts {
		fmt.Printf("Post: %2d created by %q\n", p.Number, p.CreatedBy.Name)
	}
	// Output:
	// Topic title: "What is the best JoJo opening?"
	// Poll: "What is the best JoJo opening?"
	// - Sono Chi No Sadame        23
	// - Bloody Stream             61
	// - Stand Proud               12
	// - End Of The World          14
	// - Crazy Noisy Bizarre Town  22
	// - Chase                     13
	// - Great Days                34
	// - Fighting Gold             15
	// - Traitor's Requiem         11
	// Post:  1 created by "Ringtomb"
	// Post:  2 created by "Kenzolo-folk"
	// Post:  3 created by "MechKingKillbot"
}
