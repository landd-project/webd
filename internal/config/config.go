package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	SocketPath string
	KeepHistory bool
	HistoryLimit int
}

func GetBaseDir() (string, error) {
	dir, err := os.UserConfigDir();
	if err != nil {
		return "", err;
	}

	return dir, nil;
}

func GetConfig() (Config, error) {
	var config = Config{
		SocketPath: "/tmp/webd.internal.sock",
	};

	c, err := parseConfig();
	if err != nil {
		return config, err;
	}

	if config.SocketPath != c.SocketPath && c.SocketPath != ""{
		config.SocketPath = c.SocketPath;
	}

	return config, nil;
}

func parseConfig() (Config, error) {
	var c Config;

	dir, err := GetBaseDir();
	if err != nil {
		return c, err;
	}
	path := filepath.Join(dir, "config");

	content, err := ReadFileToString(path);
	if err != nil {
		return c, err;
	}

	lines := strings.Split(content, "\n");
	for i,line := range lines {
		line = strings.TrimSpace(line);
		if line == "" || strings.HasPrefix(line, "#") {
			continue;
		}
		parts := strings.SplitN(line, ":", 2);
		if len(parts) != 2 {
			return c, fmt.Errorf("failed to parse line: %v", i);
		}
		key := parts[0];
		value := strings.TrimSpace(parts[1]);

		switch key {
		case "socket-path":
			c.SocketPath = value;
		default:
			return c, fmt.Errorf("invalid key on line: %v", i);
		}
	}

	return c, nil;
}

func ReadFileToString(path string) (string, error) {
	dat, err := os.ReadFile(path);
	if err != nil {
		fmt.Errorf("failed to read file: `%v` %v", path, err);
	}

	return string(dat), nil;
}
