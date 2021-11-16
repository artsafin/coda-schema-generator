package main

import (
	"coda-schema-generator/internal/config"
	"coda-schema-generator/internal/dto"
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

func (c *client) loadTables() (dto.Tables, error) {
	var err error

	resp, err := c.http.Get(c.endpointf("docs/%s/tables", c.opts.DocID))
	if err != nil {
		return dto.Tables{}, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	tableItems := dto.Tables{}
	err = dec.Decode(&tableItems)
	if err != nil {
		return dto.Tables{}, err
	}

	return tableItems, nil
}

func (c *client) loadColumns(tables dto.Tables) (map[string]dto.TableColumns, error) {
	var wg sync.WaitGroup

	out := make(chan dto.TableColumns)

	wg.Add(len(tables.Items))

	for _, t := range tables.Items {
		go func(tableID string) {
			resp, err := c.http.Get(c.endpointf("docs/%s/tables/%s/columns", c.opts.DocID, tableID))

			defer resp.Body.Close()

			if err != nil {
				return
			}

			dec := json.NewDecoder(resp.Body)

			columns := dto.TableColumns{
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

	cols := make(map[string]dto.TableColumns)

	for v := range out {
		cols[v.TableID] = v
	}

	return cols, nil
}
