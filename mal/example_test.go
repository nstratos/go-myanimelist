package mal_test

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
)

//go:embed testdata/*.json
var testDataJSON embed.FS

// newStubServer creates a stub server which serves some premade responses. By
// contacting this server instead of the real API we can have runnable examples
// which always produce the same output.
func newStubServer() *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	serveStubFile := func(w io.Writer, filename string) error {
		stubResponses, err := fs.Sub(testDataJSON, "testdata")
		if err != nil {
			return err
		}
		f, err := stubResponses.Open(filename)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, f); err != nil {
			return err
		}
		return nil
	}

	serveStubHandler := func(filename string) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			malError := func(err string) string {
				return fmt.Sprintf(`{"message": "", "error":"%s"}`, err)
			}
			switch r.Method {
			case http.MethodDelete:
				w.WriteHeader(http.StatusOK)
			case http.MethodGet, http.MethodPatch:
				if err := serveStubFile(w, filename); err != nil {
					http.Error(w, malError("internal"), http.StatusInternalServerError)
				}
			default:
				http.Error(w, malError("not_allowed"), http.StatusMethodNotAllowed)
			}
		}
	}

	mux.HandleFunc("/anime", serveStubHandler("animeList.json"))
	mux.HandleFunc("/anime/967", serveStubHandler("animeDetails.json"))
	mux.HandleFunc("/anime/967/my_list_status", serveStubHandler("updateMyAnimeList.json"))
	mux.HandleFunc("/anime/ranking", serveStubHandler("animeRanking.json"))
	mux.HandleFunc("/anime/season/2020/fall", serveStubHandler("animeSeasonal.json"))
	mux.HandleFunc("/anime/suggestions", serveStubHandler("animeSuggested.json"))
	mux.HandleFunc("/manga", serveStubHandler("mangaList.json"))
	mux.HandleFunc("/manga/401", serveStubHandler("mangaDetails.json"))
	mux.HandleFunc("/manga/401/my_list_status", serveStubHandler("updateMyMangaList.json"))
	mux.HandleFunc("/manga/ranking", serveStubHandler("mangaRanking.json"))
	mux.HandleFunc("/users/@me", serveStubHandler("userMyInfo.json"))
	mux.HandleFunc("/users/@me/animelist", serveStubHandler("userAnimeList.json"))
	mux.HandleFunc("/users/@me/mangalist", serveStubHandler("userMangaList.json"))
	mux.HandleFunc("/forum/boards", serveStubHandler("forumBoards.json"))
	mux.HandleFunc("/forum/topics", serveStubHandler("forumTopics.json"))
	mux.HandleFunc("/forum/topic/1877721", serveStubHandler("forumTopicDetails.json"))

	return server
}
