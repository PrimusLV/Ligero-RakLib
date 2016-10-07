package session

import (
	"raklib/utils"
	"raklib/protocol"
)

const STATE_UNCONNECTED = 0;
const STATE_CONNECTING_1 = 1;
const STATE_CONNECTING_2 = 2;
const STATE_CONNECTED = 3;

const MAX_SPLIT_SIZE = 128;
const MAX_SPLIT_COUNT = 4;

const WINDOW_SIZE = 2048;

func NewSession(sm SessionManager, addr string, p int) *Session {
	s  := Session{
		messageIndex: 0,
		channelIndex: new(map[int]int), //TODO
		sessionManager: sm,
		address: addr,
		port: p,
		mtuSize: 548,
		id: 0,
		splitID: 0,
		sendSeqNumber: 0,
		lastSeqNumber: -1,
		lastUpdate: 0,
		startTime: 0,
		isTemporal: true,
		packetToSend: make([int]protocol.DataPacket, 16*16),
		isActive: false,
		ACKQueue: make([]int),
		NACKQueue: make([]int),
		recoveryQueue: make([]protocol.DataPacket),
		splitPackets: make([]protocol.DataPacket),
		needACK: make([][]int),
		sendQueue: new(protocol.DATA_PACKET_4),
		windowStart: -1,
		receivedWindow: make(map[int]bool),
		windowEnd: WINDOW_SIZE,
		reliableWindowStart: 0,
		reliableWindowEnd: WINDOW_SIZE,
		reliableWindow: make(map[int]bool),
		lastReliableIndex: -1,	
	}
	for i := 0; i < 32; i++ { // TODO
		s.channelIndex[i] = 0
	}
	return &s
}

type Session struct {
	messageIndex int;
	channelIndex map[int]int;
    sessionManager SessionManager;
    address string;
    port int
    state int
    mtuSize int
    id int
    splitID int
	sendSeqNumber int
    lastSeqNumber int
    lastUpdate int64;
    startTime int64;
	isTemporal bool
    packetToSend map[int]protocol.DataPacket
    isActive bool
    ACKQueue []int
    NACKQueue []int
    recoveryQueue []protocol.DataPacket
	splitPackets []protocol.DataPacket;
	needACK [][]int;
	sendQueue protocol.DataPacket
    windowStart int
    receivedWindow map[int]bool;
    windowEnd int
	reliableWindowStart int
	reliableWindowEnd int
	reliableWindow map[int]bool;
	lastReliableIndex int;
}

func (this *Session) GetAddress() string {
    return this.address;
}

func (this *Session) GetPort() int {
    return this.port;
}

func (this *Session) GetID() int {
    return this.id;
}

func (this *Session) Update(time int64){
    if !this.isActive && (this.lastUpdate + 10) < time {
        this.disconnect("timeout");

        return;
    }
    this.isActive = false;

    if len(this.ACKQueue) > 0 {
        pk := new(ACK);
        pk.packets = this.ACKQueue;
        this.SendPacket(pk);
        this.ACKQueue = make([]int);
    }

    if len(this.NACKQueue) > 0 {
        pk = new(NACK);
        pk.packets = this.NACKQueue;
        this.SendPacket(pk);
        this.NACKQueue = make([]int);
    }

    if len(this.packetToSend) > 0 {
		limit := 16;
        for k, pk := range this.packetToSend {
            pk.sendTime = time;
            pk.encode();
            this.recoveryQueue[pk.seqNumber] = pk;
            delete(this.packetToSend, k);
            this.sendPacket(pk);

			if limit <= 0 {
				break;
			}
			limit--
        }

		if len(this.packetToSend) > WINDOW_SIZE {
			this.packetToSend = make(map[int]protoool.DataPacket)
		}
    }

    if len(this.needACK) > 0 {
        for identifierACK, indexes := range this.needACK {
            if len(indexes) == 0 {
                delete(this.needACK, identifierACK);
                this.sessionManager.NotifyACK(this, identifierACK);
            }
        }
    }


	for seq, pk := range this.recoveryQueue {
		if pk.sendTime < (time.Now().Unix() - 8) {
			append(this.packetToSend, pk) // TODO
			delete(this.recoveryQueue, seq)
		} else {
			break;
		}
	}

	for seq, bool := range this.receivedWindow {
		if seq < this.windowStart {
			delete(this.receivedWindow, seq);
		}else{
			break;
		}
	}

    this.sendQueue();
}

