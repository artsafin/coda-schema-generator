package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/artsafin/coda-schema-generator/dto"
	"log"
	"net/http"
	"sync"
)

type client struct {
	opts dto.APIOptions
	http *http.Client
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func NewClient(options dto.APIOptions) *client {
	c := &client{
		opts: options,
	}
	c.http = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			c.logvf("request: %s\n", req.URL.String())
			req.Header.Set("Authorization", "Bearer "+options.Token)
			return http.DefaultTransport.RoundTrip(req)
		}),
		Timeout: options.RequestTimeout,
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

func (c *client) loadEntities(res dto.ItemsContainer, entityType string) (err error) {
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

func (c *client) LoadTables() (dto.TableList, error) {
	tables := dto.TableList{}
	err := c.loadEntities(&tables, "tables")
	if err != nil {
		return dto.TableList{}, err
	}

	return tables, nil
}

func (c *client) LoadFormulas() (dto.EntityList, error) {
	el := dto.EntityList{}
	err := c.loadEntities(&el, "formulas")
	if err != nil {
		return dto.EntityList{}, err
	}
	return el, nil
}

func (c *client) LoadControls() (dto.EntityList, error) {
	el := dto.EntityList{}
	err := c.loadEntities(&el, "controls")
	if err != nil {
		return dto.EntityList{}, err
	}
	return el, nil
}

func (c *client) LoadColumns(tables dto.TableList) (cm map[string]dto.TableColumns, err error) {
	defer func() {
		c.logvf("finished loading columns for %d tables", len(tables.Items))
	}()

	var wg sync.WaitGroup

	out := make(chan dto.TableColumns)

	wg.Add(len(tables.Items))

	for _, t := range tables.Items {
		go func(tableID, tableType string) {
			defer wg.Done()

			resp, err := c.http.Get(c.endpointf("docs/%s/tables/%s/columns", c.opts.DocID, tableID))

			if err != nil {
				c.logvf("error fetching columns of %s: %v", tableID, err)
				return
			}
			defer resp.Body.Close()

			dec := json.NewDecoder(resp.Body)

			columns := dto.TableColumns{
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
		}(t.ID, t.TableType)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	cm = make(map[string]dto.TableColumns)

	for v := range out {
		cm[v.TableID] = v
	}

	if len(cm) != len(tables.Items) {
		return nil, errors.New("failed to load all tables")
	}

	return cm, nil
}
