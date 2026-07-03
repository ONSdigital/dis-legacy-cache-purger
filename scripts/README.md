# Scripts

This folder is used for useful scripts that can help when developing the `dis-legacy-cache-purger`.

Currently there is:

- [seed.js](#seed)

## Seed

This seed script will seed your local mongo database with a number of cache time entries so these can be exposed via the dp-legacy-cache-api.

To run:

```sh
make seed
```

To override parameters:

```sh
make seed SEED_COUNT=500
make seed SEED_COUNT=500 SEED_PDF_COUNT=50
make seed SEED_COUNT=500 SEED_DATA_COUNT=50
make seed SEED_COUNT=1000 SEED_COLLECTION_ID=my-collection
```

| Variable             | Default           | Description                                                                                     |
|----------------------|-------------------|-------------------------------------------------------------------------------------------------|
| `SEED_COUNT`         | `1000`            | Total number of paths to seed                                                                   |
| `SEED_PDF_COUNT`     | `100`             | Number of PDF-eligible paths to seed (e.g. `/bulletins/`, `/articles/`, `/compendium_chapter/`) |
| `SEED_DATA_COUNT`    | `100`             | Number of data-eligible paths to seed (e.g. `/datasets/`)                                       |
| `SEED_PATH`          | `/test/path`      | Path to seed for normal (non-PDF, non-data) paths                                               |
| `SEED_COLLECTION_ID` | `test-collection` | Collection ID for all seeded documents                                                          |

You can also run the script directly with `mongosh --eval` to override any variable:

```sh
mongosh mongodb://localhost:27017 --eval "var count=500; var pdfCount=50; var dataCount=50; var collectionID='test-collection'; var path='/test/path';" scripts/seed.js
```
