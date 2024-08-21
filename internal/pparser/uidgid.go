package pparser

import (
	"log"
	"os/user"
	"strconv"
	"strings"
)

func parseUIDGIDvalue(in string) (uint32, uint32) {
	ug := strings.Split(in, ":")
	if len(ug) != 2 {
		log.Fatal("Err: UID/GID format is incorrect")
	}
	uid := getUID(ug[0])
	gid := getGID(ug[1])
	return uid, gid
}

func getUID(userName string) uint32 {
	u, err := user.Lookup(userName)
	if err != nil {
		log.Fatal(ErrUserNotFound, err)
	}

	uid, err := strconv.ParseUint(u.Uid, 10, 32)
	if err != nil {
		log.Fatal("ErrUIDParse: ", userName, " ", err)
	}
	return uint32(uid)
}

func getGID(groupName string) uint32 {
	g, err := user.LookupGroup(groupName)
	if err != nil {
		log.Fatal(ErrGroupNotFound, err)
	}
	gid, err := strconv.ParseUint(g.Gid, 10, 32)
	if err != nil {
		log.Fatal("ErrGIDParse: ", groupName, " ", err)
	}
	return uint32(gid)
}
