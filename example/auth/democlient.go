package main

import (
	"context"
	"fmt"

	"github.com/nstratos/go-myanimelist/mal"
)

// demoClient has methods showcasing the usage of the different MyAnimeList API
// methods. It stores the first error it encounters so error checking only needs
// to be done once.
//
// This pattern is used for convenience and should not be used in concurrent
// code without guarding the error.
type demoClient struct {
	*mal.Client
	err error
}

func (c *demoClient) animeList(ctx context.Context) {
	if c.err != nil {
		return
	}
	anime, _, err := c.Anime.List(ctx, "galactic heroes",
		mal.Fields{"rank", "popularity", "start_season"},
		mal.Limit(5),
		mal.Offset(0),
	)
	if err != nil {
		c.err = err
		return
	}
	for _, a := range anime {
		fmt.Printf("ID: %5d, Rank: %5d, Popularity: %5d %s (%d)\n", a.ID, a.Rank, a.Popularity, a.Title, a.StartSeason.Year)
	}
}

func (c *demoClient) animeListForLoop(ctx context.Context) {
	if c.err != nil {
		return
	}
	offset := 0
	for {
		anime, resp, err := c.Anime.List(ctx, "kiseijuu",
			mal.Fields{"rank", "popularity", "start_season"},
			mal.Limit(100),
			mal.Offset(offset),
		)
		if err != nil {
			c.err = err
			return
		}
		for _, a := range anime {
			fmt.Printf("ID: %5d, Rank: %5d, Popularity: %5d %s (%d)\n", a.ID, a.Rank, a.Popularity, a.Title, a.StartSeason.Year)
		}
		fmt.Println("--------")
		fmt.Printf("Next offset: %d\n", resp.NextOffset)
		offset = resp.NextOffset
		if offset == 0 {
			break
		}
	}
}

func (c *demoClient) userAnimeList(ctx context.Context) {
	if c.err != nil {
		return
	}
	anime, _, err := c.User.AnimeList(ctx, "@me",
		mal.Fields{"list_status{start_date, end_date, priority, comments, tags}", "my_list_status{start_date, end_date, priority, comments, tags}"},
		mal.Limit(10),
		mal.Offset(0),
	)
	if err != nil {
		c.err = err
		return
	}
	for _, a := range anime {
		fmt.Printf("ID: %5d, Status: %10q, %s\n", a.Anime.ID, a.Status.Status, a.Anime.Title)
	}
}

func (c *demoClient) userMangaList(ctx context.Context) {
	if c.err != nil {
		return
	}
	manga, _, err := c.User.MangaList(ctx, "@me",
		mal.SortMangaListByListScore,
		//mal.Fields{"list_status{start_date, end_date, priority, comments, tags}", "my_list_status{start_date, end_date, priority, comments, tags}"},
		mal.Fields{"list_status{start_date, end_date, priority, comments, tags}"},
		mal.Limit(10),
		mal.Offset(0),
	)
	if err != nil {
		c.err = err
		return
	}
	for _, m := range manga {
		fmt.Printf("ID: %5d, Status: %10q, %s\n", m.Manga.ID, m.Status.Status, m.Manga.Title)
	}
}

func (c *demoClient) updateMyAnimeListStatus(ctx context.Context) {
	if c.err != nil {
		return
	}
	s, resp, err := c.Anime.UpdateMyListStatus(ctx, 820, mal.Score(8), mal.NumEpisodesWatched(4))
	if err != nil {
		c.err = err
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Printf("%#v\n", s)
}

func (c *demoClient) updateMyMangaListStatus(ctx context.Context) {
	if c.err != nil {
		return
	}
	s, resp, err := c.Manga.UpdateMyListStatus(ctx, 1, mal.MangaStatusReading, mal.Score(8), mal.NumChaptersRead(4))
	if err != nil {
		c.err = err
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Printf("%#v\n", s)
}

func (c *demoClient) deleteMyListItem(ctx context.Context) {
	if c.err != nil {
		return
	}
	resp, err := c.Anime.DeleteMyListItem(ctx, 820)
	if err != nil {
		c.err = err
		return
	}
	fmt.Println(resp.StatusCode)
}

func (c *demoClient) ranking(ctx context.Context) {
	if c.err != nil {
		return
	}
	rankings := []mal.AnimeRanking{
		mal.AnimeRankingAll,
		mal.AnimeRankingAiring,
		mal.AnimeRankingUpcoming,
		mal.AnimeRankingTV,
		mal.AnimeRankingOVA,
		mal.AnimeRankingMovie,
		mal.AnimeRankingSpecial,
		mal.AnimeRankingByPopularity,
		mal.AnimeRankingFavorite,
	}
	for _, r := range rankings {
		fmt.Println("Ranking:", r)
		anime, _, err := c.Anime.Ranking(ctx, r,
			mal.Fields{"rank", "popularity"},
		)
		if err != nil {
			c.err = err
			return
		}
		for _, a := range anime {
			fmt.Printf("Rank: %5d, Popularity: %5d %s\n", a.Rank, a.Popularity, a.Title)
		}
		fmt.Println("--------")
	}
}

func (c *demoClient) suggested(ctx context.Context) {
	if c.err != nil {
		return
	}
	anime, _, err := c.Anime.Suggested(ctx,
		mal.Limit(10),
		mal.Fields{"rank", "popularity"},
	)
	if err != nil {
		c.err = err
		return
	}
	for _, a := range anime {
		fmt.Printf("Rank: %5d, Popularity: %5d %s\n", a.Rank, a.Popularity, a.Title)
	}
}

func (c *demoClient) forumBoards(ctx context.Context) {
	if c.err != nil {
		return
	}
	forum, _, err := c.Forum.Boards(ctx)
	if err != nil {
		c.err = err
		return
	}
	for _, category := range forum.Categories {
		fmt.Printf("%s\n", category.Title)
		for _, b := range category.Boards {
			fmt.Printf("ID: %5d, Title: %5q %s\n", b.ID, b.Title, b.Description)
			for _, b := range b.Subboards {
				fmt.Printf("|-> ID: %5d, Title: %5q\n", b.ID, b.Title)
			}
		}
		fmt.Printf("-------\n")
	}
}

func (c *demoClient) forumTopics(ctx context.Context) {
	if c.err != nil {
		return
	}
	topics, _, err := c.Forum.Topics(ctx, mal.Query("kiseijuu"), mal.SortTopicsRecent)
	if err != nil {
		c.err = err
		return
	}
	for _, t := range topics {
		fmt.Printf("ID: %5d, Title: %5q created by %q\n", t.ID, t.Title, t.CreatedBy.Name)
	}
}
