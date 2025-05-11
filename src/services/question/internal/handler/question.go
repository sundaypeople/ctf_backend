package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/LainInTheWired/ctf_backend/question/model"
	"github.com/LainInTheWired/ctf_backend/question/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/xerrors"
)

type QuesionHander interface {
	CreateQuestion(c echo.Context) error
	DeleteQuestion(c echo.Context) error
	GetQuestionsInContest(c echo.Context) error
	GetQuestions(c echo.Context) error
	CloneQuestion(c echo.Context) error
	DeleteVM(c echo.Context) error
	GetQuesionIp(c echo.Context) error
	UpdateQuestion(c echo.Context) error
	ToTemplate(c echo.Context) error
}

type quesionHander struct {
	serv service.QuesionService
}
type quesionRequest struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CategoryID  int      `json:"category_id"`
	Env         string   `json:"env"`
	Sshkeys     []string `json:"sshkeys"`
	Memory      int      `json:"memory"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	CPUs        int      `json:"cpu"`
	Disk        int      `json:"disk"`
	IP          string   `json:"ip,omitempty" validate:"omitempty,cidr"`
	Gateway     string   `json:"gateway,omitempty" validate:"omitempty,ip"`
	Filename    string   `json:"filename"`
}

type updateQuestion struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Answer      string `json:"answer"`
}
type QuestionsInContestRequest struct {
	ContestID int `json:"contest_id"`
}

type CloneQuestionRequest struct {
	VMID        int      `json:"vmid"`
	ContestName string   `json:"contest_name"`
	Sshkeys     []string `json:"sshkeys"`
	Password    string   `json:"password"`
}

func NewQuestionHander(s service.QuesionService) QuesionHander {
	return &quesionHander{
		serv: s,
	}
}

func (h *quesionHander) CreateQuestion(c echo.Context) error {
	var req quesionRequest
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

	m := model.CreateQuestion{
		Name:        req.Name,
		Env:         req.Env,
		CategoryID:  req.CategoryID,
		Description: req.Description,
		Sshkeys:     req.Sshkeys,
		CPUs:        req.CPUs,
		Disk:        req.Disk,
		ID:          req.ID,
		Memory:      req.Memory,
		IP:          req.IP,
		Gateway:     req.Gateway,
		Username:    req.Username,
		Password:    req.Password,
	}

	if err := h.serv.CreateQuestion(m); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "delete contest_teams"))
}

func (h *quesionHander) DeleteQuestion(c echo.Context) error {
	// var req quesionRequest
	// if err := c.Bind(&req); err != nil {
	// 	wrappedErr := xerrors.Errorf(": %w", err)
	// 	log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	// }
	// // データをバリデーションにかける
	// if err := c.Validate(req); err != nil {
	// 	wrappedErr := xerrors.Errorf(": %w", err)
	// 	log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	// }
	sqid := c.Param("questionID")
	qid, err := strconv.Atoi(sqid)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	if err := h.serv.DeleteQuestion(qid); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "delete contest_teams"))
}

func (h *quesionHander) GetQuestions(c echo.Context) error {
	sid := c.QueryParam("id")
	fmt.Println(sid)
	if sid == "" {
		q, err := h.serv.GetQuestions()
		if err != nil {
			wrappedErr := xerrors.Errorf(": %w", err)
			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
		}
		return c.JSON(http.StatusAccepted, q)
	} else {
		id, err := strconv.Atoi(sid)
		if err != nil {
			wrappedErr := xerrors.Errorf(": %w", err)
			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
		}
		q, err := h.serv.GetQuesionByID(id)
		if err != nil {
			wrappedErr := xerrors.Errorf(": %w", err)
			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})

		}
		return c.JSON(http.StatusAccepted, q)
	}
}

func (h *quesionHander) GetQuestionsInContest(c echo.Context) error {
	sid := c.Param("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error: param")})
	}

	q, err := h.serv.GetQuestionsInContest(id)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, q)
}

func (h *quesionHander) QuestionHander(c echo.Context) error {

	return nil
}

func (h *quesionHander) CloneQuestion(c echo.Context) error {
	var req quesionRequest
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
	m := model.CreateQuestion{
		Name:        req.Name,
		Env:         req.Env,
		Description: req.Description,
		Sshkeys:     req.Sshkeys,
		CPUs:        req.CPUs,
		Disk:        req.Disk,
		ID:          req.ID,
		Memory:      req.Memory,
		IP:          req.IP,
		Password:    req.Password,
		Gateway:     req.Gateway,
	}
	vmid, err := h.serv.CloneQuestion(m)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, map[string]int{"data": vmid})
}

func (h *quesionHander) GetQuesionByID(c echo.Context) error {
	sid := c.QueryParam("questionID")
	id, err := strconv.Atoi(sid)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	question, err := h.serv.GetQuesionByID(id)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusAccepted, question)

}

func (h *quesionHander) DeleteVM(c echo.Context) error {
	var req quesionRequest
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
	// sid := c.QueryParam("questionID")
	// id, err := strconv.Atoi(sid)
	// if err != nil {
	// 	wrappedErr := xerrors.Errorf(": %w", err)
	// 	log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	// }
	if err := h.serv.DeleteVM(req.ID); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, map[string]string{"data": "delete VM"})
}

func (h *quesionHander) GetQuesionIp(c echo.Context) error {
	svmid := c.Param("vmid")
	vmid, err := strconv.Atoi(svmid)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	ip, err := h.serv.GetQuesionIp(vmid)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, map[string]interface{}{"data": *ip})
}

func (h *quesionHander) UpdateQuestion(c echo.Context) error {
	var req updateQuestion
	sid := c.Param("questionID")
	id, err := strconv.Atoi(sid)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
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

	q := model.Question{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Answer:      req.Answer,
	}

	if err := h.serv.UpdateQuestion(q); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, map[string]string{"message": "success update question"})
}

func (h *quesionHander) ToTemplate(c echo.Context) error {
	// var req quesionRequest
	// if err := c.Bind(&req); err != nil {
	// 	wrappedErr := xerrors.Errorf(": %w", err)
	// 	log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	// }
	sid := c.Param("questionID")
	id, err := strconv.Atoi(sid)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	svid := c.Param("vmid")
	vmid, err := strconv.Atoi(svid)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	// データをバリデーションにかける
	// if err := c.Validate(req); err != nil {
	// 	wrappedErr := xerrors.Errorf(": %w", err)
	// 	log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	// }

	if err := h.serv.Template(id, vmid); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, map[string]string{"data": "success"})
}
