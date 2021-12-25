/**
File: Raymond.go
Authors: Hakim Balestrieri
Date: 25.12.2021
*/

package Raymond

import (
	conf "PRR-Labo3-Balestrieri/Config"
	p "PRR-Labo3-Balestrieri/Protocol"
	"log"
)

type Raymond struct {

	//Canaux
	RayMsgOut chan p.RaymondProtocol
	RayMsgIn  chan p.RaymondProtocol
	AccessCS  chan bool

	//Variables raymond
	CurrentId int
	Status    p.RaymondStatusType
	ParentId  int
	FifoQueue []int
}

func (ray *Raymond) Run() {

	for {
		message := <-ray.RayMsgIn

		if conf.Debug {
			log.Println("Reading RayMsgIn")
		}

		switch message.ReqType {

		case p.RAYMOND_PRO_REQ:
			ray.processReq()
		case p.RAYMOND_WAIT:
			if ray.ParentId == ray.CurrentId {
				ray.wait()
			}
		case p.RAYMOND_END:
			ray.end()
		case p.RAYMOND_REQ:
			ray.req(message.ServerId)
		case p.RAYMOND_TOKEN:
			ray.token()
		default:
		}
	}

}

func (ray *Raymond) processReq() {
	if conf.Debug {
		log.Println("Raymond processReq() begin")
	}
	if ray.ParentId != ray.CurrentId {
		ray.FifoQueue = append(ray.FifoQueue, ray.CurrentId)
		if ray.Status == p.RAY_NO {
			ray.Status = p.RAY_ASK
			sendMessage(ray, p.RAYMOND_REQ)
		}
	}
	if conf.Debug {
		log.Println("Raymond processReq() terminated")
	}
}

func (ray *Raymond) wait() {
	ray.Status = p.RAY_SC
	ray.AccessCS <- true
}

func (ray *Raymond) end() {
	if conf.Debug {
		log.Println("Raymond end() begin")
	}
	ray.Status = p.RAY_NO
	if len(ray.FifoQueue) > 0 {

		firstElement := ray.FifoQueue[0]
		ray.FifoQueue = ray.FifoQueue[:1]
		ray.ParentId = firstElement

		sendMessageWithCustomId(ray, p.RAYMOND_TOKEN, firstElement)

		if len(ray.FifoQueue) > 0 {
			ray.Status = p.RAY_ASK
			sendMessage(ray, p.RAYMOND_PRO_REQ)
		}
	}
	if conf.Debug {
		log.Println("Raymond end() terminated")
	}
}

func (ray *Raymond) req(i int) {
	if conf.Debug {
		log.Println("Raymond req(i) begin")
	}
	if ray.ParentId == conf.ServerWithoutParent && ray.Status == p.RAY_NO {
		ray.ParentId = i

		sendMessageWithCustomId(ray, p.RAYMOND_TOKEN, i)

	} else {
		ray.FifoQueue = append(ray.FifoQueue, i)
		if ray.ParentId != conf.ServerWithoutParent && ray.Status == p.RAY_NO {
			ray.Status = p.RAY_ASK
			sendMessage(ray, p.RAYMOND_PRO_REQ)
		}
	}
	if conf.Debug {
		log.Println("Raymond req(i) terminated")
	}
}

func (ray *Raymond) token() {
	if conf.Debug {
		log.Println("Raymond token() begin")
	}

	firstElement := ray.FifoQueue[0]
	ray.FifoQueue = ray.FifoQueue[:1]
	ray.ParentId = firstElement

	if firstElement == ray.CurrentId {
		ray.ParentId = conf.ServerWithoutParent
	} else {
		ray.RayMsgOut <- p.RaymondProtocol{
			ReqType:  p.RAYMOND_TOKEN,
			ServerId: firstElement,
		}
		if len(ray.FifoQueue) > 0 {
			ray.Status = p.RAY_ASK
			sendMessage(ray, p.RAYMOND_REQ)
		} else {
			ray.Status = p.RAY_NO
		}
	}
	if conf.Debug {
		log.Println("Raymond token() terminated")
	}
}

func sendMessage(ray *Raymond, req p.RaymondType) {
	ray.RayMsgOut <- p.RaymondProtocol{
		ReqType:  req,
		ServerId: ray.CurrentId,
	}
}

func sendMessageWithCustomId(ray *Raymond, req p.RaymondType, id int) {
	ray.RayMsgOut <- p.RaymondProtocol{
		ReqType:  req,
		ServerId: id,
	}
}
