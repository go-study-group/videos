package youtube

import (
	"bufio"
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

func GetService(clientId string, secret string, tokenStore TokenStore) *youtube.Service {
	oauth2Config := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: secret,
		RedirectURL:  "http://localhost:8080",
		Scopes: []string{
			youtube.YoutubeUploadScope,
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://www.googleapis.com/oauth2/v3/token",
		},
	}

	ctx := context.Background()

	token := getToken(ctx, tokenStore, oauth2Config)
	client := oauth2Config.Client(ctx, token)
	service, serviceErr := youtube.New(client)
	handleError(serviceErr, "unable to get new service")

	return service
}

func getToken(ctx context.Context, tokenStore TokenStore, config *oauth2.Config) *oauth2.Token {
	token, tokenFromFileErr := tokenStore.Load()
	if tokenFromFileErr != nil {
		authorizationCode := getAuthorizationCode(config)

		var getTokenErr error
		token, getTokenErr = config.Exchange(ctx, authorizationCode)
		if getTokenErr != nil {
			log.Fatalf("Unable to retrieve token from web %v", getTokenErr)
		}

		saveErr := tokenStore.Save(token)
		if saveErr != nil {
			log.Fatalf("unable to save token: %v", saveErr)
		}
	}

	return token
}

func getAuthorizationCode(config *oauth2.Config) string {
	authorizationCode, codeFromWebErr := getAuthorizationCodeFromWeb(config)
	if codeFromWebErr != nil {
		fmt.Printf("Unable to get authorization code from web: %v\n", codeFromWebErr)
		authorizationCode = getAuthorizationCodeFromCLI(config)
	}

	return authorizationCode
}

func getAuthorizationCodeFromCLI(config *oauth2.Config) string {
	cliConfig := config
	cliConfig.RedirectURL = "urn:ietf:wg:oauth:2.0:oob"

	authCodeUrl := cliConfig.AuthCodeURL("")

	fmt.Println("Visit the URL below to get a code.",
		" This program will pause until the site is visted.")

	fmt.Println(authCodeUrl)
	fmt.Print("Enter code: ")
	scanner := bufio.NewScanner(os.Stdin)

	var code string
	for scanner.Scan() {
		line := scanner.Text()
		code = line
		break
	}

	return code
}

func getAuthorizationCodeFromWeb(config *oauth2.Config) (string, error) {
	//return "", fmt.Errorf("breaking to test cli version")

	codeCh, startWebErr := startWebServer()
	handleError(startWebErr, "Unable to start web server")

	authCodeUrl := config.AuthCodeURL("")
	openUrlErr := openURL(authCodeUrl)
	if openUrlErr != nil {
		return "", openUrlErr
	}

	fmt.Println("Your browser has been opened to an authorization URL.",
		"This program will resume once authorization has been provided.")

	code := <-codeCh

	return code, nil
}

func startWebServer() (codeCh chan string, err error) {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return nil, err
	}
	codeCh = make(chan string)
	go http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		codeCh <- code // send code to OAuth flow
		listener.Close()
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Received code: %v\r\nYou can now safely close this browser window.", code)
	}))

	return codeCh, nil
}

// openURL opens a browser window to the specified location.
// This code originally appeared at:
//   http://stackoverflow.com/questions/10377243/how-can-i-launch-a-process-that-is-not-a-file-in-go
func openURL(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost:4001/").Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Cannot open URL %s on this platform", url)
	}
	return err
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}
