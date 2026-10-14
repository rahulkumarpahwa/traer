package bootstrap

func getLinuxInstallCommand(tool string) string {
	switch tool {
	case "ffmpeg":
		return "sudo apt update && sudo apt install -y ffmpeg"
	case "yt-dlp":
		return "sudo apt update && sudo apt install -y yt-dlp"
	default:
		return ""
	}
}