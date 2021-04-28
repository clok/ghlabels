package types

import (
	"strings"

	"github.com/google/go-github/v35/github"
)

var (
	kr = k.Extend("Repo")
)

type Repo struct {
	name     string
	owner    string
	fullName string
	labels   map[string]*github.Label
}

func NewRepo(fullName string) *Repo {
	parts := strings.Split(fullName, "/")
	kr.Printf("%s -> %s / %s", fullName, parts[0], parts[1])
	return &Repo{
		name:     parts[1],
		owner:    parts[0],
		fullName: fullName,
		labels:   map[string]*github.Label{},
	}
}

func (r *Repo) Name() string {
	return r.name
}

func (r *Repo) SetName(name string) {
	r.name = name
}

func (r *Repo) Owner() string {
	return r.owner
}

func (r *Repo) SetOwner(owner string) {
	r.owner = owner
}

func (r *Repo) FullName() string {
	return r.fullName
}

func (r *Repo) SetFullName(fullName string) {
	r.fullName = fullName
}

func (r *Repo) Labels() map[string]*github.Label {
	return r.labels
}

func (r *Repo) SetLabels(labels []*github.Label) {
	for _, l := range labels {
		r.SetLabel(strings.ToLower(l.GetName()), l)
	}
}

func (r *Repo) GetLabel(name string) *github.Label {
	lower := strings.ToLower(name)
	return r.labels[lower]
}

func (r *Repo) SetLabel(name string, label *github.Label) {
	r.labels[name] = label
}

func (r *Repo) DeleteLabel(name string) {
	delete(r.labels, name)
}

func (r *Repo) HasLabel(name string) bool {
	lower := strings.ToLower(name)
	return r.Labels()[lower] != nil
}
