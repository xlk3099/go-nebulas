// Copyright (C) 2017 go-nebulas authors
//
// This file is part of the go-nebulas library.
//
// the go-nebulas library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// the go-nebulas library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-nebulas library.  If not, see <http://www.gnu.org/licenses/>.
//

package metrics

import (
	"runtime"
	"time"

	"github.com/nebulasio/go-nebulas/neblet/pb"
	"github.com/nebulasio/go-nebulas/net/p2p"
	metrics "github.com/rcrowley/go-metrics"
	influxdb "github.com/vrischmann/go-metrics-influxdb"
)

const (
	duration = 2 * time.Second
	tagName  = "nodeID"
)

// Neblet interface breaks cycle import dependency.
type Neblet interface {
	Config() nebletpb.Config
	NetService() *p2p.NetService
}

// Start metrics mornitor
func Start(neb Neblet) {
	tags := make(map[string]string)
	tags[tagName] = neb.NetService().Node().ID()
	go systemMetrics()
	influxdb.InfluxDBWithTags(metrics.DefaultRegistry, duration, neb.Config().Influxdb.Host, neb.Config().Influxdb.Db, neb.Config().Influxdb.Username, neb.Config().Influxdb.Password, tags)
}

func systemMetrics() {
	memstats := make([]*runtime.MemStats, 2)
	for i := 0; i < len(memstats); i++ {
		memstats[i] = new(runtime.MemStats)
	}

	allocs := metrics.GetOrRegisterMeter("system_allocs", nil)
	// totalAllocs := metrics.GetOrRegisterMeter("system_total_allocs", nil)
	sys := metrics.GetOrRegisterMeter("system_sys", nil)
	frees := metrics.GetOrRegisterMeter("system_frees", nil)
	heapInuse := metrics.GetOrRegisterMeter("system_heapInuse", nil)
	stackInuse := metrics.GetOrRegisterMeter("system_stackInuse", nil)

	for i := 1; ; i++ {
		runtime.ReadMemStats(memstats[i%2])
		allocs.Mark(int64(memstats[i%2].Alloc - memstats[(i-1)%2].Alloc))
		sys.Mark(int64(memstats[i%2].Sys - memstats[(i-1)%2].Sys))
		frees.Mark(int64(memstats[i%2].Frees - memstats[(i-1)%2].Frees))
		heapInuse.Mark(int64(memstats[i%2].HeapInuse - memstats[(i-1)%2].HeapInuse))
		stackInuse.Mark(int64(memstats[i%2].StackInuse - memstats[(i-1)%2].StackInuse))
		time.Sleep(2 * time.Second)
	}

}
