package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	// "github.com/qri-io/qri/core"
	// "github.com/qri-io/qri/repo"
	testrepo "github.com/qri-io/qri/repo/test"
)

func loadTestdata(path string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join("testdata", path))
}

func TestServerRoutes(t *testing.T) {
	// jobsByAutomationFile := testrepo.JobsByAutomationFile
	r, err := testrepo.NewTestRepo()
	if err != nil {
		t.Errorf("error allocating test repo: %s", err.Error())
		return
	}
	s, err := New(r, func(opt *Config) {
		opt.Online = false
		opt.MemOnly = true
	})
	if err != nil {
		t.Error(err.Error())
		return
	}

	server := httptest.NewServer(NewServerRoutes(s))
	moviesJSONData, err := loadTestdata("movies_search_query.json")
	if err != nil {
		t.Error(err.Error())
		return
	}
	sqlJSONData, err := loadTestdata("sql_statement.json")
	if err != nil {
		t.Error(err.Error())
		return
	}

	moviesPath, err := r.GetPath("movies")
	if err != nil {
		t.Errorf("error getting path: %s", err.Error())
	}
	moviesPathString := fmt.Sprintf("%s", moviesPath)

	// initJSONData, err := loadTestdata("new_dataset.json")
	// if err != nil {
	// 	t.Error(err.Error())
	// 	return
	// }
	// (get inputs for update params)
	/*
			getParams := &core.GetDatasetParams{
				Name: movies,
				Path: moviesPath,
			}
			req, err := http.NewRequest("GET", server.URL+"/datasets"+, body)

			req, err := http.NewRequest(c.method, server.URL+c.endpoint+c.queryString, bytes.NewReader(c.body))
			if err != nil {
				t.Errorf("case %d error creating request: %s", i, err.Error())
				continue
			}

			prev := &repo.DatasetRef{}
			req, err := cmd.datasetRequests(false)
			if err != nil {
				t.Error(err.Error())
				return
			}
			err = req.Get(getParams, prev)
			if err != nil {
				t.Error(err.Error())
				return
			}
			author, err := r.Profile()
			if err != nil {
				t.Error(err.Error())
				return
			}
		update := &core.UpdateParams{}
	*/

	// citiesPath, err := r.GetPath("cities")
	// if err != nil {
	// 	t.Errorf("error getting path: %s", err.Error())
	// }

	cases := []struct {
		method, endpoint string
		queryString      string
		body             []byte
		resStatus        int
	}{
		// {"GET", "/", "", nil, 200},
		{"GET", "/status", "", nil, 200},
		// {"GET", "/ipfs/", "", nil, 200},

		// Profile Routes
		{"GET", "/profile", "", nil, 200},             // getProfileHandler
		{"POST", "/profile", "", moviesJSONData, 200}, // saveProfileHandler
		// {"POST", "/profile/photo", "", nil, 200}, // SetProfilePhotoHandler
		//{"PUT", "/profile/photo", "", nil, 200}, // SetProfilePhotoHandler
		//{"POST", "/profile/poster", "", nil, 200}, // SetPosterHandler
		//{"PUT", "/profile/poster", "", nil, 200}, // SetPosterHandler

		// Search Routes
		// {"GET", "/search", "?q=movies", nil, 200}, // <-- works on curl, fails here
		// {"POST", "/search", "", moviesJSONData, 200}, // <-- works on curl, fails here

		// Peer Routes
		{"GET", "/peers", "", nil, 200}, //PeersHandler
		// {"GET", "/peers/", moviesPathString, nil, 500}, //PeerHandler
		// >>{"GET", "/connect/", "", nil, 200},       //ConnectToPeerHandler
		// {"GET", "/connections", "", nil, 200}, //ConnectionsHandler
		// >>{"GET", "/peernamespace/", "", nil, 200}, //PeerNamespaceHandler

		// Dataset Routes
		{"GET", "/datasets", "", nil, 200}, // listDatasetsHandler
		// {"POST", "/datasets", "", nil, 200},    // initDatasetsHandler
		// {"PUT", "/datasets", "", nil, 200},     // updateDatasetsHandler
		// {"OPTIONS", "/datasets/", moviesPathString, nil, 200}, // getDatasetHandler
		{"GET", "/datasets/", moviesPathString, nil, 200}, // getDatasetHandler
		// {"POST", "/init/", "", initJSONData, 200},         // InitDatasetHandler
		// {"POST", "/rename", "", nil, 200},      // RenameDatasetHandler
		// {"POST", "/add/", "", nil, 200},        // AddDatasetHandler
		// {"PUT", "/rename", "", nil, 200},       // RenameDatasetHandler
		// {"DELETE", "/datasets/", "", nil, 200}, // deleteDatasetHandler
		// {"PUT", "/datasets/", "", nil, 200},    // updateDatasetHandler
		// {"GET", "/data/ipfs/", "", nil, 200}, // StructuredDataHandler
		// {"GET", "/download/", "", nil, 200},  // ZipDatasetHandler

		// History Routes
		// {"OPTIONS", "/history/", moviesPathString, nil, 200}, // LogHandler
		{"GET", "/history/", moviesPathString, nil, 200},  // LogHandler
		{"POST", "/history/", moviesPathString, nil, 200}, // LogHandler

		// Queries Routes
		{"GET", "/queries", "", nil, 200},                // ListHandler
		{"GET", "/queries/", moviesPathString, nil, 200}, // DatasetQueriesHandler
		{"POST", "/run", "", sqlJSONData, 200},           // RunHandler
	}

	client := &http.Client{}

	// fmt.Printf("the path for movies is: %s\n", path)
	/*
		mr, err := testrepo.NewTestRepo()
		if err != nil {
			t.Errorf("error allocating test repo: %s", err.Error())
			return
		}
		path, err := mr.GetPath("movies")
		if err != nil {
			t.Errorf("error getting path: %s", err.Error())
			return
		}
	*/

	for i, c := range cases {
		req, err := http.NewRequest(c.method, server.URL+c.endpoint+c.queryString, bytes.NewReader(c.body))
		if err != nil {
			t.Errorf("case %d error creating request: %s", i, err.Error())
			continue
		}

		res, err := client.Do(req)
		if err != nil {
			t.Errorf("case %d error performing request: %s", i, err.Error())
			continue
		}

		if res.StatusCode != c.resStatus {
			t.Errorf("case %d: %s - %s status code mismatch. expected: %d, got: %d", i, c.method, c.endpoint, c.resStatus, res.StatusCode)
			continue
		}

		env := &struct {
			Meta       map[string]interface{}
			Data       interface{}
			Pagination map[string]interface{}
		}{}
		// fmt.Printf("body: '%s'\n", res.Body)
		// fmt.Println(res.Body)
		if err := json.NewDecoder(res.Body).Decode(env); err != nil {
			fmt.Println("----")
			fmt.Printf("%s\n", res.Body)
			fmt.Println("----")
			t.Errorf("case %d: %s - %s error unmarshaling json envelope: %s", i, c.method, c.endpoint, err.Error())
			continue
		}
		// fmt.Println(env)
		if env.Meta == nil {
			t.Errorf("case %d: %s - %s doesn't have a meta field", i, c.method, c.endpoint)
			continue
		}
		if env.Data == nil {
			t.Errorf("case %d: %s - %s doesn't have a data field", i, c.method, c.endpoint)
			continue
		}
	}
}
