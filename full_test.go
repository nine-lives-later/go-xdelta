package xdelta

import (
	"bytes"
	"context"
	"crypto/sha1"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
)

type testFullRoundtrip_Context struct {
	FromFilePath    string
	ToFilePath      string
	PatchFilePath   string
	AppliedFilePath string

	Header     []byte // set during seeding
	ToFileHash []byte
}

func TestFullRoundtrip(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	// get temporary directory
	tempDir, err := ioutil.TempDir("", "go-xdelta")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	defer os.RemoveAll(tempDir)

	ctx := &testFullRoundtrip_Context{
		FromFilePath:    filepath.Join(tempDir, "from"),
		ToFilePath:      filepath.Join(tempDir, "to"),
		PatchFilePath:   filepath.Join(tempDir, "patch"),
		AppliedFilePath: filepath.Join(tempDir, "to_applied"),
	}

	t.Run("Seed", func(t *testing.T) { testFullRoundtrip_Seed(t, ctx) })
	t.Run("CreatePatch", func(t *testing.T) { testFullRoundtrip_CreatePatch(t, ctx) })
	t.Run("DumpPatchInfo", func(t *testing.T) { testFullRoundtrip_DumpPatchInfo(t, ctx) })
	t.Run("ApplyPatch", func(t *testing.T) { testFullRoundtrip_ApplyPatch(t, ctx) })
	t.Run("CompareHash", func(t *testing.T) { testFullRoundtrip_CompareHash(t, ctx) })
}

func testFullRoundtrip_Seed(t *testing.T, ctx *testFullRoundtrip_Context) {
	// open the files
	fromFile, err := os.Create(ctx.FromFilePath)
	if err != nil {
		t.Fatalf("Failed to create FROM file: %v", err)
	}
	defer fromFile.Close()

	toFile, err := os.Create(ctx.ToFilePath)
	if err != nil {
		t.Fatalf("Failed to create TO file: %v", err)
	}
	defer toFile.Close()

	// determine file sizes
	rand.Seed(time.Now().UnixNano())

	buf := make([]byte, 64*1024)

	fromBlocks := int(1024 + rand.Int31n(1024))
	toBlocks := int(1024 + rand.Int31n(1024))

	t.Logf("FROM file size: %v (%v)", fromBlocks*len(buf), humanize.Bytes(uint64(fromBlocks*len(buf))))
	t.Logf("TO file size: %v (%v)", toBlocks*len(buf), humanize.Bytes(uint64(toBlocks*len(buf))))

	fromSkipMod := int(3 + rand.Int31n(10))
	toSkipMod := int(3 + rand.Int31n(10))

	// start seeding
	maxBlocks := fromBlocks
	if toBlocks > maxBlocks {
		maxBlocks = toBlocks
	}

	toHash := sha1.New()

	for block := 0; block < maxBlocks; block++ {
		_, err := rand.Read(buf)
		if err != nil {
			t.Fatalf("Failed to seed random data: %v", err)
		}

		if (block%fromSkipMod != 0) && (block < fromBlocks) {
			fromFile.Write(buf)
		}

		if (block%toSkipMod != 0) && (block < toBlocks) {
			toFile.Write(buf)
			toHash.Write(buf)
		}
	}

	// seed header
	ctx.Header = make([]byte, 7000+rand.Int31n(20000))

	rand.Read(ctx.Header)

	// done
	ctx.ToFileHash = toHash.Sum(nil)

	t.Logf("TO file hash: %x", ctx.ToFileHash)
}

func testFullRoundtrip_CreatePatch(t *testing.T, ctx *testFullRoundtrip_Context) {
	// open the files
	fromFile, err := os.Open(ctx.FromFilePath)
	if err != nil {
		t.Fatalf("Failed to open FROM file: %v", err)
	}
	defer fromFile.Close()

	toFile, err := os.Open(ctx.ToFilePath)
	if err != nil {
		t.Fatalf("Failed to open TO file: %v", err)
	}
	defer toFile.Close()

	patchFile, err := os.Create(ctx.PatchFilePath)
	if err != nil {
		t.Fatalf("Failed to open PATCH file: %v", err)
	}
	defer patchFile.Close()

	// prepare encoder
	options := EncoderOptions{
		FileID:    "TestFullRoundtrip",
		FromFile:  fromFile,
		ToFile:    toFile,
		PatchFile: patchFile,
		Header:    ctx.Header,
	}

	enc, err := NewEncoder(options)
	if err != nil {
		t.Fatalf("Failed to create encoder: %v", err)
	}
	defer enc.Close()

	// create the patch
	err = enc.Process(context.TODO())
	if err != nil {
		t.Fatalf("Failed to create patch: %v", err)
	}
}

func testFullRoundtrip_DumpPatchInfo(t *testing.T, ctx *testFullRoundtrip_Context) {
	patchFileStat, err := os.Stat(ctx.PatchFilePath)
	if err != nil {
		t.Fatalf("Failed to get patch filesize: %v", err)
	}

	t.Logf("PATCH file size: %v (%v)", patchFileStat.Size(), humanize.Bytes(uint64(patchFileStat.Size())))
}

func testFullRoundtrip_ApplyPatch(t *testing.T, ctx *testFullRoundtrip_Context) {
	// open the files
	fromFile, err := os.Open(ctx.FromFilePath)
	if err != nil {
		t.Fatalf("Failed to open FROM file: %v", err)
	}
	defer fromFile.Close()

	appliedFile, err := os.Create(ctx.AppliedFilePath)
	if err != nil {
		t.Fatalf("Failed to open APPLIED file: %v", err)
	}
	defer appliedFile.Close()

	patchFile, err := os.Open(ctx.PatchFilePath)
	if err != nil {
		t.Fatalf("Failed to open PATCH file: %v", err)
	}
	defer patchFile.Close()

	// prepare decoder
	options := DecoderOptions{
		FileID:    "TestFullRoundtrip",
		FromFile:  fromFile,
		ToFile:    appliedFile,
		PatchFile: patchFile,
	}

	dec, err := NewDecoder(options)
	if err != nil {
		t.Fatalf("Failed to apply encoder: %v", err)
	}
	defer dec.Close()

	// retrieve header
	headerChannel := make(chan []byte, 1)
	dec.Header = headerChannel

	// apply the patch
	err = dec.Process(context.TODO())
	if err != nil {
		t.Fatalf("Failed to apply patch: %v", err)
	}

	// compare the header
	readHeader := <-headerChannel

	if !bytes.Equal(ctx.Header, readHeader) {
		t.Fatalf("Header of PATCH file does not match")
	}
}

func testFullRoundtrip_CompareHash(t *testing.T, ctx *testFullRoundtrip_Context) {
	// open the files
	appliedFile, err := os.Open(ctx.AppliedFilePath)
	if err != nil {
		t.Fatalf("Failed to open APPLIED file: %v", err)
	}
	defer appliedFile.Close()

	// calculate hash
	appliedHash := sha1.New()

	_, err = io.Copy(appliedHash, appliedFile)
	if err != nil {
		t.Fatalf("Failed to hash APPLIED file: %v", err)
	}

	appliedFile.Close()

	appliedHashResult := appliedHash.Sum(nil)

	// compare
	t.Logf("APPLIED file hash: %x", appliedHashResult)

	if !bytes.Equal(ctx.ToFileHash, appliedHashResult) {
		t.Fatalf("File hash of TO and APPLIED file are different")
	}
}
