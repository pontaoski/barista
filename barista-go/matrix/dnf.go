package matrix

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/godbus/dbus"
	"github.com/olekukonko/tablewriter"
	flag "github.com/spf13/pflag"

	"github.com/matrix-org/gomatrix"
)

type Package struct {
	name        string
	desc        string
	vers        string
	downsize    string
	installsize string
}

type Distro struct {
	matches      []string
	displayName  string
	queryKitName string
}

var Distros []Distro = []Distro{
	{
		displayName:  "openSUSE Tumbleweed",
		queryKitName: "tumbleweed",
		matches:      []string{"opensuse", "os", "tumbleweed", "tw"},
	},
	{
		displayName:  "openSUSE Leap",
		queryKitName: "leap",
		matches:      []string{"leap", "opensuse-leap", "os-leap"},
	},
	{
		displayName:  "Fedora",
		queryKitName: "fedora",
		matches:      []string{"fedora"},
	},
	{
		displayName:  "Mageia",
		queryKitName: "mageia",
		matches:      []string{"mageia"},
	},
	{
		displayName:  "OpenMandriva",
		queryKitName: "openmandriva",
		matches:      []string{"openmandriva"},
	},
	{
		displayName:  "CentOS",
		queryKitName: "centos",
		matches:      []string{"centos"},
	},
}

func resolveDistro(name string) (Distro, bool) {
	var distro Distro
	set := false

	for _, dist := range Distros {
		for _, match := range dist.matches {
			if strings.ToLower(name) == strings.ToLower(match) {
				distro = dist
				set = true
			}
		}
	}

	return distro, set
}

func rqpDnf(pkg, action, distro string) ([]string, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return []string{}, err
	}
	obj := conn.Object("com.github.Appadeia.QueryKit", "/com/github/Appadeia/QueryKit")
	var reldeps []string
	err = obj.Call("com.github.Appadeia.QueryKit.QueryRepoPackage", 0, pkg, action, distro).Store(&reldeps)
	if err != nil {
		return []string{}, err
	}
	return reldeps, nil
}

func filesDnf(pkg, distro string) ([]string, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return []string{}, err
	}
	obj := conn.Object("com.github.Appadeia.QueryKit", "/com/github/Appadeia/QueryKit")
	var files []string
	err = obj.Call("com.github.Appadeia.QueryKit.ListFiles", 0, pkg, distro).Store(&files)
	if err != nil {
		return []string{}, err
	}
	return files, nil
}

func searchDnf(search, distro string) ([]Package, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return []Package{}, err
	}
	obj := conn.Object("com.github.Appadeia.QueryKit", "/com/github/Appadeia/QueryKit")
	var pkgs [][]interface{}
	err = obj.Call("com.github.Appadeia.QueryKit.SearchPackages", 0, search, distro).Store(&pkgs)
	if err != nil {
		return []Package{}, err
	}
	var packages []Package
	for _, pkg := range pkgs {
		packages = append(packages, Package{
			name:        pkg[0].(string),
			desc:        pkg[1].(string),
			vers:        pkg[2].(string),
			downsize:    humanize.Bytes(uint64(pkg[3].(int32))),
			installsize: humanize.Bytes(uint64(pkg[4].(int32))),
		})
	}
	return packages, nil
}

func PackageTable(pkgs []Package) (string, string) {
	var plaintext [][]string
	var htmlBuilder strings.Builder
	var plaintextBuilder strings.Builder
	htmlBuilder.WriteString("<table>")
	htmlBuilder.WriteString("<tr>")
	htmlBuilder.WriteString(`<th>Name</th> <th>Summary</th> <th>Version</th> <th>Download Size</th> <th>Install Size</th>`)
	htmlBuilder.WriteString("</tr>")
	for _, pkg := range pkgs {
		htmlBuilder.WriteString("<tr>")
		htmlBuilder.WriteString(fmt.Sprintf("<td>%s</td>", pkg.name))
		htmlBuilder.WriteString(fmt.Sprintf("<td>%s</td>", pkg.desc))
		htmlBuilder.WriteString(fmt.Sprintf("<td>%s</td>", pkg.vers))
		htmlBuilder.WriteString(fmt.Sprintf("<td>%s</td>", pkg.downsize))
		htmlBuilder.WriteString(fmt.Sprintf("<td>%s</td>", pkg.installsize))
		htmlBuilder.WriteString("</tr>")
		plaintext = append(plaintext, []string{pkg.name, pkg.desc, pkg.vers, pkg.downsize, pkg.installsize})
	}
	htmlBuilder.WriteString("</table>")

	table := tablewriter.NewWriter(&plaintextBuilder)
	table.SetHeader([]string{"Name", "Summary", "Version", "Download Size", "Installed Size"})
	table.AppendBulk(plaintext)
	table.Render()
	return htmlBuilder.String(), plaintextBuilder.String()
}

