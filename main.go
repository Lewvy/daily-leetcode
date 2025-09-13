package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const leetCodeAPIURL = "https://leetcode.com/graphql"

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

	friends, err := readFriendsFromFile("friends.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Please make sure a 'friends.txt' file exists in the same directory.")
		return
	}

	fmt.Println("Checking for any LeetCode submission today...")

	for _, username := range friends {
		solved, err := hasUserSolvedProblemToday(username)
		if err != nil {
			fmt.Printf("Could not check status for %s: %v\n", username, err)
			continue
		}

		if solved {
			fmt.Printf("%s solved a problem today!\n", username)
		} else {
			fmt.Printf("%s has not solved a problem yet today.\n", username)
		}
	}
}
