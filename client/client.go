package main
// This template draws inspiration from the following source:
// github.com/lucasfth/go-ass5
// Which was created by chbl, fefa and luha
// The template itself was created by chbl and luha

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	request "github.com/ChrBank/disys-2021-exam/grpc"

	"google.golang.org/grpc"
)

func main() {
	log.SetFlags(0)
	ctx := context.Background()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	c := &client{downedServerInt: 0}

	log.Printf("Enter username below:")
	fmt.Scanln(&c.name)

	// To change log location, outcomment below

	// path := fmt.Sprintf("clientlog_%s", c.name)
	// f, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()

	// log.SetOutput(f)

	c.downedServers = make (map[int32]bool) // init map of downed servers

	// Connect to all servers
	for i := 0; i < 3; i++ { // Will iterate through ports 5001, 5002, 5003
		dialNum  := int32(5001 + i)
		dialNumString := fmt.Sprintf(":%v", dialNum) 

		conn, err := grpc.Dial(dialNumString, grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		
		// create stream
		client := request.NewBiddingServiceClient(conn)
		in := &request.ClientHandshake{ClientPort: dialNum, Name: c.name} 
		//bidStream, err := client.SendBid(context.Background(), )
		stream, err := client.Handshake(ctx, in)
		if err != nil {
			log.Fatalf("open stream error %v", err)
		}
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Connected to server %v and responded %s", dialNum, resp)
		}
		if err != nil {
			log.Fatalf("Cannot receive %v", err)
		}
		c.servers = append(c.servers, client)
		c.downedServers[dialNum] = false
		time.Sleep(4 * time.Second)
	}

	c.communication()
}

func (c *client) communication() {
	amount := 0
	for { // Communication loop
		c.sendBids(int32(amount))
		log.Printf("Amount: %v", amount)
		amount++
		time.Sleep(4 * time.Second)
	}
}

func (c *client) sendBids(bid int32){
	responses := make([]string, len(c.servers))
	for i := 0; i < len(c.servers); i++ { // Send bid to all servers
		response, _ := c.sendBid(int32(i), bid)
		responses[i] = response
	}
}

func (c *client) sendBid(iteration int32, bid int32) (string, error) {
	in := &request.Bid{Name: c.name, Amount: bid}
	stream, err := c.servers[iteration].SendBid(context.Background(), in)
	if err != nil {
		serverDown(iteration, c)
		return "nil", err
	}
	resp, err := stream.Recv()
	return resp.GetResponse(), err
}

func serverDown (iteration int32, c *client) (bool){
	if (!c.downedServers[5001 + iteration]) {
		log.Printf("Server %v is down", (5001 + iteration))
		c.downedServers[5001 + iteration] = true
		return true // Server has just crashed
	}
	return false // Server was already down
}

type client struct {
	name 			string
	downedServers 	map[int32]bool
	downedServerInt int32
	servers 		[]request.BiddingServiceClient
}