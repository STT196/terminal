package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

import (
	"context"
	"crypto/md5"
	_ "embed"
	"encoding/hex"
	"errors"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/muesli/termenv"


	"github.com/terminaldotshop/terminal/go/pkg/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/charmbracelet/wish/recover"
	gossh "golang.org/x/crypto/ssh"
)

type PasswordState int

const (
	PasswordSkip PasswordState = iota
	PasswordPossible
	PasswordWaiting
	PasswordAccepted
)

const permanentHostKey = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEAvSGNhXTkqQWKlnRrTU/tpAnfzAKs4onXRZ4kQnKMcTehkDSzjl+p
Orp8mHTwtsOMIhYTGH63dUeNuE0ZPAB06Ix8oSgmFmnMBOItirlb4yRZR/mNMRtXJEU4N/
KIWQTFEpxSFXHyWxrzLD1ab8VInBa4uDZpwKFFb9uxFhZRzbTjIealJXuu+KuzbnbzarA3
AFi367Z8WxqjFJaS/mUt90j4KSZFPstcL1I9hoVkxE1Pdox/7IT/NZ3bRUv3KAcv+W9iwZ
ec8nwzpiXPOHdiHt61n5VT6/Ls92l2z1CcyiU9Gl+TLiZtUqG6VzLz2uzC67ONV2RotIbq
PBLl3g3ZPQAAA8gO33pFDt96RQAAAAdzc2gtcnNhAAABAQC9IY2FdOSpBYqWdGtNT+2kCd
/MAqziiddFniRCcoxxN6GQNLOOX6k6unyYdPC2w4wiFhMYfrd1R424TRk8AHTojHyhKCYW
acwE4i2KuVvjJFlH+Y0xG1ckRTg38ohZBMUSnFIVcfJbGvMsPVpvxUicFri4NmnAoUVv27
EWFlHNtOMh5qUle674q7NudvNqsDcAWLfrtnxbGqMUlpL+ZS33SPgpJkU+y1wvUj2GhWTE
TU92jH/shP81ndtFS/coBy/5b2LBl5zyfDOmJc84d2Ie3rWflVPr8uz3aXbPUJzKJT0aX5
MuJm1SobpXMvPa7MLrs41XZGi0huo8EuXeDdk9AAAAAwEAAQAAAQB47qhgKlM/ZCSueXhW
8gGgvxOTji5fmAXHJQxIVJhKmGi9HYWmRrKds7qRfUyhgD3tWbISGoxR+FO9Acdd32jhfV
r/bP2VnUZv5PN73XPMtGRGKmJGgRXiQkRlObZHPU6JzNyLi9WMvZm5su1NxJbd/4VTfK94
FWah1JbR6ama3rHuSELcT2rWGh0m62Sgx51GbUDAXFT2KiQADyzAYzMVjbKcG42u5saNKk
nZOjJN7EWy+haSnXwAa7K90yIMo16iak4CNtorazQ5p8AVK3ZVb+pmiYbitneVOPbea1zR
AindLx9FO9Q5d8f4fZ5sLRP7nF1D6HDUMwRNrwqY+209AAAAgQCrP3qyCuEd6CRLp9O8EO
8O810fCazuPT0IDvUaD94yH8/K0ZU5PIVsyOhdCe0egiiVw3LtIetPOMMAfW1C4Iiq8NCd
uZwDf74nu67QaPIP77+X3IZ46ipuukZnCgLxfpaZYscEZyuri2Fbb+P4BzzjQvOpkWUGzh
Y4X/M9EBK/WgAAAIEA6WkLdSLopxY600yyLT2Xg/WzzaElFBbLNPdhXVPq6CE3XVAZmyyf
snk478h5mkTD7dXyHV2vIW9/q1I/xXz6u6ZsOR5t3GXNAAI3aTYcxhoCWERMsWb1dDjSLR
CJdx8xyJu6pRYqVkRukfhgR+ZoGY2PEpX2PXa/sRIZlQRzGKMAAACBAM9vctgPgvmthLp4
KNr+qPoELwcUWxC1dJpMt4d4B4LkQV1OcKotnBEZ8jx2emqqKYCE8tvtqntsvofclExnHQ
bO+9qHLBh+bf7uqdmw1mTo6QxOEf4sFTc7nZYAki+wwsXUHCuBf+7TVTs4pmM/FREFoeMQ
6TT4GWlV8EwRBQSfAAAAEXRlcm1pbmFsLXNzaC1ob3N0AQ==
-----END OPENSSH PRIVATE KEY-----`

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		cancel()
	}()
	sshPort := os.Getenv("SSH_PORT")
	if sshPort == "" {
		sshPort = "2222"
	}

	// Load SSH host key from environment or use embedded fallback for dev
	hostKeyPEM := []byte(os.Getenv("SSH_HOST_KEY"))
	if len(hostKeyPEM) == 0 {
		log.Warn("SSH_HOST_KEY not set, using embedded key (not recommended for production)")
		hostKeyPEM = []byte(permanentHostKey)
	}

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort("0.0.0.0", sshPort)),
		wish.WithHostKeyPEM(hostKeyPEM),
		wish.WithMiddleware(
			recover.Middleware(
				bubbletea.Middleware(teaHandler),
				activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
				logging.Middleware(),
			),
		),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			hash := md5.Sum(key.Marshal())
			fingerprint := hex.EncodeToString(hash[:])
			ctx.SetValue("fingerprint", fingerprint)
			ctx.SetValue("anonymous", false)
			return true
		}),
		wish.WithKeyboardInteractiveAuth(
			func(ctx ssh.Context, challenger gossh.KeyboardInteractiveChallenge) bool {
				ctx.SetValue("fingerprint", uuid.NewString())
				ctx.SetValue("anonymous", true)
				return true
			},
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
		return
	}

	log.Info("Starting SSH server", "port", sshPort)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			cancel()
		}
	}()

	<-ctx.Done()
	s.Shutdown(ctx)
	slog.Info("Shutting down server")
}

type sshOutput struct {
	ssh.Session
}

func (s *sshOutput) Write(p []byte) (int, error) {
	return s.Session.Write(p)
}

func (s *sshOutput) Read(p []byte) (int, error) {
	return s.Session.Read(p)
}

func (s *sshOutput) Fd() uintptr {
	return 0
}

// You can wire any Bubble Tea model up to the middleware with a function that
// handles the incoming ssh.Session. Here we just grab the terminal info and
// pass it to the new model. You can also return tea.ProgramOptions (such as
// tea.WithAltScreen) on a session by session basis.
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	sessionBridge := &sshOutput{
		Session: s,
	}
	renderer := bubbletea.MakeRenderer(sessionBridge)
	fingerprint := s.Context().Value("fingerprint").(string)
	anonymous := s.Context().Value("anonymous").(bool)
	command := s.Command()
	slog.Info("got fingerprint", "fingerprint", fingerprint)
	slog.Info("got command", "command", command)

	// Get client IP address from the SSH session
	clientAddr := s.RemoteAddr().String()
	host, _, _ := net.SplitHostPort(clientAddr)
	slog.Info("client connected", "ip", host)

	if pty.Term == "xterm-ghostty" {
		renderer.SetColorProfile(termenv.TrueColor)
	}

	model, err := tui.NewModel(renderer, fingerprint, anonymous, &host, command)
	if err != nil {
		return nil, []tea.ProgramOption{}
	}
	return model, []tea.ProgramOption{tea.WithAltScreen()}
}
