package pkgman

import (
	"github.com/appadeia/barista/barista-go/log"
	"modernc.org/ql"
)

var (
	db = func() *ql.DB {
		db, err := ql.OpenFile("storage/rpm.db", &ql.Options{
			CanCreate:  true,
			FileFormat: 2,
		})
		if err != nil {
			log.Fatal(log.DatabaseFailure, "Failed to open/create database: %+v", err)
		}
		ctx := ql.NewRWCtx()
		_, _, err = db.Run(ctx, createTable)
		if err != nil {
			log.Fatal(log.DatabaseFailure, "Failed to create tables: %+v", err)
		}
		return db
	}()
)

const (
	createTable = `
BEGIN TRANSACTION;
	CREATE TABLE IF NOT EXISTS Packages (
		Distro       string NOT NULL,
		Name         string NOT NULL,
		Description  string NOT NULL,
		Version      string NOT NULL,
		DownloadSize uint64 NOT NULL,
		InstallSize  uint64 NOT NULL,
		DownloadURL  string NOT NULL,
	);
	CREATE TABLE IF NOT EXISTS Relations (
		Distro string NOT NULL,
		Name   string NOT NULL,
		Kind   string Kind IN ("Provides", "Requires", "Recommends", "Suggests", "Supplements", "Enhances", "Obsoletes"),
		Value  string NOT NULL,
	);
	CREATE TABLE IF NOT EXISTS Files (
		Distro string NOT NULL,
		Name   string NOT NULL,
		File   string NOT NULL,
	);
COMMIT;
	`
	clearDistro = `
BEGIN TRANSACTION;
	DELETE FROM Packages WHERE
		Distro == $1;
	DELETE FROM Files WHERE
		Distro == $1;
	DELETE FROM Relations WHERE
		Distro == $1;
COMMIT;
	`
	insertPackage = `
BEGIN TRANSACTION;
	INSERT INTO Packages (
		Distro, Name, Description, Version, DownloadSize, InstallSize, DownloadURL
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7
	);
COMMIT;
`
	insertRelation = `
BEGIN TRANSACTION;
	INSERT INTO Relations (
		Distro, Name, Kind, Value
	) VALUES (
		$1, $2, $3, $4
	);
COMMIT;
`
	insertFile = `
BEGIN TRANSACTION;
	INSERT INTO Files (
		Distro, Name, File
	) VALUES (
		$1, $2, $3
	);
COMMIT;
`

	searchPackages = `
SELECT * FROM Packages WHERE
	Distro == $1 &&
	Name LIKE $2;
	`

	getData = `
SELECT Value FROM Relations WHERE
	Distro == $1 &&
	Name == $2 &&
	Kind == $3;
	`

	getFile = `
SELECT Name FROM Files WHERE
	Distro == $1 &&
	File == $2;
	`

	getPackage = `
SELECT * FROM Packages WHERE
	Distro == $1 &&
	Name == $2;
	`

	getFiles = `
SELECT File FROM Files WHERE
	Distro == $1 &&
	Name == $2;
	`

	getRelation = `
SELECT Name FROM Relations WHERE
	Distro == $1 &&
	Kind == $2 &&
	Value == $3;
	`
)

// ClearDistro drops all packages belonging to a distro
func ClearDistro(name string) error {
	ctx := ql.NewRWCtx()
	_, _, err := db.Run(ctx, clearDistro, name)
	return err
}

// DBPackage represents a package in the database
type DBPackage struct {
	Distro       string
	Name         string
	Description  string
	Version      string
	DownloadSize uint64
	InstallSize  uint64
	DownloadURL  string

	Relations map[string][]string
	Files     []string
}

