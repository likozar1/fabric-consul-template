package main

import (
	"encoding/json"
	"flag"
	"log"
	"bufio"
	"fmt"
	"os"
	"time"
	"strings"
	"regexp"
	"github.com/BurntSushi/toml"
)

var (
	config        		= flag.String("c", "/go/src/fabric_consul/config/config.toml", "Path to configuration file")
	c             		= &Config{}
	URL					= ""
)

const(
	PUT 				  = "PUT"
	DELETE 				  = "DELETE"
	POST 				  = "POST"
	GET 				  = "GET"
)

type Config struct {
		Consul map[string]string
		Ssl map[string]bool
		Retry map[string]int
}

func main() {
	flag.Parse()
	// read the config
	if _, err := toml.DecodeFile(*config, c); err != nil {
		log.Panic("Failed to read the configuration file", err)
	}
	//set from config
	URL = "http://"+c.Consul["host"]+"/v1/catalog/nodes"
	//call retry
	err:=retry(c.Retry["retry"])
	//check if retry failed
	if err==0{
		log.Println("+ python file generated SUCCESSFULLY")
	}else{
		log.Println("- FAILED to generate python file")
		log.Println("- Failed after ",c.Retry["retry"]," retries")
	}
}

//get data from consul catalog
func getConsulNodes() (*ConsulNodes){
	res:=&ConsulNodes{}
	byt, _:= getRequest(URL)
    json.Unmarshal(byt, &res)
    return res
}

//return node prefix
func getPrefix(name string)string{
	splitedstring:=strings.Split(name,"-")
	if len(splitedstring)>1{
	return splitedstring[0]+"-"+getPrefixNoNumbers(splitedstring[1])
	}else{
		return name
	}
}

//return prefix without node number
func getPrefixNoNumbers(prefix string) string{
	result:=""
	for _, a := range prefix{
		if checkPrefix(string(a)){
			result+=string(a)
		}else {
			break;
		}
	}
	return result
}

//Order node by prefix
func orderNodes(catalog *ConsulNodes)*prefixedNodes{
	orderdNodes:=&prefixedNodes{}
	for _,node := range *catalog{
		orderdNodes.append(getPrefix(node.Node),node.Node)
	}
	return orderdNodes
}

//returns true if string is not numeric
func checkPrefix(prefix string) bool{
    match, _ := regexp.MatchString("([a-zA-Z]+)", prefix)
    if  !match{
     return false
    }else{
     return true
    }   

} 

//return all nodes
func getAllNodes(catalog *ConsulNodes)string{
	result:="\t'all': ["
	for _,node := range *catalog{
		result+="\n\t\t'"+node.Node+"',"
	}
	result+="\n\t],"
	return result
}

//Generate string out of ordered nodes
func getOrderedNodes(catalog *ConsulNodes)string{
	orderedMap:=orderNodes(catalog)
	result:=""
	for key, mapped := range *orderedMap{
		result+="\t'"+key+"': [\n"
		for  _,node := range mapped.nodes{
			result+="\t\t'"+node+"',\n"
		}
		result+="\t],\n"
	}
	return result
}

//Generate python file
func genFabfile(catalog *ConsulNodes){
	log.Println("Generating python file from data")
	start_of_file := ("roledefs = {")
	end_of_file := ("}")
	fileHandle, _ := os.Create(c.Consul["outfile"]+".py")
	writer := bufio.NewWriter(fileHandle)

	fmt.Fprintln(writer, start_of_file)
	fmt.Fprintln(writer,getAllNodes(catalog))
	fmt.Fprintln(writer,getOrderedNodes(catalog))

	fmt.Fprintln(writer, end_of_file)
    writer.Flush()
}

//try n times to generate data
func retry(attempts int) (err int) {
    for i := 0; ; i++ {
        consulnode:=getConsulNodes()
        if len(*consulnode) != 0 {
        	log.Println("Fetching consul data")
        	genFabfile(consulnode)
            return 0
        }else{
        	log.Println("Something went wrong")
        }

        log.Println("retrying...  retry no:",i+1,"\n")
        time.Sleep(time.Second * 2)
        if i >= (attempts - 1) {
            break
        }
        
    }
    return 127
}