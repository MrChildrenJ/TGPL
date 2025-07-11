/*
Exercise 4.12: The popular web comic xkcd has a JSON interface. For example, a request
to https://xkcd.com/571/info.0.json produces a detailed description of comic 571, one
of many favorites. Download each URL (once!) and build an offline index. Write a tool
xkcd that, using this index, prints the URL and transcript of each comic that matches
a search term provided on the command line.
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Comic struct {
	Num 		int    `json:"num"`
	Year 		string `json:"year"`
	Month 		string `json:"month"`
	Day 		string `json:"day"`
	Title 		string `json:"title"`
	SafeTitle 	string `json:"safe_title"`
	Transcript 	string `json:"transcript"`
	Alt 		string `json:"alt"`
	Img 		string `json:"img"`
	Link 		string `json:"link"`
}

type Index struct {
	Comics 	map[int]*Comic	`json:"comics"`
	LastNum int 			`json:"lastNum"`
	Updated time.Time 		`json:"updated"`
}

type SearchResult struct {
	Comic *Comic
	Score int
}

const (
	indexFile = "xkcd_index.json"		// saved json file
	baseURL   = "https://xkcd.com/"
	UserAgent = "xkcd-cli/1.0"
)

var client = http.Client{				// A custom client for more control over aspects like timeouts, 
	Timeout: 10 * time.Second,			// redirect policies, and connection pooling.
}

func fetchComic(num int) (*Comic, error) {
	var url string
	if num == 0 {
		url = baseURL + "info.0.json"
	} else {
		url = baseURL + fmt.Sprintf("%d/info.0.json", num)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Some websites block Go's default User-Agent "Go-http-client/1.1", so we need to modify it
	req.Header.Set("User-Agent", UserAgent)	

	// the most flexible method, allowing create a custom http.Request object and then execute it
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var comic Comic
	if err := json.NewDecoder(resp.Body).Decode(&comic); err != nil {
		return nil, err
	}
	return &comic, nil
}

func loadIndex() (*Index, error) {
	if _, err := os.Stat(indexFile); os.IsNotExist(err) {
		return &Index{
			Comics: make(map[int]*Comic),
			LastNum: 0,
			Updated: time.Time{},
		}, nil
	}

	data, err := os.ReadFile(indexFile)
	if err != nil {
		return nil, err
	}

	var index Index
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, err
	}

	if index.Comics == nil {
		index.Comics = make(map[int]*Comic)
	}
	
	return &index, nil
}

func saveIndex(index *Index) error {
	dir := filepath.Dir(indexFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(indexFile, data, 0644)
}

func updateIndex() error {
	fmt.Println("Loading existing index...")
	index, err := loadIndex()
	if err != nil {
		return fmt.Errorf("failed to load index: %v", err)
	}

	fmt.Println("Fetching latest comic to determine range...")
	latest, err := fetchComic(0)
	if err != nil {
		return fmt.Errorf("failed to fetch latest comic: %v", err)
	}

	fmt.Printf("Latest comic: #%d - %s\n", latest.Num, latest.Title)

	// 确定需要下载的范围
	startNum := 1
	if index.LastNum > 0 {
		startNum = index.LastNum + 1
	}

	totalToFetch := 0
	for i := startNum; i <= latest.Num; i++ {
		if _, exists := index.Comics[i]; !exists {
			totalToFetch++
		}
	}

	if totalToFetch == 0 {
		fmt.Println("Index is already up to date!")
		return nil
	}

	fmt.Printf("Need to fetch %d comics...\n", totalToFetch)

	// 下载缺失的漫画
	fetched := 0
	for i := startNum; i <= latest.Num; i++ {
		if _, exists := index.Comics[i]; exists {
			continue
		}

		fmt.Printf("Fetching comic #%d... (%d/%d)\n", i, fetched+1, totalToFetch)

		comic, err := fetchComic(i)
		if err != nil {
			fmt.Printf("Warning: failed to fetch comic #%d: %v\n", i, err)
			continue
		}

		if comic == nil {
			fmt.Printf("Warning: comic #%d does not exist\n", i)
			continue
		}

		index.Comics[i] = comic
		fetched++

		// 添加小延迟以避免过于频繁的请求
		time.Sleep(100 * time.Millisecond)

		// 每50个漫画保存一次，防止数据丢失
		if fetched%50 == 0 {
			fmt.Printf("Saving progress... (%d/%d)\n", fetched, totalToFetch)
			index.LastNum = i
			index.Updated = time.Now()
			if err := saveIndex(index); err != nil {
				fmt.Printf("Warning: failed to save progress: %v\n", err)
			}
		}
	}

	index.LastNum = latest.Num
	index.Updated = time.Now()

	fmt.Printf("Saving index with %d comics...\n", len(index.Comics))
	if err := saveIndex(index); err != nil {
		return fmt.Errorf("failed to save index: %v", err)
	}

	fmt.Printf("Successfully updated index! Fetched %d new comics.\n", fetched)
	return nil
}

func search(query string) ([]*SearchResult, error) {
	index, err := loadIndex()
	if err != nil {
		return nil, err
	}

	if len(index.Comics) == 0 {
		return nil, fmt.Errorf("index is empty. Please run 'update' first")
	}

	query = strings.ToLower(query)
	terms := strings.Fields(query)

	var results []*SearchResult

	for _, comic := range index.Comics {
		score := calculateScore(comic, terms)
		if score > 0 {
			results = append(results, &SearchResult{
				Comic: comic,
				Score: score,
			})
		}
	}

	// 按分数排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}

func calculateScore(comic *Comic, terms []string) int {
	score := 0
	
	// 将所有文本合并并转为小写
	allText := strings.ToLower(fmt.Sprintf("%s %s %s %s", 
		comic.Title, comic.SafeTitle, comic.Alt, comic.Transcript))

	for _, term := range terms {
		// 标题匹配得分更高
		if strings.Contains(strings.ToLower(comic.Title), term) {
			score += 10
		}
		if strings.Contains(strings.ToLower(comic.SafeTitle), term) {
			score += 8
		}
		// Alt文本匹配
		if strings.Contains(strings.ToLower(comic.Alt), term) {
			score += 5
		}
		// 转录文本匹配
		if strings.Contains(strings.ToLower(comic.Transcript), term) {
			score += 3
		}
		// 整体文本匹配
		if strings.Contains(allText, term) {
			score += 1
		}
	}
	return score
}

func displayComic(comic *Comic) {
	fmt.Printf("┌─ XKCD #%d ─────────────────────────────────────\n", comic.Num)
	fmt.Printf("│ Title: %s\n", comic.Title)
	fmt.Printf("│ Date:  %s-%s-%s\n", comic.Year, comic.Month, comic.Day)
	fmt.Printf("│ URL:   %s/%d/\n", baseURL, comic.Num)
	fmt.Printf("│ Image: %s\n", comic.Img)
	if comic.Link != "" {
		fmt.Printf("│ Link:  %s\n", comic.Link)
	}
	fmt.Printf("├─ Alt Text ──────────────────────────────────────\n")
	fmt.Printf("│ %s\n", wrapText(comic.Alt, 60))
	if comic.Transcript != "" {
		fmt.Printf("├─ Transcript ────────────────────────────────────\n")
		fmt.Printf("│ %s\n", wrapText(comic.Transcript, 60))
	}
	fmt.Printf("└─────────────────────────────────────────────────\n")
}

func wrapText(text string, width int) string {
	if len(text) <= width {
		return text
	}

	var lines []string
	words := strings.Fields(text)
	currentLine := ""

	for _, word := range words {
		if len(currentLine)+len(word)+1 <= width {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n│ ")
}

func showStats() error {
	index, err := loadIndex()
	if err != nil {
		return err
	}

	fmt.Printf("XKCD Index Statistics\n")
	fmt.Printf("═══════════════════════\n")
	fmt.Printf("Total comics indexed: %d\n", len(index.Comics))
	fmt.Printf("Last comic number:    %d\n", index.LastNum)
	fmt.Printf("Last updated:         %s\n", index.Updated.Format("2006-01-02 15:04:05"))
	
	if len(index.Comics) > 0 {
		fmt.Printf("\nSample comics:\n")
		// 显示前5个和后5个漫画
		var nums []int
		for num := range index.Comics {
			nums = append(nums, num)
		}
		sort.Ints(nums)

		count := 0
		for _, num := range nums {
			if count >= 5 {
				break
			}
			comic := index.Comics[num]
			fmt.Printf("  #%d: %s\n", num, comic.Title)
			count++
		}

		if len(nums) > 10 {
			fmt.Printf("  ...\n")
			for i := len(nums) - 5; i < len(nums); i++ {
				num := nums[i]
				comic := index.Comics[num]
				fmt.Printf("  #%d: %s\n", num, comic.Title)
			}
		}
	}
	return nil
}

func showRandom() error {
	index, err := loadIndex()
	if err != nil {
		return err
	}

	if len(index.Comics) == 0 {
		return fmt.Errorf("index is empty. Please run 'update' first")
	}

	// 获取随机漫画
	var nums []int
	for num := range index.Comics {
		nums = append(nums, num)
	}

	randomIndex := time.Now().UnixNano() % int64(len(nums))
	randomNum := nums[randomIndex]
	comic := index.Comics[randomNum]

	fmt.Println("Random XKCD Comic:")
	displayComic(comic)

	return nil
}

func showComic(numStr string) error {
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return fmt.Errorf("invalid comic number: %s", numStr)
	}

	index, err := loadIndex()
	if err != nil {
		return err
	}

	comic, exists := index.Comics[num]
	if !exists {
		return fmt.Errorf("comic #%d not found in index", num)
	}

	displayComic(comic)
	return nil
}

func printUsage() {
	fmt.Println("XKCD Offline Tool")
	fmt.Println("═════════════════")
	fmt.Println("Usage:")
	fmt.Println("  go run xkcd.go <command> [arguments]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  update                    - Download and update the comic index")
	fmt.Println("  search <keywords>         - Search comics by keywords")
	fmt.Println("  show <number>            - Show specific comic by number")
	fmt.Println("  random                   - Show a random comic")
	fmt.Println("  stats                    - Show index statistics")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run xkcd.go update")
	fmt.Println("  go run xkcd.go search \"programming python\"")
	fmt.Println("  go run xkcd.go show 353")
	fmt.Println("  go run xkcd.go random")
	fmt.Println("  go run xkcd.go stats")
}



func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "update":
		if err := updateIndex(); err != nil {
			log.Fatalf("Update failed: %v", err)
		}

	case "search":
		if len(os.Args) < 3 {
			log.Fatal("Search query is required")
		}
		query := strings.Join(os.Args[2:], " ")
		
		results, err := search(query)
		if err != nil {
			log.Fatalf("Search failed: %v", err)
		}

		if len(results) == 0 {
			fmt.Printf("No comics found matching '%s'\n", query)
			return
		}

		fmt.Printf("Found %d comics matching '%s':\n\n", len(results), query)
		
		maxResults := 10
		if len(results) < maxResults {
			maxResults = len(results)
		}

		for i := 0; i < maxResults; i++ {
			result := results[i]
			fmt.Printf("%d. #%d: %s (score: %d)\n", 
				i+1, result.Comic.Num, result.Comic.Title, result.Score)
			fmt.Printf("   URL: %s/%d/\n", baseURL, result.Comic.Num)
			fmt.Printf("   %s\n\n", result.Comic.Alt)
		}

		if len(results) > maxResults {
			fmt.Printf("... and %d more results\n", len(results)-maxResults)
		}

	case "show":
		if len(os.Args) < 3 {
			log.Fatal("Comic number is required")
		}
		if err := showComic(os.Args[2]); err != nil {
			log.Fatalf("Show failed: %v", err)
		}

	case "random":
		if err := showRandom(); err != nil {
			log.Fatalf("Random failed: %v", err)
		}

	case "stats":
		if err := showStats(); err != nil {
			log.Fatalf("Stats failed: %v", err)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}