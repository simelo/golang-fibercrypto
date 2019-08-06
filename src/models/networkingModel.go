package models

import (
	coin "github.com/fibercrypto/FiberCryptoWallet/src/core"
	"github.com/therecipe/qt/core"
	)

// ip, port, source, block, lastSeenIn, lastSeenOut
const (
	Ip	              = int(core.Qt__UserRole) + 1
	Port			  = int(core.Qt__UserRole) + 2
	Source            = int(core.Qt__UserRole) + 3
	Block	          = int(core.Qt__UserRole) + 4
	LastSeenIn        = int(core.Qt__UserRole) + 5
	LastSeenOut       = int(core.Qt__UserRole) + 6
)

type NetworkingModel struct {
	core.QAbstractListModel

	_ func() `constructor:"init"`

	_ map[int]*core.QByteArray `property:"roles"`
	_ []*QNetworking           `property:"networks"`

	_ func(*QNetworking)                                                        `slot:"addNetwork"`
	_ func(row int)                                                             `slot:"removeNetwork"`
	_ int                                                                       `property:"count"`
}

type QNetworking struct {
	core.QObject
// ip, port, source, block, lastSeenIn, lastSeenOut
	_ string `property:"ip"`
	_ uint16 `property:"port"`
	_ bool `property:"source"`
	_ uint64 `property:"block"`
	_ int64 `property:"lastSeenIn"`
	_ int64 `property:"lastSeenOut"`
}

func (netModel *NetworkingModel) init() {
	netModel.SetRoles(map[int]*core.QByteArray{
		Ip:             core.NewQByteArray2("ip", -1),
		Port: 			core.NewQByteArray2("port", -1),
		Source:         core.NewQByteArray2("source", -1),
		Block:          core.NewQByteArray2("block", -1),
		LastSeenIn:     core.NewQByteArray2("lastSeenIn", -1),
		LastSeenOut:    core.NewQByteArray2("lastSeenOut", -1),
	})

	netModel.ConnectData(netModel.data)
	netModel.ConnectRowCount(netModel.rowCount)
	netModel.ConnectColumnCount(netModel.columnCount)
	netModel.ConnectRoleNames(netModel.roleNames)

	netModel.ConnectAddNetwork(netModel.addNetwork)
	netModel.ConnectRemoveNetwork(netModel.removeNetwork)

}

func (netModel *NetworkingModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() {
		return core.NewQVariant()
	}

	if index.Row() >= len(netModel.Networks()) {
		return core.NewQVariant()
	}

	var w = netModel.Networks()[index.Row()]

	switch role {
	case Ip:
		{
			return core.NewQVariant1(w.Ip())
		}

	case Port:
		{
			return core.NewQVariant1(w.Port())
		}

	case Source:
		{
			return core.NewQVariant1(w.IsSource())
		}

	case Block:
		{
			return core.NewQVariant1(w.Block())
		}
	case LastSeenIn:
		{
			return core.NewQVariant1(w.LastSeenIn())
		}
	case LastSeenOut:
		{
			return core.NewQVariant1(w.LastSeenOut())
		}

	default:
		{
			return core.NewQVariant()
		}
	}
}

func (netModel *NetworkingModel) rowCount(parent *core.QModelIndex) int {
	return len(netModel.Networks())
}

func (netModel *NetworkingModel) columnCount(parent *core.QModelIndex) int {
	return 1
}

func (netModel *NetworkingModel) roleNames() map[int]*core.QByteArray {
	return netModel.Roles()
}

func (netModel *NetworkingModel) addNetwork(w *QNetworking) {
	netModel.BeginInsertRows(core.NewQModelIndex(), len(netModel.Networks()), len(netModel.Networks()))
	netModel.SetNetworks(append(netModel.Networks(), w))
	netModel.EndInsertRows()
	netModel.updateCount()

}

func (netModel *NetworkingModel) addMultipleNetworks(w []*QNetworking]) {
	for _, qnet := range w{
		netModel.addNetwork(qnet)
	}
}

func INetworkToQNetworking(net coin.INetwork) *QNetworking {
	q := NewQNetworking(nil)
	q.SetIp(net.GetIp())
	q.SetPort(net.GetPort())
	q.SetSource(net.IsTrusted())
	q.SetBlock(net.GetBlock())
	q.SetLastSeenIn(net.GetLastSeenIn())
	q.SetLastSeenOut(net.GetLastSeenOut())

	return q
}

func (netModel *NetworkingModel) removeNetwork(row int) {
	netModel.BeginRemoveRows(core.NewQModelIndex(), row, row)
	netModel.SetNetworks(append(netModel.Networks()[:row], netModel.Networks()[row+1:]...))
	netModel.EndRemoveRows()
	netModel.updateCount()

}

func (netModel *NetworkingModel) updateCount() {
	netModel.SetCount(len(netModel.Networks()))
}
