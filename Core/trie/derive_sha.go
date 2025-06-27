package trie

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/rlp"

	"FichainCore/common"
)

type Marshalable interface {
	Marshal() ([]byte, error)
}

func DeriveSha(list any) (common.Hash, error) {
	keybuf := new(bytes.Buffer)
	trie := new(Trie)

	val := reflect.ValueOf(list)
	if val.Kind() != reflect.Slice {
		return common.Hash{}, fmt.Errorf("input is not a slice")
	}

	for i := 0; i < val.Len(); i++ {
		keybuf.Reset()
		rlp.Encode(keybuf, uint(i))

		item := val.Index(i).Interface()

		m, ok := item.(Marshalable)
		if !ok {
			return common.Hash{}, fmt.Errorf("item at index %d does not implement Marshalable", i)
		}

		b, err := m.Marshal()
		if err != nil {
			return common.Hash{}, err
		}

		trie.Update(keybuf.Bytes(), b)
	}

	return trie.Hash(), nil
}
