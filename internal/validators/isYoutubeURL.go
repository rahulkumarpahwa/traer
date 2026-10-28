package validators

import "regexp"

var ytRegex = regexp.MustCompile(`^(https?://)?(www\.|music\.)?(youtube\.com|youtu\.be)/.+$`)

func IsYouTubeURL(url string) bool {
	return ytRegex.MatchString(url)
}

