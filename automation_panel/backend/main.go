package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/openconfig/gnmic/pkg/api"
	"golang.org/x/crypto/ssh"
)

type srv struct {
	be     *mux.Router
	logger *log.Logger
	Link   *LinkGroup
	Anysec *AnysecGroup
	Icmp   *IcmpGroup
}

type LinkGroup struct {
	Top    []LinkEndpoint
	Bottom []LinkEndpoint
}

type AnysecGroup struct {
	Vll  []AnysecEndpoint
	Vpls []AnysecEndpoint
	Vprn []AnysecEndpoint
}

type IcmpGroup struct {
	Vll  IcmpEndpoint
	Vpls IcmpEndpoint
	Vprn IcmpEndpoint
}

type LinkEndpoint struct {
	Host       string
	Port       string
	Reference  string
	AdminState bool
	Gnmi       GnmiConnect
}

type AnysecEndpoint struct {
	Host       string
	GroupName  string
	Peer       string
	Reference  string
	AdminState bool
	Gnmi       GnmiConnect
}

type IcmpEndpoint struct {
	SshHost     string
	SshUser     string
	SshPass     string
	Destination string
	Size        int
	Interval    float64
	AdminState  bool
	Pid         int
}

type IcmpAddons struct {
	Size     int     `json:"size"`
	Interval float64 `json:"interval"`
}

type serviceState struct {
	Vll  bool `json:"vll"`
	Vpls bool `json:"vpls"`
	Vprn bool `json:"vprn"`
}

type linkState struct {
	Top    bool `json:"top"`
	Bottom bool `json:"bottom"`
}

type allState struct {
	Icmp   serviceState `json:"icmp"`
	Link   linkState    `json:"link"`
	Anysec serviceState `json:"anysec"`
}

// API Handlers

func (s *srv) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("\n")
		s.logger.Printf("REQUEST: %s %s %s", r.RemoteAddr, r.Method, r.URL)

		next.ServeHTTP(w, r)

		corHeader := "Access-Control-Allow-Origin"
		if w.Header().Get(corHeader) != "*" {
			w.Header().Set(corHeader, "*")
		}
	})
}

func connectionOk(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("connected"))
}

func (s *srv) getAllState(w http.ResponseWriter, r *http.Request) {
	icmp := *s.Icmp
	link := *s.Link
	anysec := *s.Anysec

	response := allState{
		Icmp: serviceState{
			Vll:  icmp.Vll.AdminState,
			Vpls: icmp.Vpls.AdminState,
			Vprn: icmp.Vprn.AdminState,
		},
		Link: linkState{
			Top:    link.Top[0].AdminState && link.Top[1].AdminState,
			Bottom: link.Bottom[0].AdminState && link.Bottom[1].AdminState,
		},
		Anysec: serviceState{
			Vll:  anysec.Vll[0].AdminState && anysec.Vll[1].AdminState,
			Vpls: anysec.Vpls[0].AdminState && anysec.Vpls[1].AdminState,
			Vprn: anysec.Vprn[0].AdminState && anysec.Vprn[1].AdminState,
		},
	}

	jsonData, _ := json.Marshal(response)
	s.logger.Printf("[%d] - %s", http.StatusOK, jsonData)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(jsonData)
}

