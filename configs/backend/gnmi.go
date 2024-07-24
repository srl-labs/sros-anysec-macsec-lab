package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnmic/pkg/api"
	"github.com/openconfig/gnmic/pkg/api/path"
	"github.com/openconfig/gnmic/pkg/api/target"
)

type GnmiConnect struct {
	Host           string
	Port           string
	Username       string
	Password       string
	Target         *target.Target
	Context        context.Context
	Cancel         context.CancelFunc
	PrettyResponse bool
}

type NotificationRspMsg struct {
	Timestamp int64     `json:"timestamp,omitempty"`
	Time      time.Time `json:"time,omitempty"`
	Prefix    string    `json:"prefix,omitempty"`
	Target    string    `json:"target,omitempty"`
	Updates   []update  `json:"updates,omitempty"`
}

type update struct {
	Path   string                 `json:"path,omitempty"`
	Values map[string]interface{} `json:"values,omitempty"`
}

type setRspMsg struct {
	Source    string            `json:"source,omitempty"`
	Timestamp int64             `json:"timestamp,omitempty"`
	Time      time.Time         `json:"time,omitempty"`
	Prefix    string            `json:"prefix,omitempty"`
	Target    string            `json:"target,omitempty"`
	Results   []updateResultMsg `json:"results,omitempty"`
}

type updateResultMsg struct {
	Operation string `json:"operation,omitempty"`
	Path      string `json:"path,omitempty"`
	Target    string `json:"target,omitempty"`
}

