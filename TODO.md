
# Action Items

## Bugs

## Next

- [ ] Add the following "types" for common library and archive identifiers
  - [ ] DOI
  - [ ] ror
  - [ ] ark
  - [ ] PMCID (Pub Med Central ID)
  - [ ] ISBN
  - [ ] ISSN
- [ ] Add comments to the top of the Go files indicating authorship and copyright
- [ ] Think through model attributes and decide if they are really needed
- [X] Remove title from model
- [X] Add element generators
  - genetators in elements populate an element value
  - these include the following common types from SQL space
    - autoincrment
    - uuid 
    - timestamps (updates)
    - timestamp if now set (created timestamps)

## Someday, Maybe

- [ ] Add a Postgres generator, update modelgen to have it available
- [ ] Add a MySQL generator, update modelgen to have it available
- [ ] Add a Model YAML generator example
