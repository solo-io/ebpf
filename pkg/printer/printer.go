package printer

import (
	"fmt"
	"sort"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/rivo/tview"
	"github.com/solo-io/ebpf/pkg/internal/version"
)

const titleText = `
 ______     ______     ______   ______   ______     ______   __        
/\  ___\   /\  == \   /\  == \ /\  ___\ /\  ___\   /\__  _\ /\ \       
\ \  __\   \ \  __<   \ \  _-/ \ \  __\ \ \ \____  \/_/\ \/ \ \ \____  
 \ \_____\  \ \_____\  \ \_\    \ \_\    \ \_____\    \ \_\  \ \_____\ 
  \/_____/   \/_____/   \/_/     \/_/     \/_____/     \/_/   \/_____/ 

                              					(powered by solo.io)  `

type MapValue struct {
	Hash    uint64
	Entries []version.KvPair
	Table   *tview.Table
}

var mapOfMaps = make(map[string]MapValue)
var mapMutex = sync.RWMutex{}

type Monitor struct {
	MyChan chan version.MapEntries
	App    *tview.Application
	Flex   *tview.Flex
}

func NewMonitor() Monitor {
	app := tview.NewApplication()
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	title := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorLightCyan)
	fmt.Fprint(title, titleText)
	flex.AddItem(title, 10, 0, false)
	return Monitor{
		MyChan: make(chan version.MapEntries),
		App:    app,
		Flex:   flex,
	}
}

func (m *Monitor) Start() {
	go func() {
		if err := m.App.SetRoot(m.Flex, true).Run(); err != nil {
			panic(err)
		}
	}()
	// goroutine for updating the TUI data based on updates from loader watching maps
	go m.Watch()
}

func (m *Monitor) Watch() {
	for r := range m.MyChan {
		current := mapOfMaps[r.Name]
		newPrintHash, _ := hashstructure.Hash(r.Entries, hashstructure.FormatV2, nil)
		// Do not print if the data has not changed
		if current.Hash == newPrintHash {
			continue
		}

		// we have new entries, let's track them
		newMapVal := current
		newMapVal.Entries = r.Entries
		newMapVal.Hash = newPrintHash
		mapOfMaps[r.Name] = newMapVal

		// sort fields in key struct for consistent render order
		// TODO: use the BTF map info for this
		entry := r.Entries[0]
		theekMap := entry.Key
		keyStructKeys := []string{}
		for kk := range theekMap {
			keyStructKeys = append(keyStructKeys, kk)
		}
		sort.Strings(keyStructKeys)

		// get the instance of Table we will update
		table := newMapVal.Table
		// render the first row containing the keys
		c := 0
		for i, k := range keyStructKeys {
			cell := tview.NewTableCell(k).SetExpansion(1).SetTextColor(tcell.ColorYellow)
			table.SetCell(0, i, cell)
			c++
		}
		// last column in first row is value of the map (i.e. the counter/gauge/etc.)
		cell := tview.NewTableCell("value").SetExpansion(1).SetTextColor(tcell.ColorYellow)
		table.SetCell(0, c, cell)

		// now render each row according to the Entries we were sent by the loader
		// TODO: should we sort/order this in any specific way? right now they are
		// simply in iteration order of the underlying BTF map
		for r, entry := range newMapVal.Entries {
			r++ // increment the row index as the 0-th row is taken by the header
			ekMap := entry.Key
			eVal := entry.Value
			c := 0
			for kk, kv := range keyStructKeys {
				cell := tview.NewTableCell(ekMap[kv]).SetExpansion(1)
				table.SetCell(r, kk, cell)
				c++
			}
			cell := tview.NewTableCell(eVal).SetExpansion(1)
			table.SetCell(r, c, cell)
		}
		m.App.SetFocus(table)
		m.App.Draw()

		// print logic
		// printMap := map[string]interface{}{
		// 	"mapName": r.Name,
		// 	"entries": r.Entries,
		// }
		// byt, err := json.Marshal(printMap)
		// if err != nil {
		// 	fmt.Printf("error marshalling map data, this should never happen, %s\n", err)
		// 	continue
		// }
		// fmt.Printf("%s\n", byt)
	}
	fmt.Println("no more entries, closing")
}

func (m *Monitor) NewHashMap(name string) *tview.Table {
	table := tview.NewTable().SetFixed(1, 0)
	table.SetBorder(true).SetTitle(name)
	m.Flex.AddItem(table, 0, 1, false)
	mapMutex.Lock()
	entry := mapOfMaps[name]
	entry.Table = table
	mapOfMaps[name] = entry
	mapMutex.Unlock()
	return table
}
