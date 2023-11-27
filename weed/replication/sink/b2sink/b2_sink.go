package B2Sink

import (
	"context"
	"strings"

	"github.com/seaweedfs/seaweedfs/weed/glog"
	"github.com/seaweedfs/seaweedfs/weed/replication/repl_util"

	"github.com/kurin/blazer/b2"
	"github.com/seaweedfs/seaweedfs/weed/filer"
	"github.com/seaweedfs/seaweedfs/weed/pb/filer_pb"
	"github.com/seaweedfs/seaweedfs/weed/replication/sink"
	"github.com/seaweedfs/seaweedfs/weed/replication/source"
	"github.com/seaweedfs/seaweedfs/weed/util"
)

type B2Sink struct {
	client        *b2.Client
	bucket        string
	dir           string
	filerSource   *source.FilerSource
	isIncremental bool
	isBucketToBucket bool
}

func init() {
	sink.Sinks = append(sink.Sinks, &B2Sink{})
}

func (g *B2Sink) GetName() string {
	return "backblaze"
}

func (g *B2Sink) GetSinkToDirectory() string {
	return g.dir
}

func (g *B2Sink) IsIncremental() bool {
	return g.isIncremental
}
func (g *B2Sink) IsBucketToBucket() bool {
	return false
}


func (g *B2Sink) Initialize(configuration util.Configuration, prefix string) error {
	if configuration.GetBool(prefix + "is_bucket_to_bucket") {
		glog.Warning("is_bucket_to_bucket only works with s3.sink!, It will be ignored")
	}
	g.isIncremental = configuration.GetBool(prefix + "is_incremental")
	return g.initialize(
		configuration.GetString(prefix+"b2_account_id"),
		configuration.GetString(prefix+"b2_master_application_key"),
		configuration.GetString(prefix+"bucket"),
		configuration.GetString(prefix+"directory"),
	)
}

func (g *B2Sink) SetSourceFiler(s *source.FilerSource) {
	g.filerSource = s
}

func (g *B2Sink) initialize(accountId, accountKey, bucket, dir string) error {
	client, err := b2.NewClient(context.Background(), accountId, accountKey)
	if err != nil {
		return err
	}

	g.client = client
	g.bucket = bucket
	g.dir = dir

	return nil
}

func (g *B2Sink) DeleteEntry(key string, isDirectory, deleteIncludeChunks bool, signatures []int32) error {

	key = cleanKey(key)

	if isDirectory {
		key = key + "/"
	}

	bucket, err := g.client.Bucket(context.Background(), g.bucket)
	if err != nil {
		return err
	}

	targetObject := bucket.Object(key)

	return targetObject.Delete(context.Background())

}

func (g *B2Sink) CreateEntry(key string, entry *filer_pb.Entry, signatures []int32) error {

	key = cleanKey(key)

	if entry.IsDirectory {
		return nil
	}

	totalSize := filer.FileSize(entry)
	chunkViews := filer.ViewFromChunks(g.filerSource.LookupFileId, entry.GetChunks(), 0, int64(totalSize))

	bucket, err := g.client.Bucket(context.Background(), g.bucket)
	if err != nil {
		return err
	}

	targetObject := bucket.Object(key)
	writer := targetObject.NewWriter(context.Background())
	defer writer.Close()

	writeFunc := func(data []byte) error {
		_, writeErr := writer.Write(data)
		return writeErr
	}

	if len(entry.Content) > 0 {
		return writeFunc(entry.Content)
	}

	if err := repl_util.CopyFromChunkViews(chunkViews, g.filerSource, writeFunc); err != nil {
		return err
	}

	return nil

}

func (g *B2Sink) UpdateEntry(key string, oldEntry *filer_pb.Entry, newParentPath string, newEntry *filer_pb.Entry, deleteIncludeChunks bool, signatures []int32) (foundExistingEntry bool, err error) {
	key = cleanKey(key)
	return true, g.CreateEntry(key, newEntry, signatures)
}

func cleanKey(key string) string {
	if strings.HasPrefix(key, "/") {
		key = key[1:]
	}
	return key
}
