package session

import (
	"net"
	"bufio"
	"io"
	"errors"
	"reflect"

	"github.com/theangryangel/insim-go/pkg/protocol"
	"github.com/theangryangel/insim-go/pkg/state"
)

type InsimSession struct {
	conn net.Conn
	reader *bufio.Reader
	writer *bufio.Writer

	types map[uint8]func() (protocol.Packet)
	handlers map[reflect.Type][]reflect.Value

	GameState state.State
}

func NewInsimSession() (*InsimSession) {
	return &InsimSession{}
}

func (c *InsimSession) registerBuiltInPackets() {
	// TODO: generate this automatically from pkg/protocol/*.go
	c.RegisterPacket(protocol.ISP_MSO, protocol.NewIrpSel)
	c.RegisterPacket(protocol.ISP_TINY, protocol.NewTiny)
	c.RegisterPacket(protocol.ISP_STA, protocol.NewSta)
	c.RegisterPacket(protocol.ISP_VER, protocol.NewVer)
	c.RegisterPacket(protocol.ISP_NCN, protocol.NewNcn)
	c.RegisterPacket(protocol.ISP_CNL, protocol.NewCnl)
	c.RegisterPacket(protocol.ISP_CPR, protocol.NewCpr)
	c.RegisterPacket(protocol.ISP_NPL, protocol.NewNpl)
	c.RegisterPacket(protocol.ISP_PLL, protocol.NewPll)
}

func (c *InsimSession) RegisterPacket(ptype uint8, f func() (protocol.Packet)) {
	if c.types == nil{
		c.types = make(map[uint8]func() (protocol.Packet))
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

func (c *InsimSession) UseConn(conn net.Conn) (err error) {
	c.conn = conn
	c.reader = bufio.NewReader(c.conn)
	c.writer = bufio.NewWriter(c.conn)

	c.registerBuiltInPackets()
	c.trackGameState()
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
	verreq := protocol.NewTiny().(*protocol.Tiny)
	verreq.ReqI = 1
	verreq.SubT = protocol.TINY_VER

	err = c.Write(verreq)
	if err != nil {
		return err
	}

	connreq := protocol.NewTiny().(*protocol.Tiny)
	connreq.ReqI = 2
	connreq.SubT = protocol.TINY_NCN

	err = c.Write(connreq)
	if err != nil {
		return err
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

func (c *InsimSession) Read() (error) {
	nextByte, err := c.reader.Peek(1)
	if err != nil {
		return err
	}

  nextLen := uint8(nextByte[0])

	if c.reader.Buffered() < int(nextLen) {
		return ErrNotEnough
	}

	buf := make([]byte, nextLen)
	_, err = io.ReadFull(c.reader, buf)
	if err != nil {
		return err
	}

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
