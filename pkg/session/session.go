package session

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"time"

	"github.com/theangryangel/insim-go/pkg/protocol"
	"github.com/theangryangel/insim-go/pkg/state"
)

type InsimSession struct {
	conn    net.Conn
	writer  *bufio.Writer
	scanner *bufio.Scanner

	TimeoutDuration time.Duration

	types    map[uint8]func() protocol.Packet
	pre      map[reflect.Type][]reflect.Value
	handlers map[reflect.Type][]reflect.Value

	GameState state.GameState
}

func NewInsimSession() *InsimSession {
	return &InsimSession{
		TimeoutDuration: time.Second * 70,
	}
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
		err := payload.UnmarshalInsim(data[1:])
		if err != nil {
			return nil, err
		}
		return payload, nil
	}

	return nil, ErrUnknownType
}

func (c *InsimSession) PreOn(handler interface{}) {
	// PreOn handlers are run synchonrously before On handlers.
	// Think of it a bit like middleware in http servers.
	// It's reserved for "system" level stuff, like updating state, etc.

	// Warning: non-idomatic code ahead.
	// The "right" way to do this is to have each handler do something like
	// godiscord, or have each handler have to check the type (which would
	// be tedious to consumers of the library).
	// This also allows us to create non-packet pseudo-events dynamically.
	if c.pre == nil {
		c.pre = make(map[reflect.Type][]reflect.Value)
	}

	r := reflect.ValueOf(handler)
	t := r.Type()

	if t.Kind() != reflect.Func || t.NumIn() != 2 || t.In(0) != reflect.TypeOf(&InsimSession{}) || t.In(1).Kind() != reflect.Ptr {
		panic("bad arg")
	}

	ptype := t.In(1)
	c.pre[ptype] = append(c.pre[ptype], r)
}

func (c *InsimSession) On(handler interface{}) {
	// Warning: non-idomatic code ahead.
	// See PreOn method.
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
	conn, err := net.DialTimeout("tcp", address, time.Second*10)
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
	c.scanner = bufio.NewScanner(c.conn)
	c.scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		next := int(data[0])

		if next > len(data) {
			fmt.Println("Not enough data, only", len(data), "available")
			return 0, nil, nil
		}

		return int(next), data[0:next], nil
	})
	c.writer = bufio.NewWriter(c.conn)

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
		req.ReqI = uint8(idx + 1)
		req.SubT = subt
		err = c.Write(req)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *InsimSession) Write(packet protocol.Packet) error {
	// TODO: Do we need to find a way to do the 2 writes atomically?

	data, err := packet.MarshalInsim()
	if err != nil {
		return err
	}

	length := len(data)
	if length <= 0 {
		return errors.New("Not enough data")
	}

	// packets dont include the size or type. we need to add to it.
	length += 2

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

func (c *InsimSession) Scan(ctx context.Context) error {
	errs := make(chan error)
	lines := make(chan []byte)

	last := time.Now()
	timeout := time.NewTicker(10 * time.Second)
	defer timeout.Stop()

	go func() {
		for c.scanner.Scan() {
			lines <- c.scanner.Bytes()
		}

		if err := c.scanner.Err(); err != nil {
			fmt.Printf("Invalid input: %s", err)
			errs <- err
			return
		}
	}()

	for {
		select {
		case <-ctx.Done():
			c.Close()
			return nil
		case err := <-errs:
			if ctx.Err() != nil {
				return nil
			}

			return err
		case t := <-timeout.C:
			if ctx.Err() != nil {
				return nil
			}

			if time.Now().Sub(last) >= c.TimeoutDuration {
				fmt.Printf("Insim timeout at %s after %s\n", t, c.TimeoutDuration)
				return ErrTimeout
			}
		case buf := <-lines:
			if ctx.Err() != nil {
				return nil
			}

			last = time.Now()

			packet, err := c.Unmarshal(buf[1:])
			if err != nil {
				if err == ErrUnknownType {
					fmt.Println("Unknown Packet", buf[1])
				} else {
					return err
				}
			}

			if packet != nil {
				c.PreCall(packet)
				c.Call(packet)
				continue
			}
		}
	}
}

func (c *InsimSession) PreCall(data interface{}) {
	// Warning: non-idomatic code. See the On method for documentation.
	dtype := reflect.TypeOf(data)
	a0 := reflect.ValueOf(c)
	a1 := reflect.ValueOf(data)

	for _, handler := range c.pre[dtype] {
		handler.Call([]reflect.Value{a0, a1})
	}
}

func (c *InsimSession) Call(data interface{}) {
	// Warning: non-idomatic code. See the On method for documentation.
	dtype := reflect.TypeOf(data)
	a0 := reflect.ValueOf(c)
	a1 := reflect.ValueOf(data)

	for _, handler := range c.handlers[dtype] {
		go handler.Call([]reflect.Value{a0, a1})
	}
}

func (c *InsimSession) Conn() net.Conn {
	return c.conn
}

func (c *InsimSession) Close() error {
	return c.conn.Close()
}
