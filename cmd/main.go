package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
	"github.com/MontFerret/ferret/pkg/drivers/http"
)

type Stock struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func main() {
	listedStocks, err := getListedStocks()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, s := range listedStocks {
		fmt.Println(fmt.Sprintf("%s: %s", s.Name, s.URL))
	}
}

func getListedStocks() ([]Stock, error) {
	query := `
						LET doc = DOCUMENT('https://www.kauppalehti.fi/porssi/indeksit/OMXHPI', {driver: "cdp"})
						WAIT_ELEMENT(doc, '.stock-list', 5000)
						LET stockNameLinks = ELEMENTS(doc, '.stock-name-link')
						FOR stockNameLink IN stockNameLinks
						LET stockName = stockNameLink.innerText
						LET link =stockNameLink.attributes.href
						RETURN {
							name: stockName,
							url: link
						}
			}
	`

	comp := compiler.New()

	program, err := comp.Compile(query)

	if err != nil {
		return nil, err
	}

	// create a root context
	ctx := context.Background()

	// enable HTML drivers
	// by default, Ferret Runtime does not know about any HTML drivers
	// all HTML manipulations are done via functions from standard library
	// that assume that at least one driver is available
	ctx = drivers.WithContext(ctx, cdp.NewDriver())
	ctx = drivers.WithContext(ctx, http.NewDriver(), drivers.AsDefault())

	out, err := program.Run(ctx)

	if err != nil {
		return nil, err
	}

	res := make([]Stock, 100)

	err = json.Unmarshal(out, &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
