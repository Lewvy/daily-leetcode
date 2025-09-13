package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const leetCodeAPIURL = "https://leetcode.com/graphql"

func sendNotification(summary, body string) error {
	cmd := exec.Command("notify-send", summary, body)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to send notification: %w. Make sure 'notify-send' is installed", err)
	}
	return nil
}

func checkFriendsAndNotify() {
	fmt.Printf("[%s] Running daily LeetCode check...\n", time.Now().Format("2006-01-02 15:04:05"))

	friends, err := readFriendsFromFile("friends.txt")
	if err != nil {
		fmt.Printf("Error reading friends file: %v\n", err)
		sendNotification("LeetCode Checker Error", "Could not read friends.txt file.")
		return
	}

	var solvedList []string
	var notSolvedList []string

	for _, username := range friends {
		solved, err := hasUserSolvedProblemToday(username)
		if err != nil {
			fmt.Printf("Could not check status for %s: %v\n", username, err)
			continue
		}

		if solved {
			solvedList = append(solvedList, username)
		} else {
			notSolvedList = append(notSolvedList, username)
		}
	}

	var bodyBuilder strings.Builder
	if len(solvedList) > 0 {
		bodyBuilder.WriteString("✅ Solved:\n")
		for _, user := range solvedList {
			bodyBuilder.WriteString(fmt.Sprintf("- %s\n", user))
		}
	}
	if len(notSolvedList) > 0 {
		bodyBuilder.WriteString("\n❌ Not Solved Yet:\n")
		for _, user := range notSolvedList {
			bodyBuilder.WriteString(fmt.Sprintf("- %s\n", user))
		}
	}

	summary := fmt.Sprintf("LeetCode Status: %d/%d Solved", len(solvedList), len(friends))
	fmt.Println("Check complete. Sending notification...")
	err = sendNotification(summary, bodyBuilder.String())
	if err != nil {
		fmt.Println(err)
	}
}

func readFriendsFromFile(filename string) ([]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", filename, err)
	}

	lines := strings.Split(string(content), "\n")
	var friends []string

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			friends = append(friends, trimmedLine)
		}
	}
	return friends, nil
}

type GraphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

type UserSubmissionsResponse struct {
	Data struct {
		RecentAcSubmissionList []struct {
			TitleSlug string `json:"titleSlug"`
			Timestamp string `json:"timestamp"`
		} `json:"recentAcSubmissionList"`
	} `json:"data"`
}

func hasUserSolvedProblemToday(username string) (bool, error) {
	query := `
        query recentAcSubmissionList($username: String!, $limit: Int!) {
            recentAcSubmissionList(username: $username, limit: $limit) {
                titleSlug
                timestamp
            }
        }`

	variables := map[string]any{
		"username": username,
		"limit":    1,
	}

	reqBody, err := json.Marshal(GraphQLRequest{Query: query, Variables: variables})
	if err != nil {
		return false, fmt.Errorf("error marshalling request: %w", err)
	}

	resp, err := http.Post(leetCodeAPIURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return false, fmt.Errorf("error making http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("api returned non-200 status for %s: %s | Body: %s", username, resp.Status, string(bodyBytes))
	}

	var userResp UserSubmissionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return false, fmt.Errorf("error decoding response body: %w", err)
	}

	if len(userResp.Data.RecentAcSubmissionList) == 0 {
		return false, nil
	}

	lastSubmission := userResp.Data.RecentAcSubmissionList[0]

	now := time.Now()
	year, month, day := now.Date()

	unixTime, err := strconv.ParseInt(lastSubmission.Timestamp, 10, 64)
	if err != nil {
		return false, fmt.Errorf("could not parse timestamp: %w", err)
	}

	submissionTime := time.Unix(unixTime, 0)
	sYear, sMonth, sDay := submissionTime.Date()

	if sYear == year && sMonth == month && sDay == day {
		return true, nil
	}

	return false, nil
}

func main() {
	checkFriendsAndNotify()
}
