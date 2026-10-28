package dto


type CreateJobRequest struct {
	URLs         []string `json:"urls"`
	ContentType  string   `json:"content_type"`
	ProcessType  string   `json:"process_type"`
	AudioQuality string   `json:"audio_quality,omitempty"`
	VideoQuality string   `json:"video_quality,omitempty"`
}
