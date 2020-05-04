package cmd

import (
	"io/ioutil"
	"net/http/cookiejar"

	gabs "github.com/Jeffail/gabs/v2"
	"github.com/SAP/jenkins-library/pkg/command"
	piperhttp "github.com/SAP/jenkins-library/pkg/http"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/pkg/errors"
)

func gctsDeployCommit(config gctsDeployCommitOptions, telemetryData *telemetry.CustomData) {
	// for command execution use Command
	c := command.Command{}
	// reroute command output to logging framework
	c.Stdout(log.Entry().Writer())
	c.Stderr(log.Entry().Writer())

	// for http calls import  piperhttp "github.com/SAP/jenkins-library/pkg/http"
	// and use a  &piperhttp.Client{} in a custom system
	// Example: step checkmarxExecuteScan.go
	httpClient := &piperhttp.Client{}

	// error situations should stop execution through log.Entry().Fatal() call which leads to an os.Exit(1) in the end
	err := deployCommit(&config, telemetryData, &c, httpClient)
	if err != nil {
		log.Entry().WithError(err).Fatal("step execution failed")
	}
}

func deployCommit(config *gctsDeployCommitOptions, telemetryData *telemetry.CustomData, command execRunner, httpClient piperhttp.Sender) error {

	cookieJar, cookieErr := cookiejar.New(nil)
	if cookieErr != nil {
		return errors.Wrap(cookieErr, "deploy commit failed")
	}
	clientOptions := piperhttp.ClientOptions{
		CookieJar: cookieJar,
		Username:  config.Username,
		Password:  config.Password,
	}
	httpClient.SetOptions(clientOptions)

	url := config.Host +
		"/sap/bc/cts_abapvcs/repository/" + config.Repository +
		"/pullByCommit?sap-client=" + config.Client
	if config.Commit != "" {
		url = url + "&request=" + config.Commit
	}

	resp, httpErr := httpClient.SendRequest("GET", url, nil, nil, nil)

	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp == nil {
		return errors.New("did not retrieve a HTTP response")
	} else if httpErr != nil {
		return httpErr
	}

	bodyText, readErr := ioutil.ReadAll(resp.Body)

	if readErr != nil {
		return errors.Wrapf(readErr, "HTTP response body could not be read")
	}

	response, parsingErr := gabs.ParseJSON([]byte(bodyText))

	if parsingErr != nil {
		return errors.Wrapf(parsingErr, "HTTP response body could not be parsed as JSON: %v", string(bodyText))
	}

	log.Entry().
		WithField("repository", config.Repository).
		Infof("successfully deployed commit %v (previous commit was %v)", response.Path("toCommit").Data().(string), response.Path("fromCommit").Data().(string))
	return nil
}
