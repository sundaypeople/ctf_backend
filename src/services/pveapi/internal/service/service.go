package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/LainInTheWired/ctf-backend/pveapi/model"
	"github.com/LainInTheWired/ctf-backend/pveapi/repository"
	"github.com/cockroachdb/errors"
)

type pveService struct {
	pveRepo repository.PVERepository
}

type PVEService interface {
	CreateCloudinitVM(name string, size int, vmconf *model.VMEdit, cloneid int) (int, error)
	DeleteVMByVmid(vmid int) error
	SelectNode(cores int, memory int, disk int) (string, error)
	GenerateCloudinit(hostname string, conf []model.User, filename string, sshPwauth int) error
	TransferFileViaSCP(fname string) error
	Template(vmid int) error
	DeleteCloudinitFile(fname string) error
	GetIps(vmid int) (map[string][]string, error)
	GetClusterResource() ([]model.ClusterResources, error)
	EditVMACL() error
}

func NewPVEService(r repository.PVERepository) PVEService {
	return &pveService{
		pveRepo: r,
	}
}

// SelectLeastLoadedNode は最も負荷が低いノードを選択します
func (p *pveService) SelectNodeByCPU() (string, error) {
	nodes, err := p.pveRepo.GetNodeList()
	if err != nil {
		return "", errors.Wrap(err, "can't get nodes")
	}
	minLoad := 1.0 // CPU負荷は通常0.0〜1.0の範囲
	var selectedNode string
	for _, node := range nodes {
		// オンラインのノードのみ対象
		if node.Status != "online" {
			continue
		}

		// CPU負荷が最小のノードを選択
		if node.CPU < minLoad {
			minLoad = node.CPU
			selectedNode = node.Node
		}
	}

	if selectedNode == "" {
		return "", errors.New("not found online node")
	}
	return selectedNode, nil
}

func (p *pveService) SelectNode(cores int, memory int, disk int) (string, error) {
	nodes, err := p.pveRepo.GetNodeList()
	if err != nil {
		return "", errors.Wrap(err, "can't get nodes")
	}
	minLoad := 1.0 // CPU負荷は通常0.0〜1.0の範囲
	var selectedNode string
	for _, node := range nodes {
		// オンラインのノードのみ対象
		if node.Status != "online" {
			continue
		}

		// CPU負荷が最小のノードを選択
		if int(node.Maxmem) < memory*1048576 {
			continue
		}
		if int(node.Maxcpu) < cores {
			continue
		}
		if int64(node.Maxdisk) < int64(disk*1073741824) {
			continue
		}
		if node.CPU < minLoad {
			minLoad = node.CPU
			selectedNode = node.Node
		}

	}

	if selectedNode == "" {
		return "", errors.New("not found online node")
	}
	return selectedNode, nil

}

func (p *pveService) CreateCloudinitVM(name string, size int, vmconf *model.VMEdit, cloneid int) (int, error) {
	svmid, err := p.pveRepo.NextVMID()
	if err != nil {
		return 0, errors.Wrap(err, "can't get next VID")
	}
	vmid, err := strconv.Atoi(svmid)
	if err != nil {
		return 0, errors.Wrap(err, "can't cast vmid")
	}
	vmconf.Scsi = []string{fmt.Sprintf("vmdisk:vm-%d-disk-0,size=16G", vmid)}

	vmconf.Vmid = vmid
	cnode, err := p.SearchNodeByVmid(cloneid)
	if err != nil {
		return 0, errors.Wrap(err, "can't search vm")
	}

	err = p.pveRepo.CloneVM(name, vmid, cnode, cloneid, vmconf.Node)
	if err != nil {
		return 0, errors.Wrap(err, "can't clone vm")
	}

	// 追加 ACL
	if err := p.pveRepo.EditVMACL(vmid); err != nil {
		return 0, errors.Wrap(err, "can't edit acl")
	}

	// err = p.EditVM(*vmconf)
	err = p.fiveEditVM(vmconf)
	if err != nil {
		p.fiveDeleteVM(&model.VMDelete{Vmid: vmid, Node: vmconf.Node})
		return 0, errors.Wrap(err, "can't edit vm")
	}

	if size != 0 {
		err = p.pveRepo.ResizeDisk(vmconf.Node, "scsi0", size, vmid)
		if err != nil {
			p.fiveDeleteVM(&model.VMDelete{Vmid: vmid, Node: vmconf.Node})
			return 0, errors.Wrap(err, "can't resize vm disk")
		}
	}

	err = p.pveRepo.Boot(vmconf.Node, vmid)
	if err != nil {
		p.fiveDeleteVM(&model.VMDelete{Vmid: vmid, Node: vmconf.Node})
		return 0, errors.Wrap(err, "can't boot")
	}
	return vmid, nil
}

