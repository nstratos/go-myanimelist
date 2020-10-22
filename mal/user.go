package mal

import "time"

// User represents a MyAnimeList user.
type User struct {
	ID              int             `json:"id,omitempty"`
	Name            string          `json:"name,omitempty"`
	Location        string          `json:"location,omitempty"`
	JoinedAt        time.Time       `json:"joined_at,omitempty"`
	AnimeStatistics AnimeStatistics `json:"anime_statistics,omitempty"`
}

// AnimeStatistics about the user.
type AnimeStatistics struct {
	NumItemsWatching    int     `json:"num_items_watching,omitempty"`
	NumItemsCompleted   int     `json:"num_items_completed,omitempty"`
	NumItemsOnHold      int     `json:"num_items_on_hold,omitempty"`
	NumItemsDropped     int     `json:"num_items_dropped,omitempty"`
	NumItemsPlanToWatch int     `json:"num_items_plan_to_watch,omitempty"`
	NumItems            int     `json:"num_items,omitempty"`
	NumDaysWatched      float64 `json:"num_days_watched,omitempty"`
	NumDaysWatching     float64 `json:"num_days_watching,omitempty"`
	NumDaysCompleted    float64 `json:"num_days_completed,omitempty"`
	NumDaysOnHold       int     `json:"num_days_on_hold,omitempty"`
	NumDaysDropped      int     `json:"num_days_dropped,omitempty"`
	NumDays             float64 `json:"num_days,omitempty"`
	NumEpisodes         int     `json:"num_episodes,omitempty"`
	NumTimesRewatched   int     `json:"num_times_rewatched,omitempty"`
	MeanScore           float64 `json:"mean_score,omitempty"`
}
