# Deckard

### A framework agnostic tool for running database migrations.

#### Usage
```bash
deckard create add_login_date_to_users
# modify the created files
deckard up --host=localhost --port=5432 --user=user --password==pass --database==app
deckard down --host=localhost --port=5432 --user=user --password==pass --database==app
```

#### Managing your databases via YAML config.
Deckard also supports managing your databases via YAML.
Instead of writing
```bash
deckard up --host=localhost --port=5432 --user=user --password==pass --database==app
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