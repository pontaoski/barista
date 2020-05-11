package barista

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/appadeia/barista/barista-go/commandlib"
)

func EtyDict() string {
	return `
[
	{
		"word": "a",
		"origins": {
			"language": "A priori.",
			"word": null
		}
	},
	{
		"word": "akesi",
		"origins": {
			"language": "Dutch (Nederlands)",
			"word": "hagedis"
		}
	},
	{
		"word": "ala",
		"origins": {
			"language": "Georgian (ქართული)",
			"word": "არა (ara)"
		}
	},
	{
		"word": "ale",
		"origins": {
			"language": "Dutch (Nederlands)",
			"word": "alle"
		}
	},
	{
		"word": "anpa",
		"origins": {
			"language": "Acadian French (français acadien)",
			"word": "en bas"
		}
	},
	{
		"word": "ante",
		"origins": {
			"language": "Dutch (Nederlands)",
			"word": "ander"
		}
	},
	{
		"word": "anu",
		"origins": {
			"language": "Georgian (ქართული)",
			"word": "ან (an)"
		}
	},
	{
		"word": "awen",
		"origins": {
			"language": "Dutch (Nederlands)",
			"word": "houden"
		}
	},
	{
		"word": "e",
		"origins": {
			"language": "A priori.",
			"word": null
		}
	},
	{
		"word": "en",
		"origins": {
			"language": "Dutch (Nederlands)",
			"word": "en"
		}
	},
	{
		"word": "ijo",
		"origins": {
			"language": "Esperanto",
			"word" : "io"
		}
	},
	{
		"word": "ike",
		"origins": {
			"language": "Finnish (suomen kieli)",
			"word": "ilkeä"
		}
	},
	{
		"word": "ilo",
		"origins": {
			"language": "Esperanto",
			"word": "ilo"
		}
	},
	{
		"word": "insa",
		"origins": {
			"language": "Tok Pisin",
			"word": "insait"
		}
	},
	{
		"word": "jaki",
		"origins": {
			"language": "English",
			"word": "yucky"
		}
	},
	{
		"word": "jan",
		"origins": {
			"language": "Cantonese",
			"word": "人 /jɐn"
		}
	},
	{
		"word": "jelo",
		"origins": {
			"language": "English",
			"word": "yellow"
		}
	},
	{
		"word": "jo",
		"origins": {
			"language": "Mandarin",
			"word": "有 yǒu"
		}
	},
	{
		"word": "kala",
		"origins": {
			"language": "Finnish (suomen kieli)",
			"word": "kala"
		}
	},
	{
		"word": "kalama",
		"origins": {
			"language": "Croatian (hrvatski)",
			"word": "kalama"
		}
	},
	{
		"word": "kama",
		"origins": {
			"language": "Tok Pisin",
			"word": "kamap"
		}
	}
]
	`
}

