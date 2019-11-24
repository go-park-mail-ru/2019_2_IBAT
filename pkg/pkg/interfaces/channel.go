package interfaces

// import (
// 	"bufio"
// 	"net"

// 	"github.com/google/uuid"
// )

// //remove to handler models
// //buffer size
// const N = 10

// type Packet struct {
// 	id uuid.UUID
// }

// // Channel wraps user connection.
// type Channel struct {
// 	conn net.Conn    // WebSocket connection.
// 	send chan Packet // Outgoing packets queue.
// }

// func NewChannel(conn net.Conn) *Channel {
// 	c := &Channel{
// 		conn: conn,
// 		send: make(chan Packet, N),
// 	}

// 	go c.reader()
// 	go c.writer()

// 	return c
// }

// func (c *Channel) reader() {
// 	// We make buffered read to reduce read syscalls.
// 	buf := bufio.NewReader(c.conn)

// 	for {
// 		pkt, _ := readPacket(buf)
// 		c.handle(pkt)
// 	}
// }

// func (c *Channel) writer() {
// 	// We make buffered write to reduce write syscalls.
// 	buf := bufio.NewWriter(c.conn)

// 	for pkt := range c.send {
// 		_ := writePacket(buf, pkt)
// 		buf.Flush()
// 	}
// }
