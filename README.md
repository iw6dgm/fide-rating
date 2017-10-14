FIDE Rating parser
==================

Parse the standard XML full player rating list provided by FIDE and save data into a SQLite3 database.

Setup
-----

Create a SQLite3 database `fide.db` and create `player` table by using:

```sql
CREATE TABLE IF NOT EXISTS player (
  fideid bigint not null primary key,
  name text, country char(3),
  sex char(1),
  title text,
  w_title text,
  o_title text,
  foa_title text,
  rating integer,
  games integer,
  k smallint,
  rapid_rating integer,
  rapid_games integer,
  rapid_k smallint,
  blitz_rating integer,
  blitz_games integer,
  blitz_k smallint,
  birthday integer,
  flag text
)
```

Download XML file from FIDE website and unpack it. File name should be `players_list_xml_foa.xml`.
Download and install Go SQLite3 dependency as described in [here](https://github.com/mattn/go-sqlite3).
Compile the source code and run the executable with no parameters.

```bash
go build fideparser.go
./fideparser
```

**Warning**: The overall process is very time-consuming and the resulting database is huge.

TODO
----

FIDE Rating API: Web app to provide friendly APIs to query the FIDE dabatase created by the parser above.
