//go:build darwin

package desktopruntime

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DefaultRuntimeRoot 返回 macOS 桌面版固定的当前用户运行目录。
// LaunchAgent plist 位于 App Bundle 内，不能在构建时写入具体用户名，因此由 Core 自己解析。
func consoleUser() string {
	if output, err := exec.Command("/usr/bin/id", "-un").Output(); err == nil {
		return strings.TrimSpace(string(output))
	}
	return ""
}

func darwinUserHome() string {
	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		return home
	}

	// SMAppService launch agents intentionally start with a sparse environment and
	// may not receive HOME. Ask Directory Service for the console user's home so
	// the bundled Core and Tunnel can still resolve the same per-user runtime root.
	if output, outputErr := exec.Command("/usr/bin/dscl", ".", "-read", "/Users/"+consoleUser(), "NFSHomeDirectory").Output(); outputErr == nil {
		for _, line := range strings.Split(string(output), "\n") {
			if value, found := strings.CutPrefix(strings.TrimSpace(line), "NFSHomeDirectory:"); found {
				return strings.TrimSpace(value)
			}
		}
	}
	return ""
}

func DefaultRuntimeRoot() string {
	home := darwinUserHome()
	if home == "" {
		return ""
	}
	return filepath.Join(home, "Library", "Application Support", "AgentDock")
}
