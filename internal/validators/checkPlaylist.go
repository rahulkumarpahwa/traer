package validators

import (
	"regexp"
	"strings"
)

// PlaylistType represents the category of a YouTube playlist.
type PlaylistType string

const (
	PlaylistTypeUserCreated   PlaylistType = "user_created"
	PlaylistTypeAllUploads    PlaylistType = "all_uploads"
	PlaylistTypeVideosOnly    PlaylistType = "videos_only"
	PlaylistTypeShortsOnly    PlaylistType = "shorts_only"
	PlaylistTypeAutoMix       PlaylistType = "auto_mix"
	PlaylistTypeOfficialAlbum PlaylistType = "official_album"
	PlaylistTypeWatchLater    PlaylistType = "watch_later"
	PlaylistTypeLikedMusic    PlaylistType = "liked_music"
)

var (
	// Matches the list= query parameter value.
	listParamRe = regexp.MustCompile(`[?&]list=([A-Za-z0-9_\-]+)`)
	// Matches /playlist?list= path.
	playlistPathRe = regexp.MustCompile(`/playlist[?&]list=([A-Za-z0-9_\-]+)`)
)

// IsYouTubePlaylist checks whether the given URL is a YouTube playlist link.
// It returns (true, playlistType) if valid, or (false, "") otherwise.
func IsYouTubePlaylist(url string) (bool, PlaylistType) {
	// Extract the playlist ID from the URL.
	id := extractPlaylistID(url)
	if id == "" {
		return false, ""
	}

	// Classify by prefix.
	switch {
	case strings.HasPrefix(id, "PL"):
		return true, PlaylistTypeUserCreated
	case strings.HasPrefix(id, "UULF"):
		return true, PlaylistTypeVideosOnly
	case strings.HasPrefix(id, "UUSH"):
		return true, PlaylistTypeShortsOnly
	case strings.HasPrefix(id, "UU"):
		return true, PlaylistTypeAllUploads
	case strings.HasPrefix(id, "RD"):
		return true, PlaylistTypeAutoMix
	case strings.HasPrefix(id, "OLAK5uy_"):
		return true, PlaylistTypeOfficialAlbum
	case strings.HasPrefix(id, "WL"):
		return true, PlaylistTypeWatchLater
	case strings.HasPrefix(id, "LM"):
		return true, PlaylistTypeLikedMusic
	default:
		return false, ""
	}
}

// extractPlaylistID pulls the playlist ID from a YouTube URL.
// Supports both /playlist?list=... and /watch?v=...&list=...
func extractPlaylistID(url string) string {
	// Try /playlist?list= first.
	if m := playlistPathRe.FindStringSubmatch(url); m != nil {
		return m[1]
	}
	// Try any list= param (e.g. on /watch URLs).
	if m := listParamRe.FindStringSubmatch(url); m != nil {
		return m[1]
	}
	return ""
}   