func (this *Session) Disconnect(reason string){
    this.sessionManager.RemoveSession(this, reason);
}

func (this *Session) SendPacket(packet protocol.DataPacket){
    this.sessionManager.SendPacket(packet, this.address, this.port)
}

func (this *Session) SendQueue(){
    if len(this.sendQueue.Packets) > 0 {
    	this.sendSeqNumber += 1
        this.sendQueue.SeqNumber = this.sendQueueNumber;
		this.SendPacket(this.sendQueue);
        this.sendQueue.SendTime = time.Now().Unix() / time.Microsecond
        this.recoveryQueue[this.sendQueue->seqNumber] = this.sendQueue;
        this.sendQueue = new DATA_PACKET_4();
    }
}

/**
 * @param EncapsulatedPacket $pk
 * @param int                $flags
 */
private function addToQueue(EncapsulatedPacket $pk, $flags = RakLib::PRIORITY_NORMAL){
    $priority = $flags & 0b0000111;
    if($pk->needACK and $pk->messageIndex !== null){
        this.needACK[$pk->identifierACK][$pk->messageIndex] = $pk->messageIndex;
    }
    if($priority === RakLib::PRIORITY_IMMEDIATE){ //Skip queues
        $packet = new DATA_PACKET_0();
        $packet->seqNumber = this.sendSeqNumber++;
        if($pk->needACK){
	        $packet->packets[] = clone $pk;
	        $pk->needACK = false;
        }else{
	        $packet->packets[] = $pk->toBinary();
        }

        this.sendPacket($packet);
        $packet->sendTime = microtime(true);
        this.recoveryQueue[$packet->seqNumber] = $packet;

        return;
    }
    $length = this.sendQueue->length();
    if($length + $pk->getTotalLength() > this.mtuSize){
        this.sendQueue();
    }

    if($pk->needACK){
	    this.sendQueue->packets[] = clone $pk;
	    $pk->needACK = false;
    }else{
	    this.sendQueue->packets[] = $pk->toBinary();
    }
}

/**
 * @param EncapsulatedPacket $packet
 * @param int                $flags
 */
public function addEncapsulatedToQueue(EncapsulatedPacket $packet, $flags = RakLib::PRIORITY_NORMAL){

    if(($packet->needACK = ($flags & RakLib::FLAG_NEED_ACK) > 0) === true){
        this.needACK[$packet->identifierACK] = [];
    }

	if(
		$packet->reliability === 2 or
		$packet->reliability === 3 or
		$packet->reliability === 4 or
		$packet->reliability === 6 or
		$packet->reliability === 7
	){
		$packet->messageIndex = this.messageIndex++;

		if($packet->reliability === 3){
			$packet->orderIndex = this.channelIndex[$packet->orderChannel]++;
		}
	}

    if($packet->getTotalLength() + 4 > this.mtuSize){
        $buffers = str_split($packet->buffer, this.mtuSize - 34);
        $splitID = ++this.splitID % 65536;
        foreach($buffers as $count => $buffer){
            $pk = new EncapsulatedPacket();
            $pk->splitID = $splitID;
            $pk->hasSplit = true;
            $pk->splitCount = count($buffers);
            $pk->reliability = $packet->reliability;
            $pk->splitIndex = $count;
            $pk->buffer = $buffer;
			if($count > 0){
				$pk->messageIndex = this.messageIndex++;
			}else{
				$pk->messageIndex = $packet->messageIndex;
			}
			if($pk->reliability === 3){
				$pk->orderChannel = $packet->orderChannel;
				$pk->orderIndex = $packet->orderIndex;
			}
            this.addToQueue($pk, $flags | RakLib::PRIORITY_IMMEDIATE);
        }
    }else{
        this.addToQueue($packet, $flags);
    }
}

