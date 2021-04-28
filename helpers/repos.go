package helpers

import (
	"fmt"

	"github.com/clok/ghlabels/types"
	"github.com/google/go-github/v35/github"
	"github.com/urfave/cli/v2"
)

// GetAllRepos retrieves all repos for an Organization or User
func GetAllRepos(c *cli.Context, client *types.Client) []*github.Repository {
	var repos []*github.Repository
	switch {
	case c.String("org") != "":
		fmt.Printf("Pulling all repos for Organization %s ...\n", c.String("org"))
		return client.GetAllOrgRepos(c.String("org"))
	case c.String("user") != "":
		fmt.Printf("Pulling all repos for User %s ...\n", c.String("user"))
		return client.GetAllUserRepos(c.String("user"))
	}
	return repos
}

// ExtractQualifiedRepos ensures we only action on non-archived and non-fork repos
func ExtractQualifiedRepos(repos []*github.Repository) []*types.Repo {
	fmt.Printf("Checking %d repos ...\n", len(repos))
	var qualified []*types.Repo
	for _, r := range repos {
		switch {
		case *r.Fork:
			// skip
		case *r.Archived:
			// skip
		case *r.Private:
			qualified = append(qualified, types.NewRepo(*r.FullName))
		case *r.Visibility == "public":
			qualified = append(qualified, types.NewRepo(*r.FullName))
		default:
			// skip
		}
	}
	return qualified
}
