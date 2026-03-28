package player

type Player struct {
	name      string
	nat20     int
	nat1      int
	danoTotal int
	danoMax   int
	curaTotal int
	curaMax   int
	quedas    int
	mortes    int
	custom    string
}

func New(name string) *Player {
	return &Player{name: name}
}

func (p *Player) Name() string        { return p.name }
func (p *Player) SucessoCritico() int { return p.nat20 }
func (p *Player) FalhaCritica() int   { return p.nat1 }
func (p *Player) DanoTotal() int      { return p.danoTotal }
func (p *Player) DanoMax() int        { return p.danoMax }
func (p *Player) CuraTotal() int      { return p.curaTotal }
func (p *Player) CuraMax() int        { return p.curaMax }
func (p *Player) Quedas() int         { return p.quedas }
func (p *Player) Mortes() int         { return p.mortes }
func (p *Player) Custom() string      { return p.custom }

func (p *Player) AddNat20() { p.nat20++ }
func (p *Player) AddNat1()  { p.nat1++ }
func (p *Player) AddQueda() { p.quedas++ }
func (p *Player) AddMorte() { p.mortes++ }

func (p *Player) UpdateStats(dTotal, dMax, cTotal, cMax int) {
	p.danoTotal = dTotal
	p.danoMax = dMax
	p.curaTotal = cTotal
	p.curaMax = cMax
}

func (p *Player) SetCustom(text string) {
	p.custom = text
}

func (p *Player) LoadStats(n20, n1, dt, dm, ct, cm, q, m int, custom string) {
	p.nat20 = n20
	p.nat1 = n1
	p.danoTotal = dt
	p.danoMax = dm
	p.curaTotal = ct
	p.curaMax = cm
	p.quedas = q
	p.mortes = m
	p.custom = custom
}

type PlayerSnapshot struct {
	Nat20     int
	Nat1      int
	DanoTotal int
	DanoMax   int
	CuraTotal int
	CuraMax   int
	Quedas    int
	Mortes    int
	Custom    string
}

func (p *Player) Snapshot() PlayerSnapshot {
	return PlayerSnapshot{
		Nat20:     p.nat20,
		Nat1:      p.nat1,
		DanoTotal: p.danoTotal,
		DanoMax:   p.danoMax,
		CuraTotal: p.curaTotal,
		CuraMax:   p.curaMax,
		Quedas:    p.quedas,
		Mortes:    p.mortes,
		Custom:    p.custom,
	}
}

func (p *Player) RestoreSnapshot(s PlayerSnapshot) {
	p.LoadStats(s.Nat20, s.Nat1, s.DanoTotal, s.DanoMax, s.CuraTotal, s.CuraMax, s.Quedas, s.Mortes, s.Custom)
}