private function handleSplit(EncapsulatedPacket $packet){
	if($packet->splitCount >= self::MAX_SPLIT_SIZE or $packet->splitIndex >= self::MAX_SPLIT_SIZE or $packet->splitIndex < 0){
		return;
	}


	if(!isset(this.splitPackets[$packet->splitID])){
		if(count(this.splitPackets) >= self::MAX_SPLIT_COUNT){
			return;
		}
		this.splitPackets[$packet->splitID] = [$packet->splitIndex => $packet];
	}else{
		this.splitPackets[$packet->splitID][$packet->splitIndex] = $packet;
	}

	if(count(this.splitPackets[$packet->splitID]) === $packet->splitCount){
		$pk = new EncapsulatedPacket();
		$pk->buffer = "";
		for($i = 0; $i < $packet->splitCount; ++$i){
			$pk->buffer .= this.splitPackets[$packet->splitID][$i]->buffer;
		}

		$pk->length = strlen($pk->buffer);
		unset(this.splitPackets[$packet->splitID]);

		this.handleEncapsulatedPacketRoute($pk);
	}
}

private function handleEncapsulatedPacket(EncapsulatedPacket $packet){
	if($packet->messageIndex === null){
		this.handleEncapsulatedPacketRoute($packet);
	}else{
		if($packet->messageIndex < this.reliableWindowStart or $packet->messageIndex > this.reliableWindowEnd){
			return;
		}

		if(($packet->messageIndex - this.lastReliableIndex) === 1){
			this.lastReliableIndex++;
			this.reliableWindowStart++;
			this.reliableWindowEnd++;
			this.handleEncapsulatedPacketRoute($packet);

			if(count(this.reliableWindow) > 0){
				ksort(this.reliableWindow);

				foreach(this.reliableWindow as $index => $pk){
					if(($index - this.lastReliableIndex) !== 1){
						break;
					}
					this.lastReliableIndex++;
					this.reliableWindowStart++;
					this.reliableWindowEnd++;
					this.handleEncapsulatedPacketRoute($pk);
					unset(this.reliableWindow[$index]);
				}
			}
		}else{
			this.reliableWindow[$packet->messageIndex] = $packet;
		}
	}

}

public function getState(){
	return this.state;
}

public function isTemporal(){
	return this.isTemporal;
}

private function handleEncapsulatedPacketRoute(EncapsulatedPacket $packet){
    if(this.sessionManager === null){
        return;
    }

	if($packet->hasSplit){
		if(this.state === self::STATE_CONNECTED){
			this.handleSplit($packet);
		}
		return;
	}

	$id = ord($packet->buffer{0});
	if($id < 0x80){ //internal data packet
		if(this.state === self::STATE_CONNECTING_2){
			if($id === CLIENT_CONNECT_DataPacket::$ID){
				$dataPacket = new CLIENT_CONNECT_DataPacket;
				$dataPacket->buffer = $packet->buffer;
				$dataPacket->decode();
				$pk = new SERVER_HANDSHAKE_DataPacket;
				$pk->address = this.address;
				$pk->port = this.port;
				$pk->sendPing = $dataPacket->sendPing;
				$pk->sendPong = bcadd($pk->sendPing, "1000");
				$pk->encode();

				$sendPacket = new EncapsulatedPacket();
				$sendPacket->reliability = 0;
				$sendPacket->buffer = $pk->buffer;
				this.addToQueue($sendPacket, RakLib::PRIORITY_IMMEDIATE);
			}elseif($id === CLIENT_HANDSHAKE_DataPacket::$ID){
				$dataPacket = new CLIENT_HANDSHAKE_DataPacket;
				$dataPacket->buffer = $packet->buffer;
				$dataPacket->decode();

				if($dataPacket->port === this.sessionManager->getPort() or !this.sessionManager->portChecking){
					this.state = self::STATE_CONNECTED; //FINALLY!
					this.isTemporal = false;
					this.sessionManager->openSession($this);
				}
			}
		}elseif($id === CLIENT_DISCONNECT_DataPacket::$ID){
			this.disconnect("client disconnect");
		}elseif($id === PING_DataPacket::$ID){
			$dataPacket = new PING_DataPacket;
			$dataPacket->buffer = $packet->buffer;
			$dataPacket->decode();

			$pk = new PONG_DataPacket;
			$pk->pingID = $dataPacket->pingID;
			$pk->encode();

			$sendPacket = new EncapsulatedPacket();
			$sendPacket->reliability = 0;
			$sendPacket->buffer = $pk->buffer;
			this.addToQueue($sendPacket);
		}//TODO: add PING/PONG (0x00/0x03) automatic latency measure
	}elseif(this.state === self::STATE_CONNECTED){
		this.sessionManager->streamEncapsulated($this, $packet);

		//TODO: stream channels
	}else{
		//this.sessionManager->getLogger()->notice("Received packet before connection: " . bin2hex($packet->buffer));
	}
}