func PuDict() string {
	return `[
		{
		  "word": "a",
		  "meanings": [
			[
			  "interjection",
			  "ah! ha! oh! ooh! aw! (emotion word)"
			]
		  ]
		},
		{
		  "word": "akesi",
		  "meanings": [
			[
			  "noun",
			  "non-cute animal, reptile, amphibian, dinosaur, monster"
			]
		  ]
		},
		{
		  "word": "ala",
		  "meanings": [
			[
			  "modifier",
			  "no, not, none, un-"
			],
			[
			  "noun",
			  "nothingness, negation, zero"
			],
			[
			  "interjection",
			  "no!"
			]
		  ]
		},
		{
		  "word": "alasa",
		  "meanings": [
			[
			  "modifier",
			  "hunting-, -hunting, hunting"
			],
			[
			  "noun",
			  "hunting"
			],
			[
			  "transitive verb",
			  "to hunt, to forage"
			]
		  ]
		},
		{
		  "word": "ali",
		  "meanings": [
			[
			  "noun",
			  "everything, anything, life, the universe"
			],
			[
			  "modifier",
			  "all, every, complete, whole"
			]
		  ]
		},
		{
		  "word": "ale",
		  "meanings": [
			[
			  "noun",
			  "everything, anything, life, the universe"
			],
			[
			  "modifier",
			  "all, every, complete, whole"
			]
		  ]
		},
		{
		  "word": "anpa",
		  "meanings": [
			[
			  "noun",
			  "bottom, down, under, below, floor, beneath"
			],
			[
			  "modifier",
			  "low, lower, bottom"
			],
			[
			  "transitive verb",
			  "lower, down, defeat, overcome"
			]
		  ]
		},
		{
		  "word": "ante",
		  "meanings": [
			[
			  "noun",
			  "difference"
			],
			[
			  "modifier",
			  "different"
			],
			[
			  "context (before la)",
			  "otherwise, or else"
			],
			[
			  "transitive verb",
			  "change, alter, modify"
			]
		  ]
		},
		{
		  "word": "anu",
		  "meanings": [
			[
			  "conjunction",
			  "or"
			]
		  ]
		},
		{
		  "word": "awen",
		  "meanings": [
			[
			  "intransitive verb",
			  "stay, wait, remain"
			],
			[
			  "transitive verb",
			  "keep"
			],
			[
			  "modifier",
			  "remaining, stationary, permanent, sedentary"
			]
		  ]
		},
		{
		  "word": "e",
		  "meanings": [
			[
			  "separator",
			  "(introduces a direct object)"
			]
		  ]
		},
		{
		  "word": "en",
		  "meanings": [
			[
			  "conjunction",
			  "and (used to coordinate head nouns)"
			]
		  ]
		},
		{
		  "word": "esun",
		  "meanings": [
			[
			  "noun",
			  "market, shop, fair, bazaar, business, transaction"
			],
			[
			  "modifier",
			  "commercial, trade, marketable, for sale, salable, deductible"
			],
			[
			  "transitive verb",
			  "to buy, to sell, to barter, to swap"
			]
		  ]
		},
		{
		  "word": "ijo",
		  "meanings": [
			[
			  "noun",
			  "thing, something, stuff, anything, object"
			],
			[
			  "modifier",
			  "of something"
			],
			[
			  "transitive verb",
			  "objectify"
			]
		  ]
		},
		{
		  "word": "ike",
		  "meanings": [
			[
			  "modifier",
			  "bad, evil, wrong, evil, overly complex, (figuratively) unhealthy"
			],
			[
			  "interjection",
			  "oh dear! woe! alas!"
			],
			[
			  "noun",
			  "negativity, evil, badness"
			],
			[
			  "transitive verb",
			  "to make bad, to worsen, to have a negative effect upon"
			],
			[
			  "intransitive verb",
			  "to be bad, to suck"
			]
		  ]
		},
		{
		  "word": "ilo",
		  "meanings": [
			[
			  "noun",
			  "tool, device, machine, thing used for a specific purpose"
			]
		  ]
		},
		{
		  "word": "insa",
		  "meanings": [
			[
			  "noun",
			  "inside, inner world, centre, stomach"
			],
			[
			  "modifier",
			  "internal, central"
			]
		  ]
		},
		{
		  "word": "jaki",
		  "meanings": [
			[
			  "modifier",
			  "dirty, gross, filthy"
			],
			[
			  "noun",
			  "dirt, pollution, filth"
			],
			[
			  "transitive verb",
			  "pollute, dirty"
			],
			[
			  "interjection",
			  "ew! yuck!"
			]
		  ]
		},
		{
		  "word": "jan",
		  "meanings": [
			[
			  "noun",
			  "person, people, human, being, somebody, anybody"
			],
			[
			  "modifier",
			  "personal, human, somebody's, of people"
			],
			[
			  "transitive verb",
			  "personify, humanize, personalize"
			]
		  ]
		},
		{
		  "word": "jelo",
		  "meanings": [
			[
			  "modifier",
			  "yellow, light green"
			]
		  ]
		},
		{
		  "word":"jo",
		  "meanings":[
			[
			  "transitive verb",
			  "have, contain"
			],
			[
			  "noun",
			  "having"
			],
			[
			  "compound verb (kama)",
			  "receive, get, take, obtain"
			]
		  ]
		},
		{
		  "word":"kala",
		  "meanings":[
			[
			  "noun",
			  "fish, sea creature"
			]
		  ]
		},
		{
		  "word":"kalama",
		  "meanings":[
			[
			  "noun",
			  "sound, noise, voice"
			],
			[
			  "intransitive verb",
			  "make noise"
			],
			[
			  "transitive verb",
			  "sound, ring, play (an instrument)"
			]
		  ]
		},
		{
		  "word":"kama",
		  "meanings":[
			[
			  "noun",
			  "event, happening, chance, arrival, beginning"
			],
			[
			  "modifier",
			  "coming, future"
			],
			[
			  "intransitive verb",
			  "come, become, arrive, happen, pursue actions to arrive to (a certain state), manage to, start to"
			],
			[
			  "transitive verb",
			  "bring about, summon"
			]
		  ]
		},
		{
		  "word":"kasi",
		  "meanings":[
			[
			  "noun",
			  "plant, leaf, herb, tree, wood"
			]
		  ]
		},
		{
		  "word":"ken",
		  "meanings":[
			[
			  "noun",
			  "possibility, ability, power to do things, permission"
			],
			[
			  "intransitive verb",
			  "can, is able to, is allowed to, may, is possible"
			],
			[
			  "transitive verb",
			  "make possible, enable, allow, permit"
			],
			[
			  "context (before la)",
			  "it is possible that"
			]
		  ]
		},
		{
		  "word":"kepeken",
		  "meanings":[
			[
			  "transitive verb",
			  "use"
			],
			[
			  "preposition",
			  "with"
			]
		  ]
		},
		{
		  "word":"kili",
		  "meanings":[
			[
			  "noun",
			  "fruit, pulpy vegetable, mushroom"
			]
		  ]
		},
		{
		  "word":"kin",
		  "meanings":[
			[
			  "modifier",
			  "also, too, even, indeed (emphasizes the word that follows)"
			]
		  ]
		},
		{
		  "word":"kipisi",
		  "meanings":[
			[
			  "noun (unofficial)",
			  "section, fragment, slice"
			],
			[
			  "transitive verb (unofficial)",
			  "to cut"
			]
		  ]
		},
	  
	  
		{
		  "word":"kiwen",
		  "meanings":[
			[
			  "noun",
			  "hard thing, rock, stone, metal, mineral, clay"
			],
			[
			  "modifier",
			  "hard, solid, stone-like, made of stone or metal"
			]
		  ]
		},
		{
		  "word":"ko",
		  "meanings":[
			[
			  "noun",
			  "semi-solid or squishy substance; clay, dough, glue, paste, powder, gum"
			],
			[
			  "transitive verb",
			  "to squash, to pulverize"
			]
		  ]
		},
		{
		  "word":"kon",
		  "meanings":[
			[
			  "noun",
			  "air, wind, odour, soul"
			],
			[
			  "modifier",
			  "air-like, ethereal, gaseous"
			]
		  ]
		},
		{
		  "word":"kule",
		  "meanings":[
			[
			  "noun",
			  "color, paint"
			],
			[
			  "transitive verb",
			  "color, paint"
			],
			[
			  "modifier",
			  "colourful"
			]
		  ]
		},
		{
		  "word":"kulupu",
		  "meanings":[
			[
			  "noun",
			  "group, community, society, company"
			],
			[
			  "modifier",
			  "communal, shared, public, of the society"
			]
		  ]
		},
		{
		  "word":"kute",
		  "meanings":[
			[
			  "transitive verb",
			  "listen, hear"
			],
			[
			  "modifier",
			  "auditory, hearing"
			]
		  ]
		},
		{
		  "word":"la",
		  "meanings":[
			[
			  "separator",
			  "(between adverb or phrase of context and sentence)"
			]
		  ]
		},
		{
		  "word":"lape",
		  "meanings":[
			[
			  "n, vi",
			  "sleep, rest"
			],
			[
			  "modifier",
			  "sleeping, of sleep laso mod blue, blue-green"
			]
		  ]
		},
		{
		  "word":"laso",
		  "meanings":[
			[
			  "modifier",
			  "blue, blue-green"
			]
		  ]
		},
		{
		  "word":"lawa",
		  "meanings":[
			[
			  "noun",
			  "head, mind"
			],
			[
			  "modifier",
			  "main, leading, in charge"
			],
			[
			  "transitive verb",
			  "lead, control, rule, steer"
			]
		  ]
		},
		{
		  "word":"len",
		  "meanings":[
			[
			  "noun",
			  "clothing, cloth, fabric"
			]
		  ]
		},
		{
		  "word":"lete",
		  "meanings":[
			[
			  "noun",
			  "cold"
			],
			[
			  "modifier",
			  "cold, uncooked"
			],
			[
			  "transitive verb",
			  "cool down, chill, freeze"
			]
		  ]
		},
		{
		  "word":"li",
		  "meanings":[
			[
			  "separator",
			  "(between any subject except mi and sina and its verb; also used to introduce a new verb for the same subject)"
			]
		  ]
		},
		{
		  "word":"lili",
		  "meanings":[
			[
			  "modifier",
			  "small, little, young, a bit, short, fes, less"
			],
			[
			  "transitive verb",
			  "reduce, shorten, shrink, lessen"
			]
		  ]
		},
		{
		  "word":"linja",
		  "meanings":[
			[
			  "noun",
			  "long, very thin, floppy thing, e.g. string, rope, hair, thread, cord, chain"
			]
		  ]
		},
		{
		  "word":"lipu",
		  "meanings":[
			[
			  "noun",
			  "flat and bendable thing, e.g. paper, card, ticket"
			]
		  ]
		},
		{
		  "word":"loje",
		  "meanings":[
			[
			  "modifier",
			  "red"
			]
		  ]
		},
		{
		  "word":"lon",
		  "meanings":[
			[
			  "preposition",
			  "be located in/at/on"
			],
			[
			  "intransitive verb",
			  "be there, be present, be real/true, exist, be awake"
			]
		  ]
		},
		{
		  "word":"luka",
		  "meanings":[
			[
			  "noun",
			  "hand, arm"
			],
			[
			  "modifier",
			  "five"
			]
		  ]
		},
		{
		  "word":"lukin",
		  "meanings":[
			[
			  "transitive verb",
			  "see, look at, watch, read"
			],
			[
			  "intransitive verb",
			  "look, watch out, pay attention"
			],
			[
			  "modifier",
			  "visual, to seek to, try to, look for"
			]
		  ]
		},
		{
		  "word":"lupa",
		  "meanings":[
			[
			  "noun",
			  "hole, orifice, window, door"
			]
		  ]
		},
		{
		  "word":"ma",
		  "meanings":[
			[
			  "noun",
			  "land, earth, country, (outdoor) area"
			]
		  ]
		},
		{
		  "word":"mama",
		  "meanings":[
			[
			  "noun",
			  "parent, mother, father"
			],
			[
			  "modifier",
			  "of the parent, parental, maternal, fatherly"
			]
		  ]
		},
		{
		  "word":"mani",
		  "meanings":[
			[
			  "noun",
			  "money, material wealth, currency, dollar, capital"
			]
		  ]
		},
		{
		  "word":"meli",
		  "meanings":[
			[
			  "noun",
			  "woman, female, girl, wife, girlfriend"
			],
			[
			  "modifier",
			  "female, feminine, womanly"
			]
		  ]
		},
		{
		  "word":"mi",
		  "meanings":[
			[
			  "noun",
			  "I, we"
			],
			[
			  "modifier",
			  "my, our"
			]
		  ]
		},
		{
		  "word":"mije",
		  "meanings":[
			[
			  "noun",
			  "man, male, husband, boyfriend"
			],
			[
			  "modifier",
			  "male, masculine, manly"
			]
		  ]
		},
		{
		  "word":"moku",
		  "meanings":[
			[
			  "noun",
			  "food, meal"
			],
			[
			  "transitive verb",
			  "eat, drink, swallow, ingest, consume"
			]
		  ]
		},
		{
		  "word":"moli",
		  "meanings":[
			[
			  "noun",
			  "death"
			],
			[
			  "modifier",
			  "dead, deadly, fatal"
			],
			[
			  "transitive verb",
			  "kill"
			],
			[
			  "intransitive verb",
			  "die, be dead"
			]
		  ]
		},
		{
		  "word":"monsi",
		  "meanings":[
			[
			  "noun",
			  "back, rear end, butt, behind"
			],
			[
			  "modifier",
			  "back, rear"
			]
		  ]
		},
		{
		  "word":"monsuta",
		  "meanings":[
			[
			  "noun (unofficial)",
			  "monster, monstrosity, fearful thing, fright, mythical creatures, fear"
			]
		  ]
		},
		{
		  "word":"mu",
		  "meanings":[
			[
			  "interjection",
			  "woof! meow! moo! etc. (cute animal noise)"
			]
		  ]
		},
		{
		  "word":"mun",
		  "meanings":[
			[
			  "noun",
			  "moon"
			],
			[
			  "modifier",
			  "lunar"
			]
		  ]
		},
		{
		  "word":"musi",
		  "meanings":[
			[
			  "noun",
			  "fun, playing, game, recreation, art, entertainment"
			],
			[
			  "modifier",
			  "artful, fun, recreational"
			],
			[
			  "intransitive verb",
			  "play, have fun"
			],
			[
			  "transitive verb",
			  "amuse, entertain"
			]
		  ]
		},
		{
		  "word":"mute",
		  "meanings":[
			[
			  "modifier",
			  "many, several, very, much, a lot, abundant, numerous, more"
			],
			[
			  "noun",
			  "amount, quantity"
			],
			[
			  "transitive verb",
			  "make many or much"
			]
		  ]
		},
		{
		  "word":"namako",
		  "meanings":[
			[
			  "noun",
			  "spice, something extra, food additive, accessory"
			],
			[
			  "transitive verb",
			  "to spice, to flavor, to decorate"
			],
			[
			  "modifier",
			  "spicy, piquant"
			]
		  ]
		},
		{
		  "word":"nanpa",
		  "meanings":[
			[
			  "noun",
			  "number"
			],
			[
			  "other",
			  "-th (ordinal numbers)"
			]
		  ]
		},
		{
		  "word":"nasa",
		  "meanings":[
			[
			  "modifier",
			  "silly, crazy, foolish, drunk, strange, stupid, weird"
			],
			[
			  "transitive verb",
			  "drive crazy, make weird"
			]
		  ]
		},
		{
		  "word":"nasin",
		  "meanings":[
			[
			  "noun",
			  "way, manner, custom, road, path, doctrine, system, method"
			]
		  ]
		},
		{
		  "word":"nena",
		  "meanings":[
			[
			  "noun",
			  "bump, hill, mountain, button, nose"
			]
		  ]
		},
		{
		  "word":"ni",
		  "meanings":[
			[
			  "modifier",
			  "this, that"
			]
		  ]
		},
		{
		  "word":"nimi",
		  "meanings":[
			[
			  "noun",
			  "word, name"
			]
		  ]
		},
		{
		  "word":"noka",
		  "meanings":[
			[
			  "noun",
			  "leg, foot"
			]
		  ]
		},
		{
		  "word":"o",
		  "meanings":[
			[
			  "separator",
			  "O (vocative or imperative)"
			],
			[
			  "interjection",
			  "hey! (calling somebody's attention)"
			]
		  ]
		},
		{
		  "word":"oko",
		  "meanings":[
			[
			  "noun",
			  "eye"
			]
		  ]
		},
		{
		  "word":"olin",
		  "meanings":[
			[
			  "noun",
			  "love"
			],
			[
			  "modifier",
			  "love"
			],
			[
			  "transitive verb",
			  "to love (a person)"
			]
		  ]
		},
		{
		  "word":"ona",
		  "meanings":[
			[
			  "noun",
			  "he, she, it, they"
			],
			[
			  "modifier",
			  "his, her, its, their"
			]
		  ]
		},
		{
		  "word":"open",
		  "meanings":[
			[
			  "transitive verb",
			  "open, turn on"
			]
		  ]
		},
		{
		  "word":"pakala",
		  "meanings":[
			[
			  "noun",
			  "blunder, accident, mistake, destruction, damage, breaking"
			],
			[
			  "transitive verb",
			  "screw up, fuck up, botch, ruin, break, hurt, injure, damage, bungle, spoil, ruin"
			],
			[
			  "intransitive verb",
			  "screw up, fall apart, break"
			],
			[
			  "interjection",
			  "damn! fuck!"
			]
		  ]
		},
		{
		  "word":"pali",
		  "meanings":[
			[
			  "noun",
			  "activity, work, deed, projec"
			],
			[
			  "modifier",
			  "active, work-related, operating, working"
			],
			[
			  "transitive verb",
			  "do, make, build, create"
			],
			[
			  "intransitive verb",
			  "act, work, function"
			]
		  ]
		},
		{
		  "word":"palisa",
		  "meanings":[
			[
			  "noun",
			  "long, mostly hard object, e.g. rod, stick, branch"
			]
		  ]
		},
		{
		  "word":"pana",
		  "meanings":[
			[
			  "transitive verb",
			  "give, put, send, place, release, emit, cause"
			],
			[
			  "noun",
			  "giving, transfer, exchange"
			]
		  ]
		},
		{
		  "word":"pi",
		  "meanings":[
			[
			  "separator",
			  "of, belonging to"
			]
		  ]
		},
		{
		  "word":"pilin",
		  "meanings":[
			[
			  "noun",
			  "feelings, emotion, heart"
			],
			[
			  "intransitive verb",
			  "feel"
			],
			[
			  "transitive verb",
			  "feel, think, sense, touch"
			]
		  ]
		},
		{
		  "word":"pimeja",
		  "meanings":[
			[
			  "modifier",
			  "black, dark"
			],
			[
			  "noun",
			  "darkness, shadows"
			],
			[
			  "transitive verb",
			  "darken"
			]
		  ]
		},
		{
		  "word":"pini",
		  "meanings":[
			[
			  "noun",
			  "end, tip"
			],
			[
			  "modifier",
			  "completed, finished, past, done, ago"
			],
			[
			  "transitive verb",
			  "finish, close, end, turn off"
			]
		  ]
		},
		{
		  "word":"pipi",
		  "meanings":[
			[
			  "noun",
			  "bug, insect, spider"
			]
		  ]
		},
		{
		  "word":"poka",
		  "meanings":[
			[
			  "noun",
			  "side, hip, next to"
			],
			[
			  "modifier",
			  "neighboring"
			],
			[
			  "preposition",
			  "in the accompaniment of, with"
			]
		  ]
		},
		{
		  "word":"poki",
		  "meanings":[
			[
			  "noun",
			  "container, box, bowl, cup, glass"
			]
		  ]
		},
		{
		  "word":"pona",
		  "meanings":[
			[
			  "noun",
			  "good, simplicity, positivity"
			],
			[
			  "modifier",
			  "good, simple, positive, nice, correct, right"
			],
			[
			  "interjection",
			  "great! good! thanks! OK! cool! yay!"
			],
			[
			  "transitive verb",
			  "improve, fix, repair, make good"
			]
		  ]
		},
		{
		  "word":"pu",
		  "meanings":[
			[
			  "noun",
			  "buying and interacting with the official Toki Pona book"
			],
			[
			  "modifier",
			  "buying and interacting with the official Toki Pona book"
			],
			[
			  "transitive verb",
			  "to apply (the official Toki Pona book) to…"
			],
			[
			  "intransitive verb",
			  "to buy and to read (the official Toki Pona book)"
			]
		  ]
		},
		{
		  "word":"sama",
		  "meanings":[
			[
			  "modifier",
			  "same, similar, equal, of equal status or position"
			],
			[
			  "preposition",
			  "like, as, seem"
			],
			[
			  "context (before la)",
			  "similarly, in the same way that"
			]
		  ]
		},
		{
		  "word":"seli",
		  "meanings":[
			[
			  "noun",
			  "fire, warmth, heat"
			],
			[
			  "modifier",
			  "hot, warm, cooked"
			],
			[
			  "transitive verb",
			  "heat, warm up, cook"
			]
		  ]
		},
		{
		  "word":"selo",
		  "meanings":[
			[
			  "noun",
			  "outside, surface, skin, shell, bark, shape, peel"
			]
		  ]
		},
		{
		  "word":"seme",
		  "meanings":[
			[
			  "other",
			  "what, which, wh- (question word)"
			]
		  ]
		},
		{
		  "word":"sewi",
		  "meanings":[
			[
			  "noun",
			  "high, up, above, top, over, on"
			],
			[
			  "modifier",
			  "superior, elevated, religious, formal"
			]
		  ]
		},
		{
		  "word":"sijelo",
		  "meanings":[
			[
			  "noun",
			  "body, physical state"
			]
		  ]
		},
		{
		  "word":"sike",
		  "meanings":[
			[
			  "noun",
			  "circle, wheel, sphere, ball, cycle"
			],
			[
			  "modifier",
			  "round, cyclical"
			]
		  ]
		},
		{
		  "word":"sin",
		  "meanings":[
			[
			  "modifier",
			  "new, fresh, another, more"
			],
			[
			  "transitive verb",
			  "renew, renovate, freshen"
			]
		  ]
		},
		{
		  "word":"sina",
		  "meanings":[
			[
			  "noun",
			  "you"
			],
			[
			  "modifier",
			  "your"
			]
		  ]
		},
		{
		  "word":"sinpin",
		  "meanings":[
			[
			  "noun",
			  "front, face, chest, torso, wall"
			]
		  ]
		},
		{
		  "word":"sitelen",
		  "meanings":[
			[
			  "noun",
			  "picture, image"
			],
			[
			  "transitive verb",
			  "draw, write"
			]
		  ]
		},
		{
		  "word":"sona",
		  "meanings":[
			[
			  "noun",
			  "knowledge, wisdom, intelligence, understanding"
			],
			[
			  "transitive verb",
			  "know, understand, know how to"
			],
			[
			  "intransitive verb",
			  "know, understand"
			],
			[
			  "compound verb (kama)",
			  "learn, study"
			]
		  ]
		},
		{
		  "word":"soweli",
		  "meanings":[
			[
			  "noun",
			  "animal, especially land mammal, lovable animal"
			]
		  ]
		},
		{
		  "word":"suli",
		  "meanings":[
			[
			  "modifier",
			  "big, tall, long, adult, important"
			],
			[
			  "transitive verb",
			  "enlarge, lengthen"
			],
			[
			  "noun",
			  "size"
			]
		  ]
		},
		{
		  "word":"suno",
		  "meanings":[
			[
			  "noun",
			  "sun, light"
			]
		  ]
		},
		{
		  "word":"supa",
		  "meanings":[
			[
			  "noun",
			  "horizontal surface, e.g furniture, table, chair, pillow, floor"
			]
		  ]
		},
		{
		  "word":"suwi",
		  "meanings":[
			[
			  "noun",
			  "candy, sweet food"
			],
			[
			  "modifier",
			  "sweet, cute"
			],
			[
			  "transitive verb",
			  "sweeten"
			]
		  ]
		},
		{
		  "word":"tan",
		  "meanings":[
			[
			  "preposition",
			  "from, by, because of, since"
			],
			[
			  "noun",
			  "origin, cause"
			]
		  ]
		},
		{
		  "word":"taso",
		  "meanings":[
			[
			  "modifier",
			  "only, sole"
			],
			[
			  "conjunction",
			  "but"
			]
		  ]
		},
		{
		  "word":"tawa",
		  "meanings":[
			[
			  "preposition",
			  "to, in order to, towards, for, until"
			],
			[
			  "noun",
			  "movement, transportation"
			],
			[
			  "modifier",
			  "moving, mobile"
			],
			[
			  "intransitive verb",
			  "go, leave, walk, travel, move"
			],
			[
			  "transitive verb",
			  "move, displace"
			]
		  ]
		},
		{
		  "word":"telo",
		  "meanings":[
			[
			  "noun",
			  "water, liquid, juice, sauce"
			],
			[
			  "transitive verb",
			  "water, wash with water"
			]
		  ]
		},
		{
		  "word":"tenpo",
		  "meanings":[
			[
			  "noun",
			  "time, period of time, moment, duration, situation"
			]
		  ]
		},
		{
		  "word":"toki",
		  "meanings":[
			[
			  "noun",
			  "language, talking, speech, communication"
			],
			[
			  "transitive verb",
			  "say"
			],
			[
			  "intransitive verb",
			  "talk, chat, communicate"
			],
			[
			  "interjection",
			  "hello! hi!"
			]
		  ]
		},
		{
		  "word":"tomo",
		  "meanings":[
			[
			  "noun",
			  "indoor constructed space, e.g. house, home, room, building"
			],
			[
			  "modifier",
			  "urban, domestic, household"
			]
		  ]
		},
		{
		  "word":"tu",
		  "meanings":[
			[
			  "modifier",
			  "two"
			],
			[
			  "noun",
			  "duo, pair"
			],
			[
			  "transitive verb",
			  "double, separate/cut/divide in two"
			]
		  ]
		},
		{
		  "word":"unpa",
		  "meanings":[
			[
			  "noun",
			  "sex, sexuality"
			],
			[
			  "modifier",
			  "erotic, sexual"
			],
			[
			  "transitive verb",
			  "have sex with, sleep with, fuck"
			],
			[
			  "intransitive verb",
			  "have sex"
			]
		  ]
		},
		{
		  "word":"uta",
		  "meanings":[
			[
			  "noun",
			  "mouth"
			],
			[
			  "modifier",
			  "oral"
			]
		  ]
		},
		{
		  "word":"utala",
		  "meanings":[
			[
			  "noun",
			  "conflict, disharmony, competition, fight, war, battle, attack, blow, argument, physical or verbal violence"
			],
			[
			  "transitive verb",
			  "hit, strike, attack, compete against"
			]
		  ]
		},
		{
		  "word":"walo",
		  "meanings":[
			[
			  "modifier",
			  "white, light (colour)"
			],
			[
			  "noun",
			  "white thing/part, whiteness, lightness"
			]
		  ]
		},
		{
		  "word":"wan",
		  "meanings":[
			[
			  "modifier",
			  "one, a"
			],
			[
			  "noun",
			  "unit, element, particle, part, piece"
			],
			[
			  "transitive verb",
			  "unite, make one"
			]
		  ]
		},
		{
		  "word":"waso",
		  "meanings":[
			[
			  "noun",
			  "bird, winged animal"
			]
		  ]
		},
		{
		  "word":"wawa",
		  "meanings":[
			[
			  "noun",
			  "energy, strength, power"
			],
			[
			  "modifier",
			  "energetic, strong, fierce, intense, sure, confident"
			],
			[
			  "transitive verb",
			  "strengthen, energize, empower"
			]
		  ]
		},
		{
		  "word":"weka",
		  "meanings":[
			[
			  "modifier",
			  "away, absent, missing"
			],
			[
			  "noun",
			  "absence"
			],
			[
			  "transitive verb",
			  "throw away, remove, get rid of"
			]
		  ]
		},
		{
		  "word":"wile",
		  "meanings":[
			[
			  "transitive verb",
			  "to want, need, wish, have to, must, will, should"
			],
			[
			  "noun",
			  "desire, need, will"
			],
			[
			  "modifier",
			  "necessary"
			]
		  ]
		}
	  ]`
}

