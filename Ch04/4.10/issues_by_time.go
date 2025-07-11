package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopl.io/ch4/github"
)

func main() {
	result, err := github.SearchIssues(os.Args[1:])

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d issues:\n", result.TotalCount)

	var lessThanMonth, lessThanYear, overYear []*github.Issue

	now := time.Now()
	oneMonthAgo := now.AddDate(0, -1, 0)
	oneYearAgo := now.AddDate(-1, 0, 0)
	
	for _, item := range result.Items {
		if item.CreatedAt.After(oneMonthAgo) {
			lessThanMonth = append(lessThanMonth, item)
		} else if item.CreatedAt.After(oneYearAgo) {
			lessThanYear = append(lessThanYear, item)
		} else {
			overYear = append(overYear, item)
		}
	}

	fmt.Printf("There are %d less than a month:\n", len(lessThanMonth))
	for _, item := range lessThanMonth {
		fmt.Printf("#%-5d %9.9s %.55s %d-%02d-%04d\n", item.Number, item.User.Login, item.Title, item.CreatedAt.Year(), item.CreatedAt.Month(), item.CreatedAt.Day())
	}

	fmt.Printf("\nThere are %d less than a year:\n", len(lessThanYear))
	for _, item := range lessThanYear {
		fmt.Printf("#%-5d %9.9s %.55s %d-%02d-%04d\n", item.Number, item.User.Login, item.Title, item.CreatedAt.Year(), item.CreatedAt.Month(), item.CreatedAt.Day())
	}

	fmt.Printf("\nThere are %d over a year:\n", len(overYear))
	for _, item := range overYear {
		fmt.Printf("#%-5d %9.9s %.55s %d-%02d-%04d\n", item.Number, item.User.Login, item.Title, item.CreatedAt.Year(), item.CreatedAt.Month(), item.CreatedAt.Day())
	}
}