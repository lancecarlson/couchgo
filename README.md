couch.go
========

CouchDB Adapter for Go. Supports BulkSave and emulates couch.js API

API Overview
============

```go
c := NewClient("http://localhost:5984/myleathercouch")

c.CreateDB()

type Cat struct {
  ID string `json:"_id,omitempty"`
  Rev string `json:"_rev,omitempty"`
  Deleted bool `json:"_deleted,omitempty"`
  Name string
  Cool bool
}

cat := Cat{Name: "Octo", Cool: true}

res, err := c.Save(cat)

if err != nil {
  // Do whatever
}

lazyCat := Cat{}

err := c.Get(res.ID, lazyCat)

fmt.Println(lazyCat)

c.Delete(res.ID, res.Rev)

params := url.Values{"limit": []string{"5"}}
results, err := c.View("myapp", "all", &params, nil)
if err != nil {
   // Do whatever
}

fmt.Println(results)

for _, row := range res.Rows {
  cat := &Cat{}
  couch.Remarshal(row.Value, cat)
  fmt.Println(cat)
}
```

TODO (Top to bottom priority)
* _changes
* Attachments
