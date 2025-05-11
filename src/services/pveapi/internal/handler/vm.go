package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/LainInTheWired/ctf-backend/pveapi/model"
	"github.com/LainInTheWired/ctf-backend/pveapi/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/xerrors"
)

type PVEParam struct {
	name string
	ip   string
}

// 依存関係用の構造体
type PVEHandler struct {
	serv service.PVEService
}

//	type createVMRequest struct {
//		Cloneid  int    `json:"cloneid" validate:"required,number"`
//		Name     string `json"name,omitempty" validate:"required"`
//		Memory   int    `json"memory,omitempty" validate:"required" `
//		CPUs     int    `json"cpu,omitempty" validate:"required"`
//		Disk     int    `json"disk,omitempty" validate:"required"`
//		IP       string `json:"ip" validate:"cidr"`
//		Gateway  string `json:"gateway" validate:"ip"`
//		OS       string `json"os" validate:""`
//		Cicustom string `json:"cicustom,omitempty" validate:"required"`
//	}

type createVMRequest struct {
	Cloneid  int    `json:"cloneid" validate:"required,numeric"`
	Name     string `json:"name" validate:"required"`
	Memory   int    `json:"memory,omitempty"`
	CPUs     int    `json:"cpu,omitempty"`
	Disk     int    `json:"disk,omitempty"`
	IP       string `json:"ip" validate:"omitempty,cidr"`
	Gateway  string `json:"gateway" validate:"omitempty,ip"`
	OS       string `json:"os,omitempty"`
	Cicustom string `json:"cicustom,omitempty" validate:"required"`
}
type CloneQuestionsRequest struct {
	QID   int    `json"qid" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Teams []TeamRequest
}
type TeamRequest struct {
	ID      int      `json:"id"`
	Sshkeys []string `json:"sshkeys"`
}
type deleteVMRequest struct {
	ID int `json"id" validate:"required,number"`
}

type CloudinitRequest struct {
	Filename  string   `json:"filename" validation:"required"`
	Hostname  string   `json:"hostname" validation:"required"`
	Sshkeys   []string `json:"sshkeys"`
	SshPwauth string   `json:"ssh_pwauth"`
	Username  string   `json:"username"`
	Passwd    string   `json:"passwd"`
}

type TemplateRequest struct {
	ID int ` json:"id"  validation:"required,number"`
}

type SuccessResponse struct {
	Data string `json:"data"`
}
type DeleteCloudinit struct {
	Filename string `json:"filename"`
}

// // VMConfig は作成するVMの設定を保持します
// type VMConfig struct {
// 	VMID           string `json:"vmid"`
// 	Name           string `json:"name"`
// 	Memory         string `json:"memory"` // MB単位
// 	CPUs           string `json:"cores"`
// 	Net0           string `json:"net0"`        // 例: "virtio=DE:AD:BE:EF:00:00,bridge=vmbr0"
// 	Scsi0          string `json:"scsi0"`       // 例: "kingston_1tb:vm-200-disk-0,size=16G"
// 	Boot           string `json:"boot"`        // 例: "c"
// 	Ide2           string `json:"ide2"`        // 例: "local:iso/AlmaLinux-9.3-x86_64-boot.iso"
// 	OSType         string `json:"ostype"`      // 例: "l26" (Linux 2.6/3.x/4.x)
// 	SCSIController string `json:"scsihw"`      // 例: "virtio-scsi-single"
// 	Description    string `json:"description"` // VMの説明（オプション）
// }

func NewPVEAPI(s service.PVEService) *PVEHandler {
	return &PVEHandler{
		serv: s,
	}
}

// func (h *PVEHandler) GetVM(c echo.Context) error {
// 	_, err := h.serv.zGetVM("200")
// 	if err != nil {
// 		wrappedErr := xerrors.Errorf(": %w", err)
// 		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 	}
// 	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "create vm success"))
// }

func (h *PVEHandler) CreateCloudinitVM(c echo.Context) error {
	// リクエストから構造体にデータをコピー
	var req createVMRequest
	if err := c.Bind(&req); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	// データをバリデーションにかける
	if err := c.Validate(req); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	node, err := h.serv.SelectNode(req.CPUs, req.Memory, req.Disk)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	conf := &model.VMEdit{}
	if req.IP == "" {
		conf = &model.VMEdit{
			Memory: req.Memory,
			Cores:  req.CPUs,
			Node:   node,
			// ついか
			Ipconfig: []string{"ip=dhcp"},
			Cicustom: req.Cicustom,
		}
	} else {
		conf = &model.VMEdit{
			Memory:   req.Memory,
			Cores:    req.CPUs,
			Node:     node,
			Ipconfig: []string{fmt.Sprintf("ip=%s,gw=%s", req.IP, req.Gateway)},
			Cicustom: req.Cicustom,
		}
	}

	fmt.Printf("%+v", conf)

	vmid, err := h.serv.CreateCloudinitVM(req.Name, req.Disk, conf, req.Cloneid)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})

	}
	svmid := strconv.Itoa(vmid)
	resq := &SuccessResponse{
		Data: svmid,
	}
	return c.JSON(http.StatusOK, resq)
}

func (h *PVEHandler) DeleteVM(c echo.Context) error {
	// リクエストから構造体にデータをコピー
	var req deleteVMRequest
	if err := c.Bind(&req); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	// データをバリデーションにかける
	if err := c.Validate(req); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	if err := h.serv.DeleteVMByVmid(req.ID); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	resq := &SuccessResponse{
		Data: "success delete vm",
	}
	return c.JSON(http.StatusOK, resq)
}

func (h *PVEHandler) Cloudinit(c echo.Context) error {
	// リクエストから構造体にデータをコピー
	var req CloudinitRequest
	if err := c.Bind(&req); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	// データをバリデーションにかける
	if err := c.Validate(req); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	users := []model.User{{
		Name:              req.Username,
		Sudo:              "ALL=(ALL) NOPASSWD:ALL",
		Shell:             "/bin/bash",
		PlainTextPasswd:   req.Passwd,
		SshPwauth:         req.SshPwauth,
		SshAuthorizedKeys: req.Sshkeys,
	},
	}
	if err := h.serv.GenerateCloudinit(req.Hostname, users, req.Filename, 1); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	if err := h.serv.TransferFileViaSCP(req.Filename); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	resq := &SuccessResponse{
		Data: "success cloudinit",
	}
	return c.JSON(http.StatusOK, resq)
}
func (h *PVEHandler) DeleteCloudinit(c echo.Context) error {
	// リクエストから構造体にデータをコピー
	var req DeleteCloudinit
	if err := c.Bind(&req); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	// データをバリデーションにかける
	if err := c.Validate(req); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	if err := h.serv.DeleteCloudinitFile(req.Filename); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	resq := &SuccessResponse{
		Data: "success delete cloudinitfie",
	}
	return c.JSON(http.StatusOK, resq)
}

func (h *PVEHandler) ToTemplate(c echo.Context) error {
	// リクエストから構造体にデータをコピー
	var req TemplateRequest
	if err := c.Bind(&req); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	// データをバリデーションにかける
	if err := c.Validate(req); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	if err := h.serv.Template(req.ID); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	resq := &SuccessResponse{
		Data: "success to template",
	}
	return c.JSON(http.StatusOK, resq)
}

func (h *PVEHandler) GetIps(c echo.Context) error {
	// リクエストから構造体にデータをコピー
	svid := c.Param("vmid")
	vid, err := strconv.Atoi(svid)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	ips, err := h.serv.GetIps(vid)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	// resq := &SuccessResponse{
	// 	Data: json.Marshal(ips),
	// }
	return c.JSON(http.StatusOK, ips)
}
func (h *PVEHandler) GetClusterResource(c echo.Context) error {
	cluster, err := h.serv.GetClusterResource()
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusOK, cluster)
}

func (h *PVEHandler) EditVMACL(c echo.Context) error {
	err := h.serv.EditVMACL()
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusOK, "scesss")
}

// func (h *PVEHandler) GETTestHander(c echo.Context) error {
// 	err := h.serv.EditVM(service.VMEdit{})
// 	if err != nil {
// 		wrappedErr := xerrors.Errorf(": %w", err)
// 		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 	}

// 	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "test success"))
// }

// func (h *PVEHandler) GetNodeTestHander(c echo.Context) error {
// 	nodes, err := h.serv.GetNodeList()
// 	if err != nil {
// 		wrappedErr := xerrors.Errorf(": %w", err)
// 		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 	}
// 	for _, n := range nodes {
// 		h.serv.GetVMList(&n)
// 	}
// 	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "test success"))
// }

// func (h *PVEHandler) CloneQuestions(c echo.Context) error {
// 	// リクエストから構造体にデータをコピー
// 	var req CloneQuestionsRequest
// 	if err := c.Bind(&req); err != nil {
// 		wrappedErr := xerrors.Errorf(": %w", err)
// 		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 	}
// 	// データをバリデーションにかける
// 	if err := c.Validate(req); err != nil {
// 		wrappedErr := xerrors.Errorf(": %w", err)
// 		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 	}

// 	var wg sync.WaitGroup

// 	wg.Add(len(req.Teams))

// 	for _, team := range req.Teams {
// 		hostname := fmt.Sprintf("ctf-%s-%d", "contestName", team.ID)
// 		users := []service.User{{
// 			Name:              "user",
// 			Sudo:              "ALL=(ALL) NOPASSWD:ALL",
// 			Shell:             "/bin/bash",
// 			Passwd:            "user",
// 			SshAuthorizedKeys: team.Sshkeys,
// 		},
// 		}
// 		filename := fmt.Sprintf("ctf-%d-%d.yml", req.QID, team.ID)
// 		err := h.serv.CloudinitGenerator(filename, hostname, users)
// 		if err != nil {
// 			wrappedErr := xerrors.Errorf(": %w", err)
// 			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 		}
// 		err = h.serv.TransferFileViaSCP(filename)
// 		if err != nil {
// 			wrappedErr := xerrors.Errorf(": %w", err)
// 			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 		}

// 		nextid, err := h.serv.NextVMID()
// 		if err != nil {
// 			wrappedErr := xerrors.Errorf(": %w", err)
// 			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 		}
// 		fmt.Println(nextid)

// 		nid, err := h.serv.NextVMID()
// 		if err != nil {
// 			wrappedErr := xerrors.Errorf(": %w", err)
// 			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 		}

// 		name := fmt.Sprintf("ctf-%d-%d", req.QID, team.ID)
// 		n, err := strconv.Atoi(nid)
// 		if err != nil {
// 			wrappedErr := xerrors.Errorf(": %w", err)
// 			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 		}

// 		if err = h.serv.CloneVM(name, n, "PVE01", 9000); err != nil {
// 			wrappedErr := xerrors.Errorf(": %w", err)
// 			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 		}
// 		conf := service.VMEdit{
// 			Vmid:     n,
// 			Memory:   "2mb",
// 			Cores:    2,
// 			Node:     "PVE01",
// 			Ipconfig: []string{"ip=10.0.10.90/24,gw=10.0.10.1"},
// 			Scsi:     []string{fmt.Sprintf("vmdisk:vm-%d-disk-0,size=16G", n)},
// 			Cicustom: filename,
// 		}
// 		for i := 0; i < 5; i++ {
// 			fmt.Println(i)
// 			if err = h.serv.EditVM(conf); err != nil {
// 				if i > 4 {
// 					wrappedErr := xerrors.Errorf(": %w", err)
// 					log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 					return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
// 				}
// 			} else {
// 				break
// 			}
// 		}

// 	}

// 	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "create vm success"))
// }
