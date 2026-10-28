package types

type Playlist struct {
	ProcessType ProcessType `json:"-"`
	PlaylistUrl string      `json:"-"`
	Entries     []Entry     `json:"entries"`
	Error       []string    `json:"-"`
}

type Entry struct {
	ID  string `json:"id"`
	Url string `json:"url"`
}
