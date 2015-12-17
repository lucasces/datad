package main

// See https://godoc.org/golang.org/x/net/ipv4

import "log"
import "net"

import "golang.org/x/net/ipv4"

type MessageService struct {
	group   net.IP
	port    int
	conn    net.PacketConn
	pconn   *ipv4.PacketConn
	ifaces  []net.Interface
	ctx     *Context
	workers int
	c       chan Message
	myips   map[string]int
}

func (self *MessageService) Announce(msg Message) error {
	dst := &net.UDPAddr{IP: self.group, Port: self.port}

	data, err := msg.Encode()
	if err != nil {
		return err
	}

	for _, iface := range self.ifaces {
		err := self.pconn.SetMulticastInterface(&iface)
		if err != nil {
			return err
		}

		self.pconn.SetMulticastTTL(2)

		_, err = self.pconn.WriteTo(data, nil, dst)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *MessageService) Send(msg Message, dest string) error {
	data, err := msg.Encode()
	if err != nil {
		return err
	}
	dst := &net.UDPAddr{IP: net.ParseIP(dest), Port: self.port}

	_, err = self.pconn.WriteTo(data, nil, dst)

	return err
}

func (self *MessageService) Channel() chan Message {
	return self.c
}

func CreateMessageService(ctx *Context) (MessageService, error) {
	messageService, err := dial(ctx)
	if err != nil {
		return MessageService{}, err
	}
	ctx.MessageService = &messageService
	messageService.initProcessors()
	go messageService.listen()
	return messageService, nil
}

func filterIfaces(ifaces []net.Interface) []net.Interface {
	out := make([]net.Interface, 0)
	for _, iface := range ifaces {
		if iface.HardwareAddr == nil {
			continue
		}
		if (iface.Flags & (1 << 0)) == 0 {
			continue
		}
		if (iface.Flags & (1 << 4)) == 0 {
			continue
		}
		out = append(out, iface)
	}
	return out
}

func activeIfaces() ([]net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	return filterIfaces(ifaces), nil
}

func addressList(ifaces []net.Interface) (map[string]int, error) {
	out := make(map[string]int)
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				return nil, err
			}
			out[ip.String()] = 1
		}
	}
	return out, nil
}

func (self *MessageService) dispatchMessage(buff []byte, src net.Addr) error {
	msg, err := Decode(buff)
	if err != nil {
		return err
	}

	addr, err := net.ResolveTCPAddr("tcp4", src.String())
	if err != nil {
		return err
	}

	msg.Source = addr.IP.String()

	if self.myips[msg.Source] == 0 {
		self.c <- msg
	}

	return nil
}

func (self *MessageService) listen() {
	buff := make([]byte, 1500)
	for {
		n, cm, src, err := self.pconn.ReadFrom(buff)
		if err != nil {
			log.Fatal(err)
		}
		if cm.Dst.IsMulticast() {
			if cm.Dst.Equal(self.group) {
				self.dispatchMessage(buff[:n], src)
			} else {
				continue
			}
		} else {
			self.dispatchMessage(buff[:n], src)
		}
	}
}

func (self *MessageService) initProcessors() {
	for i := 0; i < self.workers; i++ {
		processor := NewMessageProcessor(self.ctx)
		go processor.MessageProcessorRun()
	}
}

func dial(ctx *Context) (MessageService, error) {
	group := net.ParseIP(ctx.Config.MessageService.BindGroup)
	ifaces, err := activeIfaces()
	if err != nil {
		return MessageService{}, err
	}

	myips, err := addressList(ifaces)
	if err != nil {
		return MessageService{}, err
	}

	c := make(chan Message)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", ctx.Config.MessageService.BindAddress)
	if err != nil {
		return MessageService{}, err
	}

	port := tcpAddr.Port
	conn, err := net.ListenPacket("udp4", ctx.Config.MessageService.BindAddress)
	if err != nil {
		return MessageService{}, err
	}

	pconn := ipv4.NewPacketConn(conn)
	for _, iface := range ifaces {
		if err := pconn.JoinGroup(&iface, &net.UDPAddr{IP: group}); err != nil {
			return MessageService{}, err
		}
	}

	if err := pconn.SetControlMessage(ipv4.FlagDst, true); err != nil {
		return MessageService{}, err
	}

	pconn.SetTOS(0x0)
	pconn.SetTTL(16)

	return MessageService{group, port, conn, pconn, ifaces, ctx, ctx.Config.MessageService.Workers, c, myips}, nil
}
