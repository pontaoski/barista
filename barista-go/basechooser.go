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

func Baseinator(c commandlib.Context) {
	c.SendMessage("introduction", c.I18n("welcome to the baseinator, the tool that will help you figure out the best base for your needs!"))
	c.SendMessage("largest", c.I18n("before we beign, let's set some boundaries. what is the largest base you'd be fine with using?"))
maximum:
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
	if maxNanpa > 720720 {
		c.SendMessage("largest-scold", c.I18n("i mean, it's nice that you want to use bases that large, but seriously. do you really want to use that many digits?"))
		goto maximum
	} else if maxNanpa < 2 {
		c.SendMessage("largest-scold", c.I18n("i said largest"))
		goto maximum
	}
smallest:
	c.SendMessage("smallest", c.I18n("now, let's see what the smallest base for you would be"))
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
	} else {
		c.SendMessage("farewell", c.I18n("unfortunately no bases meet your needs. what a shame!"))
	}
}
