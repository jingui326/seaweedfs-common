package operation

import (
	"context"
	"github.com/chrislusf/seaweedfs/weed/pb"
	"github.com/chrislusf/seaweedfs/weed/pb/volume_server_pb"
	"google.golang.org/grpc"
)

func GetVolumeSyncStatus(server pb.ServerAddress, grpcDialOptions []grpc.DialOption, vid uint32) (resp *volume_server_pb.VolumeSyncStatusResponse, err error) {

	WithVolumeServerClient(false, server, grpcDialOptions, func(client volume_server_pb.VolumeServerClient) error {

		resp, err = client.VolumeSyncStatus(context.Background(), &volume_server_pb.VolumeSyncStatusRequest{
			VolumeId: vid,
		})
		return nil
	})

	return
}
