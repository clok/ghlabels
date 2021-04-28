package main

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/clok/cdocs"
	"github.com/clok/ghlabels/types"
	"github.com/clok/kemba"
	"github.com/google/go-github/v35/github"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

var (
	version string
	client = types.Client{}
	k = kemba.New("ghlabels")
)

//go:embed embeds/defaults.yml
var defaultLabelsBytes []byte

// TODO: sync from Org to Repos
// TODO: update Org labels from manifest

func main() {
	k.Println("executing")

	im, err := cdocs.InstallManpageCommand(&cdocs.InstallManpageCommandInput{
		AppName: "ghlabels",
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
			Name:      "sync",
			Usage:     "sync labels - delete, rename, update",
			Subcommands: []*cli.Command{
				{
					Name: "all",
					Usage: "Sync labels across ALL repos within an org or for a user",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "org",
							Aliases:     []string{"o"},
							Usage:       "GitHub Organization to view. Cannot be used with User flag.",
						},
						&cli.StringFlag{
							Name:        "user",
							Aliases:     []string{"u"},
							Usage:       "GitHub User to view. Cannot be used with Organization flag.",
						},
						&cli.StringFlag{
							Name:        "config",
							Aliases:     []string{"c"},
							Usage:       "Path to config file withs labels to sync.",
						},
						&cli.BoolFlag{
							Name:        "merge-with-defaults",
							Aliases:     []string{"m"},
							Usage:       "Merge provided config with defaults, otherwise only use the provided config.",
						},
					},
					Action: func(c *cli.Context) error {
						kl := k.Extend("sync")
						if c.String("org") != "" && c.String("user") != "" {
							return cli.Exit(fmt.Errorf("cannot pass both organization and user flag"), 2)
						}

						var defaultLabels types.Config
						err := yaml.Unmarshal(defaultLabelsBytes, &defaultLabels)
						if err != nil {
							return cli.Exit(err, 2)
						}

						var userConfig types.Config
						var config types.Config
						if c.String("config") != "" {
							yamlFile, err := ioutil.ReadFile(c.String("config"))
							if err != nil {
								return cli.Exit(err, 2)
							}
							err = yaml.Unmarshal(yamlFile, &userConfig)
							if err != nil {
								return cli.Exit(err, 2)
							}
							if c.Bool("merge-with-defaults") {
								defaultLabels.MergeLeft(userConfig)
								config = defaultLabels
							} else {
								config = userConfig
							}
						} else {
							config = defaultLabels
						}
						kl.Log(config)

						client.GenerateClient()

						// get all pages of results
						var repos []*github.Repository
						switch {
						case c.String("org") != "":
							repos = client.GetAllOrgRepos(c.String("org"))
						case c.String("user") != "":
							repos = client.GetAllUserRepos(c.String("user"))
						}

						k.Printf("Found %d repos", len(repos))
						// Only action on non-archived and non-fork repos
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

						if len(qualified) == 0 {
							fmt.Println("No repos to update")
							return nil
						}

						// Confirm?
						fmt.Printf("Found %d repos that qualify\n", len(qualified))
						confirm := false
						prompt := &survey.Confirm{
							Message: "Sync dem labels????",
						}
						err = survey.AskOne(prompt, &confirm)
						if err != nil {
							return cli.Exit(err, 3)
						}

						if !confirm {
							return cli.Exit("No action. Goodbye!", 0)
						}

						// Do the thing!
						// NOTE: This is expensive. On a large org, it will many K API calls.
						actions := 0
						for _, repo := range qualified {
							repo.SetLabels(client.GetLabels(repo))

							// TODO: Clean up duplicate code
							// Rename labels
							// 1. Check if has the label to rename
							// 2. Check if new label already exists
							// 3. If NOT
							//     - Rename the label
							// 4. If New Label DOES EXISTS
							//     - Create new label, if it does not exist
							//     - Swap out old Label for new Label on ALL
							//       Issues that have the old Label
							for _, l := range defaultLabels.Rename {
								if repo.HasLabel(l.From) {
									k.Extend("rename").Printf("%s -> %s", l.From, l.To)

									if sync := defaultLabels.FindSyncLabel(l.To); sync != nil {
										if !repo.HasLabel(l.To) {
											label := client.CreateLabel(repo, *sync)
											repo.SetLabel(label.GetName(), label)
										}

										// If there are issues with the original label
										issues, err := client.GetRepoIssues(repo, l.From)
										if err != nil {
											return cli.Exit(err, 2)
										}

										// Update Issues, swapping From for To
										for _, issue := range issues {
											kl.Printf("updating issue %d labels. Swapping %s for %s", *issue.Number, l.From, l.To)
											err = client.UpdateIssueLabels(repo, issue, l.From, l.To)
											if err != nil {
												return cli.Exit(err, 2)
											}
										}

										// Delete From label from the Repo on GH
										client.DeleteLabel(repo, l.From)
										// local copy
										repo.DeleteLabel(l.From)

										// Update the To
										_ = client.UpdateLabel(repo, l.To, *sync)
									} else {
										client.RenameLabel(repo, l)
										updatedLabel := client.GetLabel(repo, l.To)
										repo.SetLabel(updatedLabel.GetName(), updatedLabel)
										repo.DeleteLabel(l.From)
									}
								}
							}

							// Update labels
							for _, l := range defaultLabels.Sync {
								if repo.HasLabel(l.Name) {
									label := repo.GetLabel(l.Name)
									// is the label different
									if l.Color != *label.Color || l.Description != *label.Description {
										k.Extend("update").Printf("%s: Color or description is different.", l.Name)
										client.UpdateLabel(repo, l.Name, l)
									}
								} else {
									k.Extend("create").Printf("%s: Does not exist. Creating...", l.Name)
									label := client.CreateLabel(repo, l)
									repo.SetLabel(label.GetName(), label)
								}
							}

							// Delete labels
							for _, l := range defaultLabels.Remove {
								if repo.HasLabel(l) {
									k.Extend("remove").Println(l)
									client.DeleteLabel(repo, l)
									repo.DeleteLabel(l)
								}
							}

							fmt.Printf("Label sync complete for %s", repo.Name())
						}

						fmt.Printf("TOTAL ACTIONS: %d\n", actions)
						return nil
					},
				},
				{
					Name: "repo",
					Usage: "Sync labels for a single repo",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "repo",
							Aliases:     []string{"r"},
							Usage:       "Repo name including owner. Examlple: clok/ghlabels",
							Required: true,
						},
						&cli.StringFlag{
							Name:        "config",
							Aliases:     []string{"c"},
							Usage:       "Path to config file withs labels to sync.",
						},
						&cli.BoolFlag{
							Name:        "merge-with-defaults",
							Aliases:     []string{"m"},
							Usage:       "Merge provided config with defaults, otherwise only use the provided config.",
						},
					},
					Action: func(c *cli.Context) error {
						kl := k.Extend("sync:repo")

						var defaultLabels types.Config
						err := yaml.Unmarshal(defaultLabelsBytes, &defaultLabels)
						if err != nil {
							return cli.Exit(err, 2)
						}

						var userConfig types.Config
						var config types.Config
						if c.String("config") != "" {
							yamlFile, err := ioutil.ReadFile(c.String("config"))
							if err != nil {
								return cli.Exit(err, 2)
							}
							err = yaml.Unmarshal(yamlFile, &userConfig)
							if err != nil {
								return cli.Exit(err, 2)
							}
							if c.Bool("merge-with-defaults") {
								defaultLabels.MergeLeft(userConfig)
								config = defaultLabels
							} else {
								config = userConfig
							}
						} else {
							config = defaultLabels
						}
						kl.Log(config)

						// get all pages of results
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
						// Do the thing!
						repo.SetLabels(client.GetLabels(repo))

						// TODO: Clean up duplicate code
						// Rename labels
						// 1. Check if has the label to rename
						// 2. Check if new label already exists
						// 3. If NOT
						//     - Rename the label
						// 4. If New Label DOES EXISTS
						//     - Create new label, if it does not exist
						//     - Swap out old Label for new Label on ALL
						//       Issues that have the old Label
						for _, l := range defaultLabels.Rename {
							if repo.HasLabel(l.From) {
								k.Extend("rename").Printf("%s -> %s", l.From, l.To)

								if sync := defaultLabels.FindSyncLabel(l.To); sync != nil {
									if !repo.HasLabel(l.To) {
										label := client.CreateLabel(repo, *sync)
										repo.SetLabel(label.GetName(), label)
									}

									// If there are issues with the original label
									issues, err := client.GetRepoIssues(repo, l.From)
									if err != nil {
										return cli.Exit(err, 2)
									}

									// Update Issues, swapping From for To
									for _, issue := range issues {
										kl.Printf("updating issue %d labels. Swapping %s for %s", *issue.Number, l.From, l.To)
										err = client.UpdateIssueLabels(repo, issue, l.From, l.To)
										if err != nil {
											return cli.Exit(err, 2)
										}
									}

									// Delete From label from the Repo on GH
									client.DeleteLabel(repo, l.From)
									// local copy
									repo.DeleteLabel(l.From)

									// Update the To
									_ = client.UpdateLabel(repo, l.To, *sync)
								} else {
									client.RenameLabel(repo, l)
									updatedLabel := client.GetLabel(repo, l.To)
									repo.SetLabel(updatedLabel.GetName(), updatedLabel)
									repo.DeleteLabel(l.From)
								}
							}
						}

						// Update labels
						for _, l := range defaultLabels.Sync {
							if repo.HasLabel(l.Name) {
								label := repo.GetLabel(l.Name)
								// is the label different
								if l.Color != *label.Color || l.Description != *label.Description {
									k.Extend("update").Printf("%s: Color or description is different.", l.Name)
									client.UpdateLabel(repo, l.Name, l)
								}
							} else {
								k.Extend("create").Printf("%s: Does not exist. Creating...", l.Name)
								label := client.CreateLabel(repo, l)
								repo.SetLabel(label.GetName(), label)
							}
						}

						// Delete labels
						for _, l := range defaultLabels.Remove {
							if repo.HasLabel(l) {
								k.Extend("remove").Println(l)
								client.DeleteLabel(repo, l)
								repo.DeleteLabel(l)
							}
						}

						fmt.Printf("Label sync complete for %s\n", repo.FullName())
						return nil
					},
				},
			},
		},
		{
			Name: "dump-defaults",
			Usage: "print default labels yaml to STDOUT",
			Action: func(c *cli.Context) error {
				fmt.Println(string(defaultLabelsBytes))
				return nil
			},
		},
		{
			Name:      "stats",
			Usage:     "prints out repo stats",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "org",
					Aliases:     []string{"o"},
					Usage:       "GitHub Organization to view. Cannot be used with User flag.",
				},
				&cli.StringFlag{
					Name:        "user",
					Aliases:     []string{"u"},
					Usage:       "GitHub User to view. Cannot be used with Organization flag.",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("org") != "" && c.String("user") != "" {
					return cli.Exit(fmt.Errorf("cannot pass both organization and user flag"), 2)
				}

				client.GenerateClient()

				// get all pages of results
				var repos []*github.Repository
				switch {
				case c.String("org") != "":
					repos = client.GetAllOrgRepos(c.String("org"))
				case c.String("user") != "":
					repos = client.GetAllUserRepos(c.String("user"))
				}

				k.Printf("Found %d repos", len(repos))
				counts := map[string]int{"archived": 0, "forks": 0, "public": 0, "private": 0, "other": 0}
				for _, r := range repos {
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