type wordOrigin struct {
	Language   string `json:"language"`
	SourceWord string `json:"word"`
}

type puEty struct {
	Word   string     `json:"word"`
	Origin wordOrigin `json:"origins"`
}

type puWord struct {
	Word     string     `json:"word"`
	Meanings [][]string `json:"meanings"`
}

var pu []puWord
var puEtym []puEty

func init() {
	err := json.Unmarshal([]byte(PuDict()), &pu)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(EtyDict()), &puEtym)
	if err != nil {
		panic(err)
	}
	commandlib.RegisterCommand(commandlib.Command{
		Name:  "toki pona vocabulary lookup",
		Usage: "Look up words in pu",
		ID:    "pu-lookup",
		Match: [][]string{
			{"o", "pu"},
		},
		Flags: []commandlib.Flag{
			commandlib.BoolFlag{
				LongFlag: "browse",
			},
		},
		Action: PuSearch,
	})
	commandlib.RegisterCommand(commandlib.Command{
		Name:  "toki pona etymolygy lookup",
		Usage: "Look up words in etymolygy",
		ID:    "ety-lookup",
		Match: [][]string{
			{"o", "ety"},
		},
		Flags: []commandlib.Flag{
			commandlib.BoolFlag{
				LongFlag: "browse",
			},
		},
		Action: EtymologySearch,
	})
	commandlib.RegisterCommand(commandlib.Command{
		Name:  "toki pona quiz",
		Usage: "Quiz yourself on toki pona",
		ID:    "tpo-quiz",
		Match: [][]string{
			{"o", "quiz"},
		},
		Action: Quiz,
	})
}

