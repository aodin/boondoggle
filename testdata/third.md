# Third
[tags]: <> (markdown,sql,go)

Text before.

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, ❤")
}
```

and then some SQL:

```sql
SELECT
    name AS "Name"
FROM "table"
    JOIN "other" ON "other".table_id = "table".id
WHERE name ILIKE 'a%'
ORDER BY LOWER(name)
```

Text after.
