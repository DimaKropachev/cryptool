package progressbar

import "github.com/cheggaaa/pb/v3"

const (
	tmpl = `{{string . "prefix"}} {{bar . (white "[") (green "=") (green ">") (white "_") (white "]")}} {{etime . }}`

	PrefixEncrypt = "Encrypting"
	PrefixDecrypt = "Decrypting"
)

type ProgressBar struct {
	prefix string
	pb     *pb.ProgressBar
}

func New(prefix string, total int64) *ProgressBar {
	p := pb.New64(total)
	p.SetTemplateString(tmpl)

	p.Set("prefix", prefix)
	return &ProgressBar{
		prefix: prefix,
		pb:     p,
	}
}

func (p *ProgressBar) Start() {
	p.pb.Start()
}

func (p *ProgressBar) Finish() {
	p.pb.Finish()
}

func (p *ProgressBar) Add(n int) {
	p.pb.Add(n)
}

type ProgressBarsPull struct {
	pbs  []*ProgressBar
	pool *pb.Pool
}

func NewPull() *ProgressBarsPull {
	return &ProgressBarsPull{
		pbs: []*ProgressBar{},
	}
}

func (p *ProgressBarsPull) Add(prefix string, total int64) *ProgressBar {
	pb := New(prefix, total)
	p.pbs = append(p.pbs, pb)
	return pb
}

func (p *ProgressBarsPull) Start() error {
	bars := []*pb.ProgressBar{}
	for _, bar := range p.pbs {
		bars = append(bars, bar.pb)
	}

	pool, err := pb.StartPool(bars...)
	if err != nil {
		return err
	}

	p.pool = pool

	return nil
}

func (p *ProgressBarsPull) Stop() {
	p.pool.Stop()
}
