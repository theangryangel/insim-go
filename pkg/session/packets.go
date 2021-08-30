package session

import (
	"github.com/theangryangel/insim-go/pkg/protocol"
)

func useBuiltInPackets(c *InsimSession) {
	// TODO: generate this automatically from pkg/protocol/*.go
	c.RegisterPacket(protocol.IrpSel, protocol.NewIrpSel)
	c.RegisterPacket(protocol.IrpErr, protocol.NewIrpErr)
	c.RegisterPacket(protocol.IrpHos, protocol.NewIrpHos)
	c.RegisterPacket(protocol.IrpHlr, protocol.NewIrpHlr)

	c.RegisterPacket(protocol.IspTiny, protocol.NewTiny)
	c.RegisterPacket(protocol.IspSmall, protocol.NewSmall)
	c.RegisterPacket(protocol.IspSta, protocol.NewSta)
	c.RegisterPacket(protocol.IspVer, protocol.NewVer)
	c.RegisterPacket(protocol.IspNcn, protocol.NewNcn)
	c.RegisterPacket(protocol.IspToc, protocol.NewToc)
	c.RegisterPacket(protocol.IspCnl, protocol.NewCnl)
	c.RegisterPacket(protocol.IspCpr, protocol.NewCpr)
	c.RegisterPacket(protocol.IspNpl, protocol.NewNpl)
	c.RegisterPacket(protocol.IspPll, protocol.NewPll)
	c.RegisterPacket(protocol.IspPlp, protocol.NewPlp)
	c.RegisterPacket(protocol.IspMci, protocol.NewMci)
	c.RegisterPacket(protocol.IspMso, protocol.NewMso)
	c.RegisterPacket(protocol.IspCch, protocol.NewCch)

	c.RegisterPacket(protocol.IspIsm, protocol.NewIsm)
	c.RegisterPacket(protocol.IspRst, protocol.NewRst)
	c.RegisterPacket(protocol.IspReo, protocol.NewReo)
	c.RegisterPacket(protocol.IspSpx, protocol.NewSpx)
	c.RegisterPacket(protocol.IspLap, protocol.NewLap)
	c.RegisterPacket(protocol.IspFin, protocol.NewFin)
	c.RegisterPacket(protocol.IspRes, protocol.NewRes)
	c.RegisterPacket(protocol.IspFlg, protocol.NewFlg)
	c.RegisterPacket(protocol.IspPla, protocol.NewPla)
	c.RegisterPacket(protocol.IspCon, protocol.NewCon)

	c.RegisterPacket(protocol.IspVtn, protocol.NewVtn)
}
