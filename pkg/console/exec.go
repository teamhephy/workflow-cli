package console

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"io/ioutil"
	"github.com/gorilla/websocket"
	"k8s.io/client-go/rest"
)

// ExecOptions describe a execute request args.
type ExecOptions struct {
	Namespace string
	Pod       string
	Container string
	Command   []string
	TTY       bool
	Stdin     bool
}

type RoundTripCallback func(c *websocket.Conn) error

type WebsocketRoundTripper struct {
	TLSConfig *tls.Config
	Callback  RoundTripCallback
}

var cacheBuff bytes.Buffer

var protocols = []string{
	"v4.channel.k8s.io",
	"v3.channel.k8s.io",
	"v2.channel.k8s.io",
	"channel.k8s.io",
}

const (
	stdin = iota
	stdout
	stderr
)

func WebsocketCallback(c *websocket.Conn) error {
	errChan := make(chan error, 3)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		buf := make([]byte, 1025)
		for {
			n, err := os.Stdin.Read(buf[1:])
			if err != nil {
				if err == io.EOF {
					continue
				}
				errChan <- err
				return
			}

			cacheBuff.Write(buf[1:n])
			cacheBuff.Write([]byte{13, 10}) // == 13->\r, 10->\n
			if err := c.WriteMessage(websocket.BinaryMessage, buf[:n+1]); err != nil {
				errChan <- err
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			_, buf, err := c.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}

			if len(buf) > 1 {
				var w io.Writer
				switch buf[0] {
				case stdout:
					w = os.Stdout
				case stderr:
					w = os.Stderr
				}

				if w == nil {
					continue
				}
				s := strings.Replace(string(buf[1:]), cacheBuff.String(), "", -1)
				_, err = w.Write([]byte(s))
				if err != nil {
					errChan <- err
					return
				}
			}
			cacheBuff.Reset()
		}
	}()

	cc := make(chan os.Signal, 1)
	signal.Notify(cc, os.Interrupt)
	go func(){
		for _ = range cc {
			// sig is a ^C, handle it
			if err := c.WriteMessage(websocket.BinaryMessage, []byte{13, 10,13,10}); err != nil {
				errChan <- err
				return
			}
		}
	}()

	wg.Wait()
	close(errChan)
	err := <-errChan
	return err
}

func (wrt *WebsocketRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	dialer := &websocket.Dialer{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: wrt.TLSConfig,
		Subprotocols:    protocols,
	}
	log.Printf("[RoundTrip] Url: %s",r.URL.String())
	tokenLength := len(myKubeApiAccess.Token)
	if tokenLength > 11 {
		log.Printf("[RoundTrip] Decoded token: %s ... %s",myKubeApiAccess.Token[0:10],myKubeApiAccess.Token[tokenLength-10:])
	}else{
		log.Println(myKubeApiAccess.Token)
	}
	conn, resp, err := dialer.Dial(r.URL.String(), http.Header{"Authorization": []string{"Bearer "+myKubeApiAccess.Token} })
	if resp.StatusCode != 200 &&  resp.StatusCode != 101 {
		bodyBytes, resp_err := ioutil.ReadAll(resp.Body)
		if resp_err != nil {
			log.Fatal(resp_err)
		}
		bodyString := string(bodyBytes)
		log.Printf("[RoundTrip] Error connecting to remote!\n")
		log.Printf("[RoundTrip] HTTP Status: %d\n",resp.StatusCode)
		log.Printf("[RoundTrip] HTTP Protocol: %s\n",resp.Proto)
		log.Printf("[RoundTrip] HTTP headers: %#v\n",resp.Header)
		log.Printf("[RoundTrip] HTTP body: %s\n",bodyString)
	}
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return resp, wrt.Callback(conn)
}

func ExecRoundTripper(config *rest.Config, f RoundTripCallback) (http.RoundTripper, error) {
	rt := &WebsocketRoundTripper{
		Callback:  f,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return rest.HTTPWrappersForConfig(config, rt)
}

func ExecRequest(config *rest.Config, opts *ExecOptions) (*http.Request, error) {
	u, err := url.Parse(config.Host)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "https":
		u.Scheme = "wss"
	case "http":
		u.Scheme = "ws"
	default:
		return nil, fmt.Errorf("Unrecognised URL scheme in %v", u)
	}

	u.Path = fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/exec", opts.Namespace, opts.Pod)

	rawQuery := "stdout=true&tty=true"
	for _, c := range opts.Command {
		rawQuery += "&command=" + c
	}

	if opts.Container != "" {
		rawQuery += "&container=" + opts.Container
	}

	if opts.Stdin {
		rawQuery += "&stdin=true"
	}
	u.RawQuery = rawQuery

	return &http.Request{
		Method: http.MethodGet,
		URL:    u,
	}, nil
}