func (p *pveService) fiveEditVM(conf *model.VMEdit) error {
	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Second)
		if err := p.pveRepo.EditVM(*conf); err != nil {
			fmt.Println(i)
			if i > 3 {
				return errors.Wrap(err, "can't error")
			}
		} else {
			break
		}
	}
	return nil
}
func (p *pveService) fiveTemplateVM(node string, vmid int) error {
	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Second)
		if err := p.pveRepo.Template(node, vmid); err != nil {
			fmt.Println(i)
			if i > 3 {
				return errors.Wrap(err, "can't error")
			}
		} else {
			break
		}
	}
	return nil
}

func (p *pveService) fiveDeleteVM(conf *model.VMDelete) error {
	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Second)
		if err := p.pveRepo.DeleteVM(conf); err != nil {
			if i > 3 {
				return errors.Wrap(err, "can't error")
			}
		} else {
			break
		}
	}
	return nil
}

func (p *pveService) SearchNodeByVmid(vmid int) (string, error) {
	res, err := p.pveRepo.GetClusterResourcesList()
	if err != nil {
		return "", err
	}
	for _, v := range res {
		if vmid == v.Vmid {
			return v.Node, nil
		}
	}
	return "", errors.New("not found this vmid in cluster")
}
func (p *pveService) DeleteVMByVmid(vmid int) error {
	n, err := p.SearchNodeByVmid(vmid)
	fmt.Println("search error", err)
	if err != nil {
		return errors.Wrap(err, "can't search node")
	}
	conf := &model.VMDelete{
		Vmid: vmid,
		Node: n,
	}
	fmt.Println("vmid", conf.Vmid, "node", n)
	if err := p.pveRepo.Shutdown(n, vmid); err != nil {
		return err
	}

	time.Sleep(30 * time.Second)

	if err := p.fiveDeleteVM(conf); err != nil {
		return err
	}
	return nil

}

func (p *pveService) GenerateCloudinit(hostname string, conf []model.User, filename string, sshPwauth int) error {
	if err := p.pveRepo.CloudinitGenerator(filename, hostname, hostname, sshPwauth, conf); err != nil {
		return errors.Wrap(err, "can't create cloudinit")
	}
	return nil
}

func (p *pveService) TransferFileViaSCP(fname string) error {
	if err := p.pveRepo.TransferFileViaSCP(fname); err != nil {
		return errors.Wrap(err, "can't send file")
	}
	return nil
}

func (p *pveService) DeleteCloudinitFile(fname string) error {
	if err := p.pveRepo.DeleteFile(fname); err != nil {
		return errors.Wrap(err, "can't to template")
	}
	return nil
}

func (p *pveService) Template(vmid int) error {
	fmt.Println("template: ", vmid)
	node, err := p.SearchNodeByVmid(vmid)

	if err != nil {
		return errors.Wrap(err, "can't found err")
	}
	if err := p.pveRepo.Shutdown(node, vmid); err != nil {
		return errors.Wrap(err, "can't stop vm")
	}

	if err := p.fiveTemplateVM(node, vmid); err != nil {
		return errors.Wrap(err, "can't to template")
	}
	return nil
}

func (p *pveService) GetIps(vmid int) (map[string][]string, error) {
	ips := map[string][]string{}
	node, err := p.SearchNodeByVmid(vmid)
	if err != nil {
		return nil, errors.Wrap(err, "can't found err")
	}
	GetNetIntFormQumeAgent, err := p.pveRepo.GetNetIntFormQumeAgent(node, vmid)
	if err != nil {
		return nil, errors.Wrap(err, "can't get int")
	}
	fmt.Printf("ips:   %+v", GetNetIntFormQumeAgent)
	for _, results := range GetNetIntFormQumeAgent {
		for _, reips := range results.IPAddresses {
			ips[results.Name] = append(ips[results.Name], reips.IPAddress)
		}
	}
	return ips, nil
}

func (p *pveService) GetClusterResource() ([]model.ClusterResources, error) {
	res, err := p.pveRepo.GetClusterResourcesList()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *pveService) EditVMACL() error {
	err := p.pveRepo.EditVMACL(137)
	if err != nil {
		return err
	}
	return nil
}