func (s *srv) setAdminState(w http.ResponseWriter, r *http.Request) {
	module := mux.Vars(r)["module"]
	service := mux.Vars(r)["service"]
	state := mux.Vars(r)["state"]

	switch module {
	case "icmp":
		var body IcmpAddons
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			s.logger.Printf("ERROR decoding request body: %s", err.Error())
		}
		//fmt.Println(body)
		switch service {
		case "vll":
			s.Icmp.Vll.Size = body.Size
			s.Icmp.Vll.Interval = body.Interval
			err := s.sshIcmpPing(&s.Icmp.Vll, state)
			if err != nil {
				s.logger.Println(err)
			}
		case "vpls":
			s.Icmp.Vpls.Size = body.Size
			s.Icmp.Vpls.Interval = body.Interval
			err := s.sshIcmpPing(&s.Icmp.Vpls, state)
			if err != nil {
				s.logger.Println(err)
			}
		case "vprn":
			s.Icmp.Vprn.Size = body.Size
			s.Icmp.Vprn.Interval = body.Interval
			err := s.sshIcmpPing(&s.Icmp.Vprn, state)
			if err != nil {
				s.logger.Println(err)
			}
		}
	case "link":
		switch service {
		case "top":
			s.setLinkAdminState(s.Link.Top, state)
		case "bottom":
			s.setLinkAdminState(s.Link.Bottom, state)
		}
	case "anysec":
		switch service {
		case "vll":
			s.setAnysecAdminState(s.Anysec.Vll, state)
		case "vpls":
			s.setAnysecAdminState(s.Anysec.Vpls, state)
		case "vprn":
			s.setAnysecAdminState(s.Anysec.Vprn, state)
		}
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte("success"))
}

