package table

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)


type Table struct {
	t *tablewriter.Table
}

func New() *Table {
	return &Table{tablewriter.NewTable(os.Stdout,
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
			Settings: tw.Settings{Separators: tw.Separators{BetweenRows: tw.On}},
		})),
		)}
}

func (t *Table) SetHeader(headlines []string) {
	t.t.Header(headlines)
}

func (t *Table) SetContent(content [][]string) error {
	return t.t.Bulk(content)
}

func (t *Table) Render() error {
	return t.t.Render()
}