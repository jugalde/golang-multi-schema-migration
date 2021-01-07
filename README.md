# Golang Multi Schema Migration

## Use case
When creating a B2B app or more broadly an app where multiple users are split across multiple organizations, one way of securely segmenting data by organization is to have one schema per organization. This repo aims to solve the problem of database migrations for that multi-schema use case.


## Why not use golang-migrate?
[golang-migrate](https://github.com/golang-migrate/migrate) solves for a general use case well. But for my multi-schema use case, I prefer a clear segmentation in how they are created, stored, and deployed. Also, I prefer to store my migrations in .go files in case there may be a need to add go code into my migration (rare but has happened before).

## How to use?
### Create public and schema-specific migrations
The migrate binary can be used to create your migrations. The output directory is defined in cmd/migrate/main, as the variable `const directory = "/migrations"`. Usage is shown below.
```
  $ ./cmd/migrate/migrate -h
    Usage of ./cmd/migrate/migrate:
      -isPublic
          set true if this migration is ran at the public level
      -name string
          the name of the migration (only lowercase and underscore)
```

### How to migrate up and down


Use `CreateMigrator(hostName, port, database, user, password string, lg *log.Logger) (*Migrator, error)` in `migrator.go`. Migrator then exposes the following API.

```
// Public API
MigratePublicUpToLatest(ctx context.Context) error 
MigratePublicUpToDateString(ctx context.Context, dateString string) error 
MigratePublicDownFully(ctx context.Context) error
MigratePublicDownToDateString(ctx context.Context, dateString string) error 

// Schema API
MigrateSchemaUpToLatest(ctx context.Context, schema string) error 
MigrateSchemaUpToDateString(ctx context.Context, schema string, dateString string) error 
MigrateSchemaDownFully(ctx context.Context, schema string) error 
MigrateSchemaDownToDateString(ctx context.Context, schema string, dateString string) error
```

To summarize the functionality, for public and schema level migrations, Migrator lets you migrate up to the latest migration, up to a specific datestring, down to a specific date string, for down fully (undoing every migration including the first).

I intentionally don't expose running migrations on all schema in one go. I've seen that lead to service outages before, so if this is functionality you really want, implement it yourself with full knowledge of your infrastructure.

## Testing migrations
Currently, `migrations/public` and `migrations/schema` have sample migrations. These migrations are used to test all of the API of Migrator in `migration_test.go`.

Run tests using 
``
    $make integration_test
``
Running tests will require docker and docker-compose. I'd personally recommend testing migrations using int tests like this, but rather than one test for all migrations, one test per migration.

## Miscellaneous Notes
golang-multi-schema-migration is built with postgres in mind. It uses pgxpool from [pgx](https://github.com/jackc/pgx/) to manage migrations. If you wish to work with other databases, you should be able to modify the code pretty easily.





