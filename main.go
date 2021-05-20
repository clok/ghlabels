package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/clok/cdocs"
	"github.com/clok/ghlabels/helpers"
	"github.com/clok/ghlabels/types"
	"github.com/clok/kemba"
	"github.com/urfave/cli/v2"
)

var (
	version string
	client  = types.Client{}
	k       = kemba.New("ghlabels")
)

func main() {
	k.Println("executing")

	im, err := cdocs.InstallManpageCommand(&cdocs.InstallManpageCommandInput{
		AppName: "ghlabels",
		Hidden:  true,
	})
	if err != nil {
		log.Fatal(err)
	}

	app := cli.NewApp()
	app.Name = "ghlabels"
	app.Version = version
	app.Usage = "label sync for repos and organizations"
	app.Commands = []*cli.Command{
		{
			Name:  "sync",
			Usage: "sync labels - rename, sync, delete",
			Subcommands: []*cli.Command{
				{
					Name:      "all",
					Usage:     "Sync labels across ALL repos within an org or for a user",
					UsageText: helpers.SyncAllUsageText,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "org",
							Aliases: []string{"o"},
							Usage:   "GitHub Organization to view. Cannot be used with User flag.",
						},
						&cli.StringFlag{
							Name:    "user",
							Aliases: []string{"u"},
							Usage:   "GitHub User to view. Cannot be used with Organization flag.",
						},
						&cli.StringFlag{
							Name:    "config",
							Aliases: []string{"c"},
							Usage:   "Path to config file withs labels to sync.",
						},
						&cli.BoolFlag{
							Name:    "merge-with-defaults",
							Aliases: []string{"m"},
							Usage:   "Merge provided config with defaults, otherwise only use the provided config.",
						},
					},
					Action: func(c *cli.Context) error {
						if err := helpers.ValidateOrgUserArgs(c); err != nil {
							return cli.Exit(err, 2)
						}

						config, err := helpers.DetermineConfig(c)
						if err != nil {
							return cli.Exit(err, 2)
						}

						client.GenerateClient()
						repos := helpers.GetAllRepos(c, &client)

						k.Printf("Found %d repos", len(repos))
						qualified := helpers.ExtractQualifiedRepos(repos)

						if len(qualified) == 0 {
							fmt.Println("No repos to update")
							return nil
						}

						// Confirm?
						fmt.Printf("Found %d repos that qualify\n", len(qualified))
						confirm := false
						prompt := &survey.Confirm{
							Message: fmt.Sprintf("Sync labels on %d qualifying repos?", len(qualified)),
						}
						err = survey.AskOne(prompt, &confirm)
						if err != nil {
							return cli.Exit(err, 3)
						}

						if !confirm {
							return cli.Exit("No action. Goodbye!", 0)
						}

						// NOTE: This is expensive. On a large org, it will make many API calls.
						for i, repo := range qualified {
							err = helpers.SyncLabels(repo, &client, config, i)
							if err != nil {
								return cli.Exit(err, 2)
							}
						}

						return nil
					},
				},
				{
					Name:      "repo",
					Usage:     "Sync labels for a single repo",
					UsageText: helpers.SyncRepoUsageText,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "repo",
							Aliases:  []string{"r"},
							Usage:    "Repo name including owner. Examlple: clok/ghlabels",
							Required: true,
						},
						&cli.StringFlag{
							Name:    "config",
							Aliases: []string{"c"},
							Usage:   "Path to config file withs labels to sync.",
						},
						&cli.BoolFlag{
							Name:    "merge-with-defaults",
							Aliases: []string{"m"},
							Usage:   "Merge provided config with defaults, otherwise only use the provided config.",
						},
					},
					Action: func(c *cli.Context) error {
						kl := k.Extend("sync:repo")

						config, err := helpers.DetermineConfig(c)
						if err != nil {
							return cli.Exit(err, 2)
						}

						repo := types.NewRepo(c.String("repo"))
						kl.Log(repo)

						// Confirm?
						confirm := false
						prompt := &survey.Confirm{
							Message: fmt.Sprintf("Sync labels on the %s repository?", repo.FullName()),
						}
						err = survey.AskOne(prompt, &confirm)
						if err != nil {
							return cli.Exit(err, 3)
						}

						if !confirm {
							return cli.Exit("No action. Goodbye!", 0)
						}

						client.GenerateClient()
						err = helpers.SyncLabels(repo, &client, config, 0)
						if err != nil {
							return cli.Exit(err, 2)
						}

						return nil
					},
				},
			},
		},
		{
			Name:  "config",
			Usage: "commands for viewing or generating configuration",
			Subcommands: []*cli.Command{
				{
					Name:  "defaults",
					Usage: "print default labels yaml to STDOUT",
					Action: func(c *cli.Context) error {
						fmt.Println(helpers.GetDefaultConfig())
						return nil
					},
				},
				{
					Name:  "pull",
					Usage: "pull labels from a Repo and print to STDOUT",
					UsageText: `
$ ghlabels config generate repo [<owner>/<repo_name>]

Example:

	$ ghlabels config generate repo clok/ghlabels
`,
					Action: func(c *cli.Context) error {
						kl := k.Extend("config:generate:repo")
						if c.NArg() <= 0 {
							_ = cli.ShowSubcommandHelp(c)
							return cli.Exit("ERROR: no repo provided", 2)
						}

						client.GenerateClient()
						repo := types.NewRepo(c.Args().Get(0))

						repo.SetLabels(client.GetLabels(repo))
						kl.Log(repo)

						var labels []types.Label
						for _, l := range repo.Labels() {
							labels = append(labels, types.Label{
								Name:        l.GetName(),
								Color:       l.GetColor(),
								Description: l.GetDescription(),
							})
						}
						kl.Log(labels)

						config := types.Config{Sync: labels}
						kl.Log(config)

						output, err := yaml.Marshal(&config)
						if err != nil {
							return cli.Exit(err, 2)
						}
						fmt.Printf("\n---\n%s\n", string(output))

						return nil
					},
				},
			},
		},
		{
			Name:  "stats",
			Usage: "prints out repository stats",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "org",
					Aliases: []string{"o"},
					Usage:   "GitHub Organization to view. Cannot be used with User flag.",
				},
				&cli.StringFlag{
					Name:    "user",
					Aliases: []string{"u"},
					Usage:   "GitHub User to view. Cannot be used with Organization flag.",
				},
			},
			Action: func(c *cli.Context) error {
				if err := helpers.ValidateOrgUserArgs(c); err != nil {
					return cli.Exit(err, 2)
				}

				client.GenerateClient()
				repos := helpers.GetAllRepos(c, &client)

				k.Printf("Found %d repos", len(repos))
				counts := map[string]int{"archived": 0, "forks": 0, "public": 0, "private": 0, "other": 0}
				for _, r := range repos {
					switch {
					case *r.Archived:
						counts["archived"]++
					case *r.Fork:
						counts["forks"]++
					case *r.Private:
						counts["private"]++
					case *r.Visibility == "public":
						counts["public"]++
					default:
						counts["other"]++
					}
				}
				fmt.Printf("public: %d | private: %d | forks: %d | archived: %d\n", counts["public"], counts["private"], counts["forks"], counts["archived"])
				return nil
			},
		},
		im,
	}

	if os.Getenv("DOCS_MD") != "" {
		docs, err := cdocs.ToMarkdown(app)
		if err != nil {
			panic(err)
		}
		fmt.Println(docs)
		return
	}

	if os.Getenv("DOCS_MAN") != "" {
		docs, err := cdocs.ToMan(app)
		if err != nil {
			panic(err)
		}
		fmt.Println(docs)
		return
	}

	err = app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
