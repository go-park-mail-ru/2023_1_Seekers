package smtp

import (
	"fmt"
	"math/rand"
	"net"
)

func MxRecordSendMostPriority(to string, funcSend func(hostNameMostPriority string) (isSend bool, err error),
) error {
	mxRecords := &MXRecords{}
	err := mxRecords.Set(to)
	if err != nil {
		return err
	}

	return mxRecords.SendToMostPriortyRecord(funcSend)
}

//---------------------------------------------------------------------------------------------------//
// mx records

type MXRecords struct {
	Mxs []*net.MX
}

func (t *MXRecords) Init() {
	t.Mxs = make([]*net.MX, 0, 10)
}

func (t *MXRecords) Set(_to string) error {
	if t.Mxs == nil {
		t.Init()
	}

	domain, err := ParseDomain(_to)
	if err != nil {
		return err
	}
	t.Mxs, err = net.LookupMX(domain)
	if err != nil {
		return err
	}

	return nil
}

func (t *MXRecords) SendToMostPriortyRecord(fnCbSendmail func(hostMostpriorty string) (isSend bool, err error)) error {
	if len(t.Mxs) == 0 {
		return fmt.Errorf("mx records is empty")
	}

	var posStart, posEnd int
	minPref := t.Mxs[0].Pref
	lenMxs := len(t.Mxs)

	for i := 0; i < lenMxs; i++ {
		if i == lenMxs-1 || t.Mxs[i+1].Pref > minPref {
			posEnd = i + 1
			randBetweenstartend := rand.Intn(posEnd - posStart)
			randBetweenstartend += posStart

			for j := randBetweenstartend; j < posEnd; j++ {
				host := t.Mxs[j].Host
				isSend, err := fnCbSendmail(host)
				if err != nil {
					return err
				} else if isSend == true {
					return nil
				}
			}

			for j := posStart; j < randBetweenstartend; j++ {
				sHost := t.Mxs[j].Host
				isSend, err := fnCbSendmail(sHost)
				if err != nil {
					return err
				} else if isSend == true {
					return nil
				}
			}

			posStart = i + 1
			if i != lenMxs-1 {
				minPref = t.Mxs[i+1].Pref
			}
		}
	}
	return fmt.Errorf("failed send_mail | every mx records refuse to send mail")
}
