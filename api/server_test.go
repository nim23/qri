package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qri-io/qri/repo/test"
)

func TestServerRoutes(t *testing.T) {
	cases := []struct {
		method, endpoint string
		body             []byte
		resStatus        int
	}{
		// {"GET", "/", nil, 200},
		{"GET", "/status", nil, 200},
		// {"GET", "/ipfs/", nil, 200},

		// Profile Routes
		{"GET", "/profile", nil, 200},
		//{"POST", "/profile", nil, 200},
		//{"POST", "/profile/photo", nil, 200},
		//{"PUT", "/profile/photo", nil, 200},
		//{"POST", "/profile/poster", nil, 200},
		//{"PUT", "/profile/poster", nil, 200},

		// Search Routes
		// >>{"GET", "/search", nil, 500},
		//{"GET", "/search", nil, 200},
		// {"POST", "/search", nil, 200},

		// Peer Routes
		{"GET", "/peers", nil, 200}, //PeersHandler
		// >>{"GET", "/peers/", nil, 200},         //PeerHandler
		// >>{"GET", "/connect/", nil, 200},       //ConnectToPeerHandler
		// >>{"GET", "/connections", nil, 200},    //ConnectionsHandler
		// >>{"GET", "/peernamespace/", nil, 200}, //PeerNamespaceHandler

		// Dataset Routes
		{"GET", "/datasets", nil, 200}, // listDatasetsHandler
		// {"POST", "/datasets", nil, 200},    // initDatasetsHandler
		// {"PUT", "/datasets", nil, 200},     // updateDatasetsHandler
		{"GET", "/datasets/", nil, 500}, // getDatasetHandler
		// {"GET", "/datasets/", nil, 200}, // getDatasetHandler
		// {"PUT", "/datasets/", nil, 200},    // updateDatasetHandler
		// {"DELETE", "/datasets/", nil, 200}, // deleteDatasetHandler
		// {"POST", "/add/", nil, 200},        // AddDatasetHandler
		// {"POST", "/init/", nil, 200},       // InitDatasetHandler
		// {"POST", "/rename", nil, 200},      // RenameDatasetHandler
		// {"PUT", "/rename", nil, 200},       // RenameDatasetHandler
		// {"GET", "/data/ipfs/", nil, 200}, // StructuredDataHandler
		// {"GET", "/download/", nil, 200},  // ZipDatasetHandler

		// History Routes
		// {"GET", "/history/", nil, 500}, // LogHandler
		// {"GET", "/history/", nil, 200}, // LogHandler

		// Queries Routes
		{"GET", "/queries", nil, 200}, // ListHandler
		// >>{"GET", "/queries/", nil, 200}, // DatasetQueriesHandler
		// {"POST", "/run", nil, 200},     // RunHandler
	}

	client := &http.Client{}

	r, err := test.NewTestRepo()
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

	for i, c := range cases {
		req, err := http.NewRequest(c.method, server.URL+c.endpoint, bytes.NewReader(c.body))
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

		if err := json.NewDecoder(res.Body).Decode(env); err != nil {
			t.Errorf("case %d: %s - %s error unmarshaling json envelope: %s", i, c.method, c.endpoint, err.Error())
			continue
		}

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
