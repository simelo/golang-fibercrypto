import QtQuick 2.12
import QtQuick.Controls 2.12
import QtQuick.Controls.Material 2.12
import QtQuick.Layouts 1.12

Item {
    id: root

    property date date: "2000-01-01 00:00"
    property int type: TransactionDetails.Type.Send
    property int status: TransactionDetails.Status.Preview
    property var statusString: [ qsTr("Confirmed"), qsTr("Pending"), qsTr("Preview") ]
    property real amount: 0
    property int hoursReceived: 0
    property int hoursBurned: 0
    property string transactionID

    enum Status {
        Confirmed,
        Pending,
        Preview
    }
    enum Type {
        Send,
        Receive
    }

    implicitHeight: 400
    implicitWidth: 650
    clip: true

    ColumnLayout {
        id: columnLayoutRoot
        anchors.fill: parent
        spacing: 20

        RowLayout {
            Layout.fillWidth: true

            ColumnLayout {
                Layout.fillWidth: true

                Label {
                    text: qsTr("Transaction")
                    font.bold: true
                    Layout.fillWidth: true
                }

                GridLayout {
                    id: gridLayoutBasicInfo
                    Material.foreground: Material.Grey
                    columns: 2
                    columnSpacing: 10

                    Layout.fillWidth: true

                    Label {
                        text: qsTr("Date:")
                        font.bold: true
                    }
                    Label {
                        text: Qt.formatDateTime(root.date, Qt.DefaultLocaleShortDate)
                    }

                    Label {
                        text: qsTr("Status:")
                        font.bold: true
                    }
                    Label {
                        text: statusString[root.status]
                    }

                    Label {
                        text: qsTr("Hours:")
                        font.bold: true
                    }
                    Label {
                        text: root.hoursReceived + ' ' + qsTr("received") + ' | ' + hoursBurned + ' ' + qsTr("burned")
                    }

                    Label {
                        text: qsTr("Tx ID:")
                        font.bold: true
                    }
                    Label {
                        text: root.transactionID
                        Layout.fillWidth: true
                    }
                } // GridLayout
            }

            ColumnLayout {
                Layout.alignment: Qt.AlignTop
                Layout.rightMargin: 20
                Image {
                    source: "qrc:/images/send-" + (type === TransactionDetails.Type.Receive ? "blue" : "amber") + ".svg"
                    sourceSize: "96x96"
                    fillMode: Image.PreserveAspectFit
                    mirror: type === TransactionDetails.Type.Receive
                    Layout.fillWidth: true
                }
                Label {
                    text: (type === TransactionDetails.Type.Receive ? "Receive" : "Send") + ' ' + amount + ' ' + qsTr("SKY")
                    font.bold: true
                    font.pointSize: 14
                    horizontalAlignment: Label.AlignHCenter
                    Layout.fillWidth: true
                }
            }
        } // RowLayout

        Rectangle {
            height: 1
            color: Material.color(Material.Grey)
            Layout.fillWidth: true
        }

        RowLayout {
            Layout.fillWidth: true

            ColumnLayout {
                Layout.fillWidth: true

                Label {
                    text: qsTr("Inputs")
                    font.pointSize: Qt.application.font.pointSize + 2
                    font.bold: true
                    font.italic: true
                    Layout.fillWidth: true
                }

                ListView {
                    id: listViewInputs
                    model: listModelOutputsInputs
                    implicitHeight: 120 * count
                    interactive: false
                    Layout.fillWidth: true
                    delegate: InputOutputDelegate {
                        width: ListView.view.width
                    }
                }
            }

            ColumnLayout {
                Layout.fillWidth: true

                Label {
                    text: qsTr("Outputs")
                    font.pointSize: Qt.application.font.pointSize + 3
                    font.bold: true
                    font.italic: true
                    Layout.fillWidth: true
                }

                ListView {
                    id: listViewOutputs
                    model: listModelOutputsInputs
                    implicitHeight: 120 * count
                    interactive: false
                    Layout.fillWidth: true
                    delegate: InputOutputDelegate {
                        width: ListView.view.width
                    }
                }
            }
        } // RowLayout
    } // ColumnLayout (root)

    // Roles: address, addressSky, addressCoinHours
    // Use listModel.append( { "address": value, "addressSky": value, "addressCoinHours": value } )
    // Or implement the model in the backend (a more recommendable approach)
    ListModel {
        id: listModelOutputsInputs
        ListElement { address: "qrxw7364w8xerusftaxkw87ues"; addressSky: 30; addressCoinHours: 1049 }
        ListElement { address: "8745yuetsrk8tcsku4ryj48ije"; addressSky: 12; addressCoinHours: 16011 }
    }
}