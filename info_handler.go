package podstalk

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mhutter/podstalk/k8s"
)

type InfoHandler struct {
	info PodInfo
	kc   *k8s.Client
	t    *template.Template
}

// PodInfo contains all interesting information about a pod
type PodInfo struct {
	Name           string
	Namespace      string
	IP             string
	ServiceAccount string
	NodeName       string
	NodeIP         string
	Info           map[string]string
	Now            string
	Title          string
	Siblings       []string
	Refresh        int
	BasePath       string
}

func NewInfoHandler() InfoHandler {
	h := InfoHandler{
		info: collectPodInfo(),
		kc:   k8s.NewClient(),
		t:    template.Must(template.New("index").Parse(htmlTemplate)),
	}

	return h
}

func (h InfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := h.info
	p.Now = time.Now().Format(time.RFC3339Nano)

	if h.kc != nil {
		for _, pod := range h.kc.ListPods().Items {
			p.Siblings = append(p.Siblings, pod.Metadata.Name)
		}
	}

	var err error
	p.Refresh, err = strconv.Atoi(r.URL.Query().Get("refresh"))
	if err != nil {
		p.Refresh = 0
	}

	err = h.t.Execute(w, p)
	if err != nil {
		log.Fatal(err)
	}
}

func collectPodInfo() PodInfo {
	return PodInfo{
		Name:           os.Getenv("POD_NAME"),
		Namespace:      os.Getenv("POD_NAMESPACE"),
		IP:             os.Getenv("POD_IP"),
		ServiceAccount: os.Getenv("POD_SA"),
		NodeName:       os.Getenv("NODE_NAME"),
		NodeIP:         os.Getenv("NODE_IP"),
		Title:          GetEnvOr("TITLE", "Podstalk"),
		Info:           collectEnv(),
		BasePath:       os.Getenv("BASE_PATH"),
	}
}

func collectEnv() (info map[string]string) {
	pattern := regexp.MustCompile("^info_([^=]+)=(.+)$")

	for _, line := range os.Environ() {
		l := strings.ToLower(line)
		if m := pattern.FindStringSubmatch(l); m != nil {
			info[m[1]] = m[2]
		}
	}
	return
}
