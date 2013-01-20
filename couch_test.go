package couch

import (
	"testing"
	"fmt"
	"net/url"
)

const (
	Host = "http://localhost:5984"
	CreateDeleteDB = Host + "/couch-go-delete-db"
	DB = Host + "/couch-go-testdb-data"
)

func TestAllDBs(t *testing.T) {
	c, _ := NewClientURL(Host)
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
	c, _ := NewClientURL(CreateDeleteDB)
	
	res, _, err := c.CreateDB()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)

	res, _, err = c.DeleteDB()
	if res.Ok != true {
		t.Error("Problem deleting")
	}
	fmt.Println(res)

	if err != nil {
		t.Error(err)
	}
}

func TestSave(t *testing.T) {
	c, _ := NewClientURL(DB)

	res, err := c.Save(map[string]string{"test1": "value1", "test2": "value2"})
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Save")
	fmt.Println(res)
}

type Cow struct {
	ID string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	Name string
}

func TestSaveWithId(t *testing.T) {
	c, _ := NewClientURL(DB)
	
	cow := new(Cow)
	c.Get("testcow", cow)
	c.Delete("testcow", cow.Rev)
	fmt.Println(cow)
	res, err := c.Save(Cow{ID: "testcow", Name: "Fred"})
	if err != nil {
		t.Error(err)
	}
	fmt.Println("SaveWithId")
	fmt.Println(res)
}

func TestGetAndSave(t *testing.T) {
	c, _ := NewClientURL(DB)

	doc := map[string]interface{}{}
	err := c.Get("explicit", &doc)
	if err != nil {
		t.Error(err)
	}

	doc["updatekey"] = "updated"

	res, err := c.Save(doc)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("GetAndSave")
	fmt.Println(res)
}

func TestDelete(t *testing.T) {
	c, _ := NewClientURL(DB)
	
	doc := map[string]string{"_id": "deleteme"}
	res, err := c.Save(doc)
	if err != nil {
		t.Error(err)
	}
	err = c.Delete(res.ID, res.Rev)
	if err != nil {
		t.Error(err)
	}
}

type Cat struct {
	ID string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	Name string
	Cool bool
}

func TestBulkSave(t *testing.T) {
	c, _ := NewClientURL(DB)
	
	cat1 := Cat{Name: "Hakki", Cool: true}
	cat2 := Cat{Name: "Farb", Cool: false}
	cats := []interface{}{}
	cats = append(cats, cat1, cat2)

	resp, _, err := c.BulkSave(cats...)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("BulkSave")
	fmt.Println(resp)
	fmt.Println("/BulkSave")
}

type Dog struct {
	ID string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	Type string `json:"type"`
	Name string
}

func TestView(t *testing.T) {
	c, _ := NewClientURL(DB)

	c.Save(&Dog{Name: "Savannah", Type: "dog"})

	params := url.Values{"limit": []string{"5"}}
	res, err := c.View("dog", "dog", &params, nil)
	if err != nil {
		t.Error(err)
	}

	for _, row := range res.Rows {
		dog := &Dog{}
		Remarshal(row.Value, dog)
		fmt.Println(dog)
	}

	fmt.Println("View")
	fmt.Println(res)
	fmt.Println("View")

	dog1 := Dog{ID: "dog1", Type: "dog"}
	dog2 := Dog{ID: "dog2", Type: "dog"}
	c.BulkSave(dog1, dog2)
	
	resp, err := c.View("dog", "dog", nil, &[]string{"dog1", "dog2"})
	if err != nil {
		t.Error(err)
	}

	fmt.Println("ViewWithKeys")
	fmt.Println(resp)
	fmt.Println("ViewWithKeys")
}

func TestCopy(t *testing.T) {
	c, _ := NewClientURL(DB)
	
	res, _, err := c.Copy("explicit", "explicit-copy", nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}