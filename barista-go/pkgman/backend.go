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
	Conflicts PackageDataKind = iota
	Provides
	Requires
	Recommends
	Suggests
	Supplements
	Enhances
	Obsoletes
)

type PackageRelationKind int

const (
	WhatConflicts PackageRelationKind = iota
	WhatProvides
	WhatRequires
	WhatRecommends
	WhatSuggests
	WhatSupplements
	WhatEnhances
	WhatObsoletes
)

type Pkg struct {
	Name         string
	Description  string
	Version      string
	DownloadSize uint64
	InstallSize  uint64
	URL          string
}

// Backend represents a backend that can search packages
type Backend interface {
	Distros() []string
	Refresh(string) error
}

var backends = []Backend{}
