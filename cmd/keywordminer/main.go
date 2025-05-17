package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xshoji/go-keywordminer/config"
	"github.com/xshoji/go-keywordminer/internal/language/english"
	"github.com/xshoji/go-keywordminer/pkg/analyzer"
)

const (
	UsageRequiredPrefix = "\u001B[33m" + "(REQ)" + "\u001B[0m "
	UsageDummy          = "########"
	TimeFormat          = "2006-01-02 15:04:05.0000 [MST]"
)

var (
	// Command options ( the -h, --help option is defined by default in the flag package )
	commandDescription     = "A tool for extracting and analyzing keywords from web pages. \n  Fetches titles, meta tags, and identifies top keywords with their relevance scores."
	commandOptionMaxLength = 0
	optionUrl              = defineFlagValue("u", "url" /*    */, UsageRequiredPrefix+"URL" /*   */, "").(*string)
)

func init() {
	formatUsage(commandDescription, &commandOptionMaxLength, new(bytes.Buffer))
}

// Build:
// $ GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -trimpath ./cmd/keywordminer
func main() {
	flag.Parse()
	if *optionUrl == "" {
		flag.Usage()
		os.Exit(0)
	}

	cfg := config.DefaultConfig()
	anlz, err := analyzer.NewAnalyzer(*optionUrl, cfg)
	if err != nil {
		handleError(err, "NewAnalyzer")
		os.Exit(1)
	}

	title, err := anlz.FetchTitle()
	handleError(err, "FetchTitle")
	if err == nil && title != "" {
		fmt.Println("[Title]")
		fmt.Println(title)
	} else {
		fmt.Println("[Title] None")
	}

	meta, err := anlz.FetchMetaTags()
	handleError(err, "FetchMetaTags(meta)")
	if err == nil && len(meta) > 0 {
		fmt.Println("\n[Meta Tags]")
		for k, v := range meta {
			fmt.Printf("%s: %s\n", k, v)
		}
	} else {
		fmt.Println("\n[Meta Tags] None")
	}

	// ストップワード・正規化関数はConfigから取得
	stopWords := cfg.EnglishStopWords
	pluralSingularMap := cfg.PluralSingularMap
	invariantWords := cfg.InvariantWords
	normalize := func(word string) string {
		return english.NormalizeEnglishKeyword(word, pluralSingularMap, invariantWords)
	}

	keywordsWithScores, kerr := anlz.GetTopKeywords(20, stopWords, normalize)
	if kerr != nil {
		handleError(kerr, "GetTopKeywords")
	}
	fmt.Printf("\n[Top Keywords]\n")
	if len(keywordsWithScores) > 0 {
		for _, kws := range keywordsWithScores {
			fmt.Printf("%s (score: %d), ", kws.Keyword, kws.Score)
		}
		fmt.Println()
	} else {
		fmt.Println("No keywords found")
	}
}

// =======================================
// Common Utils
// =======================================

func handleError(err error, prefixErrMessage string) {
	if err != nil {
		fmt.Printf("%s [ERROR %s]: %v\n", time.Now().Format(TimeFormat), prefixErrMessage, err)
	}
}

// Helper function for flag
func defineFlagValue(short, long, description string, defaultValue any) (f any) {
	flagUsage := short + UsageDummy + description
	switch v := defaultValue.(type) {
	case string:
		f = flag.String(short, "", UsageDummy)
		flag.StringVar(f.(*string), long, v, flagUsage)
	case int:
		f = flag.Int(short, 0, UsageDummy)
		flag.IntVar(f.(*int), long, v, flagUsage)
	case bool:
		f = flag.Bool(short, false, UsageDummy)
		flag.BoolVar(f.(*bool), long, v, flagUsage)
	case float64:
		f = flag.Float64(short, 0.0, UsageDummy)
		flag.Float64Var(f.(*float64), long, v, flagUsage)
	default:
		panic("unsupported flag type")
	}
	return
}

func formatUsage(description string, maxLength *int, buffer *bytes.Buffer) {
	func() { flag.CommandLine.SetOutput(buffer); flag.Usage(); flag.CommandLine.SetOutput(os.Stderr) }()
	usageOption := regexp.MustCompile("(-\\S+)( *\\S*)+\n*\\s+"+UsageDummy+"\n\\s*").ReplaceAllString(buffer.String(), "")
	re := regexp.MustCompile("\\s(-\\S+)( *\\S*)( *\\S*)+\n\\s+(.+)")
	usageFirst := strings.Replace(strings.Replace(strings.Split(usageOption, "\n")[0], ":", " [OPTIONS] [-h, --help]", -1), " of ", ": ", -1) + "\n\nDescription:\n  " + description + "\n\nOptions:\n"
	usageOptions := re.FindAllString(usageOption, -1)
	for _, v := range usageOptions {
		*maxLength = max(*maxLength, len(re.ReplaceAllString(v, " -$1")+re.ReplaceAllString(v, "$2"))+2)
	}
	usageOptionsRep := make([]string, 0)
	for _, v := range usageOptions {
		usageOptionsRep = append(usageOptionsRep, fmt.Sprintf("  -%-1s,%-"+strconv.Itoa(*maxLength)+"s%s", strings.Split(re.ReplaceAllString(v, "$4"), UsageDummy)[0], re.ReplaceAllString(v, " -$1")+re.ReplaceAllString(v, "$2"), strings.Split(re.ReplaceAllString(v, "$4"), UsageDummy)[1]+"\n"))
	}
	sort.SliceStable(usageOptionsRep, func(i, j int) bool {
		return strings.Count(usageOptionsRep[i], UsageRequiredPrefix) > strings.Count(usageOptionsRep[j], UsageRequiredPrefix)
	})
	flag.Usage = func() { _, _ = fmt.Fprint(flag.CommandLine.Output(), usageFirst+strings.Join(usageOptionsRep, "")) }
}
