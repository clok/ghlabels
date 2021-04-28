package helpers

import (
	"fmt"

	"github.com/clok/ghlabels/types"
	"github.com/clok/kemba"
)

var (
	k   = kemba.New("ghlabels:helpers")
	ksl = k.Extend("SyncLabels")
)

// SyncLabels will apply the provided Labels Config to a given repo.
func SyncLabels(repo *types.Repo, client *types.Client, labels *types.Config, id int) error {
	fmt.Printf("\rUpdating labels %s: START", repo.FullName())

	// Do the thing!
	repo.SetLabels(client.GetLabels(repo))

	// Rename labels
	// 1. Check if has the label to rename
	// 2. Check if new label already exists
	// 3. If NOT
	//     - Rename the label
	// 4. If New Label DOES EXISTS
	//     - Create new label, if it does not exist
	//     - Swap out old Label for new Label on ALL
	//       Issues that have the old Label
	for _, l := range labels.Rename {
		fmt.Printf("\r[%d] Updating labels %s: RENAME", id, repo.FullName())
		if repo.HasLabel(l.From) {
			ksl.Printf("[rename] %s -> %s", l.From, l.To)

			if sync := labels.FindSyncLabel(l.To); sync != nil {
				if !repo.HasLabel(l.To) {
					label := client.CreateLabel(repo, *sync)
					repo.SetLabel(label.GetName(), label)
				}

				// If there are issues with the original label
				issues, err := client.GetRepoIssues(repo, l.From)
				if err != nil {
					return err
				}

				// Update Issues, swapping From for To
				for _, issue := range issues {
					ksl.Printf("[rename] updating issue %d labels. Swapping %s for %s", *issue.Number, l.From, l.To)
					err = client.UpdateIssueLabels(repo, issue, l.From, l.To)
					if err != nil {
						return err
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
	for _, l := range labels.Sync {
		fmt.Printf("\r[%d] Updating labels %s: SYNC", id, repo.FullName())
		if repo.HasLabel(l.Name) {
			label := repo.GetLabel(l.Name)
			// is the label different
			if l.Color != *label.Color || l.Description != *label.Description {
				ksl.Printf("[update] %s Color or description is different.", l.Name)
				client.UpdateLabel(repo, l.Name, l)
			}
		} else {
			ksl.Printf("[create] %s: Does not exist. Creating...", l.Name)
			label := client.CreateLabel(repo, l)
			repo.SetLabel(label.GetName(), label)
		}
	}

	// Delete labels
	for _, l := range labels.Remove {
		fmt.Printf("\r[%d] Updating labels %s: DELETE", id, repo.FullName())
		if repo.HasLabel(l) {
			ksl.Printf("[remove] deleting label %s", l)
			client.DeleteLabel(repo, l)
			repo.DeleteLabel(l)
		}
	}

	fmt.Printf("\r[%d] Updating labels %s: COMPLETE\n", id, repo.FullName())
	return nil
}
