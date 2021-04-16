package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigMergeLeft(t *testing.T) {
	is := assert.New(t)

	t.Run("merge abd override - Rename", func(t *testing.T) {
		c := Config{
			Rename: []Rename{
				{
					From: "test",
					To: "yolo",
				},
				{
					From: "test2",
					To: "yolo2",
				},
			},
		}

		u := Config{
			Rename: []Rename{
				{
					From: "test",
					To:   "not yolo",
				},
				{
					From: "new",
					To:   "tag",
				},
			},
		}
		c.MergeLeft(u)

		type test struct {
			expect Rename
			index int
		}

		table := []test{
			{ expect: Rename{"test", "not yolo"}, index: 0 },
			{ expect: Rename{"new", "tag"}, index: 1 },
			{ expect: Rename{"test2", "yolo2"}, index: 2 },
		}

		for _, tc := range table {
			is.Equal(tc.expect, c.Rename[tc.index])
		}
	})

	t.Run("merge abd override - Remove", func(t *testing.T) {
		c := Config{
			Remove: []string{"wont", "do a thing"},
		}

		u := Config{
			Remove: []string{"will", "do a thing"},
		}
		c.MergeLeft(u)

		type test struct {
			expect string
			index int
		}

		table := []test{
			{ expect: "will", index: 0 },
			{ expect: "do a thing", index: 1 },
			{ expect: "wont", index: 2 },
		}

		for _, tc := range table {
			is.Equal(tc.expect, c.Remove[tc.index])
		}
	})

	t.Run("merge abd override - Sync", func(t *testing.T) {
		c := Config{
			Sync: []Label{
				{
					Name: "all that matters",
					Color: "FFFFFF",
					Description: "a thing",
				},
				{
					Name: "is the name",
					Color: "000000",
					Description: "another thing",
				},
			},
		}

		u := Config{
			Sync: []Label{
				{
					Name: "all that matters",
					Color: "123456",
					Description: "different",
				},
				{
					Name: "new-tag",
					Color: "987654",
					Description: "yolo",
				},
			},
		}
		c.MergeLeft(u)

		type test struct {
			expect Label
			index int
		}

		table := []test{
			{ expect: Label{"all that matters", "123456", "different"}, index: 0 },
			{ expect: Label{"new-tag", "987654", "yolo"}, index: 1 },
			{ expect: Label{"is the name", "000000", "another thing"}, index: 2 },
		}

		for _, tc := range table {
			is.Equal(tc.expect, c.Sync[tc.index])
		}
	})
}