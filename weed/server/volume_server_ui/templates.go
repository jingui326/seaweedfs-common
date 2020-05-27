package master_ui

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
)

const unit = 1024

// BytesToHumanReadable exposes the bytesToHumanReadble.
// It returns the converted human readable representation of the bytes
// and the conversion state that true for conversion realy occurred or false for no conversion.
func BytesToHumanReadable(b uint64) (string, bool) {
	return bytesToHumanReadble(b), b >= unit
}

func bytesToHumanReadble(b uint64) string {
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

func percentFrom(total uint64, part_of uint64) string {
	return fmt.Sprintf("%.2f", (float64(part_of)/float64(total))*100)
}

func join(data []int64) string {
	var ret []string
	for _, d := range data {
		ret = append(ret, strconv.Itoa(int(d)))
	}
	return strings.Join(ret, ",")
}

var funcMap = template.FuncMap{
	"join":                join,
	"bytesToHumanReadble": bytesToHumanReadble,
	"percentFrom":         percentFrom,
}

var StatusTpl = template.Must(template.New("status").Funcs(funcMap).Parse(`<!DOCTYPE html>
<html>
  <head>
    <title>SeaweedFS {{ .Version }}</title>
	<link rel="stylesheet" href="/seaweedfsstatic/bootstrap/3.3.1/css/bootstrap.min.css">
	<script type="text/javascript" src="/seaweedfsstatic/javascript/jquery-2.1.3.min.js"></script>
	<script type="text/javascript" src="/seaweedfsstatic/javascript/jquery-sparklines/2.1.2/jquery.sparkline.min.js"></script>
    <script type="text/javascript">
    $(function() {
		var periods = ['second', 'minute', 'hour', 'day'];
		for (i = 0; i < periods.length; i++) { 
		    var period = periods[i];
	        $('.inlinesparkline-'+period).sparkline('html', {
				type: 'line', 
				barColor: 'red', 
				tooltipSuffix:' request per '+period,
			}); 
		}
    });
    </script>
	<style>
	#jqstooltip{
		height: 28px !important;
		width: 150px !important;
	}
	</style>
  </head>
  <body>
    <div class="container">
      <div class="page-header">
	    <h1>
          <a href="https://github.com/chrislusf/seaweedfs"><img src="/seaweedfsstatic/seaweed50x50.png"></img></a>
          SeaweedFS <small>{{ .Version }}</small>
	    </h1>
      </div>

      <div class="row">
        <div class="col-sm-6">          
    	  <h2>Disk Stats</h2>
          <table class="table table-striped">
            <thead>
              <tr>
              <th>Path</th>
              <th>Total</th>
              <th>Free</th>
              <th>Usage</th>
              </tr>
            </thead>
          <tbody>          
          {{ range .DiskStatuses }}
            <tr>
              <td>{{ .Dir }}</td>
              <td>{{ bytesToHumanReadble .All }}</td>
              <td>{{ bytesToHumanReadble .Free  }}</td>
              <td>{{ percentFrom .All .Used}}%</td>
            </tr>
          {{ end }}
          </tbody>
          </table>
        </div>

        <div class="col-sm-6">
          <h2>System Stats</h2>
          <table class="table table-condensed table-striped">
            <tr>
              <th>Masters</th>
              <td>{{.Masters}}</td>
            </tr>
            <tr>
              <th>Weekly # ReadRequests</th>
              <td><span class="inlinesparkline-day">{{ .Counters.ReadRequests.WeekCounter.ToList | join }}</span></td>
            </tr>
            <tr>
              <th>Daily # ReadRequests</th>
              <td><span class="inlinesparkline-hour">{{ .Counters.ReadRequests.DayCounter.ToList | join }}</span></td>
            </tr>
            <tr>
              <th>Hourly # ReadRequests</th>
              <td><span class="inlinesparkline-minute">{{ .Counters.ReadRequests.HourCounter.ToList | join }}</span></td>
            </tr>
            <tr>
              <th>Last Minute # ReadRequests</th>
              <td><span class="inlinesparkline-second">{{ .Counters.ReadRequests.MinuteCounter.ToList | join }}</span></td>
            </tr>
          {{ range $key, $val := .Stats }}
            <tr>
              <th>{{ $key }}</th>
              <td>{{ $val }}</td>
            </tr>
          {{ end }}
          </table>
        </div>
      </div>

      <div class="row">
        <h2>Volumes</h2>
        <table class="table table-striped">
          <thead>
            <tr>
              <th>Id</th>
              <th>Collection</th>
              <th>Data Size</th>
              <th>Files</th>
              <th>Trash</th>
              <th>TTL</th>
              <th>ReadOnly</th>
            </tr>
          </thead>
          <tbody>
          {{ range .Volumes }}
            <tr>
              <td><code>{{ .Id }}</code></td>
              <td>{{ .Collection }}</td>
              <td>{{ bytesToHumanReadble .Size }}</td>
              <td>{{ .FileCount }}</td>
              <td>{{ .DeleteCount }} / {{bytesToHumanReadble .DeletedByteCount}}</td>
              <td>{{ .Ttl }}</td>
              <td>{{ .ReadOnly }}</td>
            </tr>
          {{ end }}
          </tbody>
        </table>
      </div>

      <div class="row">
        <h2>Remote Volumes</h2>
        <table class="table table-striped">
          <thead>
            <tr>
              <th>Id</th>
              <th>Collection</th>
              <th>Size</th>
              <th>Files</th>
              <th>Trash</th>
              <th>Remote</th>
              <th>Key</th>
            </tr>
          </thead>
          <tbody>
          {{ range .RemoteVolumes }}
            <tr>
              <td><code>{{ .Id }}</code></td>
              <td>{{ .Collection }}</td>
              <td>{{ bytesToHumanReadble .Size }}</td>
              <td>{{ .FileCount }}</td>
              <td>{{ .DeleteCount }} / {{bytesToHumanReadble .DeletedByteCount}}</td>
              <td>{{ .RemoteStorageName }}</td>
              <td>{{ .RemoteStorageKey }}</td>
            </tr>
          {{ end }}
          </tbody>
        </table>
      </div>

      <div class="row">
        <h2>Erasure Coding Shards</h2>
        <table class="table table-striped">
          <thead>
            <tr>
              <th>Id</th>
              <th>Collection</th>
              <th>Shard Size</th>
              <th>Shards</th>
              <th>CreatedAt</th>
            </tr>
          </thead>
          <tbody>
          {{ range .EcVolumes }}
            <tr>
              <td><code>{{ .VolumeId }}</code></td>
              <td>{{ .Collection }}</td>
              <td>{{ bytesToHumanReadble .ShardSize }}</td>
              <td>{{ .ShardIdList }}</td>
              <td>{{ .CreatedAt.Format "02 Jan 06 15:04 -0700" }}</td>
            </tr>
          {{ end }}
          </tbody>
        </table>
      </div>

    </div>
  </body>
</html>
`))
