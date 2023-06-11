package dns

import (
	"crypto/rand"
	"math/big"
	"bufio"
	"fmt"
	"net"

	"golang.org/x/net/dns/dnsmessage"
)

const ROOT_SERVERS = "198.41.0.4,199.9.14.201,192.33.4.12,199.7.91.13,192.203.230.10,192.5.5.241,192.112.36.4,198.97.190.53"

func handlePacket(pc net.PacketConn, addr net.Addr, buf []byte) error {
	return fmt.Errorf("not implemented yet")
}

func outgoingDnsQuery(servers []net.IP, question dnsmessage.Question) (*dnsmessage.Parser, *dnsmessage.Header, error) {
	max := uint16(^uint16(0))

	// generate a random number max to unit16
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return nil, nil, err
	}

	// convert random number type from bigint to uint16
	id := uint16(randomNumber.Int64())

	// define a UDP message (question)
	msg := dnsmessage.Message{
		Header: dnsmessage.Header{
			ID: id,
			Response: false,
			OpCode: dnsmessage.OpCode(0),
		},
		Questions: []dnsmessage.Question{question},
	}
	
	// pack message to send to root server
	buf, err := msg.Pack()
	if err != nil {
		return nil, nil, err
	}

	// make connection with servers and send message (for each servers)
	//* connections is with udp protocol
	var conn net.Conn
	for _, server := range servers{
		conn, err = net.Dial("udp", server.String() + ":53" /*dns port*/)
		if err == nil {
			break
		}
	}
	if conn == nil{
		return nil, nil, fmt.Errorf("faild to make connection to servers: %s", err)
	}
	// here we have the connection!
	_, err = conn.Write(buf)
	if err != nil {
		return nil, nil, err
	}

	// read answer  from connection
	answer := make([]byte, 512)
	n, err := bufio.NewReader(conn).Read(answer)
	if err != nil {
		return nil, nil, err
	}

	conn.Close()

	var p dnsmessage.Parser
	header, err := p.Start(answer[:n])
	if err != nil {
		return nil, nil, fmt.Errorf("parser start error: %s", err)
	}

	questions, err := p.AllQuestions()
	if err != nil {
		return nil, nil, err
	}
	if len(questions) != len(msg.Questions){
		return nil, nil, fmt.Errorf("answer packet dosen't have the same amount of questions")
	}
	
	err = p.SkipAllQuestions()
	if err != nil {
		return nil, nil, err
	}

	return &p, &header, nil
}