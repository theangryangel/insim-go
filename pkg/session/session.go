package session

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"time"

	"github.com/theangryangel/insim-go/pkg/protocol"
	"github.com/theangryangel/insim-go/pkg/state"
)

type InsimSession struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer

	types    map[uint8]func() protocol.Packet
	handlers map[reflect.Type][]reflect.Value

	GameState state.GameState

	mangled uint8
	discard uint8
}

func NewInsimSession() *InsimSession {
	return &InsimSession{}
}

func (c *InsimSession) RegisterPacket(ptype uint8, f func() protocol.Packet) {
	if c.types == nil {
		c.types = make(map[uint8]func() protocol.Packet)
	}

	c.types[ptype] = f
}

func (c *InsimSession) Unmarshal(data []byte) (packet protocol.Packet, err error) {
	ptype := uint8(data[0])

	if v, found := c.types[ptype]; found {
		payload := v()
		err := payload.Unmarshal(data[1:])
		if err != nil {
			return nil, err
		}
		return payload, nil
	}

	return nil, ErrUnknownType
}

func (c *InsimSession) On(handler interface{}) {
	// Warning: non-idomatic code ahead.
	// The "right" way to do this is to have each handler do something like
	// godiscord, or have each handler have to check the type (which would
	// be tedious to consumers of the library).
	// This also allows us to create non-packet pseudo-events dynamically.
	if c.handlers == nil {
		c.handlers = make(map[reflect.Type][]reflect.Value)
	}

	r := reflect.ValueOf(handler)
	t := r.Type()

	if t.Kind() != reflect.Func || t.NumIn() != 2 || t.In(0) != reflect.TypeOf(&InsimSession{}) || t.In(1).Kind() != reflect.Ptr {
		panic("bad arg")
	}

	ptype := t.In(1)
	c.handlers[ptype] = append(c.handlers[ptype], r)
}

func (c *InsimSession) Dial(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	return c.UseConn(conn)
}

func (c *InsimSession) Use(f func(*InsimSession)) {
	f(c)
}

func (c *InsimSession) UseConn(conn net.Conn) (err error) {
	c.conn = conn
	c.reader = bufio.NewReader(c.conn)
	c.writer = bufio.NewWriter(c.conn)
	c.mangled = 0
	c.discard = 0

	c.Use(useBuiltInPackets)
	c.Use(usePing)
	c.Use(useGameState)
	return nil
}

func (c *InsimSession) Init() (err error) {
	isi := protocol.NewInit()

	err = c.Write(isi)
	if err != nil {
		return err
	}
	return nil
}

func (c *InsimSession) SelectRelayHost(hostname string) (err error) {
	sel := protocol.NewIrpSel().(*protocol.IrpSel)
	sel.HName = hostname
	return c.Write(sel)
}

func (c *InsimSession) RequestState() (err error) {
	subts := [...]uint8{
		protocol.TINY_VER,
		protocol.TINY_SST,
		protocol.TINY_NCN,
		protocol.TINY_NPL,
	}
	req := protocol.NewTiny().(*protocol.Tiny)

	for idx, subt := range subts {
		req.ReqI = uint8(idx)
		req.SubT = subt
		err = c.Write(req)
		fmt.Printf("Req State %d\n", subt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *InsimSession) Write(packet protocol.Packet) error {
	// TODO: Do we need to find a way to do the 2 writes atomically?

	data, err := packet.Marshal()
	if err != nil {
		return err
	}

	length := len(data)
	if length <= 0 {
		return errors.New("Not enough data")
	}

	length += 2 // packets dont include the size or type. we need to add to it.

	err = c.writer.WriteByte(byte(length))
	if err != nil {
		return err
	}

	ptype := packet.Type()
	err = c.writer.WriteByte(byte(ptype))
	if err != nil {
		return err
	}
	_, err = c.writer.Write(data)
	if err != nil {
		panic(err)
	}

	c.writer.Flush()
	return err
}

func (c *InsimSession) Read() error {
	// TODO Remove. This is to aide debugging a potential WSL2 packetloss issue
	if c.discard > 0 {
		c.reader.Discard(int(c.discard))
		c.discard = 0
	}

	nextByte, err := c.reader.Peek(1)
	if err != nil {
		return err
	}

	nextLen := uint8(nextByte[0])

	if c.reader.Buffered() < int(nextLen) {
		fmt.Printf("Needed %d, got %d, Buffer Cap %d", nextLen, c.reader.Buffered(), c.reader.Size())
		c.mangled++
		time.Sleep(time.Duration(c.mangled) * time.Second)
		if c.mangled > 4 {
			buf := make([]byte, c.reader.Buffered())
			_, err = io.ReadFull(c.reader, buf)

			fmt.Printf("Buffer: %v", buf)
			//panic(ErrNotEnough)
			fmt.Printf("resetting...\n")
			c.discard = nextLen - uint8(len(buf))
			c.mangled = 0
			c.reader.Reset(bufio.NewReader(c.conn))
		}
		return ErrNotEnough
	}

	buf := make([]byte, nextLen)
	_, err = io.ReadFull(c.reader, buf)
	if err != nil {
		return err
	}

	fmt.Printf("Got: %d\n", buf[1])

	packet, err := c.Unmarshal(buf[1:])
	if err != nil {
		return err
	}

	if packet != nil {
		c.Call(packet)
		return nil
	}

	return ErrNoPacket
}

func (c *InsimSession) ReadLoop() error {
	for {
		err := c.Read()
		if err != nil {

			if err == ErrUnknownType {
				continue
			}

			if err == ErrNotEnough {
				fmt.Printf("Not enough data to read full packet\n")
				continue
			}

			return err
		}
	}
}

func (c *InsimSession) Call(data interface{}) {
	// Warning: non-idomatic code. See the On method for documentation.
	dtype := reflect.TypeOf(data)
	a0 := reflect.ValueOf(c)
	a1 := reflect.ValueOf(data)

	for _, handler := range c.handlers[dtype] {
		handler.Call([]reflect.Value{a0, a1})
	}
}

func (c *InsimSession) Conn() net.Conn {
	return c.conn
}

func (c *InsimSession) Close() error {
	return c.conn.Close()
}
