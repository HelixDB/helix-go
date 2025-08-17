package internal

import "github.com/HelixDB/helix-go"

var HelixClient *helix.Client

func ConfigHelix() {
	HelixClient = helix.NewClient("http://localhost:6969")
}
