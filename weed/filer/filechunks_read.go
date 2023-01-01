package filer

import (
	"container/list"
	"github.com/seaweedfs/seaweedfs/weed/pb/filer_pb"
	"golang.org/x/exp/slices"
)

// readResolvedChunks returns a container.List of VisibleInterval
func readResolvedChunks(chunks []*filer_pb.FileChunk, startOffset int64, stopOffset int64) (visibles *IntervalList[VisibleInterval]) {

	var points []*Point
	for _, chunk := range chunks {
		if chunk.IsChunkManifest {
			println("This should not happen! A manifest chunk found:", chunk.GetFileIdString())
		}
		start, stop := max(chunk.Offset, startOffset), min(chunk.Offset+int64(chunk.Size), stopOffset)
		if start >= stop {
			continue
		}
		points = append(points, &Point{
			x:       chunk.Offset,
			ts:      chunk.ModifiedTsNs,
			chunk:   chunk,
			isStart: true,
		})
		points = append(points, &Point{
			x:       chunk.Offset + int64(chunk.Size),
			ts:      chunk.ModifiedTsNs,
			chunk:   chunk,
			isStart: false,
		})
	}
	slices.SortFunc(points, func(a, b *Point) bool {
		if a.x != b.x {
			return a.x < b.x
		}
		if a.ts != b.ts {
			return a.ts < b.ts
		}
		return !a.isStart
	})

	var prevX int64
	queue := list.New() // points with higher ts are at the tail
	visibles = NewIntervalList[VisibleInterval]()
	var prevPoint *Point
	for _, point := range points {
		if queue.Len() > 0 {
			prevPoint = queue.Back().Value.(*Point)
		} else {
			prevPoint = nil
		}
		if point.isStart {
			if prevPoint != nil {
				if point.x != prevX && prevPoint.ts < point.ts {
					addToVisibles(visibles, prevX, prevPoint, point)
					prevX = point.x
				}
			}
			// insert into queue
			if prevPoint == nil || prevPoint.ts < point.ts {
				queue.PushBack(point)
				prevX = point.x
			} else {
				for e := queue.Front(); e != nil; e = e.Next() {
					if e.Value.(*Point).ts > point.ts {
						queue.InsertBefore(point, e)
						break
					}
				}
			}
		} else {
			isLast := true
			for e := queue.Back(); e != nil; e = e.Prev() {
				if e.Value.(*Point).ts == point.ts {
					queue.Remove(e)
					break
				}
				isLast = false
			}
			if isLast && prevPoint != nil {
				addToVisibles(visibles, prevX, prevPoint, point)
				prevX = point.x
			}
		}
	}

	return
}

func addToVisibles(visibles *IntervalList[VisibleInterval], prevX int64, startPoint *Point, point *Point) {
	if prevX < point.x {
		chunk := startPoint.chunk
		visible := VisibleInterval{
			start:        prevX,
			stop:         point.x,
			fileId:       chunk.GetFileIdString(),
			modifiedTsNs: chunk.ModifiedTsNs,
			chunkOffset:  prevX - chunk.Offset,
			chunkSize:    chunk.Size,
			cipherKey:    chunk.CipherKey,
			isGzipped:    chunk.IsCompressed,
		}
		visibles.AppendInterval(&Interval[VisibleInterval]{
			StartOffset: visible.start,
			StopOffset:  visible.stop,
			TsNs:        visible.modifiedTsNs,
			Value:       visible,
		})
	}
}

type Point struct {
	x       int64
	ts      int64
	chunk   *filer_pb.FileChunk
	isStart bool
}
