# go-keywordminer

A command-line tool that extracts and analyzes keywords from web pages. This tool helps you identify the most relevant keywords, retrieve page titles and meta tags, making it useful for SEO analysis and content optimization.

## Features

- Extract and analyze keywords from any web page
- Retrieve page titles and meta tags
- Calculate keyword relevance scores
- Display top keywords ranked by importance

## Installation

### Build from source

```
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -trimpath ./cmd/keywordminer
```

### Using Go

```
go install github.com/xshoji/go-keywordminer/cmd/keywordminer@latest
```

## Usage

Basic usage requires providing a URL to analyze:

```
keywordminer -u https://example.com
```

or using the long option format:

```
keywordminer --url https://example.com
```

### Example output

```
[ Command options ]
  -u, --url http://example.com       URL

[Title]
Example Domain

[Meta Tags]
description: This is an example website

[Top Keywords]
example (score: 15), domain (score: 12), website (score: 8), ...
```

## Important Considerations

When using this tool, please be aware of the following:

- Always respect the terms of service of the websites you analyze
- Do not use this tool to extract personal information or copyrighted content
- Be considerate of the website's server load by limiting request frequency

## License

This project is licensed under the MIT License - see the LICENSE file for details.