func lookupPu(word string) (puWord, bool) {
	for _, val := range pu {
		if val.Word == word {
			return val, true
		}
	}
	return puWord{}, false
}

func lookupEty(word string) (puEty, bool) {
	for _, val := range puEtym {
		if val.Word == word {
			return val, true
		}
	}
	return puEty{}, false
}

func toEmbed(v interface{}) commandlib.Embed {
	switch v.(type) {
	case puWord:
		word := v.(puWord)
		embed := commandlib.Embed{}
		embed.Title.Text = word.Word
		for _, meaning := range word.Meanings {
			embed.Fields = append(embed.Fields, commandlib.EmbedField{
				Title: meaning[0],
				Body:  meaning[1],
			})
		}
		return embed
	case puEty:
		ety := v.(puEty)
		embed := commandlib.Embed{}
		embed.Title.Text = ety.Word
		if ety.Origin.SourceWord != "" {
			embed.Fields = append(embed.Fields, commandlib.EmbedField{
				Title: ety.Origin.Language,
				Body:  ety.Origin.SourceWord,
			})
		} else {
			embed.Body = ety.Origin.Language
		}
		return embed
	}
	return commandlib.Embed{}
}

func EtymologySearch(c commandlib.Context) {
	if c.IsFlagSet("browse") {
		embedlist := commandlib.EmbedList{
			ItemTypeName: c.I18n("Word"),
		}
		for _, word := range puEtym {
			embedlist.Embeds = append(embedlist.Embeds, toEmbed(word))
		}
		c.SendMessage("primary", embedlist)
		return
	}
	if c.NArgs() < 1 {
		c.SendMessage(
			"primary",
			commandlib.ErrorEmbed(c.I18n("Please provide a word from pu to look up the etymolygy of.")),
		)
		return
	}
	if word, ok := lookupEty(c.Arg(0)); ok {
		c.SendMessage("primary", toEmbed(word))
		return
	} else {
		c.SendMessage(
			"primary",
			commandlib.ErrorEmbed(c.I18n("Sorry, I don't think that's a word.")),
		)
		return
	}
}

