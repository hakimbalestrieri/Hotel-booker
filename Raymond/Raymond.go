package Raymond

import (
	p "PRR-Labo3-Balestrieri/Protocol"
)

type Raymond struct {
	me                  int
	status              p.RaymondStatus
	parent              int
	queue               []int
	RaymondMsgBroadcast chan p.RaymondProtocol
}

//Init raymond algorithm
func Init(ray *Raymond) {
	ray.RaymondMsgBroadcast = make(chan p.RaymondProtocol)
	ray.status = p.RAY_NO
}

func (ray *Raymond) Run() {

	ray.waitForReady()

	for {
		select {
		case msg := <-ray.RaymondMsgBroadcast:
			switch msg.ReqType {
			case p.RAY_REQ:
				handleRequest(ray, msg)
			case p.RAY_WAIT:
				handleWait(ray, msg)
			case p.RAY_END:
				handleEnd(ray, msg)
			case p.RAY_REQ_SENDER:
				handleRequestSender(ray, msg)
			case p.RAY_TOKEN:
				handleToken(ray, msg)
			}
		}
	}
}

func (ray *Raymond) waitForReady() {
	for {
		if ray.parent != nil {
			break
		}
	}
}

func handleRequest(ray *Raymond, msg p.RaymondProtocol) {
	if ray.parent !== nil {
		ray.queue = append(ray.queue, msg.ServerId)
		if ray.status == p.RAY_NO {
			ray.status = p.RAY_REQ_SENDER
			ray.RaymondMsgBroadcast <- p.RaymondProtocol{
				ReqType:  p.RAY_REQ,
				ServerId: ray.me,
				ParentId: ray.parent,
			}
		}
	}
}

func handleWait(ray *Raymond, msg p.RaymondProtocol) {
	ray.status = p.RAY_SC
	ray.RaymondMsgBroadcast <- p.RaymondProtocol{
		ReqType:  p.RAY_END,
		ServerId: ray.me,
	}
}

func handleEnd(ray *Raymond, msg p.RaymondProtocol) {
	ray.status = p.RAY_NO
	if len(ray.queue) > 0 {
		ray.queue = ray.queue[1:]

		//appeler  sendToRoodNode(msg, indexDuServeur)
		if len(ray.queue) > 0 {
			ray.status = p.RAY_REQ_SENDER
			ray.RaymondMsgBroadcast <- p.RaymondProtocol{
				ReqType:  p.RAY_REQ,
				ServerId: ray.me,
				ParentId: ray.parent,
			}
		}
	}
}

func handleRequestSender(ray *Raymond, msg p.RaymondProtocol) {
	if ray.parent == nil && ray.status == p.RAY_NO {
		ray.parent = msg.ServerId
		ray.RaymondMsgBroadcast <- p.RaymondProtocol{
			ReqType:  p.RAY_TOKEN,
			ServerId: ray.parent,
		}
		ray.RaymondMsgBroadcast <- p.RaymondProtocol{
			ReqType:  p.RAY_REQ,
			ServerId: ray.me,
			ParentId: ray.parent,
		}
	} else if ray.parent != nil && ray.status == p.RAY_NO {
		ray.status = p.RAY_REQ_SENDER
		ray.RaymondMsgBroadcast <- p.RaymondProtocol{
			ReqType:  p.RAY_REQ,
			ServerId: ray.me,
			ParentId: ray.parent,
		}
	}
}

func handleToken(ray *Raymond, msg p.RaymondProtocol) {
	ray.RaymondMsgBroadcast <- p.RaymondProtocol{
		ReqType:  p.RAY_TOKEN,
		ServerId: msg.ServerId,
	}
	if ray.parent == msg.ServerId {
		ray.parent = nil
	} else if ray.queue != 0 {
		ray.status = p.RAY_REQ_SENDER
		ray.RaymondMsgBroadcast <- p.RaymondProtocol{
			ReqType:  p.RAY_REQ,
			ServerId: ray.me,
			ParentId: ray.parent,
		}
	} else {
		ray.status = p.RAY_NO
	}
}
