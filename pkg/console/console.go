package console

import (
	"github.com/sirupsen/logrus"
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"k8s.io/client-go/tools/clientcmd"
	"errors"
	"encoding/base64"
	deis "github.com/teamhephy/controller-sdk-go"
)

var log = logrus.New()

func init() {
	log.SetLevel(logrus.PanicLevel)
	log.SetOutput(os.Stdout)
}

type KubeApiAccess struct {
	ApiEndpoint string
	Token 		string
	Error 		bool
	Msg			string
}

var myKubeApiAccess KubeApiAccess

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
	if (httpResp.StatusCode == 401 || httpResp.StatusCode == 403) {
		return errors.New("\nPermission denied. Please ensure that you have access to the application '"+application+"'")
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

	containerName := applicationName+"-"+procType
	opts := &ExecOptions{}
	opts.Namespace = applicationName
	opts.Pod = podName
	opts.Container = containerName
	opts.TTY = true 
	opts.Stdin = true
	opts.Command = []string{execCommand}

	//log.Printf("[Start] Controller Url: %s",deis.Client.ControllerURL)
	//log.Printf("[Start] User token: %s...",deis.Client.Token[0:5])
	log.Printf("[Start] Container name: %s",opts.Container)

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