func (s *srv) sshIcmpPing(endpoint *IcmpEndpoint, action string) error {
	logFilePath := "/tmp/pingPids.txt"
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	config := &ssh.ClientConfig{
		User: endpoint.SshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(endpoint.SshPass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", endpoint.SshHost+":22", config)
	if err != nil {
		return fmt.Errorf("failed to dial ssh: %s", err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create ssh session: %s", err)
	}
	defer session.Close()

	switch action {
	case "enable":
		cmd := fmt.Sprintf("nohup ping -s %d -i %f %s &> /dev/null & echo $!", endpoint.Size, endpoint.Interval, endpoint.Destination)
		var stdoutBuf bytes.Buffer
		session.Stdout = &stdoutBuf

		if err := session.Run(cmd); err != nil {
			return fmt.Errorf("failed to initiate icmp ping: %s", err)
		}

		pid := strings.Trim(stdoutBuf.String(), "\n")
		s.logger.Printf("[%d] - Ping command is running on %s with PID: %s", http.StatusOK, endpoint.SshHost, pid)
		num, err := strconv.Atoi(pid)

		if err != nil {
			return fmt.Errorf("failed to convert PID string to integer: %s", err)
		}
		endpoint.Pid = num
		endpoint.AdminState = true

		logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open %s log file: %v", logFilePath, err)
		}
		defer logFile.Close()

		logEntry := fmt.Sprintf("%v\n", endpoint)
		if _, err := logFile.WriteString(logEntry); err != nil {
			return fmt.Errorf("failed to write to %s log file: %v", logFilePath, err)
		}
	case "disable":
		if endpoint.Pid != 0 {
			cmdKill := fmt.Sprintf("kill -9 %d", endpoint.Pid)
			cmdSed := exec.Command("sed", "-i", fmt.Sprintf("/%v/d", endpoint), logFilePath)

			if err := session.Run(cmdKill); err != nil {
				return fmt.Errorf("failed to terminate process: %s", err)
			}

			s.logger.Printf("[%d] - Ping command on %s with PID: %d terminated", http.StatusOK, endpoint.SshHost, endpoint.Pid)
			endpoint.Pid = 0
			endpoint.AdminState = false

			stdoutStderr, err := cmdSed.CombinedOutput()
			if err != nil {
				return fmt.Errorf("error when removing PID %d reference from log file %s: %s", endpoint.Pid, logFilePath, stdoutStderr)
			}
		}
	}

	return nil
}

func (s *srv) setLinkAdminState(endpoint []LinkEndpoint, action string) {
	for i := range endpoint {
		xpath := fmt.Sprintf("/configure/port[port-id=%s]/admin-state", endpoint[i].Port)
		jsonData, _ := setRequest(endpoint[i].Gnmi, xpath, action)
		s.logger.Printf("[%d] - %s", http.StatusOK, jsonData)
	}
}

func (s *srv) setAnysecAdminState(endpoint []AnysecEndpoint, action string) {
	for i := range endpoint {
		xpath := fmt.Sprintf("/configure/anysec/tunnel-encryption/encryption-group[group-name=%s]/peer[peer-ip-address=%s]/admin-state", endpoint[i].GroupName, endpoint[i].Peer)
		jsonData, _ := setRequest(endpoint[i].Gnmi, xpath, action)
		s.logger.Printf("[%d] - %s", http.StatusOK, jsonData)
	}
}

// Link Subscribe Handler

func linkSubscribeResponse(v *LinkEndpoint) {
	var mu sync.Mutex
	subRspChan, subErrChan := v.Gnmi.Target.ReadSubscriptions()
	for {
		select {
		case rsp := <-subRspChan:
			state := formatSubscribeResponse(rsp.Response)
			mu.Lock()
			if state == "enable" {
				v.AdminState = true
			} else if state == "disable" {
				v.AdminState = false
			}
			mu.Unlock()
		case tgErr := <-subErrChan:
			log.Fatalf("subscription %q stopped: %v", tgErr.SubscriptionName, tgErr.Err)
		}
	}
}

func linkSubscribeRequest(group *[]LinkEndpoint) {
	endpoint := *group
	for i := range endpoint {
		endpoint[i].Gnmi = GnmiConnect{
			Host:     endpoint[i].Host,
			Port:     "57400",
			Username: "admin",
			Password: "admin",
		}

		endpoint[i].Gnmi.Context, endpoint[i].Gnmi.Cancel = context.WithCancel(context.Background())
		endpoint[i].Gnmi.Target = initGnmic(endpoint[i].Gnmi)

		xpath := fmt.Sprintf("/configure/port[port-id=%s]/admin-state", endpoint[i].Port)
		subReq, err := api.NewSubscribeRequest(
			api.Encoding("json_ietf"),
			api.SubscriptionListMode("stream"),
			api.Subscription(
				api.Path(xpath),
				api.SubscriptionMode("on-change"),
			))

		if err != nil {
			log.Fatal(err)
		}

		go endpoint[i].Gnmi.Target.Subscribe(endpoint[i].Gnmi.Context, subReq, fmt.Sprintf("link-sub-%d", i))
		go linkSubscribeResponse(&endpoint[i])
	}
}

// Anysec Subscribe Handler

func anysecSubscribeResponse(v *AnysecEndpoint) {
	var mu sync.Mutex
	subRspChan, subErrChan := v.Gnmi.Target.ReadSubscriptions()
	for {
		select {
		case rsp := <-subRspChan:
			state := formatSubscribeResponse(rsp.Response)
			mu.Lock()
			if state == "enable" {
				v.AdminState = true
			} else if state == "disable" {
				v.AdminState = false
			}
			mu.Unlock()
		case tgErr := <-subErrChan:
			log.Fatalf("subscription %q stopped: %v", tgErr.SubscriptionName, tgErr.Err)
		}
	}
}

func anysecSubscribeRequest(group *[]AnysecEndpoint) {
	endpoint := *group
	for i := range endpoint {
		endpoint[i].Gnmi = GnmiConnect{
			Host:     endpoint[i].Host,
			Port:     "57400",
			Username: "admin",
			Password: "admin",
		}

		endpoint[i].Gnmi.Context, endpoint[i].Gnmi.Cancel = context.WithCancel(context.Background())
		endpoint[i].Gnmi.Target = initGnmic(endpoint[i].Gnmi)

		xpath := fmt.Sprintf("/configure/anysec/tunnel-encryption/encryption-group[group-name=%s]/peer[peer-ip-address=%s]/admin-state", endpoint[i].GroupName, endpoint[i].Peer)
		subReq, err := api.NewSubscribeRequest(
			api.Encoding("json_ietf"),
			api.SubscriptionListMode("stream"),
			api.Subscription(
				api.Path(xpath),
				api.SubscriptionMode("on-change"),
			))

		if err != nil {
			log.Fatal(err)
		}

		go endpoint[i].Gnmi.Target.Subscribe(endpoint[i].Gnmi.Context, subReq, fmt.Sprintf("anysec-sub-%d", i))
		go anysecSubscribeResponse(&endpoint[i])
	}
}

// Subscribe Trigger

func (s *srv) subscribeTrigger() {
	link := *s.Link
	linkSubscribeRequest(&link.Top)
	linkSubscribeRequest(&link.Bottom)

	anysec := *s.Anysec
	anysecSubscribeRequest(&anysec.Vll)
	anysecSubscribeRequest(&anysec.Vpls)
	anysecSubscribeRequest(&anysec.Vprn)
}

// Main

func main() {
	s := &srv{
		be:     mux.NewRouter(),
		logger: log.New(os.Stderr, "", log.LstdFlags),
		Link: &LinkGroup{
			Top: []LinkEndpoint{
				{
					Host:       "pe1",
					Port:       "1/1/c1/1",
					Reference:  "src",
					AdminState: false,
				},
				{
					Host:       "p3",
					Port:       "1/1/c2/1",
					Reference:  "dest",
					AdminState: false,
				},
			},
			Bottom: []LinkEndpoint{
				{
					Host:       "pe1",
					Port:       "1/1/c2/1",
					Reference:  "src",
					AdminState: false,
				},
				{
					Host:       "p4",
					Port:       "1/1/c2/1",
					Reference:  "dest",
					AdminState: false,
				},
			},
		},
		Anysec: &AnysecGroup{
			Vll: []AnysecEndpoint{
				{
					Host:       "pe1",
					GroupName:  "EG_VLL-1001",
					Peer:       "10.0.0.21",
					Reference:  "end-a",
					AdminState: false,
				},
				{
					Host:       "pe2",
					GroupName:  "EG_VLL-1001",
					Peer:       "10.0.0.11",
					Reference:  "end-b",
					AdminState: false,
				},
			},
			Vpls: []AnysecEndpoint{
				{
					Host:       "pe1",
					GroupName:  "EG_VPLS-1002",
					Peer:       "10.0.0.22",
					Reference:  "end-a",
					AdminState: false,
				},
				{
					Host:       "pe2",
					GroupName:  "EG_VPLS-1002",
					Peer:       "10.0.0.12",
					Reference:  "end-b",
					AdminState: false,
				},
			},
			Vprn: []AnysecEndpoint{
				{
					Host:       "pe1",
					GroupName:  "EG_VPRN-1003",
					Peer:       "10.0.0.2",
					Reference:  "end-a",
					AdminState: false,
				},
				{
					Host:       "pe2",
					GroupName:  "EG_VPRN-1003",
					Peer:       "10.0.0.1",
					Reference:  "end-b",
					AdminState: false,
				},
			},
		},
		Icmp: &IcmpGroup{
			Vll: IcmpEndpoint{
				SshHost:     "client7",
				SshUser:     "user",
				SshPass:     "multit00l",
				Destination: "192.168.51.8",
				Size:        2000,
				Interval:    0.01,
				AdminState:  false,
				Pid:         0,
			},
			Vpls: IcmpEndpoint{
				SshHost:     "client7",
				SshUser:     "user",
				SshPass:     "multit00l",
				Destination: "192.168.52.8",
				Size:        2000,
				Interval:    0.01,
				AdminState:  false,
				Pid:         0,
			},
			Vprn: IcmpEndpoint{
				SshHost:     "client7",
				SshUser:     "user",
				SshPass:     "multit00l",
				Destination: "192.168.63.8",
				Size:        2000,
				Interval:    0.01,
				AdminState:  false,
				Pid:         0,
			},
		},
	}

	s.subscribeTrigger()
	s.be.Use(s.logMiddleware)
	s.be.HandleFunc("/", connectionOk).Methods("GET")
	s.be.HandleFunc("/get_all_state", s.getAllState).Methods("GET")
	s.be.HandleFunc("/set/{module}/{service}/{state}", s.setAdminState).Methods("POST")

	s.logger.Printf("Launch Web UI via http://localhost:8080")
	err := http.ListenAndServe(":8080", s.be)
	if err != nil {
		s.logger.Printf("ListenAndServer Error: %s", err.Error())
	}
}
