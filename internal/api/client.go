package api

import (
	"coda-schema-generator/internal/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type client struct {
	opts config.APIOptions
	http *http.Client
}

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (fn RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func NewClient(options config.APIOptions) *client {
	return &client{
		opts: options,
		http: &http.Client{
			Transport: RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
				if options.Verbose {
					log.Printf("request: %s\n", req.URL.String())
				}
				req.Header.Set("Authorization", "Bearer "+options.Token)
				return http.DefaultTransport.RoundTrip(req)
			}),
		},
	}
}

func (c *client) endpointf(endpoint string, params ...interface{}) string {
	return fmt.Sprintf("%s/%s", c.opts.Endpoint, fmt.Sprintf(endpoint, params...))
}

func (c *client) loadEntities(entityType string) (EntityList, error) {
	var err error

	resp, err := c.http.Get(c.endpointf("docs/%s/%s", c.opts.DocID, entityType))
	if err != nil {
		return EntityList{}, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	items := EntityList{}
	err = dec.Decode(&items)
	if err != nil {
		return EntityList{}, err
	}

	return items, nil
}

func (c *client) LoadTables() (EntityList, error) {
	return c.loadEntities("tables")
}

func (c *client) LoadFormulas() (EntityList, error) {
	return c.loadEntities("formulas")
}

func (c *client) LoadControls() (EntityList, error) {
	return c.loadEntities("controls")
}

func (c *client) LoadColumns(tables EntityList) (map[string]TableColumns, error) {
	var wg sync.WaitGroup

	out := make(chan TableColumns)

	wg.Add(len(tables.Items))

	for _, t := range tables.Items {
		go func(tableID string) {
			resp, err := c.http.Get(c.endpointf("docs/%s/tables/%s/columns", c.opts.DocID, tableID))

			defer resp.Body.Close()

			if err != nil {
				return
			}

			dec := json.NewDecoder(resp.Body)

			columns := TableColumns{
				TableID: tableID,
			}
			err = dec.Decode(&columns)
			if err != nil {
				return
			}
			//columns.TableID = tableID

			out <- columns
			wg.Done()
		}(t.ID)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	cols := make(map[string]TableColumns)

	for v := range out {
		cols[v.TableID] = v
	}

	return cols, nil
}