func getValue(updValue *gnmi.TypedValue) (interface{}, error) {
	if updValue == nil {
		return nil, nil
	}
	var value interface{}
	var jsondata []byte
	switch updValue.Value.(type) {
	case *gnmi.TypedValue_AsciiVal:
		value = updValue.GetAsciiVal()
	case *gnmi.TypedValue_BoolVal:
		value = updValue.GetBoolVal()
	case *gnmi.TypedValue_BytesVal:
		value = updValue.GetBytesVal()
	case *gnmi.TypedValue_DecimalVal:
		//lint:ignore SA1019 still need DecimalVal for backward compatibility
		value = updValue.GetDecimalVal()
	case *gnmi.TypedValue_FloatVal:
		//lint:ignore SA1019 still need GetFloatVal for backward compatibility
		value = updValue.GetFloatVal()
	case *gnmi.TypedValue_DoubleVal:
		value = updValue.GetDoubleVal()
	case *gnmi.TypedValue_IntVal:
		value = updValue.GetIntVal()
	case *gnmi.TypedValue_StringVal:
		value = updValue.GetStringVal()
	case *gnmi.TypedValue_UintVal:
		value = updValue.GetUintVal()
	case *gnmi.TypedValue_JsonIetfVal:
		jsondata = updValue.GetJsonIetfVal()
	case *gnmi.TypedValue_JsonVal:
		jsondata = updValue.GetJsonVal()
	case *gnmi.TypedValue_LeaflistVal:
		value = updValue.GetLeaflistVal()
	case *gnmi.TypedValue_ProtoBytes:
		value = updValue.GetProtoBytes()
	case *gnmi.TypedValue_AnyVal:
		value = updValue.GetAnyVal()
	}
	if value == nil && len(jsondata) != 0 {
		err := json.Unmarshal(jsondata, &value)
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}

func formatGetResponse(m *gnmi.GetResponse) interface{} {
	notifications := make([]NotificationRspMsg, 0, len(m.GetNotification()))
	for _, notif := range m.GetNotification() {

		msg := NotificationRspMsg{
			Prefix:    path.GnmiPathToXPath(notif.GetPrefix(), false),
			Updates:   make([]update, 0, len(notif.GetUpdate())),
			Timestamp: notif.Timestamp,
			Time:      time.Unix(0, notif.Timestamp),
			Target:    notif.GetPrefix().GetTarget(),
		}

		for i, upd := range notif.GetUpdate() {
			pathElems := make([]string, 0, len(upd.GetPath().GetElem()))
			for _, pElem := range upd.GetPath().GetElem() {
				pathElems = append(pathElems, pElem.GetName())
			}

			value, err := getValue(upd.GetVal())
			if err != nil {
				return nil
			}
			msg.Updates = append(msg.Updates,
				update{
					Path:   path.GnmiPathToXPath(upd.GetPath(), false),
					Values: make(map[string]interface{}),
				})
			msg.Updates[i].Values[strings.Join(pathElems, "/")] = value
		}
		notifications = append(notifications, msg)
	}

	// return only xpath value
	result := make([]interface{}, 0, len(notifications))
	for _, n := range notifications {
		for _, u := range n.Updates {
			for _, v := range u.Values {
				result = append(result, v)
			}
		}
	}

	//fmt.Println(notifications, result)
	return result[0]
}

func getRequest(c GnmiConnect, xpath string) interface{} {
	getReq, err := api.NewGetRequest(
		api.Path(xpath),
		api.Encoding("json_ietf"))
	if err != nil {
		log.Fatal(err)
	}

	getResp, err := c.Target.Get(c.Context, getReq)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(prototext.Format(getResp))
	return formatGetResponse(getResp)
}

func formatSubscribeResponse(m *gnmi.SubscribeResponse) interface{} {
	switch m := m.GetResponse().(type) {
	case *gnmi.SubscribeResponse_Update:
		msg := NotificationRspMsg{
			Timestamp: m.Update.Timestamp,
		}

		msg = NotificationRspMsg{
			Prefix:    path.GnmiPathToXPath(m.Update.GetPrefix(), false),
			Updates:   make([]update, 0, len(m.Update.GetUpdate())),
			Timestamp: m.Update.Timestamp,
			Time:      time.Unix(0, m.Update.Timestamp),
			Target:    m.Update.GetPrefix().GetTarget(),
		}

		for i, upd := range m.Update.Update {
			if upd.Path == nil {
				upd.Path = new(gnmi.Path)
			}
			pathElems := make([]string, 0, len(upd.Path.Elem))
			for _, pElem := range upd.Path.Elem {
				pathElems = append(pathElems, pElem.GetName())
			}
			value, err := getValue(upd.Val)
			if err != nil {
				return nil
			}
			msg.Updates = append(msg.Updates,
				update{
					Path:   path.GnmiPathToXPath(upd.Path, false),
					Values: make(map[string]interface{}),
				})
			msg.Updates[i].Values[strings.Join(pathElems, "/")] = value
		}

		// fmt.Println(msg)
		// return only xpath value
		for _, u := range msg.Updates {
			for _, v := range u.Values {
				return v
			}
		}
	}
	return nil
}

func formatSetResponse(m *gnmi.SetResponse, pretty bool) ([]byte, error) {
	responses := setRspMsg{
		Prefix:    path.GnmiPathToXPath(m.GetPrefix(), false),
		Target:    m.GetPrefix().GetTarget(),
		Timestamp: m.Timestamp,
		Time:      time.Unix(0, m.Timestamp),
	}

	responses.Results = make([]updateResultMsg, 0, len(m.GetResponse()))
	for _, u := range m.GetResponse() {
		responses.Results = append(responses.Results, updateResultMsg{
			Operation: u.Op.String(),
			Path:      path.GnmiPathToXPath(u.GetPath(), false),
			Target:    u.GetPath().GetTarget(),
		})
	}

	if pretty {
		return json.MarshalIndent(responses, "", "  ")
	}
	return json.Marshal(responses)
}

func setRequest(c GnmiConnect, xpath string, value string) ([]byte, error) {
	setReq, err := api.NewSetRequest(
		api.Update(
			api.Path(xpath),
			api.Value(value, "json_ietf")),
	)
	if err != nil {
		log.Fatal(err)
	}

	setResp, err := c.Target.Set(c.Context, setReq)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(prototext.Format(setResp))
	return formatSetResponse(setResp, c.PrettyResponse)
}

func initGnmic(c GnmiConnect) *target.Target {
	tg, err := api.NewTarget(
		api.Address(fmt.Sprintf("%s:%s", c.Host, c.Port)),
		api.Username(c.Username),
		api.Password(c.Password),
		api.Insecure(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = tg.CreateGNMIClient(c.Context)
	if err != nil {
		log.Fatal(err)
	}

	return tg
}
