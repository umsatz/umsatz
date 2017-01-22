package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

func migrations_20141019153443_create_bank_accounts_sql() ([]byte, error) {
	return bindata_read([]byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x00, 0xff, 0x8c, 0x91,
		0x41, 0x4f, 0xc2, 0x40, 0x10, 0x85, 0xef, 0xf3, 0x2b, 0xe6, 0x48, 0x23,
		0x18, 0x83, 0xc6, 0x0b, 0xa7, 0xb5, 0x6c, 0xb4, 0xb1, 0x16, 0xb2, 0x16,
		0x13, 0x4e, 0x64, 0x76, 0xa9, 0x64, 0xa2, 0x9d, 0xc5, 0x76, 0x37, 0xfe,
		0x7d, 0x17, 0xe2, 0x41, 0x53, 0x42, 0xd8, 0xf3, 0xf7, 0xde, 0xb7, 0x93,
		0x37, 0x99, 0xe0, 0x55, 0xcb, 0xbb, 0x8e, 0x42, 0x83, 0xab, 0x3d, 0xe4,
		0x46, 0xab, 0x5a, 0x63, 0xad, 0x1e, 0x4a, 0x8d, 0xfb, 0x68, 0x3f, 0xd9,
		0x5d, 0x5b, 0x92, 0x8f, 0x0d, 0x39, 0xe7, 0xa3, 0x84, 0x1e, 0x47, 0x80,
		0xc8, 0x5b, 0x3c, 0xf1, 0x5e, 0xb5, 0x29, 0x54, 0x89, 0x4b, 0x53, 0xbc,
		0x28, 0xb3, 0xc6, 0x67, 0xbd, 0x1e, 0x43, 0x82, 0x93, 0x81, 0xbe, 0x0e,
		0x1d, 0x2c, 0x3b, 0xa4, 0x10, 0x3a, 0xb6, 0x31, 0xc9, 0xfc, 0x3b, 0x86,
		0x8e, 0xa4, 0x27, 0x17, 0xd8, 0x4b, 0x9f, 0xc0, 0xa3, 0xc7, 0xf9, 0x6d,
		0xf3, 0xaf, 0x35, 0x7f, 0x52, 0x46, 0xe5, 0xb5, 0x36, 0xf8, 0x96, 0x5a,
		0x8b, 0xea, 0x71, 0x74, 0x3b, 0xcd, 0xb0, 0x5a, 0xd4, 0x58, 0xad, 0xca,
		0x72, 0x9c, 0x72, 0xbf, 0x5f, 0xdb, 0x48, 0x6c, 0x6d, 0xd3, 0x5d, 0x9e,
		0xe3, 0x24, 0x1c, 0x5e, 0x31, 0xcc, 0xdd, 0xdf, 0x65, 0x07, 0xdc, 0xb2,
		0x3b, 0x71, 0xf4, 0x10, 0x9f, 0xde, 0x1c, 0x71, 0xa1, 0xb6, 0xb9, 0xb0,
		0x1d, 0xb2, 0x19, 0xc0, 0xdf, 0x21, 0xe6, 0xfe, 0x5b, 0x60, 0x6e, 0x16,
		0xcb, 0x33, 0x43, 0xcc, 0xe0, 0x27, 0x00, 0x00, 0xff, 0xff, 0x2c, 0x44,
		0xa9, 0xfb, 0xbb, 0x01, 0x00, 0x00,
	},
		"migrations/20141019153443-create-bank_accounts.sql",
	)
}

