package bootstrap

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func installTool(tool string) error {
	switch runtime.GOOS {

	case "linux":
		return installLinux(tool)

	case "darwin":
		return installMac(tool)

	case "windows":
		return installWindows(tool)

	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func installLinux(tool string) error {
	cmd := exec.Command("bash", "-c", getLinuxInstallCommand(tool))
	log.Println("installing:", tool)
	return cmd.Run()
}

func installMac(tool string) error {
	cmd := exec.Command("brew", "install", tool)
	log.Println("installing:", tool)
	return cmd.Run()
}

func installWindows(tool string) error {
	cmd := exec.Command("winget", "install", tool)
	log.Println("installing:", tool)
	return cmd.Run()
}