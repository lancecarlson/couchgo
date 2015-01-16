package couch

import (
	"fmt"
	"net/url"
	"testing"
)

const (
	Host           = "http://localhost:5984"
	ReplicateSrcDB = Host + "/couch-go-replicate-src"
	ReplicateTarDB = Host + "/couch-go-replicate-tar"
	CreateDeleteDB = Host + "/couch-go-delete-db"
	DB             = Host + "/couch-go-testdb-data"
)

func newClientURL(url string) (client *Client, err error) {
	client, _ = NewClientURL(url)
	if url != Host {
		client.DeleteDB()
		client.CreateDB()
	}
	return
}

func TestAllDBs(t *testing.T) {
	c, _ := newClientURL(Host)
	if results, err := c.AllDBs(); err != nil || len(results) == 0 {
		t.Fatal(err, results)
	}

	c, _ = newClientURL(DB)
	if results, err := c.AllDBs(); err != nil || len(results) == 0 {
		t.Fatal(err, results)
	}
}

func TestAllDocs(t *testing.T) {
	c, _ := newClientURL(DB)
	if ids, err := c.AllDocs(); err != nil {
		t.Fatal(ids, err)
	} else if len(ids) != 0 {
		t.Fatal(ids)
	}

	id := "a0000000"
	obj := map[string]interface{}{
		"id":   id,
		"name": "me",
	}
	c.Add(id, obj)
	if ids, err := c.AllDocs(); err != nil {
		t.Fatal(ids, err)
	} else if len(ids) != 1 {
		t.Fatal(ids)
	} else if ids[0] != id {
		t.Fatal(ids)
	}
}

func TestCreateAndDeleteDB(t *testing.T) {
	c, _ := newClientURL(CreateDeleteDB)
	c.DeleteDB()
	if res, code, err := c.CreateDB(); err != nil {
		t.Fatal(err, code, res)
	}

	if res, code, err := c.DeleteDB(); err != nil || !res.Ok {
		t.Fatal(res, code, err)
	}
}

func TestSave(t *testing.T) {
	c, _ := newClientURL(DB)
	if res, err := c.Save(map[string]string{"test1": "value1", "test2": "value2"}); err != nil {
		t.Fatal(res, err)
	}
}

type Cow struct {
	ID   string `json:"_id,omitempty"`
	Rev  string `json:"_rev,omitempty"`
	Name string
}

func TestSaveWithId(t *testing.T) {
	c, _ := newClientURL(DB)

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
	c, _ := newClientURL(DB)
	ids := map[string]interface{}{
		"_id":      "explicit",
		"explicit": "explicite",
	}
	c.Save(ids)
	if err := c.Get("explicit", &ids); err != nil {
		t.Fatal(err)
	}

	ids["updatekey"] = "updated"

	if res, err := c.Save(ids); err != nil {
		t.Fatal(res, err)
	}
}

func TestDelete(t *testing.T) {
	c, _ := newClientURL(DB)

	ids := map[string]string{"_id": "deleteme"}
	res, err := c.Save(ids)
	if err != nil {
		t.Error(err)
	}
	err = c.Delete(res.ID, res.Rev)
	if err != nil {
		t.Error(err)
	}
}

type Cat struct {
	ID   string `json:"_id,omitempty"`
	Rev  string `json:"_rev,omitempty"`
	Name string
	Cool bool
}

func TestBulkSave(t *testing.T) {
	c, _ := newClientURL(DB)

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
	ID   string `json:"_id,omitempty"`
	Rev  string `json:"_rev,omitempty"`
	Type string `json:"type"`
	Name string
}

func TestView(t *testing.T) {
	c, _ := newClientURL(DB)

	c.Save(&Dog{Name: "Savannah", Type: "dog"})

	params := url.Values{"limit": []string{"5"}}
	if res, err := c.View("dog", "dog", &params, nil); err != nil {
		t.Fatal(res, err)
	}
	dog1 := Dog{ID: "dog1", Type: "dog"}
	dog2 := Dog{ID: "dog2", Type: "dog"}
	c.BulkSave(dog1, dog2)

	if res, err := c.View("dog", "dog", nil, &[]string{"dog1", "dog2"}); err != nil {
		t.Fatal(res, err)
	}
}

func TestCopy(t *testing.T) {
	c, _ := newClientURL(DB)

	if res, _, err := c.Copy("explicit", "explicit-copy", nil); err == nil {
		t.Fatal(res)
	}
}

func TestReplicate(t *testing.T) {
	c, _ := newClientURL(Host)
	newClientURL(ReplicateSrcDB)
	newClientURL(ReplicateTarDB)
	req := &ReplicateRequest{
		Source: ReplicateSrcDB,
		Target: ReplicateTarDB,
	}
	if res, _, err := c.Replicate(req); err != nil {
		t.Fatal(res, err)
	}
}
