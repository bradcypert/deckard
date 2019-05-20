# Deckard

### A framework agnostic tool for running database migrations.

#### Usage
```bash
deckard create add_login_date_to_users
# modify the created files
deckard up --host=localhost --port=5432 --user=user --password=pass --database=app
deckard down --host=localhost --port=5432 --user=user --password=pass --database=app
```

#### Managing your databases via YAML config.
Deckard also supports managing your databases via YAML.
Instead of writing
```bash
deckard up --host=localhost --port=5432 --user=user --password=pass --database=app
```

You can simply write
```bash
deckard up --key=prod
```
when you have a a `.deckard.yml` in your home directory. In this instance, your YAML should look like this:
```yaml
prod:
    host: localhost
    port: 5432
    user: user
    password: pass
    database: app
```

Alternatively, you can provide deckard the path to the configuration value you want to use.
```bash
deckard up --config=/usr/app/deckard.yml --key=prod
```

#### Verifying a migration was ran against the database.
Sometimes, you may find yourself curious if the migration was ran against the DB. You certainly can fire up your favorite database client and query for the metadata entry (or the schema change!), but Deckard also allows you to verify that a given migration has been ran against a given database. Simply use `deckard verify ~/path/to/my/migration.up.sql` and deckard will verify that that migration has been ran. A word of warning: We simply check to ensure the metadata table contains a matching entry for the migration. Basically, deckard is only verifying that the migration has been applied in the "UP" position. It won't verify that the schema is currently matching the changes that were introduced in that migration.