func PuSearch(c commandlib.Context) {
	if c.IsFlagSet("browse") {
		embedlist := commandlib.EmbedList{
			ItemTypeName: c.I18n("Word"),
		}
		for _, word := range pu {
			embedlist.Embeds = append(embedlist.Embeds, toEmbed(word))
		}
		c.SendMessage("primary", embedlist)
		return
	}
	if c.NArgs() < 1 {
		c.SendMessage(
			"primary",
			commandlib.ErrorEmbed(c.I18n("Please provide a word from pu to look up the meaning of.")),
		)
		return
	}
	if word, ok := lookupPu(c.Arg(0)); ok {
		c.SendMessage("primary", toEmbed(word))
		return
	} else {
		c.SendMessage(
			"primary",
			commandlib.ErrorEmbed(c.I18n("Sorry, I don't think that's a word.")),
		)
		return
	}
}

const nextRespDuration = 3 * time.Second

type meaningWordList []string

func (m meaningWordList) contains(s string) bool {
	lower := strings.ToLower(s)
	for _, word := range strings.Fields(lower) {
		for _, meaning := range m {
			if strings.TrimSpace(word) == strings.TrimSpace(meaning) {
				return true
			}
		}
	}
	return false
}

func (p puWord) toMeaningWordList() (ret meaningWordList) {
	for _, meaning := range p.Meanings {
		for _, split := range strings.Fields(strings.ReplaceAll(strings.ReplaceAll(meaning[1], ",", " "), "!", " ")) {
			ret = append(ret, strings.ToLower(split))
		}
	}
	return
}

