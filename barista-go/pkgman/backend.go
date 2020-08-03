package pkgman

import "errors"

// ErrNotSupported represents an unsupported operation
var ErrNotSupported = errors.New("This operation is not supported")

// Distro represents a distro
type Distro interface {
	isDistro()
}

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
	Distros() []Distro
	Refresh(distro Distro) error
	Search(query string) (packages []Pkg, err error)
}
