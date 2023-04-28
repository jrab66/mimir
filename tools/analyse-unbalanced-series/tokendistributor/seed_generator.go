package tokendistributor

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"golang.org/x/exp/slices"
)

type SeedGenerator interface {
	generateSeedByZone(zones []Zone, tokensPerInstance int, maxTokenValue Token) map[Zone][]Token
}

type TestSeedGenerator struct{}

func (g TestSeedGenerator) generateSeedByZone(zones []Zone, tokensPerInstance int, maxTokenValue Token) map[Zone][]Token {
	allSeeds := []Token{17, 101, 185, 269, 353, 437, 513, 597, 681, 765, 849, 933}
	slices.Sort(allSeeds)
	fmt.Println(allSeeds)

	return distributeSortedTokensByZone(allSeeds, zones, tokensPerInstance)
}

type PerfectlySpacedSeedGenerator struct{}

func (g PerfectlySpacedSeedGenerator) generateSeedByZone(zones []Zone, tokensPerInstance int, maxTokenValue Token) map[Zone][]Token {
	allSeeds := make([]Token, 0, tokensPerInstance*len(zones))
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tokensCount := tokensPerInstance * len(zones)
	offset := uint32(math.Ceil(float64(maxTokenValue) / float64(tokensCount)))
	curr := uint32(r.Intn(int(maxTokenValue)))
	for i := 0; i < tokensCount; i++ {
		allSeeds = append(allSeeds, Token(curr))
		if uint32(maxTokenValue)-offset < curr {
			curr -= uint32(maxTokenValue) - offset
		} else {
			curr += offset
		}
	}

	slices.Sort(allSeeds)
	fmt.Println(allSeeds)

	return distributeSortedTokensByZone(allSeeds, zones, tokensPerInstance)
}

type RandomSeedGenerator struct{}

func (g RandomSeedGenerator) generateSeedByZone(zones []Zone, tokensPerInstance int, maxTokenValue Token) map[Zone][]Token {
	allSeeds := make([]Token, 0, tokensPerInstance*len(zones))
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tokensCount := tokensPerInstance * len(zones)
	for {
		token := Token(r.Intn(int(maxTokenValue)))
		if slices.Contains(allSeeds, token) {
			continue
		}
		allSeeds = append(allSeeds, token)
		if len(allSeeds) == tokensCount {
			break
		}
	}

	slices.Sort(allSeeds)

	return distributeSortedTokensByZone(allSeeds, zones, tokensPerInstance)
}

func distributeSortedTokensByZone(allSeeds []Token, zones []Zone, tokensPerInstance int) map[Zone][]Token {
	seedByZone := make(map[Zone][]Token, len(zones))
	for i, token := range allSeeds {
		zone := zones[i%len(zones)]
		currSeed, ok := seedByZone[zone]
		if !ok {
			currSeed = make([]Token, 0, tokensPerInstance)
		}
		currSeed = append(currSeed, token)
		seedByZone[zone] = currSeed
	}
	return seedByZone
}
