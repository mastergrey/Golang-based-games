package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

const simulations = 100000

type Card struct {
	Rank int // 2–14 (Ace = 14)
	Suit int // 0–3
}

type Hand []Card

func main() {
	rand.Seed(time.Now().UnixNano())

	// Example input (can be replaced by user input)
	players := []Hand{
		{{14, 0}, {14, 1}}, // AA
		{{13, 2}, {13, 3}}, // KK
	}

	board := Hand{
		{2, 0}, {7, 1}, {9, 2}, // Flop
		// Turn, River unknown
	}

	calcOdds(players, board)
}

func calcOdds(players []Hand, board Hand) {
	wins := make([]int, len(players))
	ties := make([]int, len(players))

	for i := 0; i < simulations; i++ {
		deck := newDeck()
		removeKnown(&deck, players, board)

		shuffle(deck)

		// Complete board
		fullBoard := append(Hand{}, board...)
		fullBoard = append(fullBoard, deck[:5-len(board)]...)
		deck = deck[5-len(board):]

		bestRanks := make([]int, len(players))

		for p := range players {
			hand := append(players[p], fullBoard...)
			bestRanks[p] = evaluate(hand)
		}

		max := maxInt(bestRanks)
		count := 0
		for _, r := range bestRanks {
			if r == max {
				count++
			}
		}

		for p := range bestRanks {
			if bestRanks[p] == max {
				if count > 1 {
					ties[p]++
				} else {
					wins[p]++
				}
			}
		}
	}

	fmt.Println("Results:")
	for i := range players {
		fmt.Printf(
			"Player %d | Win: %.2f%% | Tie: %.2f%%\n",
			i+1,
			100*float64(wins[i])/simulations,
			100*float64(ties[i])/simulations,
		)
	}
}

func newDeck() []Card {
	deck := make([]Card, 0, 52)
	for s := 0; s < 4; s++ {
		for r := 2; r <= 14; r++ {
			deck = append(deck, Card{r, s})
		}
	}
	return deck
}

func removeKnown(deck *[]Card, players []Hand, board Hand) {
	known := append(board, flatten(players)...)

	filtered := (*deck)[:0]
	for _, c := range *deck {
		if !contains(known, c) {
			filtered = append(filtered, c)
		}
	}
	*deck = filtered
}

func flatten(players []Hand) []Card {
	var all []Card
	for _, p := range players {
		all = append(all, p...)
	}
	return all
}

func shuffle(deck []Card) {
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
}

/*
HAND RANKING (higher is better):
8 = Straight Flush
7 = Four of a Kind
6 = Full House
5 = Flush
4 = Straight
3 = Three of a Kind
2 = Two Pair
1 = One Pair
0 = High Card
*/
func evaluate(cards []Card) int {
	best := 0
	comb := combinations(cards, 5)

	for _, c := range comb {
		r := rankHand(c)
		if r > best {
			best = r
		}
	}
	return best
}

func rankHand(hand []Card) int {
	sort.Slice(hand, func(i, j int) bool {
		return hand[i].Rank > hand[j].Rank
	})

	isFlush := true
	for i := 1; i < 5; i++ {
		if hand[i].Suit != hand[0].Suit {
			isFlush = false
			break
		}
	}

	isStraight := true
	for i := 1; i < 5; i++ {
		if hand[i].Rank != hand[0].Rank-i {
			isStraight = false
			break
		}
	}

	counts := make(map[int]int)
	for _, c := range hand {
		counts[c.Rank]++
	}

	switch {
	case isStraight && isFlush:
		return 8
	case hasN(counts, 4):
		return 7
	case hasN(counts, 3) && hasN(counts, 2):
		return 6
	case isFlush:
		return 5
	case isStraight:
		return 4
	case hasN(counts, 3):
		return 3
	case pairs(counts) == 2:
		return 2
	case pairs(counts) == 1:
		return 1
	default:
		return 0
	}
}

func hasN(m map[int]int, n int) bool {
	for _, v := range m {
		if v == n {
			return true
		}
	}
	return false
}

func pairs(m map[int]int) int {
	c := 0
	for _, v := range m {
		if v == 2 {
			c++
		}
	}
	return c
}

func combinations(cards []Card, k int) [][]Card {
	var res [][]Card
	var comb []Card

	var dfs func(int, int)
	dfs = func(i, left int) {
		if left == 0 {
			tmp := make([]Card, len(comb))
			copy(tmp, comb)
			res = append(res, tmp)
			return
		}
		for j := i; j <= len(cards)-left; j++ {
			comb = append(comb, cards[j])
			dfs(j+1, left-1)
			comb = comb[:len(comb)-1]
		}
	}

	dfs(0, k)
	return res
}

func contains(arr []Card, c Card) bool {
	for _, v := range arr {
		if v == c {
			return true
		}
	}
	return false
}

func maxInt(arr []int) int {
	m := arr[0]
	for _, v := range arr {
		if v > m {
			m = v
		}
	}
	return m
}
