package main

// #cgo CFLAGS: -I/opt/halon/include
// #cgo LDFLAGS: -Wl,--unresolved-symbols=ignore-all
// #include <HalonMTA.h>
// #include <stdlib.h>
// #include <syslog.h>
// static void c_syslog(int priority, const char *message) {
//     syslog(priority, "%s", message);
// }
import "C"
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

type Project struct {
	Id          string   `json:"id"`
	Topics      []string `json:"topics"`
	Credentials string   `json:"credentials"`
}

type Config struct {
	Projects []Project `json:"projects"`
}

type Topic struct {
	Id    string
	Topic *pubsub.Topic
}

type Client struct {
	Id     string
	Client *pubsub.Client
	Topics []Topic
}

var clients []Client

//export Halon_version
func Halon_version() C.int {
	return C.HALONMTA_PLUGIN_VERSION
}

//export Halon_init
func Halon_init(hic *C.HalonInitContext) C.bool {
	var config *C.HalonConfig
	if !C.HalonMTA_init_getinfo(hic, C.HALONMTA_INIT_CONFIG, nil, 0, unsafe.Pointer(&config), nil) {
		Syslog(C.LOG_CRIT, "pubsub: Could not get init config")
		return false
	}

	cfg, err := GetConfigAsJSON(config)
	if err == nil {
		var parsedConfig Config
		json.Unmarshal([]byte(cfg), &parsedConfig)

		for _, project := range parsedConfig.Projects {
			if project.Id == "" {
				Syslog(C.LOG_CRIT, "pubsub: Missing required \"id\" setting for project")
				return false
			}
			if len(project.Topics) == 0 {
				Syslog(C.LOG_CRIT, "pubsub: Missing or empty \"topics\" setting for project")
				return false
			}

			var c *pubsub.Client
			ctx := context.Background()
			if project.Credentials != "" {
				c, err = pubsub.NewClient(ctx, project.Id, option.WithCredentialsFile(project.Credentials))
			} else {
				c, err = pubsub.NewClient(ctx, project.Id)
			}
			if err != nil {
				Syslog(C.LOG_CRIT, "pubsub: "+err.Error())
				return false
			}
			client := Client{Id: project.Id, Client: c}
			for _, topicId := range project.Topics {
				topic := Topic{Id: topicId, Topic: c.Topic(topicId)}
				client.Topics = append(client.Topics, topic)
			}
			clients = append(clients, client)
		}
	}

	return true
}

//export pubsub_publish
func pubsub_publish(hhc *C.HalonHSLContext, args *C.HalonHSLArguments, ret *C.HalonHSLValue) {
	projectId, err := GetArgumentAsString(args, 0, true)
	if err != nil {
		SetException(hhc, "pubsub: "+err.Error())
		return
	}

	topicId, err := GetArgumentAsString(args, 1, true)
	if err != nil {
		SetException(hhc, "pubsub: "+err.Error())
		return
	}

	message, err := GetArgumentAsString(args, 2, true)
	if err != nil {
		SetException(hhc, "pubsub: "+err.Error())
		return
	}

	var client *pubsub.Client
	var topic *pubsub.Topic

	for _, c := range clients {
		if c.Id == projectId {
			client = c.Client
			for _, t := range c.Topics {
				if t.Id == topicId {
					topic = t.Topic
				}
			}
		}
	}

	if client == nil {
		SetException(hhc, "pubsub: No project matched the argument at position 0")
		return
	}
	if topic == nil {
		SetException(hhc, "pubsub: No topic matched the argument at position 1")
		return
	}

	ctx := context.Background()
	res := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(message),
	})
	go getTask(hhc, ret, ctx, res)
	C.HalonMTA_hsl_suspend_return(hhc)
}

func getTask(hhc *C.HalonHSLContext, ret *C.HalonHSLValue, ctx context.Context, res *pubsub.PublishResult) {
	id, err := res.Get(ctx)
	if err != nil {
		value := map[string]interface{}{"id": id, "error": err.Error()}
		SetReturnValueToAny(ret, value)
		C.HalonMTA_hsl_schedule(hhc)
		return
	}
	value := map[string]interface{}{"id": id}
	SetReturnValueToAny(ret, value)
	C.HalonMTA_hsl_schedule(hhc)
}

//export Halon_hsl_register
func Halon_hsl_register(hhrc *C.HalonHSLRegisterContext) C.bool {
	pubsub_publish_cs := C.CString("pubsub_publish")
	C.HalonMTA_hsl_register_function(hhrc, pubsub_publish_cs, nil)
	C.HalonMTA_hsl_module_register_function(hhrc, pubsub_publish_cs, nil)
	return true
}

func GetConfigAsJSON(cfg *C.HalonConfig) (string, error) {
	var x *C.char
	y := C.HalonMTA_config_to_json(cfg, &x, nil)
	defer C.free(unsafe.Pointer(x))
	if y {
		return C.GoString(x), nil
	} else {
		if x != nil {
			return "", errors.New(C.GoString(x))
		} else {
			return "", errors.New("failed to get config")
		}
	}
}

func GetArgumentAsString(args *C.HalonHSLArguments, pos uint64, required bool) (string, error) {
	var x = C.HalonMTA_hsl_argument_get(args, C.ulong(pos))
	if x == nil {
		if required {
			return "", fmt.Errorf("missing argument at position %d", pos)
		} else {
			return "", nil
		}
	}
	var y *C.char
	if C.HalonMTA_hsl_value_get(x, C.HALONMTA_HSL_TYPE_STRING, unsafe.Pointer(&y), nil) {
		return C.GoString(y), nil
	} else {
		return "", fmt.Errorf("invalid argument at position %d", pos)
	}
}

func SetReturnValueToAny(ret *C.HalonHSLValue, val interface{}) error {
	x, err := json.Marshal(val)
	if err != nil {
		return err
	}
	y := C.CString(string(x))
	defer C.free(unsafe.Pointer(y))
	var z *C.char
	if !(C.HalonMTA_hsl_value_from_json(ret, y, &z, nil)) {
		if z != nil {
			err = errors.New(C.GoString(z))
			C.free(unsafe.Pointer(z))
		} else {
			err = errors.New("failed to parse return value")
		}
		return err
	}
	return nil
}

func SetException(hhc *C.HalonHSLContext, msg string) {
	x := C.CString(msg)
	y := unsafe.Pointer(x)
	defer C.free(y)
	exception := C.HalonMTA_hsl_throw(hhc)
	C.HalonMTA_hsl_value_set(exception, C.HALONMTA_HSL_TYPE_EXCEPTION, y, 0)
}

// A wrapper for the "syslog" C function
func Syslog(pri int, str string) {
	x := C.CString(str)
	defer C.free(unsafe.Pointer(x))
	C.c_syslog(C.int(pri), x)
}

func main() {}
