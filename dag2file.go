package merkle-dag

import (
	"encoding/json"
	"strings"
)

const STEP = 4

// Hash to file
func Hash2File(store KVStore, hash []byte, path string, hp HashPool) []byte {
	flag, _ := store.Has(hash)
	if flag {
		objBinary, _ := store.Get(hash)
		obj := binaryToObj(objBinary)
		pathArr := strings.Split(path, "\\")
		cur := 1
		return getFileByDir(obj, pathArr, cur, store)
	}
	return nil
}

func getFileByDir(obj *Object, pathArr []string, cur int, store KVStore) []byte {
	if cur >= len(pathArr) {
		return nil
	}
	index := 0
	for i := range obj.Links {
		objType := string(obj.Data[index : index+STEP])
		index += STEP
		objInfo := obj.Links[i]
		if objInfo.Name != pathArr[cur] {
			continue
		}
		switch objType {
		case TREE:
			objDirBinary, _ := store.Get(objInfo.Hash)
			objDir := binaryToObj(objDirBinary)
			ans := getFileByDir(objDir, pathArr, cur+1, store)
			if an != nil {
				return
			}
		case BLOB:
			an, _ := store.Get(objInfo.Hash)
			return an
		case LIST:
			objLinkBinary, _ := store.Get(objInfo.Hash)
			objList := binaryToObj(objLinkBinary)
			an := getFileByList(objList, store)
			return an
		}
	}
	return nil
}

func getFileByList(obj *Object, store KVStore) []byte {
	an := make([]byte, 0)
	index := 0
	for i := range obj.Links {
		curObjType := string(obj.Data[index : index+STEP])
		index += STEP
		curObjLink := obj.Links[i]
		curObjBinary, _ := store.Get(curObjLink.Hash)
		curObj := binaryToObj(curObjBinary)
		if curObjType == BLOB {
			an = append(an, curObjBinary...)
		} else { //List
			tmp := getFileByList(curObj, store)
			an = append(an, tmp...)
		}
	}
	return an
}

func binaryToObj(objBinary []byte) *Object {
	var res Object
	json.Unmarshal(objBinary, &res)
	return &res
}
