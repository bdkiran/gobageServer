# gobageServer
Go Server implementation for use with a postgresSQL server.
Required dependencies:
github.com/gorilla/mux
github.com/lib/pq

Database table example has a structure of:
TABLE users (
 	id SERIAL PRIMARY KEY,
 	age INT,
 	first_name TEXT,
 	last_name TEXT,
  email TEXT UNIQUE NOT NULL
)

Please create a tablebase with same structure or make nessecery changes to match your table.
