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
  Name string
  Cool bool
}

cat := Cat{Name: "Octo", Cool: true}

resp, err := lc.Save(cat)

if err != nil {
  // Do whatever
}

lazyCat := Cat{}

err := lc.Get(res.ID, lazyCat)

fmt.Println(lazyCat)

lc.Delete(res.ID, res.Rev)

params := url.Values{"limit": []string{"5"}}
results, err := c.ViewRaw("myapp", "all", &params)
if err != nil {
   // Do whatever
}

fmt.Println(results)
```

TODO (Top to bottom priority)
* Replication
* _changes
* Attachments
