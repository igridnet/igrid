package mq

import (
	"fmt"
	"github.com/igridnet/users/models"
	"strings"
)

const (
	PublishOperation   = 1
	SubscribeOperation = 2
)

// Authorize check if a certain node is allowed to make either publish or
// subscribe operation in a certain region comm channel
// Actuator are allowed to subscribe only, Sensors can only publish. Controller
// nodes can do both they can publish or subscribe to all communication channels
// of the region
func Authorize(node models.Node,topic string, operation int)(bool,error){
	var (
		region string
		nodeId string
	)

	split := strings.Split(topic, "/")
	if len(split) > 2{
		return false, fmt.Errorf("topic should be specified as -t <region-id>/<node-id> or just <region>")
	}
	nodeType := models.NodeType(node.Type)

	if nodeType == models.ControllerNode{
		if strings.TrimSpace(topic) == ""{
			return false, fmt.Errorf("can not publish to null topic")
		}

		region = split[0]
		if node.Region != region{
			return false, fmt.Errorf("nodes are not allowed to publish/subscribe outsied their region")
		}

		return true,nil
	}
	if nodeType != models.ControllerNode && len(split)<2{
		return false, fmt.Errorf("not allowed to perform any operation in this topic use format <region-id>/<node-id>")
	}

	if nodeType == models.ActuatorNode{
		if operation == PublishOperation{
			return false,fmt.Errorf("actuators are not allowed to publish, they can only subscribe")
		}

		if operation == SubscribeOperation{
			region, nodeId = split[0],split[1]
			if node.Region != region || node.UUID != nodeId{
				return false,fmt.Errorf("nodes should only subscribes in the topic(node-id) of their region")
			}
		}

		return true,nil
	}

	if nodeType == models.SensorNode && operation == SubscribeOperation{
		if operation == SubscribeOperation{
			return false,fmt.Errorf("sensor nodes are not allowed to subscribe, they can only publish")
		}

		if operation == PublishOperation{
			region, nodeId = split[0],split[1]
			if node.Region != region || node.UUID != nodeId{
				return false,fmt.Errorf("nodes should only publish in the topic(node-id) of their region")
			}
		}
		return true,nil
	}

	return true,nil
}

