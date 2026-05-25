package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const configFile = "config.yaml"

func main() {
	// 找到 config.yaml（先看当前目录，再看可执行文件同目录）
	cfgPath := findConfig()
	cfg, err := loadConfig(cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// 首次运行 — 生成配置
	if cfg == nil {
		cfg2 := defaultConfig()
		if err := saveConfig(cfgPath, &cfg2); err != nil {
			log.Fatalf("create config: %v", err)
		}
		fmt.Printf(`
╔══════════════════════════════════════════════════╗
║          cc-gateway — 首次运行                    ║
╠══════════════════════════════════════════════════╣
║                                                   ║
║  配置文件已生成：%s
║                                                   ║
║  请编辑此文件，填入你的 Command Code API Key：     ║
║    commandcode:                                    ║
║      api_key: "user_xxx"                           ║
║                                                   ║
║  API Key 获取：https://commandcode.ai/studio       ║
║                                                   ║
║  填完后重新运行即可启动。                           ║
║  对外 API Key（自动生成）：%s
║                                                   ║
╚══════════════════════════════════════════════════╝
`, cfgPath, cfg2.APIKey)
		os.Exit(0)
	}

	// 检查是否填了 CC API Key
	if cfg.CommandCode.APIKey == "" {
		fmt.Printf("❌ 请先在 %s 中填入 commandcode.api_key\n", cfgPath)
		fmt.Println("   获取：https://commandcode.ai/studio")
		os.Exit(1)
	}

	if cfg.Port == 0 {
		cfg.Port = 11434
	}
	if cfg.CommandCode.BaseURL == "" {
		cfg.CommandCode.BaseURL = "https://api.commandcode.ai"
	}

	cc := NewCCClient(cfg.CommandCode.APIKey, cfg.CommandCode.BaseURL)

	if err := runServer(cc, cfg); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func findConfig() string {
	// 1. 当前目录
	if _, err := os.Stat(configFile); err == nil {
		return configFile
	}
	// 2. 可执行文件同目录
	exe, err := os.Executable()
	if err == nil {
		p := filepath.Join(filepath.Dir(exe), configFile)
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	// 3. 默认当前目录（用于首次生成）
	return configFile
}
