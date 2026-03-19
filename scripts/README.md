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

There are variables inside the script that can be used to configure it - these may be made available at the command line in the future.
