# coda-schema-generator
Schema generator for coda.io documents

# Usage

## Generate schema for your Coda document
```
docker run --rm ghcr.io/artsafin/coda-schema-generator/coda-schema-generator:main $CODA_TOKEN $CODA_DOCUMENT > internal/codaschema/ids.go
```

where:

`$CODA_TOKEN` is an API token for Coda
`$CODA_DOCUMENT` is Coda document ID

## Use in your code

Query some `Users` grid filtering by `LastName` column:

```
params := coda.ListRowsParameters{
    Query: fmt.Sprintf("\"%s\":\"%s\"", codaschema.ID.Table.Users.Cols.LastName, lastName),
}
resp, err := coda.ListTableRows(codaDocumentId, codaschema.ID.Table.Users.ID, params)
```
