package console

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	deis "github.com/teamhephy/controller-sdk-go"
	"io/ioutil"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"time"
)

var log = logrus.New()

func init() {
	log.SetLevel(logrus.PanicLevel)
	log.SetOutput(os.Stdout)
}

type KubeApiAccess struct {
	ApiEndpoint      string
	Token            string
	WebsocketTimeout int
	Error            bool
	Msg              string
}

var myKubeApiAccess KubeApiAccess
var timeOutContext context.Context
var timeOutContextCancelation context.CancelFunc

func getKubernetesApiAndToken(c *deis.Client, application string) error {
	u := fmt.Sprintf("/v2/apps/%s/console-token", application)
	httpResp, httpErr := c.Request("GET", u, nil)
	if httpErr != nil {
		log.Println(httpErr)
		return httpErr
	}
	defer httpResp.Body.Close()

	body, httpErr := ioutil.ReadAll(httpResp.Body)
	if httpErr != nil {
		log.Println(httpErr)
		return httpErr
	}
	// No hephy token provided == 401; Wrong hephy token provided == 403
	if httpResp.StatusCode == 401 || httpResp.StatusCode == 403 {
		return errors.New("\nPermission denied. Please ensure that you have access to the application '" + application + "'")
	}
	json.Unmarshal([]byte(body), &myKubeApiAccess)
	if myKubeApiAccess.Error {
		return errors.New(myKubeApiAccess.Msg)
	}

	return nil
}

func Start(c *deis.Client, applicationName string, podName string, procType string, execCommand string, debug bool) error {
	if debug {
		log.SetLevel(logrus.DebugLevel)
	}
	collectParametersError := getKubernetesApiAndToken(c, applicationName)
	if collectParametersError != nil {
		return collectParametersError
	}

	containerName := applicationName + "-" + procType
	opts := &ExecOptions{}
	opts.Namespace = applicationName
	opts.Pod = podName
	opts.Container = containerName
	opts.TTY = true
	opts.Stdin = true
	opts.Command = []string{execCommand}

	log.Printf("[Start] Container name: %s", opts.Container)

	sDec, err := base64.StdEncoding.DecodeString(myKubeApiAccess.Token)
	if err != nil {
		log.Println(err)
		return err
	}
	myKubeApiAccess.Token = string(sDec)
	config, err := clientcmd.BuildConfigFromFlags(myKubeApiAccess.ApiEndpoint, "")
	if err != nil {
		log.Println(err)
		return err
	}

	timeOutContext, timeOutContextCancelation = context.WithTimeout(context.TODO(), time.Duration(myKubeApiAccess.WebsocketTimeout)*time.Second)
	defer timeOutContextCancelation()

	wrt, err := ExecRoundTripper(config, WebsocketCallback)
	if err != nil {
		log.Println(err)
		return err
	}

	req, err := ExecRequest(config, opts)
	if err != nil {
		log.Println(err)
		return err
	}

	if _, err = wrt.RoundTrip(req); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
