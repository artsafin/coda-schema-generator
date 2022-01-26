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
	c := &client{
		opts: options,
	}
	c.http = &http.Client{
		Transport: RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			c.logvf("request: %s\n", req.URL.String())
			req.Header.Set("Authorization", "Bearer "+options.Token)
			return http.DefaultTransport.RoundTrip(req)
		}),
	}

	return c
}

func (c *client) logvf(format string, params ...interface{}) {
	if !c.opts.Verbose {
		return
	}
	log.Printf(format, params...)
}

func (c *client) endpointf(endpoint string, params ...interface{}) string {
	return fmt.Sprintf("%s/%s", c.opts.Endpoint, fmt.Sprintf(endpoint, params...))
}

func (c *client) loadEntities(res ItemsContainer, entityType string) (err error) {
	defer func() {
		c.logvf("finished loading %s: %d items", entityType, res.Count())
	}()

	resp, err := c.http.Get(c.endpointf("docs/%s/%s", c.opts.DocID, entityType))
	if err != nil {
		c.logvf("error fetching %s entities: %v", entityType, err)
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	err = dec.Decode(res)
	if err != nil {
		c.logvf("error decoding json of %s entities: %v", entityType, err)
		return err
	}

	return nil
}

func (c *client) LoadTables() (TableList, error) {
	tables := TableList{}
	err := c.loadEntities(&tables, "tables")
	if err != nil {
		return TableList{}, err
	}

	return tables, nil
}

func (c *client) LoadFormulas() (EntityList, error) {
	el := EntityList{}
	err := c.loadEntities(&el, "formulas")
	if err != nil {
		return EntityList{}, err
	}
	return el, nil
}

func (c *client) LoadControls() (EntityList, error) {
	el := EntityList{}
	err := c.loadEntities(&el, "controls")
	if err != nil {
		return EntityList{}, err
	}
	return el, nil
}

func (c *client) LoadColumns(tables TableList) (cm map[string]TableColumns, err error) {
	defer func() {
		c.logvf("finished loading columns for %d tables", len(tables.Items))
	}()

	var wg sync.WaitGroup

	out := make(chan TableColumns)

	wg.Add(len(tables.Items))

	for _, t := range tables.Items {
		go func(tableID, tableType string) {
			resp, err := c.http.Get(c.endpointf("docs/%s/tables/%s/columns", c.opts.DocID, tableID))

			defer resp.Body.Close()

			if err != nil {
				c.logvf("error fetching columns of %s: %v", tableID, err)
				return
			}

			dec := json.NewDecoder(resp.Body)

			columns := TableColumns{
				TableID:   tableID,
				TableType: tableType,
			}
			err = dec.Decode(&columns)
			if err != nil {
				c.logvf("error decoding columns json of %s: %v", tableID, err)
				return
			}

			c.logvf("finished loading columns for %s: %d items", tableID, len(columns.Items))

			out <- columns
			wg.Done()
		}(t.ID, t.TableType)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	cm = make(map[string]TableColumns)

	for v := range out {
		cm[v.TableID] = v
	}

	return cm, nil
}
