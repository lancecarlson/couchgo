package couch

import (
	"testing"
	"fmt"
	"net/url"
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

type Cow struct {
	ID string `json:"_id"`
	Rev string `json:"_rev"`
	Name string
}

func TestSaveWithId(t *testing.T) {
	c, _ := NewClientURL(DB)
	
	cow := new(Cow)
	c.Get("testcow", cow)
	c.Delete("testcow", cow.Rev)
	fmt.Println(cow)
	_, _, err := c.Save(Cow{ID: "testcow", Name: "Fred"})
	if err != nil {
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
	ID string `json:"_id"`
	Rev string `json:"_rev"`
	Name string
	Cool bool
}

func TestBulkSave(t *testing.T) {
	c, _ := NewClientURL(DB)
	
	cat1 := Cat{Name: "Hakki", Cool: true}
	cat2 := Cat{Name: "Farb", Cool: false}
	cats := []interface{}{}
	cats = append(cats, cat1, cat2)

	resp, err := c.BulkSave(cats...)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)
}

/*
type Dog struct {
	Name string
}

func TestView(t *testing.T) {
	c, _ := NewClientURL(DB)

	dogs := []Row{{Value: Dog{}}}
	err := c.View("dog", "all", nil, dogs)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(dogs)
}*/

func TestViewRaw(t *testing.T) {
	c, _ := NewClientURL(DB)

	params := url.Values{"limit": []string{"2"}}
	res, err := c.ViewRaw("dog", "all", &params)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}

func TestCopy(t *testing.T) {
	c, _ := NewClientURL(DB)
	
	res, err := c.Copy("explicit", "explicit-copy", nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}