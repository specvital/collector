package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/hibiken/asynq"
	handler "github.com/specvital/collector/internal/handler/queue"
)

func main() {
	redisURL := flag.String("redis", os.Getenv("REDIS_URL"), "Redis URL")
	flag.Parse()

	if flag.NArg() < 1 {
		printUsage()
		os.Exit(1)
	}

	if *redisURL == "" {
		fmt.Fprintln(os.Stderr, "Error: Redis URL is required (use -redis flag or set REDIS_URL)")
		os.Exit(1)
	}

	owner, repo, err := ParseGitHubURL(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := enqueue(*redisURL, owner, repo); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to enqueue task: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: enqueue [flags] <github-url>")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Arguments:")
	fmt.Fprintln(os.Stderr, "  <github-url>  GitHub repository URL (e.g., github.com/owner/repo)")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintln(os.Stderr, "  enqueue github.com/octocat/Hello-World")
	fmt.Fprintln(os.Stderr, "  enqueue -redis redis://localhost:6379 github.com/owner/repo")
	fmt.Fprintln(os.Stderr, "  enqueue https://github.com/owner/repo.git")
}

func enqueue(redisURL, owner, repo string) error {
	opt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		return fmt.Errorf("parse redis URI: %w", err)
	}

	client := asynq.NewClient(opt)
	defer client.Close()

	payload, err := json.Marshal(handler.AnalyzePayload{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	task := asynq.NewTask(handler.TypeAnalyze, payload)
	info, err := client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("enqueue task: %w", err)
	}

	slog.Info("task enqueued",
		"id", info.ID,
		"queue", info.Queue,
		"owner", owner,
		"repo", repo,
	)
	return nil
}
