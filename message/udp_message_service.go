package message

// See https://godoc.org/golang.org/x/net/ipv4

import "log"
import "net"

import "golang.org/x/net/ipv4"

import "datad/context"
import "datad/defs"

type UDPMessageService struct {
	group   net.IP
	port    int
	conn    net.PacketConn
	pconn   *ipv4.PacketConn
	ifaces  []net.Interface
	ctx     context.Context
	workers int
	c       chan defs.Message
	myips   map[string]int
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

func Dial(ctx context.Context) (UDPMessageService, error) {
	group := net.ParseIP(ctx.Config.MessageService.BindGroup)
	ifaces, err := activeIfaces()
	if err != nil {
		return UDPMessageService{}, err
	}

	myips, err := addressList(ifaces)
	if err != nil {
		return UDPMessageService{}, err
	}

	c := make(chan defs.Message)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", ctx.Config.MessageService.BindAddress)
	if err != nil {
		return UDPMessageService{}, err
	}

	port := tcpAddr.Port
	conn, err := net.ListenPacket("udp4", ctx.Config.MessageService.BindAddress)
	if err != nil {
		return UDPMessageService{}, err
	}

	pconn := ipv4.NewPacketConn(conn)
	for _, iface := range ifaces {
		if err := pconn.JoinGroup(&iface, &net.UDPAddr{IP: group}); err != nil {
			return UDPMessageService{}, err
		}
	}

	if err := pconn.SetControlMessage(ipv4.FlagDst, true); err != nil {
		return UDPMessageService{}, err
	}

	pconn.SetTOS(0x0)
	pconn.SetTTL(16)

	return UDPMessageService{group, port, conn, pconn, ifaces, ctx, ctx.Config.MessageService.Workers, c, myips}, nil
}

func dispatchMessage(handler UDPMessageService, buff []byte, src net.Addr) error {
	msg, err := Decode(buff)
	if err != nil {
		return err
	}

	addr, err := net.ResolveTCPAddr("tcp4", src.String())
	if err != nil {
		return err
	}

	msg.Source = addr.IP.String()

	if handler.myips[msg.Source] == 0 {
		handler.c <- msg
	}

	return nil
}

func Listen(handler UDPMessageService) {
	buff := make([]byte, 1500)
	for {
		n, cm, src, err := handler.pconn.ReadFrom(buff)
		if err != nil {
			log.Fatal(err)
		}
		if cm.Dst.IsMulticast() {
			if cm.Dst.Equal(handler.group) {
				dispatchMessage(handler, buff[:n], src)
			} else {
				continue
			}
		} else {
			dispatchMessage(handler, buff[:n], src)
		}
	}
}

func (self UDPMessageService) Announce(msg defs.Message) error {
	dst := &net.UDPAddr{IP: self.group, Port: self.port}

	data, err := Encode(msg)
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

func (self UDPMessageService) Send(msg defs.Message, dest string) error {
	data, err := Encode(msg)
	if err != nil {
		return err
	}
	dst := &net.UDPAddr{IP: net.ParseIP(dest), Port: self.port}

	_, err = self.pconn.WriteTo(data, nil, dst)

	return err
}

func (self UDPMessageService) Channel() chan defs.Message {
	return self.c
}

func (self UDPMessageService) initProcessors() {
	for i := 0; i < self.workers; i++ {
		processor := NewMessageProcessor(self.ctx)
		self.ctx.WaitGroup.Add(1)
		go processor.MessageProcessorRun()
	}
}

func NewUDPMessageService(ctx context.Context) (UDPMessageService, error) {
	discoHandler, err := Dial(ctx)
	if err != nil {
		return UDPMessageService{}, err
	}
	discoHandler.ctx.MessageService = discoHandler
	discoHandler.initProcessors()
	ctx.WaitGroup.Add(1)
	go Listen(discoHandler)
	return discoHandler, nil
}
