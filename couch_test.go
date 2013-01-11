package couch

import (
	"testing"
	"fmt"
)

const (
	URL = "http://localhost:5984"
	DB = "http://localhost:5984/couch-go-testdb-data"
)

func TestAllDBs(t *testing.T) {
	c, _ := NewClientURL(URL)
	results, err := c.AllDBs()

	fmt.Println(results)

	if err != nil {
		t.Error(err)
	}
	
	if results == nil {
		t.Error("no results")
	}
}

func TestCreateAndDeleteDB(t *testing.T) {
	c, _ := NewClientURL(URL)
	
	res, err := c.CreateDB("couch-go-testdb")
	if err != nil {
		t.Error(err)
	}

	res, err = c.DeleteDB("couch-go-testdb")
	if res["ok"] != true {
		t.Error("Problem deleting")
	}

	if err != nil {
		t.Error(err)
	}
}

func TestSave(t *testing.T) {
	c, _ := NewClientURL(DB)

	id, rev, err := c.Save(map[string]string{"test1": "value1", "test2": "value2"})
	if err != nil {
		t.Error(err)
	}

	if id == "" {
		t.Error(err)
	}

	if rev == "" {
		t.Error(err)
	}
		
}

func TestGetAndSave(t *testing.T) {
	c, _ := NewClientURL(DB)

	doc := map[string]interface{}{}
	err := c.Get("explicit", &doc)
	if err != nil {
		t.Error(err)
	}

	doc["updatekey"] = "updated"

	_, _, err = c.Save(doc)
	if err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	c, _ := NewClientURL(DB)
	
	doc := map[string]string{"_id": "deleteme"}
	id, rev, err := c.Save(doc)
	if err != nil {
		t.Error(err)
	}

	err = c.Delete(id, rev)
	if err != nil {
		t.Error(err)
	}
}

type Cat struct {
	Name string
	Cool bool
}

func TestBulkSave(t *testing.T) {
	c, _ := NewClientURL(DB)
	
	cat1 := Cat{Name: "Hakki", Cool: true}
	cat2 := Cat{Name: "Farb", Cool: false}
	cats := []interface{}{}
	cats = append(cats, cat1, cat2)

	err := c.BulkSave(cats...)
	if err != nil {
		t.Error(err)
	}
}

