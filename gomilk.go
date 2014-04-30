package main
 
import (
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
)

func readURL(url string) (string, error) {
	response, err := http.Get(url)
	
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	
	if err != nil {
		return "", err
	}

	return string(contents), nil
}
 
func main() {
	contents, err := readURL("http://127.0.0.1:9292/gmilk?package=milkode&query=def+test")
	// contents, err := readURL("http://www.yahoo.co.jp")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	fmt.Println(contents)
}
