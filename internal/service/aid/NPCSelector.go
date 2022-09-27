package aid

import "math/rand"

type NPCSelector struct {
	npcPhones []string
}

func newNPCSelector() *NPCSelector {
	return &NPCSelector{
		npcPhones: []string{
			"15548720906",
			"18349162361",
			"13693520907",
			"18512303122",
			"18487540201",
			"13610239417",
			"16602833635",
			"15110007178",
			"17150305361",
			"18516597958",
			"13616067770",
		},
	}
}

func (g NPCSelector) RandomPhone() string {
	return g.npcPhones[rand.Int()%len(g.npcPhones)]
}
