package models

import (
	"context"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/performance"
)

//MetricPerf vsphere performance
type MetricPerf struct {
	Metric   string `json:"metric"`
	Value    int64  `json:"value"`
	Instance string `json:"instance"`
}

//DatastoreWithURL datastore with url
type DatastoreWithURL struct {
	Datastore string
	URL       string
}

//CounterWithID get all counter list(Key:counter name,Value:counter id)
func CounterWithID(ctx context.Context, c *govmomi.Client) (map[string]int32, error) {
	CounterNameID := make(map[string]int32)
	m := performance.NewManager(c.Client)
	p, err := m.CounterInfoByKey(ctx)
	if err != nil {
		return nil, err
	}
	for _, cc := range p {
		CounterNameID[cc.Name()] = cc.Key
	}
	return CounterNameID, nil
}

//DsWithURL get all datastore with URL
func DsWithURL(ctx context.Context, c *govmomi.Client) ([]DatastoreWithURL, error) {
	dss, err := VsphereDatastores(ctx, c)
	if err != nil {
		return nil, err
	}
	datastoreWithURL := []DatastoreWithURL{}
	if dss != nil {
		for _, ds := range dss {
			datastoreWithURL = append(datastoreWithURL, DatastoreWithURL{Datastore: ds.Summary.Name, URL: ds.Summary.Url})
		}
	}
	return datastoreWithURL, nil
}

//CounterIDByName get counter key by counter name
func CounterIDByName(CounterNameID map[string]int32, Name []string) []int32 {
	IDList := make([]int32, 0)
	for _, eachName := range Name {
		ID, exit := CounterNameID[eachName]
		if exit {
			IDList = append(IDList, ID)
		}
	}
	return IDList
}
