package server

import (
	"sync"

	"github.com/FactomProject/FactomCode/common"
	"github.com/FactomProject/FactomCode/wire"
)

// ProcessListMgr contains a list of valid confirmation messages
// and is used for consensus building
type ProcessListMgr struct {
	sync.RWMutex
	MyProcessList    *ProcessList
	NextDBlockHeight uint32
	serverPrivKey    common.PrivateKey
}

// NewProcessListMgr create sa new process list
func NewProcessListMgr(height uint32, plSizeHint uint, privKey common.PrivateKey) *ProcessListMgr {
	plMgr := new(ProcessListMgr)
	plMgr.MyProcessList = NewProcessList(plSizeHint)
	plMgr.NextDBlockHeight = height
	plMgr.serverPrivKey = privKey
	return plMgr
}

// AddToFollowersProcessList creates a new process list item and add it to the MyProcessList
func (plMgr *ProcessListMgr) AddToFollowersProcessList(msg wire.Message, ack *wire.MsgAck, hash *wire.ShaHash) error {
	err := ack.Sign(&plMgr.serverPrivKey)
	if err != nil {
		return err
	}
	plMgr.MyProcessList.nextIndex++

	plItem := &ProcessListItem{
		Ack:     ack,
		Msg:     msg,
		MsgHash: hash,
	}
	plMgr.MyProcessList.AddToProcessList(plItem)
	return nil
}

// AddToLeadersProcessList creates a new process list item and add it to the MyProcessList
func (plMgr *ProcessListMgr) AddToLeadersProcessList(msg wire.FtmInternalMsg, hash *wire.ShaHash,
	msgType byte, dirBlockTimestamp uint32, coinbaseTimestamp int64, sid string) (ack *wire.MsgAck, err error) {
	
	ack = wire.NewMsgAck(plMgr.NextDBlockHeight, uint32(plMgr.MyProcessList.nextIndex), hash, msgType,
		dirBlockTimestamp, uint64(coinbaseTimestamp), sid)
	// Sign the ack using server private keys
	err = ack.Sign(&plMgr.serverPrivKey)
	if err != nil {
		return nil, err
	}
	plMgr.MyProcessList.nextIndex++

	plItem := &ProcessListItem{
		Ack:     ack,
		Msg:     msg,
		MsgHash: hash,
	}
	plMgr.MyProcessList.AddToProcessList(plItem)
	return ack, nil
}

// IsMyPListExceedingLimit checks if the number of process list items is exceeding the size limit
func (plMgr *ProcessListMgr) IsMyPListExceedingLimit() bool {
	return (plMgr.MyProcessList.totalItems >= common.MAX_PLIST_SIZE)
}