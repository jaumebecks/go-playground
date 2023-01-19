# My personal Go Playground

## File Generator

### Problem statement

Given a dummy feed data source, create a CSV from its data iteration

### Sources

`file-generator-db` sqlite3 database contains a table `main.feed` which includes

| id_item                 | id_offer | price | title  | brand  | category | in_promo |
|-------------------------|----------|-------|--------|--------|----------|----------|
| int (PK auto-increment) | int      | float | string | string | string   | boolean  |