func Quiz(c commandlib.Context) {
	score := 0
	rand.Seed(time.Now().Unix())
	c.SendMessage("starting", fmt.Sprintf(c.I18n("Starting quiz... Type '%s' to cancel"), c.I18n("cancel")))
Wait:
	for i := 0; i < 10; i++ {
		word := pu[rand.Intn(len(pu))]
		c.SendMessage(fmt.Sprintf("primary-%d", i), commandlib.Embed{
			Title: commandlib.EmbedHeader{
				Text: word.Word,
			},
			Body: c.I18n("What does this word mean?"),
			Footer: commandlib.EmbedHeader{
				Text: fmt.Sprintf("Word %d out of %d", i+1, 10),
			},
		})
		timeoutChan := make(chan struct{})
		go func() {
			time.Sleep(7 * time.Second)
			timeoutChan <- struct{}{}
		}()
		for {
			select {
			case msg := <-c.NextResponse():
				if strings.Contains(msg, c.I18n("cancel")) {
					c.SendMessage(fmt.Sprintf("primary-%d", i), commandlib.ErrorEmbed(c.I18n("Quiz cancelled.")))
					return
				}
				if word.toMeaningWordList().contains(msg) {
					grats := toEmbed(word)
					grats.Colour = 0x00ff00
					c.SendMessage(fmt.Sprintf("primary-%d", i), grats)
					c.SendMessage(fmt.Sprintf("congrats-%d", i), "Correct!")
					score++
					time.Sleep(nextRespDuration)
					continue Wait
				}
			case <-timeoutChan:
				wrong := toEmbed(word)
				wrong.Colour = 0xff0000
				c.SendMessage(fmt.Sprintf("primary-%d", i), wrong)
				time.Sleep(nextRespDuration)
				continue Wait
			}
		}
	}
	c.SendMessage("primary", commandlib.Embed{
		Title: commandlib.EmbedHeader{
			Text: c.I18n("Quiz Results"),
		},
		Fields: []commandlib.EmbedField{
			{
				Title: c.I18n("Correct Results"),
				Body:  strconv.Itoa(score),
			},
			{
				Title: c.I18n("Incorrect Results"),
				Body:  strconv.Itoa(10 - score),
			},
		},
	})
}