public function handlePacket(Packet $packet){
    this.isActive = true;
    this.lastUpdate = microtime(true);
    if(this.state === self::STATE_CONNECTED or this.state === self::STATE_CONNECTING_2){
        if($packet::$ID >= 0x80 and $packet::$ID <= 0x8f and $packet instanceof DataPacket){ //Data packet
            $packet->decode();

			if($packet->seqNumber < this.windowStart or $packet->seqNumber > this.windowEnd or isset(this.receivedWindow[$packet->seqNumber])){
				return;
			}

			$diff = $packet->seqNumber - this.lastSeqNumber;

			unset(this.NACKQueue[$packet->seqNumber]);
			this.ACKQueue[$packet->seqNumber] = $packet->seqNumber;
			this.receivedWindow[$packet->seqNumber] = $packet->seqNumber;

			if($diff !== 1){
				for($i = this.lastSeqNumber + 1; $i < $packet->seqNumber; ++$i){
					if(!isset(this.receivedWindow[$i])){
						this.NACKQueue[$i] = $i;
					}
				}
			}

			if($diff >= 1){
				this.lastSeqNumber = $packet->seqNumber;
				this.windowStart += $diff;
				this.windowEnd += $diff;
			}

			foreach($packet->packets as $pk){
				this.handleEncapsulatedPacket($pk);
			}
		}else{
            if($packet instanceof ACK){
                $packet->decode();
                foreach($packet->packets as $seq){
                    if(isset(this.recoveryQueue[$seq])){
                        foreach(this.recoveryQueue[$seq]->packets as $pk){
                            if($pk instanceof EncapsulatedPacket and $pk->needACK and $pk->messageIndex !== null){
                                unset(this.needACK[$pk->identifierACK][$pk->messageIndex]);
                            }
                        }
                        unset(this.recoveryQueue[$seq]);
                    }
                }
            }elseif($packet instanceof NACK){
                $packet->decode();
                foreach($packet->packets as $seq){
                    if(isset(this.recoveryQueue[$seq])){
						$pk = this.recoveryQueue[$seq];
						$pk->seqNumber = this.sendSeqNumber++;
                        this.packetToSend[] = $pk;
						unset(this.recoveryQueue[$seq]);
                    }
                }
            }
        }

    }elseif($packet::$ID > 0x00 and $packet::$ID < 0x80){ //Not Data packet :)
        $packet->decode();
        if($packet instanceof OPEN_CONNECTION_REQUEST_1){
            $packet->protocol; //TODO: check protocol number and refuse connections
            $pk = new OPEN_CONNECTION_REPLY_1();
            $pk->mtuSize = $packet->mtuSize;
            $pk->serverID = this.sessionManager->getID();
            this.sendPacket($pk);
            this.state = self::STATE_CONNECTING_1;
        }elseif(this.state === self::STATE_CONNECTING_1 and $packet instanceof OPEN_CONNECTION_REQUEST_2){
            this.id = $packet->clientID;
            if($packet->serverPort === this.sessionManager->getPort() or !this.sessionManager->portChecking){
                this.mtuSize = min(abs($packet->mtuSize), 1464); //Max size, do not allow creating large buffers to fill server memory
                $pk = new OPEN_CONNECTION_REPLY_2();
                $pk->mtuSize = this.mtuSize;
                $pk->serverID = this.sessionManager->getID();
				$pk->clientAddress = this.address;
                $pk->clientPort = this.port;
                this.sendPacket($pk);
                this.state = self::STATE_CONNECTING_2;
            }
        }
    }
}

public function close(){
	$data = "\x00\x00\x08\x15";
    this.addEncapsulatedToQueue(EncapsulatedPacket::fromBinary($data), RakLib::PRIORITY_IMMEDIATE); //CLIENT_DISCONNECT packet 0x15
    this.sessionManager = null;
}