// InsertPackage inserts a package into the database
func InsertPackage(pkg DBPackage) error {
	ctx := ql.NewRWCtx()

	_, _, err := db.Run(ctx, insertPackage, pkg.Distro, pkg.Name, pkg.Description, pkg.Version, pkg.DownloadSize, pkg.InstallSize, pkg.DownloadURL)
	if err != nil {
		return err
	}

	for rel, vals := range pkg.Relations {
		for _, val := range vals {
			_, _, err := db.Run(ctx, insertRelation, pkg.Distro, pkg.Name, rel, val)
			if err != nil {
				return err
			}
		}
	}
	for _, file := range pkg.Files {
		_, _, err := db.Run(ctx, insertFile, pkg.Distro, pkg.Name, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func Search(distro, query string) (packages []Pkg, err error) {
	ctx := ql.NewRWCtx()

	data, _, err := db.Run(ctx, searchPackages, distro, query)

	if err != nil {
		return
	}

	for _, row := range data {
		err = row.Do(false, func(data []interface{}) (bool, error) {
			packages = append(packages, Pkg{
				data[1].(string),
				data[2].(string),
				data[3].(string),
				data[4].(uint64),
				data[5].(uint64),
				data[6].(string),
			})
			return true, nil
		})
		if err != nil {
			return
		}
	}

	return
}

func GetData(distro, pkg string, kind PackageDataKind) (data []string, err error) {
	val := map[PackageDataKind]string{
		Conflicts:   "Conflicts",
		Provides:    "Provides",
		Requires:    "Requires",
		Recommends:  "Recommends",
		Suggests:    "Suggests",
		Supplements: "Supplements",
		Enhances:    "Enhances",
		Obsoletes:   "Obsoletes",
	}[kind]

	ctx := ql.NewRWCtx()

	ret, _, err := db.Run(ctx, getData, distro, pkg, val)

	for _, row := range ret {
		err = row.Do(false, func(data []interface{}) (bool, error) {
			data = append(data, data[0].(string))
			return true, nil
		})
	}

	return
}

func WhatOwnsFile(distro, file string) (pkg Pkg, err error) {
	ctx := ql.NewRWCtx()

	ret, _, err := db.Run(ctx, getFile, distro, file)

	if err != nil {
		return
	}

	data, err := ret[0].FirstRow()

	if err != nil {
		return
	}

	dbpkg, _, err := db.Run(ctx, getPackage, distro, data[0])

	if err != nil {
		return
	}

	retpkg, err := dbpkg[0].FirstRow()

	return Pkg{
		retpkg[1].(string),
		retpkg[2].(string),
		retpkg[3].(string),
		retpkg[4].(uint64),
		retpkg[5].(uint64),
		retpkg[6].(string),
	}, nil
}

func GetFiles(distro, pkg string) (files []string, err error) {
	ctx := ql.NewRWCtx()

	ret, _, err := db.Run(ctx, getFiles, distro, pkg)
	if err != nil {
		return
	}

	for _, row := range ret {
		err = row.Do(false, func(data []interface{}) (bool, error) {
			files = append(files, data[0].(string))
			return true, nil
		})
		if err != nil {
			return
		}
	}

	return
}

func GetPackagesWithRelation(distro string, kind PackageRelationKind, value string) (pkgs []Pkg, err error) {
	ctx := ql.NewRWCtx()

	kindString := map[PackageRelationKind]string{
		WhatConflicts:   "Conflicts",
		WhatProvides:    "Provides",
		WhatRequires:    "Requires",
		WhatRecommends:  "Recommends",
		WhatSuggests:    "Suggests",
		WhatSupplements: "Supplements",
		WhatEnhances:    "Enhances",
		WhatObsoletes:   "Obsoletes",
	}[kind]

	ret, _, err := db.Run(ctx, getRelation, distro, kindString, value)
	if err != nil {
		return
	}

	var names []string

	for _, row := range ret {
		err = row.Do(false, func(data []interface{}) (bool, error) {
			names = append(names, data[0].(string))
			return true, nil
		})
		if err != nil {
			return
		}
	}

	for _, name := range names {
		var data []ql.Recordset

		data, _, err = db.Run(ctx, getPackage, distro, name)
		if err != nil {
			return
		}

		for _, row := range data {
			err = row.Do(false, func(data []interface{}) (bool, error) {
				pkgs = append(pkgs, Pkg{
					data[1].(string),
					data[2].(string),
					data[3].(string),
					data[4].(uint64),
					data[5].(uint64),
					data[6].(string),
				})
				return true, nil
			})
			if err != nil {
				return
			}
		}
	}

	return
}
