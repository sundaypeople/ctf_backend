package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/LainInTheWired/ctf_backend/team/model"
	"github.com/LainInTheWired/ctf_backend/team/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/xerrors"
)

type TeamHander interface {
	CreateTeam(c echo.Context) error
	DeleteTeam(c echo.Context) error
	ListTeamByContest(c echo.Context) error
	ListTeamUserByContest(c echo.Context) error
	ListUsers(c echo.Context) error
	EditTeam(c echo.Context) error
}

type teamHander struct {
	serv service.TeamService
}

type createTeamRequest struct {
	Name    string `json:"name" validate:"required"`
	UserIDs []int  `json:"user_ids"`
}

type deleteTeamRequest struct {
	ID int `json:"id" validate:"required"`
}

type listTeamByContestRequest struct {
	ID   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type listTeamByContestResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type joinContestRequest struct {
	ContestID int `json:"contest_id" validate:"required"`
	TeamID    int `json:"team_id" validate:"required"`
}

func NewTeamHandler(sv service.TeamService) TeamHander {
	return &teamHander{
		serv: sv,
	}
}

func (t *teamHander) CreateTeam(c echo.Context) error {
	var req createTeamRequest
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

	m := model.Team{
		Name: req.Name,
	}
	tid, err := t.serv.CreateTeam(m, req.UserIDs)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusAccepted, map[string]string{"team_id": strconv.Itoa(tid)})
}

func (t *teamHander) DeleteTeam(c echo.Context) error {
	// var req deleteTeamRequest
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
	stid := c.Param("teamID")
	tid, err := strconv.Atoi(stid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error: param")})
	}

	m := model.Team{
		ID: tid,
	}

	if err := t.serv.DeleteTeam(m); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "Delete teams "))
}
func (t *teamHander) EditTeam(c echo.Context) error {
	var req createTeamRequest
	stid := c.Param("teamID")
	tid, err := strconv.Atoi(stid)
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
	m := model.Team{
		ID:   tid,
		Name: req.Name,
	}

	if err := t.serv.EditTeam(m, req.UserIDs); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	response := map[string]string{
		"message": "success Edit Team",
	}
	return c.JSON(http.StatusOK, response)
}

func (t *teamHander) ListTeamByContest(c echo.Context) error {
	sid := c.Param("contestID")
	id, err := strconv.Atoi(sid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error: param")})
	}

	teams, err := t.serv.ListTeamByContest(id)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusOK, teams)
}

func (t *teamHander) JoinContest(c echo.Context) error {
	var req joinContestRequest
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

	m := model.ContestTeams{
		ContestID: req.ContestID,
		TeamID:    req.TeamID,
	}
	if err := t.serv.JoinContest(m); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "Join contest_teams"))
}

func (t *teamHander) ListTeamUserByContest(c echo.Context) error {
	var teams []model.Team
	scid := c.Param("contestID")
	cid, err := strconv.Atoi(scid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error: param")})
	}
	suid := c.QueryParam("userid")
	uid, err := strconv.Atoi(suid)
	if err != nil {
		teams, err = t.serv.ListTeamUsersByContest(cid)
		if err != nil {
			wrappedErr := xerrors.Errorf(": %w", err)
			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
		}
	} else {
		teams, err = t.serv.ListTeamInContestByUserID(cid, uid)
		if err != nil {
			wrappedErr := xerrors.Errorf(": %w", err)
			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
		}
	}

	return c.JSON(http.StatusOK, teams)
}

func (t *teamHander) ListUsers(c echo.Context) error {
	users, err := t.serv.ListUsers()
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusOK, users)
}
