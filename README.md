# kemono-sync

`kemono.su` Patreon etc. archiver scraper.

Scrapes post files and attachments into the working directory.

Files are saved to a `.cache` directory and then linked into a creator's directory by

- Post ID
- Post title
- Post added date

for lookup in different ways.

The scraping mechanism isn't particularly sophisticated, but is hardcoded to 8 simultaneous
downloads to avoid a deliberate attempt to hit any rate limits.

Cached files won't be re-downloaded when re-running the same command, so it can be used to
incrementally sync files e.g. from a specific creator.

## Usage

Scrape files from specific post

```shell
(cd cmd/kemono-get-post && go install)
# You're Wrong About podcast:
# Extended Cut: George Michael With Marcus McCann, Part 2
kemono-get-post kemono.su patreon 19142444 103374306
```

Scrape files from creator

```shell
(cd cmd/kemono-get-creator && go install)
# You're Wrong About podcast
kemono-get-post kemono.su patreon 19142444
```

## Todo

- Lookup from service by name instead of just by ID.
- Link into library by name instead of by ID.
- Retry for failed downloads - currently only headers are retried, not body.
- More sophisticated sync options e.g. after X date or post.
- Sync options by library config file e.g. list creators to sync and do several
  simultaneously.
