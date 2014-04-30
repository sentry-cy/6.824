package pbservice

import "viewservice"
import "net/rpc"
//import "fmt"
import (
	"time"
	"strconv"
)

// You'll probably need to uncomment these:
// import "time"
// import "crypto/rand"
// import "math/big"



type Clerk struct {
	vs *viewservice.Clerk
	// Your declarations here

	view viewservice.View
	me string
}


func MakeClerk(vshost string, me string) *Clerk {
	ck := new(Clerk)
	ck.vs = viewservice.MakeClerk(me, vshost)
	// Your ck.* initializations here
	ck.view = viewservice.View{}

	// note: argument me could not garantuee unique, this ck.me field would be put along
	// with PutArgs to show client's identification
	ck.me = strconv.FormatInt(nrand(), 10)
	return ck
}


//
// call() sends an RPC to the rpcname handler on server srv
// with arguments args, waits for the reply, and leaves the
// reply in reply. the reply argument should be a pointer
// to a reply structure.
//
// the return value is true if the server responded, and false
// if call() was not able to contact the server. in particular,
// the reply's contents are only valid if call() returned true.
//
// you should assume that call() will time out and return an
// error after a while if it doesn't get a reply from the server.
//
// please use call() to send all RPCs, in client.go and server.go.
// please don't change this function.
//
func call(srv string, rpcname string,
	args interface{}, reply interface{}) bool {
	c, errx := rpc.Dial("unix", srv)
	if errx != nil {
		return false
	}
	defer c.Close()

	err := c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

//	fmt.Println(err)
	return false
}


func (ck *Clerk) UpdateView() {
	view, err := ck.vs.Ping(ck.view.Viewnum)
	if err != nil {
//		fmt.Println("[client.UpdateView] fails: ", view, err)
		return
	}
	ck.view = view
}

//
// fetch a key's value from the current primary;
// if they key has never been set, return "".
// Get() must keep trying until it either the
// primary replies with the value or the primary
// says the key doesn't exist (has never been Put().
//
func (ck *Clerk) Get(key string) string {

	// Your code here.
	if ck.view.Viewnum == 0{
		ck.UpdateView()
	}
	args := &GetArgs{key, nrand()}
	var reply GetReply

	for {
		ok := call(ck.view.Primary, "PBServer.Get", args, &reply)
		if ok {
			return reply.Value
		}
		time.Sleep(viewservice.PingInterval)
		ck.UpdateView()
	}
	return "???"
}

//
// tell the primary to update key's value.
// must keep trying until it succeeds.
//
func (ck *Clerk) PutExt(key string, value string, dohash bool) string {

	// Your code here.
	if ck.view.Viewnum == 0 {
		ck.UpdateView()
	}
//	fmt.Println("[Clerk.Put]: key/value", key, "/", value)
//	num := nrand()
//	fmt.Println("uid: ", num, strconv.FormatInt(nrand(), 10))
	args := &PutArgs{key, value, dohash, false, strconv.FormatInt(nrand(), 10), ck.me}
//	fmt.Println("args ", args)
	var reply PutReply

	for {
		ok := call(ck.view.Primary, "PBServer.Put", args, &reply)
		if ok {
//			fmt.Println("client.PutExt returns")
			return reply.PreviousValue
		}
		//just always update view, because call may return false, reply empty when server got killed
		time.Sleep(viewservice.PingInterval)
		ck.UpdateView()
	}
	return "???"
}

func (ck *Clerk) Put(key string, value string) {
	ck.PutExt(key, value, false)
}
func (ck *Clerk) PutHash(key string, value string) string {
	v := ck.PutExt(key, value, true)
	return v
}
