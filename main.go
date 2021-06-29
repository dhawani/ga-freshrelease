package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	base           string
	token          string
	from           string
	regex          string
	str            string
	statusRegex    string
)
func init() {
	flag.StringVar(&base, "base", "", "Freshrelease Base URL")
	flag.StringVar(&token, "token", "", "Freshrelease API Token")
	flag.StringVar(&regex, "regex", "([a-zA-Z0-9]+-[0-9]+)", "Issue Key Regex")
	flag.StringVar(&statusRegex, "status-regex", "", "Regex to match status label")
	flag.StringVar(&from, "from", "string", "Find from predefined place (should be either 'string', 'branch', or 'commits')")
	flag.StringVar(&str, "str", "", "Provide a string to extract issue key from")
	flag.Parse()
}

func main() {
	if from == "string" {
		issueKey := regexp.MustCompile(regex).FindString(str)
		if issueKey != "" {
			fmt.Printf("Issue key matches.\n")
			if token != "" && base != "" {
				split := strings.Split(issueKey, "-")
				if len(split) == 2 {
					projectKey := split[0]
					issue, err := GetIssue(base, token, projectKey, issueKey)
					if err == nil {
						if !issue.Issue.Deleted {
							if statusRegex != "" {
								label := ""
								for _, s := range issue.Statuses {
									if s.ID == issue.Issue.StatusID {
										label = s.Label
										break
									}
								}
								if matched, err1 := regexp.MatchString(statusRegex, label); err1 == nil && matched {
									fmt.Printf("%s found with %s status.", issueKey, label)
								} else {
									fmt.Printf("Invalid status: %s for issue %s.", label, issueKey)
									os.Exit(1)
								}
							} else {
								fmt.Printf("%s found.", issueKey)
							}
						} else {
							fmt.Printf("%s is deleted.", issueKey)
							os.Exit(1)
						}
					} else {
						fmt.Printf("Error: %s\n", err.Error())
						os.Exit(1)
					}
				} else {
					fmt.Printf("Project key not found from %s\n", issueKey)
					os.Exit(1)
				}
			}
		} else {
			fmt.Printf("Issue key not found.\n")
			os.Exit(1)
		}
	} else {
		fmt.Printf("Unknown from: %s\n", from)
		os.Exit(1)
	}
}
