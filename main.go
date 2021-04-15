package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"log"
	"os"

	"github.com/clok/kemba"
	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

var (
	k = kemba.New("sync-labels")
)

type Label struct {
	Name string `yaml:"name"`
	Color string `yaml:"color"`
	Description string `yaml:"description"`
}

//go:embed defaults.yml
var defaultLabelsBytes []byte

func main() {
	var defaultLabels []Label
	err := yaml.Unmarshal(defaultLabelsBytes, &defaultLabels)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	k.Log(defaultLabels)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "TOKEN"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{
		Type: "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, "GoodwayGroup", opt)
		if err != nil {
			log.Fatal(err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			k.Println("done fetching")
			break
		}
		opt.Page = resp.NextPage
		k.Printf("Next Page: %d", opt.Page)
	}

	k.Printf("Found %d repos", len(allRepos))
	counts := map[string]int{"archived": 0, "forks": 0, "public": 0, "private": 0, "other": 0}
	for _, r := range allRepos {
		switch {
		case *r.Archived:
			counts["archived"] += 1
		case *r.Fork:
			counts["forks"] += 1
		case *r.Private:
			counts["private"] += 1
		case *r.Visibility == "public":
			counts["public"] += 1
		default:
			counts["other"] += 1
		}
	}
	k.Log(counts)
}

// printJSON prints v as JSON encoded with indent to stdout. It panics on any error.
func printJSON(v interface{}) {
	w := json.NewEncoder(os.Stdout)
	w.SetIndent("", "\t")
	err := w.Encode(v)
	if err != nil {
		panic(err)
	}
}
