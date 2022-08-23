package ws

import (
	"fmt"
	"strings"

	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const NODE_TYPE_RESP = "resp"
const NODE_TYPE_P2P = "p2p"
const NODE_TYPE_GROUP = "group"
const NODE_TYPE_BROAD = "broad"

type NodeWrite struct {
	Type string //resp , p2p,  group, broad
	Uid  string
	Gid  string
	Msg  interface{}
}

type NodeRead struct {
	Uid string
	Gid string
	Msg interface{}
}

type WsIo interface {
	Init(adapater string, hostip string) bool // chan  redis   kafka
	ReadFromws()
	Write2ws()
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WsIoBase  for  chan
type WsIoBase struct {
	WriteDataCh chan NodeWrite
	ReadDataCh  chan NodeRead
}

func (this *WsIoBase) Init(adapater string, hostip string) bool {
	if adapater != "chan" {
		return false
	}
	this.WriteDataCh = make(chan NodeWrite, 256)
	this.ReadDataCh = make(chan NodeRead, 1024)
	return true
}

func (this *WsIoBase) ReadFromws(doMsg func(msg NodeRead)) {
	for {
		node := <-this.ReadDataCh
		if doMsg != nil {
			doMsg(node)
		}
	}
}

func (this *WsIoBase) NewMsg(id, group string, msg []byte) {

	var node NodeRead
	node.Uid = id
	node.Gid = group
	node.Msg = msg
	this.ReadDataCh <- node
}

func (this *WsIoBase) Write2ws(node NodeWrite) {
	if node.Type == NODE_TYPE_RESP || node.Type == NODE_TYPE_P2P {
		bs, err := json.Marshal(node.Msg)
		if err == nil {
			WebsocketManager.Send(node.Uid, node.Gid, bs)
		} else {

		}
	}
	if node.Type == NODE_TYPE_GROUP {
		bs, err := json.Marshal(node.Msg)
		if err == nil {
			WebsocketManager.SendGroup(node.Gid, bs)
		} else {

		}
	}
	if node.Type == NODE_TYPE_BROAD {
		bs, err := json.Marshal(node.Msg)
		if err == nil {
			WebsocketManager.SendAll(bs)
		} else {

		}
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WsIoBase  for  redis    未完成
type WsIoRedis struct {
	__redisClient *redis.Client
	ReadDataCh    chan NodeRead
}

func (this *WsIoRedis) Init(adapater string, hostip string) bool {
	if adapater != "redis" {
		return false
	}
	ip := ""
	redispwd := ""
	fields := strings.Split(hostip, "@")
	if len(fields) == 1 {
		ip = hostip
	}
	if len(fields) == 2 {
		ip = fields[1]
		redispwd = fields[0]
	}
	// redispwd := strings.Split(redisip, "@")[0]
	// ip := strings.Split(redisip, "@")[1]
	this.__redisClient = redis.NewClient(&redis.Options{
		Addr:     ip,
		Password: redispwd, // no password set   "!test123"
		DB:       0,        // use default DB
	})

	_, err := this.__redisClient.Ping().Result()
	if err != nil {
		fmt.Println("err", err)
		return false
	}
	return true
}

// XADD
// XADD mystream * field1 value1 field2 value2 field3 value3
func (this *WsIoRedis) ReadFromws(doMsg func(msg NodeRead)) {

	for {
		node := <-this.ReadDataCh
		doMsg(node)
	}

}

// XGROUP CREATE - 创建消费者组
// XREADGROUP GROUP - 读取消费者组中的消息
// XACK - 将消息标记为"已处理"
func (this *WsIoRedis) Write2ws(node NodeWrite) {
	if node.Type == NODE_TYPE_RESP || node.Type == NODE_TYPE_P2P {
		bs, err := json.Marshal(node.Msg)
		if err == nil {
			WebsocketManager.Send(node.Uid, node.Gid, bs)
		} else {

		}

	}
	if node.Type == NODE_TYPE_GROUP {
		bs, err := json.Marshal(node.Msg)
		if err == nil {
			WebsocketManager.SendGroup(node.Gid, bs)
		} else {

		}
	}
	if node.Type == NODE_TYPE_BROAD {
		bs, err := json.Marshal(node.Msg)
		if err == nil {
			WebsocketManager.SendAll(bs)
		} else {

		}
	}
}

// func publishTicketReceivedEvent(client *redis.Client) error {
// 	log.Println("Publishing event to Redis")
// 	err := client.XAdd(&redis.XAddArgs{
// 	  Stream:       "tickets",
// 	  MaxLen:       0,
// 	  MaxLenApprox: 0,
// 	  ID:           "",
// 	  Values: map[string]interface{}{
// 		"whatHappened": string("ticket received"),
// 		"ticketID":     int(rand.Intn(100000000)),
// 		"ticketData":   string("some ticket data"),
// 	  },
// 	}).Err()
// 	return err
//   }

//消费

// subject := "tickets"
// consumersGroup := "tickets-consumer-group"
// err = redisClient.XGroupCreate(subject, consumersGroup, "0").Err()
// if err != nil {
// log.Println(err)
// }

// 该程序没有退出，仍然在侦听消息。这是因为我们使用BLOCK = 0参数调用XREADGROUP，这意味着程序执行将被阻塞无限长的时间，直到新消息到来。
// uniqueID := xid.New().String()
// for {
//   entries, err := redisClient.XReadGroup(&redis.XReadGroupArgs{
//     Group:    consumersGroup,
//     Consumer: uniqueID,
//     Streams:  []string{subject, ">"},
//     Count:    2,
//     Block:    0,
//     NoAck:    false,
//   }).Result()
//   if err != nil {
//     log.Fatal(err)
//   }
// for i := 0; i < len(entries[0].Messages); i++ {
//   messageID := entries[0].Messages[i].ID
//   values := entries[0].Messages[i].Values
//   eventDescription := fmt.Sprintf("%v", values["whatHappened"])
//   ticketID := fmt.Sprintf("%v", values["ticketID"])
//   ticketData := fmt.Sprintf("%v", values["ticketData"])
//   if eventDescription == "ticket received" {
//     err := handleNewTicket(ticketID, ticketData)
//     if err != nil {
//       log.Fatal(err)
//     }
//     redisClient.XAck(subject, consumersGroup, messageID)
//   }
// }
