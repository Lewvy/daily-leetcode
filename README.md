# LeetCode Daily Submission Checker

A simple command-line tool written in Go that checks if your friends (or a list of specified LeetCode users) have solved at least one problem on LeetCode today. It reads a list of usernames from a local file and queries the LeetCode GraphQL API to check their recent submission status.

-----

## Features

  - **Daily Tracking**: Quickly see which users from your list have made a successful submission on the current date.
  - **Simple Configuration**: Usernames are managed in a straightforward `friends.txt` file.
  - **Direct API Interaction**: Communicates directly with the official LeetCode GraphQL API.
  - **Lightweight**: Built with only the Go standard library, requiring no external dependencies.

-----

## How It Works

The application performs the following steps:

1.  **Reads Usernames**: It reads a list of LeetCode usernames from a `friends.txt` file located in the same directory.
2.  **Iterates and Queries**: For each username, it constructs a GraphQL query to fetch the user's most recent accepted submission.
3.  **API Request**: It sends an HTTP POST request to the LeetCode API endpoint (`https://leetcode.com/graphql`).
4.  **Parses Response**: It parses the JSON response to extract the timestamp of the last submission.
5.  **Compares Dates**: The submission timestamp is converted to a date and compared with the current system date.
6.  **Prints Status**: It prints a status message to the console for each user, indicating whether they have solved a problem today or not.

-----

## Prerequisites

To run this project, you need to have **Go** installed on your system. You can download it from the [official Go website](https://go.dev/dl/).

-----

## Installation & Usage

Follow these steps to get the checker up and running:

1.  **Get the code**:
    Save the Go code into a file named `main.go`.

2.  **Create the friends list**:
    In the same directory as `main.go`, create a file named `friends.txt`. Add the LeetCode usernames you want to track, with each username on a new line.

    **Example `friends.txt`:**

    ```text
    leetcode_username1
    another_user
    friend_coder
    ```

3.  **Run the application**:
    Open your terminal, navigate to the project directory, and run the following command:

    ```sh
    go run main.go
    ```

-----

## Example Output

After running the command, you will see an output similar to this in your terminal:

```
Checking for any LeetCode submission today...
leetcode_username1 solved a problem today!
another_user has not solved a problem yet today.
friend_coder solved a problem today!
```

If the application encounters an issue (e.g., an invalid username or an API error), it will print a descriptive error message for that specific user.
