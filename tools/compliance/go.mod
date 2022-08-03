<<<<<<< HEAD   (0ff8c6 Merge "Merge empty history for sparse-8814173-L7490000095568)
=======
module android/soong/tools/compliance

require google.golang.org/protobuf v0.0.0

replace google.golang.org/protobuf v0.0.0 => ../../../../external/golang-protobuf

require (
	android/soong v0.0.0
	github.com/google/blueprint v0.0.0
)

require golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f // indirect

replace android/soong v0.0.0 => ../../../soong

replace github.com/google/blueprint => ../../../blueprint

// Indirect deps from golang-protobuf
exclude github.com/golang/protobuf v1.5.0

replace github.com/google/go-cmp v0.5.5 => ../../../../external/go-cmp

// Indirect dep from go-cmp
exclude golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543

go 1.18
>>>>>>> BRANCH (deebb4 Merge "Version bump to TKB1.220803.001.A1 [core/build_id.mk])
