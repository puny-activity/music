package util

func ToPointer[V any](v V) *V {
	return &v
}
