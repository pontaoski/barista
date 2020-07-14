package barista

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/appadeia/barista/barista-go/bases"
	"github.com/appadeia/barista/barista-go/commandlib"
)

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Name:  "Baseinator",
		Usage: "Play the baseinator game",
		ID:    "baseinator",
		Matches: []string{
			"sudo baseinator",
			"o nasin nanpa",
		},
		Action: Baseinator,
	})
}

var sarcasm = map[bases.Base]string{
	2:      I18n("this is very simple. very, very simple. takes a lot of digits to write stuff."),
	3:      I18n("some math says this is the most economic base out there? dunno what that means."),
	4:      I18n("better binary, i guess."),
	5:      I18n("base five, for five-fingered creatures."),
	6:      I18n("this is the favourite of the guy who became a meme for a rant on rjienrlwey."),
	7:      I18n("wow, you look like you're into weird numbers."),
	8:      I18n("octal, the base that people only know because of UNIX. my favourite, though."),
	9:      I18n("weirdly fine."),
	10:     I18n("well uh, it's decimal. booring."),
	11:     I18n("not decimal."),
	12:     I18n("apparently people really like this."),
	13:     I18n("only takes one digit to represent a baker's dozen of bread."),
	14:     I18n("how do you even get to these bases?!"),
	15:     I18n("great base if you have a third hand. does anyone have a third hand?"),
	16:     I18n("what is this bases's job? being less ugly binary."),
	17:     I18n("the suboptimal base for suboptimal people."),
	18:     I18n("i'm pretty thirsy"),
	19:     I18n("base 19, the 19th base!"),
	20:     I18n("a decent base for decent people."),
	24:     I18n("only takes a single digit to write a day"),
	36:     I18n("6 * 6 gets you this?"),
	42:     I18n("the lojban word for 42 is vore."),
	48:     I18n("1/3 gross, 100% inefficient"),
	60:     I18n("the seconds are ticking"),
	69:     I18n("jajajaja sex number kkk kekeke"),
	144:    I18n("quite a gross base."),
	360:    I18n("if you need to write all possible integer degrees with only one number..."),
	420:    I18n("smonk weed everyday"),
	1337:   I18n("l33t speak is annoying ngl"),
	2540:   I18n("a superior highly composite numbers"),
	5040:   I18n("plato time"),
	9000:   I18n("well, you could have gone over 9000..."),
	9001:   I18n("you went over 9000!"),
	55440:  I18n("plato, you're drunk."),
	720720: I18n("you should learn how to say no."),
}

func Baseinator(c commandlib.Context) {
	if c.Type() == commandlib.EditCommand {
		return
	}
	c.SendMessage("introduction", c.I18n("welcome to the baseinator, the tool that will help you figure out the best base for your needs!"))
	c.SendMessage("largest", c.I18n("before we begin, let's set some boundaries. what is the largest base you'd be fine with using?"))
	goto postMaximum
maximum:
	c.SendMessage("largest-post", c.I18n("give me a better base"))
postMaximum:
	nanpa, ok := c.AwaitResponse(time.Second * 30)
	if !ok {
		c.SendMessage("farewell", c.I18n("mmm, sorry, but you were too slow. see you next time!"))
		return
	}
	maxNanpa, err := strconv.ParseInt(nanpa, 10, 64)
	if err != nil {
		c.SendMessage("largest-scold", c.I18n("that's not a number, you bonk"))
		goto maximum
	}
	if maxNanpa > 1000000 {
		c.SendMessage("largest-scold", c.I18n("i mean, it's nice that you want to use bases that large, but seriously. do you really want to use that many digits?"))
		goto maximum
	} else if maxNanpa < 2 {
		c.SendMessage("largest-scold", c.I18n("i said largest"))
		goto maximum
	}
	c.SendMessage("smallest", c.I18n("now, let's see what the smallest base for you would be"))
	goto postSmallest
smallest:
	c.SendMessage("smallest-post", c.I18n("got a better base than that?"))
postSmallest:
	nanpa, ok = c.AwaitResponse(time.Second * 30)
	if !ok {
		c.SendMessage("farewell", c.I18n("mmm, sorry, but you were too slow. see you next time!"))
		return
	}
	minNanpa, err := strconv.ParseInt(nanpa, 10, 64)
	if err != nil {
		c.SendMessage("smallest-scold", c.I18n("ain't a number, twat"))
		goto smallest
	}
	if minNanpa > 720720 {
		c.SendMessage("smallest-scold", c.I18n("y tho"))
		goto smallest
	} else if minNanpa > maxNanpa {
		c.SendMessage("smallest-scold", c.I18n("i said smallest, mate"))
		goto smallest
	} else if minNanpa == maxNanpa {
		c.SendMessage("smallest-scold", c.I18n("really? i guess it's good that you like that base, but seriously. pick a smaller base than your maximum."))
		goto smallest
	} else if minNanpa < 2 {
		c.SendMessage("smallest-scold", c.I18n("how do you even use a base that small"))
		goto smallest
	}
	pool := bases.PoolForRange(minNanpa, maxNanpa)
	denom := 2
	for len(pool) > 1 {
		status := ""
		amt, length, recurring := pool.LargestExpansionFor(denom)
		if len(amt) == len(pool) {
			denom++
			continue
		}
		if recurring {
			status = fmt.Sprintf(c.I18n("%d bases in your pool represent 1/%d with a repeating decimal"), len(amt), denom)
		} else {
			status = fmt.Sprintf(c.I18n("%d bases in your pool represent 1/%d with %d digit(s)"), len(amt), denom, length)
		}
		smallest := pool.Smallest()
		c.SendMessage(
			fmt.Sprintf("base-count-%d", denom),
			fmt.Sprintf(
				c.I18n("there are %d bases in your pool. the smallest base in your pool is base %d, %s.\n %s.\ndo you want to remove these? (Y/n)"),
				len(pool),
				smallest,
				smallest.Name(),
				status,
			),
		)
		resp, ok := c.AwaitResponse(time.Second * 30)
		if !ok {
			c.SendMessage(fmt.Sprintf("base-too-slow-%d", denom), c.I18n("i guess you do, then. next base!"))
			pool.RemoveLongestExpansionFor(denom)
		} else if strings.HasPrefix(strings.ToLower(resp), "y") {
			c.SendMessage(fmt.Sprintf("base-yes-%d", denom), c.I18n("aight, removed that. next base!"))
			pool.RemoveLongestExpansionFor(denom)
		} else {
			c.SendMessage(fmt.Sprintf("base-no-%d", denom), c.I18n("ok then, next base!"))
		}
		denom++
	}
	if len(pool) == 1 {
		best := pool.Smallest()
		c.SendMessage("farewell", fmt.Sprintf(c.I18n("oh my, looks like base %d, %s, is the best base for your needs!"), best, best.Name()))
		if val, ok := sarcasm[best]; ok {
			c.SendMessage("sarcasm", fmt.Sprintf(c.I18n(val), best, best.Name()))
		}
	} else {
		c.SendMessage("farewell", c.I18n("unfortunately no bases meet your needs. what a shame!"))
	}
}