func DnfSearch(client *gomatrix.Client, ev *gomatrix.Event, args []string) {
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	distro := flags.StringP("distro", "d", "", "")
	flags.Parse(args)
	if *distro == "" {
		SendMessage(client, ev.RoomID, "You need to provide a distro with the `-d`/`--distro` flag.")
		return
	}
	dist, valid := resolveDistro(*distro)
	if !valid {
		SendMessage(client, ev.RoomID, "You need to provide a distro that Barista supports.")
		SendMessage(client, ev.RoomID, "Barista supports the following distros: `tumbleweed`, `leap`, `fedora`, `openmandriva`, and `centos`.")
		return
	}
	query := strings.Join(flags.Args(), " ")
	if query == "" {
		SendMessage(client, ev.RoomID, "You need to provide a query.")
		return
	}
	packages, err := searchDnf(query, dist.queryKitName)
	if err != nil {
		SendMessage(client, ev.RoomID, "There was an error: "+err.Error())
		return
	}
	if len(packages) == 0 {
		SendMessage(client, ev.RoomID, "No packages were found.")
		return
	}
	html, plaintext := PackageTable(packages)
	SendHTMLMessage(client, ev.RoomID, html, plaintext)
}

func DnfRepoquery(client *gomatrix.Client, ev *gomatrix.Event, args []string) {
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	distro := flags.StringP("distro", "d", "", "")
	flagVals := map[string]*string{
		"conflicts":   flags.String("conflicts", "", "List packages that this package is conflicting."),
		"obsoletes":   flags.String("obsoletes", "", "List packages that this package is obsoleting."),
		"provides":    flags.String("provides", "", "List packages that this package is providing."),
		"recommends":  flags.String("recommends", "", "List packages that this package is recommending."),
		"enhances":    flags.String("enhances", "", "List packages that this package is recommending."),
		"supplements": flags.String("supplements", "", "List packages that this package is supplementing."),
		"suggests":    flags.String("suggests", "", "List packages that this package is suggesting."),
		"requires":    flags.String("requires", "", "List packages that this package is requiring."),
		"list":        flags.StringP("list", "l", "", "List files of a package."),
	}
	flags.Parse(args)
	if *distro == "" {
		SendMessage(client, ev.RoomID, "You need to provide a distro with the `-d`/`--distro` flag.")
		SendHTMLMessage(client, ev.RoomID, WrapCode(flags.FlagUsages()), flags.FlagUsages())
		return
	}
	for key, val := range flagVals {
		if *val != "" {
			if key != "files" && key != "list" {
				reldeps, err := rqpDnf(*val, key, *distro)
				if err != nil {
					SendMessage(client, ev.RoomID, "There was an error: "+err.Error())
					return
				}
				SendHTMLMessage(client, ev.RoomID, WrapCode(strings.Join(reldeps, "\n")), strings.Join(reldeps, "\n"))
				return
			}
			if key == "files" {
				files, err := filesDnf(key, *distro)
				if err != nil {
					SendMessage(client, ev.RoomID, "There was an error: "+err.Error())
					return
				}
				SendHTMLMessage(client, ev.RoomID, WrapCode(strings.Join(files, "\n")), strings.Join(files, "\n"))
				return
			}
		}
	}
	SendHTMLMessage(client, ev.RoomID, WrapCode(flags.FlagUsages()), flags.FlagUsages())
}
