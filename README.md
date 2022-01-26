# coda-schema-generator
Schema generator for coda.io documents.

If you have a huge Coda document with lots of tables and columns and would like to fetch data programmatically you have several options:
1) Hardcode plain names of the tables, columns or other objects.
  This approach is rather brittle because Coda allows to change things too easily for to keep the pace of updating hardcoded values in the code.
2) Hardcode IDs of the objects. Better option but also subject to suffer from changing Coda objects, less often but still.
3) Introspect Coda document dynamically and generate code according to the most recent structure (is what this project do).
  This approach allows to be sure that no changes in Coda structure affect your app logic: your app will just fail to build unless all the references to Coda objects are correct

In essence the schema is just a bunch of Coda tables/columns/formulas/controls identifiers packed into a (opinionated) convenient structure.

Example of generated file:
```
...
    Users: _usersTable{
        ID: "grid-BoqFoXRyHd",
        Cols: _usersTableColumns{
            FirstName: _column {
                ID:   "c-Kauq98as09",
                Name: "First name",
            },
            LastName: _column {
                ID:   "c-Aas98Taka-",
                Name: "Last name",
            },
            Location: _column {
                ID:   "c-OIaiunafia",
                Name: "Location",
            },
        },
    },
...
```

# Usage for users

## Generate schema for your Coda document
```
docker run --rm ghcr.io/artsafin/coda-schema-generator/coda-schema-generator:main $CODA_TOKEN $CODA_DOCUMENT > internal/codaschema/ids.go
```

where:

- `$CODA_TOKEN` is an API token for Coda
- `$CODA_DOCUMENT` is Coda document ID (`XXXXXXXXXX` part in the `https://coda.io/d/YOURHUMANDOCNAME_dXXXXXXXXXX` url)

## Use in your code

Query some `Users` grid filtering by `LastName` column:

(note using `codaschema.ID.Table.Users.Cols.LastName` and `codaschema.ID.Table.Users.ID` instead of plain names or)

```
params := coda.ListRowsParameters{
    Query: fmt.Sprintf("\"%s\":\"%s\"", codaschema.ID.Table.Users.Cols.LastName, lastName),
}
resp, err := coda.ListTableRows(codaDocumentId, codaschema.ID.Table.Users.ID, params)
```

# Usage for developers

Build it:

```
make
```

Run it:

```
./build/csg $CODA_TOKEN $CODA_DOCUMENT > internal/codaschema/ids.go
```