func migrations_20141019153624_create_transactions_sql() ([]byte, error) {
	return bindata_read([]byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x00, 0xff, 0x94, 0x93,
		0xcd, 0xae, 0xd3, 0x30, 0x10, 0x85, 0xf7, 0x7e, 0x8a, 0x59, 0x36, 0xa2,
		0x65, 0x01, 0x88, 0x4d, 0x57, 0x26, 0x75, 0xdb, 0x88, 0xfc, 0x54, 0x8e,
		0x8b, 0x28, 0x9b, 0xc8, 0x75, 0x5c, 0x14, 0x35, 0xb1, 0x23, 0xc7, 0x01,
		0xf1, 0xf6, 0xd8, 0x41, 0x48, 0x75, 0x48, 0xb9, 0xba, 0x5e, 0x7a, 0xbe,
		0x39, 0xe7, 0x78, 0x34, 0xde, 0x6c, 0xe0, 0x4d, 0xd7, 0x7c, 0x37, 0xdc,
		0x4a, 0x38, 0xf7, 0x28, 0xa6, 0x04, 0x33, 0x02, 0x0c, 0x7f, 0x4a, 0x09,
		0xf4, 0xe3, 0xb5, 0x6d, 0xc4, 0x5b, 0x6b, 0xb8, 0x1a, 0xb8, 0xb0, 0x8d,
		0x56, 0x03, 0xac, 0x10, 0x40, 0x53, 0xc3, 0xc2, 0x29, 0x09, 0x4d, 0x70,
		0x0a, 0x27, 0x9a, 0x64, 0x98, 0x5e, 0xe0, 0x33, 0xb9, 0xac, 0x1d, 0xdb,
		0x8f, 0xa6, 0xd7, 0x83, 0x9c, 0xb1, 0x8c, 0x7c, 0x65, 0xbe, 0x5a, 0xcb,
		0x41, 0x98, 0xa6, 0xf7, 0xd2, 0x0b, 0xd5, 0xc1, 0x72, 0x3b, 0x0e, 0x73,
		0x9f, 0xf8, 0x88, 0x29, 0x8e, 0x19, 0xa1, 0xf0, 0xc5, 0xf9, 0x24, 0xf9,
		0x61, 0xf5, 0xfe, 0x5d, 0x04, 0x90, 0x17, 0x0c, 0xf2, 0x73, 0x9a, 0xfa,
		0x46, 0x61, 0xa4, 0x7b, 0x4f, 0x5d, 0x71, 0x1b, 0x98, 0x26, 0x19, 0x29,
		0x19, 0xce, 0x4e, 0xec, 0x5b, 0x40, 0xff, 0xe0, 0xed, 0x68, 0x79, 0x55,
		0xfb, 0x11, 0x2c, 0xd2, 0x1e, 0xea, 0xb8, 0xf2, 0x40, 0x65, 0xe4, 0x4d,
		0x1a, 0xa9, 0x84, 0x7c, 0x92, 0xe5, 0xe3, 0x87, 0x28, 0x8c, 0x32, 0x0e,
		0x56, 0x77, 0xd2, 0x84, 0x8d, 0x2f, 0xf7, 0x59, 0x6d, 0x79, 0x5b, 0xf1,
		0x4e, 0x8f, 0xca, 0x56, 0x42, 0x2a, 0x3b, 0xcd, 0x21, 0xc9, 0xd9, 0x7f,
		0xa8, 0xd1, 0x78, 0x83, 0x5f, 0x4b, 0x13, 0x0a, 0xc5, 0x6f, 0x52, 0xfe,
		0x23, 0x3d, 0x89, 0xcf, 0x8b, 0x7f, 0x15, 0x17, 0xa7, 0x1e, 0xad, 0x91,
		0xc3, 0x5b, 0x2d, 0x5c, 0x84, 0x2b, 0x57, 0xf7, 0xea, 0x61, 0x2f, 0xe6,
		0x41, 0x8d, 0xec, 0xb4, 0x9b, 0xde, 0x0c, 0x0b, 0x29, 0x87, 0xc5, 0x45,
		0x5e, 0x32, 0x8a, 0xfd, 0xfd, 0xa4, 0xeb, 0xf9, 0xdb, 0x1d, 0x79, 0x76,
		0x5f, 0x50, 0x92, 0x1c, 0x72, 0xbf, 0x57, 0xb0, 0x0a, 0x4c, 0xa3, 0xa9,
		0x4e, 0xc9, 0x9e, 0x50, 0x92, 0xc7, 0xa4, 0x84, 0xe9, 0x9e, 0x0b, 0xe1,
		0xdf, 0xe0, 0x36, 0xd6, 0x11, 0x90, 0x61, 0x16, 0x1f, 0x61, 0xbf, 0xe0,
		0xf3, 0x27, 0xd9, 0x33, 0xa3, 0x30, 0xf7, 0xeb, 0x9c, 0x50, 0xb4, 0x45,
		0xe8, 0xf1, 0x7f, 0xed, 0xf4, 0x4f, 0x85, 0x76, 0xb4, 0x38, 0x3d, 0xff,
		0x5f, 0x5b, 0xf4, 0x3b, 0x00, 0x00, 0xff, 0xff, 0xb3, 0xed, 0x8b, 0x9d,
		0x91, 0x03, 0x00, 0x00,
	},
		"migrations/20141019153624-create-transactions.sql",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"migrations/20141019153443-create-bank_accounts.sql": migrations_20141019153443_create_bank_accounts_sql,
	"migrations/20141019153624-create-transactions.sql": migrations_20141019153624_create_transactions_sql,
}
// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"migrations": &_bintree_t{nil, map[string]*_bintree_t{
		"20141019153443-create-bank_accounts.sql": &_bintree_t{migrations_20141019153443_create_bank_accounts_sql, map[string]*_bintree_t{
		}},
		"20141019153624-create-transactions.sql": &_bintree_t{migrations_20141019153624_create_transactions_sql, map[string]*_bintree_t{
		}},
	}},
}}
