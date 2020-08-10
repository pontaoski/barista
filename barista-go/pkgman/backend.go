package pkgman

import "errors"

var (
	// ErrNotSupported represents an unsupported operation
	ErrNotSupported = errors.New("This operation is not supported")
	// ErrBadDistro represents a distro that doesn't exist
	ErrBadDistro = errors.New("This distro doesn't exist")
)

type PackageDataKind int

const (
	Provides PackageDataKind = iota
	Requires
	Recommends
	Suggests
	Supplements
	Enhances
	Obsoletes
)

type PackageRelationKind int

const (
	WhatOwnsFile PackageRelationKind = iota
	WhatConflicts
	WhatRequires
	WhatObsoletes
	WhatProvides
	WhatRecommends
	WhatEnhances
	WhatSupplements
	WhatSuggests
)

type Pkg struct {
	Name         string
	Description  string
	Version      string
	DownloadSize string
	InstallSize  string
	URL          string
}

// Backend represents a backend that can search packages
type Backend interface {
	Distros() []string
	Refresh(string) error
	Search(query string) (packages []Pkg, err error)
	GetData(pkg string, kind PackageDataKind) (data []string, err error)
	GetRelation(pkg string, relation PackageRelationKind) (pkgs []Pkg, err error)
}

var backends = []Backend{}
