# coda-schema-generator
Represents a Coda.io document as Go code structs and methods (see [# Generated code](#generated-code) below).
Generated code is assumed to be used as part of a more complicated Go application that extend Coda.io capabilities.

It allows to use Coda document programmatically:
- Use autogenerated data structs reflecting tables and columns of the Coda document
- Load data from your document via Coda REST API
- Closely reflects Coda types, supports Lookups columns and allows to fetch data with deep loading of Lookup relations
- Refer to IDs of the Coda entities

Generated code very closely reflects the contents of the document -- tables, views, controls, formulas and provides a way to validate your app against current version of the Coda document.
This approach allows to be sure that no changes in Coda structure affect your app logic: your app will just fail to build unless the symbols from the generated code are used properly.

Additionally for each formula column there is a comment generated with the formula text.
This creates an orthogonal usage of the generated file -- finding/keeping records of usages of a certain column among other columns' formulas.

# Usage as standalone CLI app

## Generate schema for your Coda document
```
docker run --rm ghcr.io/artsafin/coda-schema-generator:v1.0.8 $CODA_TOKEN $CODA_DOCUMENT > internal/codaschema/codaschema.go
```

where:

- `$CODA_TOKEN` is an API token for Coda
- `$CODA_DOCUMENT` is a Coda document ID (`XXXXXXXXXX` part in the `https://coda.io/d/YOURHUMANDOCNAME_dXXXXXXXXXX` url)

# Usage as library

It's possible to include this package as library to have access to intermediate data and to have more control over generation.

```
go get -d github.com/artsafin/coda-schema-generator@latest
```

Use in your code:

```
import (
	...
	"github.com/artsafin/coda-schema-generator/schema"
	"github.com/artsafin/coda-schema-generator/generator"
	...
)

func main() {
    ...
    sch, err := schema.Get(opts)
	if err != nil {
		panic(err)
	}

	fmt.Println(sch.Tables.Items)
	
	err = generator.Generate("codaschema", "codaapi", sch, os.Stdout)
	if err != nil {
		panic(err)
	}
    ...
}
```

For example it allows to exclude some tables or columns from generation.

# Generated code

This is the most interesting part!

The output file will contain the following code:

## `codaschema.ID` variable and related strongly-typed structs to store IDs

There are tables, columns, formulas and controls type structs reflecting _metadata_ of the entities contained in the document.
They are useful only as formula comments holders and as helper types for populating the `codaschema.ID` (see next section)

### `codaschema.ID` variable

The `codaschema.ID` global variable contains an enumeration of document entities metadata - their IDs and human names.

Mostly this is useful only to _refer_ to tables, table columns, etc by their human names for future use in calling Coda APIs.
Otherwise one would need not only to hardcode them manually but also to keep them up to date along with changes in the Coda document.

For example if there is a `All Users` table in Coda with a `first name` and `last name` columns the generated code will include:

(note how the names are converted to Go symbols)

```
codaschema.ID.Table.AllUsers // Structure related to the [All Users] table

codaschema.ID.Table.AllUsers.ID // A string field containing ID of the [All Users] table for use API calls.
                                // Example value: "grid-RcSeMRomST"

codaschema.ID.Table.AllUsers.Name // A string field containing the name of the [All Users] table.
                                  // Example value: "All Users"

codaschema.ID.Table.AllUsers.Cols // Structure related to the [All Users] columns

codaschema.ID.Table.AllUsers.Cols.FirstName // Structure related to the [All Users].[first name] column

codaschema.ID.Table.AllUsers.Cols.FirstName.ID   // (see below)
codaschema.ID.Table.AllUsers.Cols.FirstName.Name // String fields containing ID and name of the column (similar to table)
```

Changing name of the table/column in Coda document and regenerating the code **will change** the schema -- thus possibly making your app impossible to build because it may have used fields that are undefined now.

The reason behind that is the difference of typical lifecycles of the Go application and Coda document: thanks to the team behind the Coda it is very easy to rapidly change a structure of the Coda document and overturn the whole way how the document is modelled. But it's not so easy to keep the same pace of the changes on the Go app side.  

So this is considered a good practice as it encourages to keep the mental model of the document with the app code up to date.

## Data structs reflecting the document data model and their constructors

The data structs will be created with `<Table Name>` type names and represent one row of table/view data. Fields of the data struct will reflect the columns of the table.
The types of the struct fields will be as close to the column types of the source table as possible (see [Type Mapping](#type-mapping) below).

Additionally there will be a `New<Table Name>` constructor generated per each table that creates an instance of the table row from a `Valuer` interface.

For example if there is a `All Users` table in Coda with a `first name`, `last name` and `location` columns there will be:

```
type Valuer interface {
	GetValue(key string) (value interface{}, ok bool)
}


type AllUsers struct {
    FirstName string
    LastName  string
    Location  LocationLookup // see more below on Lookup types
}

func NewAllUsers(row codaschema.Valuer) (codaschema.AllUsers, error) {
    ...
}

type Location struct {
    City string
    Zip  string
}

func NewLocation(row codaschema.Valuer) (codaschema.Location, error) {
    ...
}
```

These structs and constructors are useful on their own as they assume nothing about the way how the data has been loaded (by means of `Valuer`) and return clean strongly-typed structs containing the minimal data without implications on how to work with them.


### Type mapping

| Coda Type                 | Go Type                                                                        |
|---------------------------|--------------------------------------------------------------------------------|
| Date and Time, Date, Time | `time.Time`                                                                    |
| Scale                     | `uint8`                                                                        |
| Number, Slider            | `float64`                                                                      |
| Checkbox                  | `bool`                                                                         |
| Person, Reaction          | [`[]codaschema.Person`](internal/templates/coda_types.go.tmpl)       |
| Image, Attachment         | [`[]codaschema.Attachment`](internal/templates/coda_types.go.tmpl)   |
| Currency                  | [`codaschema.MonetaryAmount`](internal/templates/coda_types.go.tmpl) |
| Lookup                    | `codaschema.<Table Name>Lookup` (see below more on Lookup types)               |
| (the rest)                | `string`                                                                       |

### Lookup types

#### Lookup type `<Table Name>Lookup`

For each Coda table that is referenced by other Coda tables or views in Lookup columns the generator yields a `<Table Name>Lookup` type.

The lookup type represents possible multiple lookup values in the document as a slice of row references to another table (to which the Lookup in Coda is made).

Example of `<Table Name>Lookup` struct definition:

```
type LocationLookup struct {
    Values []LocationRowRef
}
```
Lookup type's `Values` is a slice because generally Coda allows to store multiple references in one table cell.

Lookup types also have handy methods to fetch the first row reference and the data of the first row reference.

A lookup type is only a container of the references to some another table rows, holding data no more than Coda provides it in API when the table rows are fetched.

#### Row reference type `<Table Name>RowRef`

A reference to a row is a `<Table Name>RowRef` type.

Example of `<Table Name>RowRef`:

```
type LocationRowRef struct {
    Name  string
    RowID string
    Data  *Location
}
```
When the table data is loaded this struct is guaranteed to contain `Name` and `RowID`.

The `Data` will be empty because it requires an extra HTTP request to be populated. The `Data` field has a type of the target referenced table. It will be used by deep data loaders, see [Deep loading](#deep-loading) section below.

`Name` is a value of Display column of the referenced table. In some usecases having just this value may be enough for the purposes of the app even without deep-loading of the row data. However using Name in cases except printing to the user is not recommended as the semantic use of the `Name` in app cannot be reliably verified against changes of Display column in Coda.

`Row ID` is a unique ID of the row, for example - `"i-4p2VnnLVjo"`. It is not what the `RowId(thisRow)` Coda formula returns.

## Loading data into structs

The generator creates functions that leverage previously described data structures and constructors for fetching table data directly from the Coda REST API using the https://github.com/artsafin/coda-go-client client library.

### CodaDocument type

`CodaDocument` is an abstraction of the Coda Document to avoid passing around common parameters every time.

Methods for data loading are all having the `CodaDocument` receiver.

`CodaDocument` has a common `ListAllRows` method that can be used separately to build your own logic:

```
func (d *CodaDocument) ListAllRows(ctx context.Context, tableID string, extraParams ...codaapi.ListRowsParam) ([]codaapi.Row, error)
```

Where:
- `codaapi.Row` type satisfies the `codaschema.Valuer` interface and therefore can be used with `New<Table Name>` data constructors
- `tableID` can be taken from `codaschema.ID.Tables.<Table Name>.ID`
- `codaapi.ListRowsParam` is a way to set the Query Parameters for the [ListRows API endpoint](https://coda.io/developers/apis/v1#operation/listRows)

Example of usage:
```
    token := "..."
    docID := "..."
    doc, err := codaschema.NewCodaDocument("https://coda.io/apis/v1", token, docID)
    rows, err := doc.ListAllRows(
        ctx,
        codaschema.ID.Table.AllUsers.ID,
        codaapi.ListRows.SortBy(codaapi.RowsSortByNatural),
        codaapi.ListRows.Query(codaschema.ID.Table.AllUsers.Cols.FirstName.ID, "donald"),
    )
    if err != nil {
        ...
    }

    for _, row := range rows {
        item, err := NewAllUsers(&row)
        if err != nil {
            ...
        }
        fmt.Println(item)
    }
```


`CodaDocument` also contains functions specific to your document tables with two flavors:
- shallow loading of the data - i.e. loading minimal data with just one HTTP request per table ([see note](#pagination));
- deep loading - routines to enrich shallow data with the nested data of the Lookup columns. They load data into the `codaschema.<Table Name>Lookup[N].Data` fields.

### Shallow loading

The methods here are just syntactic sugar on top of the `ListAllRows` method.

**`List<TableName>` methods** -- load data as a slice of `<Table Name>` structs.

Example signature:

```
func (d *CodaDocument) ListAllUsers(ctx context.Context, extraParams ...codaapi.ListRowsParam) ([]codaschema.AllUsers, error)
```

The slice will be ordered in the way Coda returns rows via the API. You can use the `codaapi.ListRows.SortBy(...)` parameter with the `codaapi.RowsSortBy*` constants to control the order.

**`MapOf<TableName>` methods** -- load data as a map of `<Table Name>` structs keyed by Coda row ID. The method also returns a slice of row IDs that conveys the order in which the rows were returned from the Coda API (because Go doesn't maintain an order in `map`s). 

Example signature:

```
func (d *CodaDocument) MapOfAllUsers(ctx context.Context, extraParams ...codaapi.ListRowsParam) (map[RowID]AllUsers, []RowID, error)
```

`MapOf` methods were created to deal with lookup relations - in such scenario you usually know the Row ID and would like to address rows by ID directly.

### Deep loading

Deep loading methods require an already loaded map of data using `MapOf` methods.
They go through the rows and populate the `Data` field of the `RowRef` structs of the lookup values.

The method mutates the map passed to the function.

Example signature:
```
func (doc *CodaDocument) LoadRelationsAllUsers(ctx context.Context, shallow map[RowID]codaschema.AllUsers, rels codaschema.Tables) (err error)
```

The last argument is a `codaschema.Tables` struct which is an enumeration of the document tables. This argument specifies which Lookup relations should be loaded by the `LoadRelations` method.

`LoadRelations<Table Name>` methods load _all related_ table data into an internal `CodaDocument` cache.

This approach has it's pros and cons:

- Pro: relation loading is mostly a special routine in the app when a lot of deeply-loaded entities are needed at once, therefore raising HIT/MISS cache ratio.
  Caching data of each kind of the referenced table at CodaDocument struct optimizes application by CPU and time while sacrificing RAM;
- Pro: every kind of data is loaded and parsed only once; even for further calls to other `LoadRelations<Table Name>` methods no extra requests will be made
- Pro: solves the N+1 problem
- Cons: the volume of the referenced data may be high and it is not limited in any way ([see note](#pagination)). Unfortunately Coda API doesn't provide any way to fetch only needed rows in batch by their IDs.
  Extra RAM usage can be an issue
- Cons: cache is not invalidated nor cleared; `CodaDocument` instance lifecycle should be as short as possible

### Pagination

The `ListAllRows` routine loads all the pages of the table data, so as methods depending on it - `ListRows<Table Name>`, `MapOf<Table Name>` etc.

Since each load of the page makes a separate HTTP request loading the big table will issue several HTTP requests depending on the page size.

# Example of application using codaschema

```
package main

import (
	"context"
	"fmt"
	"test-app/codaschema"
	"github.com/artsafin/coda-go-client/codaapi"
	"os"
)

func main() {
    token := os.Args[1]
    docID := os.Args[2]

    doc, err := codaschema.NewCodaDocument("https://coda.io/apis/v1", token, docID)
	if err != nil {
		panic(err)
	}

	ctx := context.Background() // On real apps it should include timeout at least

	usersMap, usersOrder, err := doc.MapOfAllUsers(
		ctx,
		codaapi.ListRows.SortBy(codaapi.RowsSortByNatural), // Returns the order that is seen in the browser 
	)
	if err != nil {
		panic(err)
	}
	err = doc.LoadRelationsAllUsers(ctx, usersMap, codaschema.Tables{
		Location: true,
	})
	if err != nil {
		panic(err)
	}

	for idx, rowid := range usersOrder {
		locationRef, ok := usersMap[rowid].Location.FirstRef()
		if !ok {
		    // Even if there is no first ref in Location lookup the locationRef will be an empty struct
		}

		fmt.Printf(
			"#%03d | %v: %20s | %v (%v) | %v %v\n",
			idx+1,
			locationRef.RowID,
			usersMap[rowid].Location.FirstData().City, // Or the same: locationRef.Data.City
			                                           // Or the same: usersMap[rowid].Location.Values[0].Data.City
			usersMap[rowid].FirstName,
			usersMap[rowid].LastName,
		)
	}
}
```

## Other generated stuff

This includes supporting functions and types for all of the above:

- `Valuer` interface that decouples an API library https://github.com/artsafin/coda-go-client from the generated code
- Various structs reflecting some Coda complex types: Person, MonetaryAmount, Attachment/ImageObject and a Structured Value.
- Routines for parsing of basic internal types (strings, dates, numbers etc)
- Aggregate errors container

Mostly these are useful only for other `codaschema` code.

# Known shortcomings

- One may have security concerns about feeding the generator an API key exposing access to sensitive data
- Coda allows to violate column types via Formulas. E.g. if there is a Lookup column but the column formula yields a string onto the cell value the loading routines will result in error. Types are generated according to the column metadata and the data conversion routine (that is part of the loading) will fail unless the data fits the declared type
- There is no way to limit what code do you need - medium-sized document (~30 tables, ~20 views, ~20 canvas formulas) can yield ~10k LOC
- This project is targeted to work with Coda in readonly mode
- It doesn't include APIs to work with Packs because the client library has an issue with OpenAPI schemas when they have `Response` suffixes (which is the case for all Packs APIs)
- The project itself has a minimal test coverage - I'm not yet certain on how to test generated code yet
