package hander

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/LainInTheWired/ctf_backend/contest/model"
	"github.com/LainInTheWired/ctf_backend/contest/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"golang.org/x/xerrors"
)

type ContestHander interface {
	CreateContest(c echo.Context) error
	DeleteContest(c echo.Context) error
	JoinTeamsinContest(c echo.Context) error
	DeleteTeamContest(c echo.Context) error
	ListContest(c echo.Context) error
	ListContestByTeams(c echo.Context) error
	StartContest(c echo.Context) error
	GetPoints(c echo.Context) error
	CheckAnswer(c echo.Context) error
	ListQuestionsByContestID(c echo.Context) error
	JoinContestQuestions(c echo.Context) error
	UpdateContestQuestions(c echo.Context) error
	StopContest(c echo.Context) error
	GetCloudinit(c echo.Context) error
	GetClusterResource(c echo.Context) error
	AllVMDelete(c echo.Context) error
}

type contestHander struct {
	serv service.ContestService
}

type createContestRequest struct {
	Name      string `json"name" validate:"required"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type deleteContestRequest struct {
	ID int `json:"id" validate:"required"`
}
type contestJoinTeamRquest struct {
	ContestID int `json:"contest_id"`
	TeamID    int `json:"team_id" validate:"required"`
}
type listContestResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
type joinContestQuesiontsRequest struct {
	QID   int `json:"qid" validate:"required"`
	Point int `json:"point" validate:"required,min=0"`
}
type updateContestQuesionts struct {
	Point int `json:"point" validate:"required,min=0"`
}

type startContestRequest struct {
	contestID int `json:"contest_id"`
}

type CheckAnswerResponse struct {
	// TeamID     int    `json:"team_id"`
	Answer     string `json:"answer"`
	QuestionID int    `json:"question_id"`
}

func NewContestHander(cr service.ContestService) ContestHander {
	return &contestHander{
		serv: cr,
	}
}

func (h *contestHander) CreateContest(c echo.Context) error {
	// リクエストから構造体にデータをコピー
	var req createContestRequest
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

	layout := "2006-01-02 15:04:05"

	st, err := time.Parse(layout, req.StartDate)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	et, err := time.Parse(layout, req.EndDate)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	m := model.Contest{
		Name:      req.Name,
		StartDate: st,
		EndDate:   et,
	}
	fmt.Println(m)

	if err := h.serv.CreateContest(m); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "create contest"))
}

func (h *contestHander) DeleteContest(c echo.Context) error {
	// var req deleteContestRequest
	scid := c.Param("contestID")
	cid, err := strconv.Atoi(scid)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
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

	m := model.Contest{
		ID: cid,
	}
	if err := h.serv.DeleteContest(m); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "delete contest "))

}

func (h *contestHander) JoinTeamsinContest(c echo.Context) error {
	// var req contestJoinTeamRquest
	scid := c.Param("contestID")
	cid, err := strconv.Atoi(scid)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	stid := c.Param("teamID")
	tid, err := strconv.Atoi(stid)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
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

	m := model.ContestsTeam{
		ContestID: cid,
		TeamID:    tid,
	}

	if err := h.serv.CreateTeamContest(m); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "join contest_teams"))
}
func (h *contestHander) DeleteTeamContest(c echo.Context) error {
	// var req contestJoinTeamRquest
	scid := c.Param("contestID")
	cid, err := strconv.Atoi(scid)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	stid := c.Param("teamID")
	tid, err := strconv.Atoi(stid)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
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

	m := model.ContestsTeam{
		ContestID: cid,
		TeamID:    tid,
	}

	if err := h.serv.DeleteTeamContest(m); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "delete contest_teams"))
}

func (h *contestHander) ListContest(c echo.Context) error {

	contests, err := h.serv.ListContest()
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusAccepted, contests)
}

func (h *contestHander) ListContestByTeams(c echo.Context) error {
	sid := c.Param("contestID")
	id, err := strconv.Atoi(sid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error: param")})
	}

	contests, err := h.serv.ListContestByTeams(id)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusAccepted, contests)
}

func (h *contestHander) JoinContestQuestions(c echo.Context) error {
	reqs := make([]joinContestQuesiontsRequest, 0)
	scid := c.Param("contestID")
	cid, err := strconv.Atoi(scid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error: param")})
	}
	if err := c.Bind(&reqs); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	// データをバリデーションにかける
	for _, req := range reqs {
		if err := c.Validate(req); err != nil {
			wrappedErr := xerrors.Errorf(": %w", err)
			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
		}
	}
	cqs := []model.ContestQuestions{}

	for _, req := range reqs {
		cq := model.ContestQuestions{
			QuestionID: req.QID,
			ContestID:  cid,
			Point:      req.Point,
		}
		cqs = append(cqs, cq)
	}

	// qcids := []map[string]int{}
	// for _, req := range *req {
	// 	t := map[string]int{"qid": req.QID, "cid": cid}
	// 	qcids = append(qcids, t)
	// }

	if err := h.serv.JoinListContestQuesionts(cqs); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "join contests_quesions"))
}
func (h *contestHander) UpdateContestQuestions(c echo.Context) error {
	var req updateContestQuesionts
	scid := c.Param("contestID")
	cid, err := strconv.Atoi(scid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error: param")})
	}
	sqid := c.Param("questionID")
	qid, err := strconv.Atoi(sqid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error: param")})
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

	cq := &model.ContestQuestions{
		ContestID:  cid,
		QuestionID: qid,
		Point:      req.Point,
	}
	if err := h.serv.UpdateContestQuesionts(cq); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "join contests_quesions"))
}
func (h *contestHander) StartContest(c echo.Context) error {
	sid := c.Param("contestID")
	id, err := strconv.Atoi(sid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error: param")})
	}
	var req startContestRequest
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

	if err := h.serv.StartContest(id); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "join contests_quesions"))
}
func (h *contestHander) StopContest(c echo.Context) error {
	scid := c.Param("contestID")
	cid, err := strconv.Atoi(scid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error: param")})
	}
	var req startContestRequest
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

	if err := h.serv.StopContest(cid); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "join contests_quesions"))
}
func (h *contestHander) GetPoints(c echo.Context) error {
	sid := c.Param("contestID")
	id, err := strconv.Atoi(sid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error: param")})
	}
	points, err := h.serv.GetPoints(id)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, points)
}

func (h *contestHander) CheckAnswer(c echo.Context) error {
	scid := c.Param("contestID")
	cid, err := strconv.Atoi(scid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error: param")})
	}

	suid := c.Request().Header.Get("X-User-ID")
	uid, err := strconv.Atoi(suid)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User ID not found",
		})
	}
	teams, err := h.serv.GetTeamByUserID(cid, uid)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	var req CheckAnswerResponse
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
	ans, err := h.serv.CheckQuestion(cid, req.QuestionID, teams[0].ID, req.Answer)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	if !ans {
		return c.JSON(http.StatusAccepted, map[string]any{
			"message": "fail your answe",
			"correct": false,
		})
	} else {
		return c.JSON(http.StatusAccepted, map[string]any{
			"message": "correct your answer",
			"correct": true,
		})
	}
}

func (h *contestHander) ListQuestionsByContestID(c echo.Context) error {
	suid := c.Request().Header.Get("X-User-ID")
	fmt.Println(suid)
	uid, err := strconv.Atoi(suid)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User ID not found",
		})
	}
	// cid := 1
	scid := c.Param("contestID")
	cid, err := strconv.Atoi(scid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error: param")})
	}
	//  adminやった場合
	if uid == 1 {
		ques, err := h.serv.ListQuestionsByContestIDAdmin(cid)
		if err != nil {
			wrappedErr := xerrors.Errorf(": %w", err)
			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})

		}
		return c.JSON(http.StatusOK, ques)
	}
	teams, err := h.serv.GetTeamByUserID(cid, uid)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	fmt.Printf("%+v", teams)

	if len(teams) == 0 {
		wrappedErr := errors.Wrap(err, "request bind error")
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusForbidden, map[string]string{"message": "not join team"})
	}

	fmt.Printf("%+v", teams)
	points, err := h.serv.ListQuestionsByContestID(cid, teams[0].ID)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	fmt.Printf("%+v", points)

	return c.JSON(http.StatusAccepted, points)
}

func (h *contestHander) GetCloudinit(c echo.Context) error {
	scid := c.Param("contestID")
	cid, err := strconv.Atoi(scid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error: param")})
	}
	sqid := c.Param("questionID")
	qid, err := strconv.Atoi(sqid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error: param")})
	}
	suid := c.Request().Header.Get("X-User-ID")
	fmt.Println(c.Request().Header)

	uid, err := strconv.Atoi(suid)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User ID not found",
		})
	}
	teams, err := h.serv.GetTeamByUserID(cid, uid)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	cloudinit, err := h.serv.GetCloudinit(cid, teams[0].ID, qid)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusOK, cloudinit)
}

func (h *contestHander) GetClusterResource(c echo.Context) error {
	cluster, err := h.serv.GetClusterResource()
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusOK, cluster)
}

func (h *contestHander) AllVMDelete(c echo.Context) error {
	err := h.serv.AllDeleteVM()
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusOK, nil)
}
