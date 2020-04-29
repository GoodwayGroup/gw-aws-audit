package lib

type Metrics struct {
	Volumes         int
	SumVolumeSize   int64
	VolumeCosts     int
	Snapshots       int
	SumSnapshotSize int64
	Processed       int64
	Skipped         int64
	Modified        int64
}
