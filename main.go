package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"sub2sing/generator"
	"sub2sing/parser"
)

func main() {
	urlFlag := flag.String("url", "", "v2ray 订阅链接 URL")
	fileFlag := flag.String("file", "", "本地订阅文件路径")
	outputFlag := flag.String("o", "", "输出文件路径（默认 stdout）")
	flag.Parse()

	if *urlFlag == "" && *fileFlag == "" {
		fmt.Fprintln(os.Stderr, "用法: sub2sing -url <订阅链接> -o config.json")
		fmt.Fprintln(os.Stderr, "      sub2sing -file <订阅文件> -o config.json")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// 1. 获取订阅数据
	var data []byte
	var err error

	if *urlFlag != "" {
		data, err = fetchURL(*urlFlag)
	} else {
		data, err = os.ReadFile(*fileFlag)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取订阅失败: %v\n", err)
		os.Exit(1)
	}

	// 2. 解析订阅
	nodes, err := parser.ParseSubscription(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "解析订阅失败: %v\n", err)
		os.Exit(1)
	}

	if len(nodes) == 0 {
		fmt.Fprintln(os.Stderr, "未解析到任何节点，请检查订阅格式")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "解析到 %d 个节点\n", len(nodes))
	for i, n := range nodes {
		fmt.Fprintf(os.Stderr, "  [%d] %s (%s) - %s:%d\n", i+1, n.Tag, n.Type, n.Server, n.ServerPort)
	}

	// 3. 生成配置
	config := generator.Generate(nodes)
	jsonData, err := config.ToJSON()
	if err != nil {
		fmt.Fprintf(os.Stderr, "生成配置失败: %v\n", err)
		os.Exit(1)
	}

	// 4. 输出
	if *outputFlag != "" {
		err = os.WriteFile(*outputFlag, jsonData, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "写入文件失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "配置已写入 %s\n", *outputFlag)
	} else {
		fmt.Println(string(jsonData))
	}
}

func fetchURL(url string) ([]byte, error) {
	// 清理 URL 中的换行和空格
	url = strings.TrimSpace(url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP 状态码: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	return data, nil
}
