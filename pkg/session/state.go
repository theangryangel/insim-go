package session

import (
	"github.com/theangryangel/insim-go/pkg/protocol"
)

func useGameState(c *InsimSession) {
	c.PreOn(func(c *InsimSession, ncn *protocol.Ncn) {
		c.GameState.FromNcn(ncn)
	})

	c.PreOn(func(c *InsimSession, cnl *protocol.Cnl) {
		c.GameState.FromCnl(cnl)
	})

	c.PreOn(func(c *InsimSession, toc *protocol.Toc) {
		c.GameState.FromToc(toc)
	})

	c.PreOn(func(c *InsimSession, npl *protocol.Npl) {
		c.GameState.FromNpl(npl)
	})

	c.PreOn(func(c *InsimSession, pll *protocol.Pll) {
		c.GameState.FromPll(pll)
	})

	c.PreOn(func(c *InsimSession, pll *protocol.Plp) {
		c.GameState.FromPlp(pll)
	})

	c.PreOn(func(c *InsimSession, mci *protocol.Mci) {
		c.GameState.FromMci(mci)
	})

	c.PreOn(func(c *InsimSession, flg *protocol.Flg) {
		c.GameState.FromFlg(flg)
	})

	c.PreOn(func(c *InsimSession, pla *protocol.Pla) {
		c.GameState.FromPla(pla)
	})

	c.PreOn(func(c *InsimSession, rst *protocol.Rst) {
		c.GameState.FromRst(rst)
	})

	c.PreOn(func(c *InsimSession, spx *protocol.Spx) {
		c.GameState.FromSpx(spx)
	})

	c.PreOn(func(c *InsimSession, lap *protocol.Lap) {
		c.GameState.FromLap(lap)
	})

	c.PreOn(func(c *InsimSession, sta *protocol.Sta) {
		c.GameState.FromSta(sta)
	})

	c.PreOn(func(c *InsimSession, vtn *protocol.Vtn) {
		c.GameState.FromVtn(vtn)
	})

	c.PreOn(func(c *InsimSession, vtn *protocol.Tiny) {
		c.GameState.FromTiny(vtn)
	})

	c.PreOn(func(c *InsimSession, vtn *protocol.Small) {
		c.GameState.FromSmall(vtn)
	})

	c.PreOn(func(c *InsimSession, res *protocol.Res) {
		c.GameState.FromRes(res)
	})

	c.PreOn(func(c *InsimSession, fin *protocol.Fin) {
		c.GameState.FromFin(fin)
	})
}
