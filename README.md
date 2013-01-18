couch.go
========

CouchDB Adapter for Go. Supports BulkSave and emulates couch.js API

API Overview
============

```go
c := NewClient("http://localhost:5984")

c.CreateDB("myleathercouch")

lc := NewClient("http://localhost:5984/myleathcouch")

type Cat struct {
  ID string `json:"_id,omitempty"`
  Rev string `json:"_rev,omitempty"`
  Deleted bool `json:"_deleted,omitempty"`
  Name string
  Cool bool
}

cat := Cat{Name: "Octo", Cool: true}

res, err := lc.Save(cat)

if err != nil {
  // Do whatever
}

lazyCat := Cat{}

err := lc.Get(res.ID, lazyCat)

fmt.Println(lazyCat)

lc.Delete(res.ID, res.Rev)

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
* Replication
* _changes
* Attachments
