package docdbuild

import (
	"fmt"
	"gopkg.in/go-playground/webhooks.v5/github"
	"net/http"
	"os/exec"
)

// AutoDeploy : Pulls latest commit from remote master and deploys
func AutoDeploy(res http.ResponseWriter, req *http.Request) {
	hook, _ := github.New(github.Options.Secret(""))

	payload, err := hook.Parse(req, github.PushEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			fmt.Println("Unknown event")
		}
	}

	switch payload.(type) {

	case github.PushPayload:
		push := payload.(github.PushPayload)
		fmt.Println("Change detected on", push.Ref, "- deploying...")
		cmd := exec.Command("git", "pull")
		err := cmd.Run()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("Deployment complete")
	}
	fmt.Fprintf(res, "Hello, %s!", req.URL.Path[1:])
}
