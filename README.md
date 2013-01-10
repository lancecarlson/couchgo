couch.go
========

CouchDB Adapter for Go. Supports BulkSave and emulates couch.js API

API Overview
============

```
c := NewClient("http://localhost:5984")

c.CreateDB("myleathercouch")

lc := NewClient("http://localhost:5984/myleathcouch")

type Cat struct {
  Name string
  Cool bool
}

cat := Cat{Name: "Octo", Cool: true}

id, rev, err := lc.Save(cat)

if err != nil {
  // Do whatever
}

lazyCat := Cat{}

err := lc.Get(id, lazyCat)

fmt.Println(lazyCat)

lc.Delete(id, rev)
````

TODO: BulkSave :)
