package snowflake

import (
	cryptorand "crypto/rand"
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"math/big"
	"net"
	"strings"

	"github.com/bwmarrin/snowflake"
)

const (
	defaultNodeID int64 = 1
	maxNodeID     int64 = 1<<10 - 1
)

var node *snowflake.Node

func init() {
	// Custom Epoch (Sunday, December 1, 2019 1:00:00 AM)
	snowflake.Epoch = 1575162000000

	var err error
	node, err = snowflake.NewNode(createNodeID())
	if err != nil {
		log.Println("Error on initing a new snowflake node:", err.Error())
	}
}

// NextID generates an unique int64 of the snowflake ID
// Ex: 587159296175034368
func NextID() (int64, error) {
	if node == nil {
		return 0, errors.New("snowflake: fail to init snowflake node")
	}

	return node.Generate().Int64(), nil
}

func createNodeID() int64 {
	var sb strings.Builder

	ifaces, err := net.Interfaces()
	if err != nil {
		// valid node ID: random 0 -> 1023
		return randomValidNodeID()
	}

	// calculate node ID from hardware addresses
	for _, iface := range ifaces {
		var mac []byte = iface.HardwareAddr
		if len(mac) > 0 {
			for _, macPort := range mac {
				sb.WriteString(fmt.Sprintf("%02X", macPort))
			}
		}
	}

	h := fnv.New64a()
	h.Write([]byte(sb.String()))

	// node ID: hash(addresses) & 1023
	return int64(h.Sum64() & uint64(maxNodeID))
}

func randomValidNodeID() int64 {
	val, err := cryptorand.Int(cryptorand.Reader, big.NewInt(maxNodeID))
	if err != nil {
		return defaultNodeID
	}

	return val.Int64()
}
