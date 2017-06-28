package gtrace

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	testLogN_M(t, 1, 10)
	testLogN_M(t, 3, 1)
	testLogN_M(t, 6, 1)
	testLogN_M(t, 8, 1)
	testLogN_M(t, 20, 1)

	testLogN_M(t, 3, 1)
	testLogN_M(t, 3, 2)
	testLogN_M(t, 6, 3)
	testLogN_M(t, 8, 4)
	testLogN_M(t, 20, 5)
}

var rollTestPath = "/tmp/rollogTest"

func testLogN_M(t *testing.T, chunkBits, keepn uint) {

	os.RemoveAll(rollTestPath)

	//----------------------------------------
	// 打开新文件

	f, err := OpenRolloger(rollTestPath, chunkBits, keepn) // 8 bytes, rollsize 32 bytes
	assert.Equal(t, nil, err, "open fail")

	text := []byte("123456789")

	totalLen := int64(0)
	for i := 0; i < 1000; i++ {
		err := f.Log(text)
		assert.Equal(t, nil, err, "write error")
		totalLen += int64(len(text) + 1)
	}
	assert.Equal(t, nil, f.Close(), "close fail")

	//----------------------------------------
	// 关闭后重新打开

	ff, err := os.OpenFile(rollTestPath+"/___", os.O_RDWR|os.O_CREATE, 0666)
	assert.Equal(t, nil, err, "create file fail")
	ff.Close()

	f, err = OpenRolloger(rollTestPath, chunkBits, keepn)
	assert.Equal(t, nil, err, "open fail")

	for i := 0; i < 1000; i++ {
		err := f.Log(text)
		assert.Equal(t, nil, err, "write error")
		totalLen += int64(len(text) + 1)
	}
	assert.Equal(t, nil, f.Close(), "close fail")

	//----------------------------------------
	// 检查当前文件内容是否正确

	l := len(text) + 1 // append '\n'
	fsize := ((1 << chunkBits / l) + 1) * l
	count := 2000 * l / fsize

	type x struct {
		name string
		size int64
	}
	want := []x{}
	for i := count - 1; i >= 0 && i > count-int(keepn); i-- {
		want = append(want, x{
			name: strconv.FormatInt(int64(i), 36),
			size: int64(fsize),
		})
	}
	want = append([]x{
		x{
			name: strconv.FormatInt(int64(count), 36),
			size: int64((2000 * l) % fsize),
		}}, want...)

	list, err := listSortFile(rollTestPath)
	assert.Equal(t, nil, err, "list file error")

	got := []x{}
	for _, f := range list {
		got = append(got, x{
			name: f.Name(),
			size: f.Size(),
		})
	}
	assert.Equal(t, want, got, "file remains wrong")
}

// --------------------------------------------------------------------
