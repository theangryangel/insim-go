package session

import(
	"github.com/theangryangel/insim-go/pkg/protocol"
)

func useBuiltInPackets(c *InsimSession) {
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
	c.RegisterPacket(protocol.ISP_PLP, protocol.NewPlp)
	c.RegisterPacket(protocol.ISP_MCI, protocol.NewMci)
}
