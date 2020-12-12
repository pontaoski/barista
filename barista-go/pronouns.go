package barista

import (
	"fmt"
	"strings"

	"github.com/appadeia/barista/barista-go/commandlib"
)

const pronounsData = `she	her	her	hers	herself
he	him	his	his	himself
they	them	their	theirs	themselves
ze	hir	hir	hirs	hirself
ze	zir	zir	zirs	zirself
xey	xem	xyr	xyrs	xemself
ae	aer	aer	aers	aerself
e	em	eir	eirs	emself
ey	em	eir	eirs	eirself
fae	faer	faer	faers	faerself
fey	fem	feir	feirs	feirself
hu	hum	hus	hus	humself
it	it	its	its	itself
jee	jem	jeir	jeirs	jemself
kit	kit	kits	kits	kitself
ne	nem	nir	nirs	nemself
peh	pehm	peh's	peh's	pehself
per	per	per	pers	perself
sie	hir	hir	hirs	hirself
se	sim	ser	sers	serself
shi	hir	hir	hirs	hirself
si	hyr	hyr	hyrs	hyrself
they	them	their	theirs	themself
thon	thon	thons	thons	thonself
ve	ver	vis	vis	verself
ve	vem	vir	virs	vemself
vi	ver	ver	vers	verself
vi	vim	vir	virs	vimself
vi	vim	vim	vims	vimself
xie	xer	xer	xers	xerself
xe	xem	xyr	xyrs	xemself
xey	xem	xeir	xeirs	xemself
yo	yo	yos	yos	yosself
ze	zem	zes	zes	zirself
ze	mer	zer	zers	zemself
zee	zed	zeta	zetas	zedself
zie	zir	zir	zirs	zirself
zie	zem	zes	zes	zirself
zie	hir	hir	hirs	hirself
zme	zmyr	zmyr	zmyrs	zmyrself`

type pronoun struct {
	subject              string
	object               string
	possessiveDeterminer string
	possessivePronoun    string
	reflexive            string
}

var pronouns []pronoun

func init() {
	for _, str := range strings.Split(pronounsData, "\n") {
		data := strings.Split(str, "\t")
		pronouns = append(pronouns, pronoun{
			subject:              data[0],
			object:               data[1],
			possessiveDeterminer: data[2],
			possessivePronoun:    data[3],
			reflexive:            data[4],
		})
	}
	commandlib.RegisterCommand(commandlib.Command{
		Name:     "English Pronouns",
		Usage:    "Shows example sentences using English pronouns",
		Examples: "o pronoun ze",
		ID:       "pronouns-en",
		Matches:  []string{"o pronoun"},
		Action: func(c commandlib.Context) {
			if c.Arg(0) == "" {
				c.SendMessage("primary", commandlib.ErrorEmbed("Please provide a pronoun"))
				return
			}

			for _, pronoun := range pronouns {
				if pronoun.subject == c.Arg(0) {
					c.SendMessage("primary", fmt.Sprintf(
						"%s went to the park, and I went with %s. %s had brought %s frisbee. Well, I think it was %s. %s threw the frisbee to %s.",
						pronoun.subject, pronoun.object, pronoun.subject, pronoun.possessiveDeterminer, pronoun.possessivePronoun, pronoun.subject, pronoun.reflexive,
					))
					return
				}
			}

			c.SendMessage("primary", commandlib.ErrorEmbed("Sorry, I don't recognise that pronoun."))
		},
	})
}
