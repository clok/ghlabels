package types

import (
	"context"
	"log"
	"os"

	"github.com/clok/kemba"
	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

var (
	k  = kemba.New("ghlabels:types")
	kc = k.Extend("Client")
)

type Client struct {
	client *github.Client
}

func (c *Client) Client() *github.Client {
	return c.client
}

func (c *Client) GenerateClient() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	c.client = github.NewClient(tc)
}

func (c Client) GetAllUserRepos(user string) (allRepos []*github.Repository) {
	kl := kc.Extend("GetAllUserRepos")
	opt := &github.RepositoryListOptions{
		Type:        "all",
		ListOptions: github.ListOptions{PerPage: 50},
	}

	ctx := context.Background()
	for {
		repos, resp, err := c.Client().Repositories.List(ctx, user, opt)
		if err != nil {
			log.Fatal(err)
		}
		allRepos = append(allRepos, repos...)

		opt.Page = resp.NextPage
		kl.Printf("Next Page: %d", opt.Page)
		if resp.NextPage == 0 {
			kl.Println("done fetching")
			break
		}
	}
	return
}

func (c Client) GetAllOrgRepos(org string) (allRepos []*github.Repository) {
	kl := kc.Extend("GetAllOrgRepos")
	opt := &github.RepositoryListByOrgOptions{
		Type:        "all",
		ListOptions: github.ListOptions{PerPage: 50},
	}

	ctx := context.Background()
	for {
		repos, resp, err := c.Client().Repositories.ListByOrg(ctx, org, opt)
		if err != nil {
			log.Fatal(err)
		}
		allRepos = append(allRepos, repos...)

		opt.Page = resp.NextPage
		kl.Printf("Next Page: %d", opt.Page)
		if resp.NextPage == 0 {
			kl.Println("done fetching")
			break
		}
	}
	return
}

func (c Client) GetRepository(owner string, name string) *github.Repository {
	kl := kc.Extend("GetRepository")

	ctx := context.Background()
	repo, _, err := c.Client().Repositories.Get(ctx, owner, name)
	if err != nil {
		log.Fatal(err)
	}
	kl.Log(repo)

	return repo
}

func (c Client) GetLabel(repo *Repo, name string) *github.Label {
	kl := kc.Extend(repo.FullName())

	ctx := context.Background()
	label, _, err := c.Client().Issues.GetLabel(ctx, repo.Owner(), repo.Name(), name)
	if err != nil {
		log.Fatal(err)
	}
	kl.Printf("GetLabel - updated label: %# v", label)

	return label
}

func (c Client) GetLabels(repo *Repo) []*github.Label {
	kl := kc.Extend(repo.FullName())

	opt := &github.ListOptions{
		PerPage: 100,
	}
	ctx := context.Background()
	labels, _, err := c.Client().Issues.ListLabels(ctx, repo.Owner(), repo.Name(), opt)
	if err != nil {
		log.Fatal(err)
	}
	kl.Printf("GetLabels - found %d labels", len(labels))

	return labels
}

func (c Client) GetRepoIssues(repo *Repo, label string) ([]*github.Issue, error) {
	kl := kc.Extend(repo.FullName())

	opts := github.IssueListByRepoOptions{State: "all"}
	if label != "" {
		opts.Labels = []string{label}
	}
	kl.Printf("options: %# v", opts)

	ctx := context.Background()
	issues, _, err := c.Client().Issues.ListByRepo(ctx, repo.Owner(), repo.Name(), &opts)
	kl.Printf("GetRepoIssues: %# v", issues)
	if err != nil {
		return nil, err
	}
	return issues, nil
}

func (c Client) UpdateIssueLabels(repo *Repo, issue *github.Issue, from string, to string) error {
	kl := kc.Extend(repo.FullName())

	var labels []string
	// remove From
	for _, l := range issue.Labels {
		if *l.Name != from {
			labels = append(labels, *l.Name)
		}
	}
	// add To
	labels = append(labels, to)
	kl.Printf("UpdateIssueLabels: %# v", labels)

	ctx := context.Background()
	_, resp, err := c.Client().Issues.AddLabelsToIssue(ctx, repo.Owner(), repo.Name(), *issue.Number, labels)
	if err != nil {
		return err
	}
	kl.Printf("UpdateIssueLabels: %s: %d", resp.Status, resp.StatusCode)

	return nil
}

func (c Client) RenameLabel(repo *Repo, label Rename) {
	kl := kc.Extend(repo.FullName())

	kl.Printf("RenameLabel: %# v", label)
	ctx := context.Background()
	_, _, err := c.Client().Issues.EditLabel(ctx, repo.Owner(), repo.Name(), label.From, &github.Label{
		Name: &label.To,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (c Client) UpdateLabel(repo *Repo, name string, label Label) *github.Label {
	kl := kc.Extend(repo.FullName())

	opts := &github.Label{
		Name:        &label.Name,
		Color:       &label.Color,
		Description: &label.Description,
	}
	kl.Printf("UpdateLabel: label '%s' with opts: %# v", name, opts)
	ctx := context.Background()
	updatedLabel, _, err := c.Client().Issues.EditLabel(ctx, repo.Owner(), repo.Name(), name, opts)
	if err != nil {
		log.Fatal(err)
	}
	return updatedLabel
}

func (c Client) CreateLabel(repo *Repo, label Label) *github.Label {
	kl := kc.Extend(repo.FullName())

	opts := &github.Label{
		Name:        &label.Name,
		Color:       &label.Color,
		Description: &label.Description,
	}
	kl.Printf("CreateLabel: %# v", opts)
	ctx := context.Background()
	updatedLabel, _, err := c.Client().Issues.CreateLabel(ctx, repo.Owner(), repo.Name(), opts)
	if err != nil {
		log.Fatal(err)
	}
	return updatedLabel
}

func (c Client) DeleteLabel(repo *Repo, label string) {
	ctx := context.Background()
	_, err := c.Client().Issues.DeleteLabel(ctx, repo.Owner(), repo.Name(), label)
	if err != nil {
		log.Fatal(err)
	}
}
