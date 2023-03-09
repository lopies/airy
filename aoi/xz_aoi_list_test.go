package aoi

import (
	"fmt"
	"github.com/airy/player"
	"math/rand"
	"testing"
)

var xzm = NewXZListAOIManager()
var players = make([]*player.Player, 0, 20)

func TestAddPlayer(t *testing.T) {
	p := player.NewPlayer(uint32(0), "")
	xzm.Enter(p, 0, 0)
	players = append(players, p)
	for i := 1; i < 20; i++ {
		p := player.NewPlayer(uint32(i), "")
		players = append(players, p)
		xzm.Enter(p, rand.Float32()*100, rand.Float32()*100)
	}
	xzm.String()
}

func TestMovePlayer(t *testing.T) {
	TestAddPlayer(t)
	xzm.Move(players[0], 500, 500)
	xzm.String()
}

func TestExit(t *testing.T) {
	TestAddPlayer(t)
	xzm.Leave(players[0])
	xzm.String()
}

func TestNeighbor(t *testing.T) {
	TestAddPlayer(t)
	for i := 0; i < len(players); i++ {
		fmt.Printf("player%d : X = %f, Z= %f\n", i, players[i].X, players[i].Z)
		neighbors := xzm.Neighbors(players[i])
		for p, _ := range neighbors {
			fmt.Printf("neighbor: X = %f, Z = %f\n", p.X, p.Z)
		}
		fmt.Println()
		fmt.Println()
	}
}
