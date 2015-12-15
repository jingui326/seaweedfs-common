package operation

import (
	"errors"
	"strconv"
	"time"

	"github.com/chrislusf/seaweedfs/go/glog"
)

type VidInfo struct {
	Locations       Locations
	NextRefreshTime time.Time
}
type VidCache struct {
	cache []VidInfo
}

func (vc *VidCache) Get(vid string) (Locations, error) {
	id, err := strconv.Atoi(vid)
	if err != nil {
		glog.V(1).Infof("Unknown volume id %s", vid)
		return nil, err
	}
	if 0 < id && id <= len(vc.cache) {
		if vc.cache[id-1].Locations == nil {
			return nil, errors.New("Not Set")
		}
		if vc.cache[id-1].NextRefreshTime.Before(time.Now()) {
			return nil, errors.New("Expired")
		}
		return vc.cache[id-1].Locations, nil
	}
	return nil, errors.New("Not Found")
}
func (vc *VidCache) Set(vid string, locations Locations, duration time.Duration) {
	id, err := strconv.Atoi(vid)
	if err != nil {
		glog.V(1).Infof("Unknown volume id %s", vid)
		return
	}
	if id > len(vc.cache) {
		for i := id - len(vc.cache); i > 0; i-- {
			vc.cache = append(vc.cache, VidInfo{})
		}
	}
	if id > 0 {
		vc.cache[id-1].Locations = locations
		vc.cache[id-1].NextRefreshTime = time.Now().Add(duration)
	}
}
