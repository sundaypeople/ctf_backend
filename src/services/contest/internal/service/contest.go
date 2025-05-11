package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"

	"github.com/LainInTheWired/ctf_backend/contest/model"
	"github.com/LainInTheWired/ctf_backend/contest/repository"
	"github.com/cockroachdb/errors"
)

type ContestService interface {
	CreateContest(c model.Contest) error
	DeleteContest(c model.Contest) error
	CreateTeamContest(c model.ContestsTeam) error
	DeleteTeamContest(c model.ContestsTeam) error
	ListContest() ([]model.Contest, error)
	ListContestByTeams(tid int) ([]model.Contest, error)
	JoinListContestQuesionts(ContestQuestions []model.ContestQuestions) error
	StartContest(cid int) error
	GetPoints(cid int) ([]model.ResponsePoints, error)
	CheckQuestion(cid int, qid int, tid int, ans string) (bool, error)
	ListQuestionsByContestID(cid int, tid int) (*model.Contest, error)
	GetTeamByUserID(cid int, uid int) ([]model.Team, error)
	UpdateContestQuesionts(cq *model.ContestQuestions) error
	StopContest(cid int) error
	GetCloudinit(cid, tid, qid int) (*model.Cloudinit, error)
	GetClusterResource() ([]model.ClusterResources, error)
	AllDeleteVM() error
	ListQuestionsByContestIDAdmin(cid int) (*model.Contest, error)
}

type contestService struct {
	pveRepo   repository.PVEAPIRepository
	mysqlRepo repository.MysqlRepository
	teamRepo  repository.TeamRepository
	quesRepo  repository.QuestionRepository
}

func NewContestService(pveRepo repository.PVEAPIRepository, mysqlRepo repository.MysqlRepository, teamRepo repository.TeamRepository, quesRepo repository.QuestionRepository) ContestService {
	return &contestService{
		pveRepo:   pveRepo,
		mysqlRepo: mysqlRepo,
		teamRepo:  teamRepo,
		quesRepo:  quesRepo,
	}
}

func (r *contestService) CreateContest(c model.Contest) error {
	// モデルの構造体に移し替えてから、repositoryに渡す
	if err := r.mysqlRepo.InsertContest(c); err != nil {
		return errors.Wrap(err, "can't create contest")
	}
	return nil
}

func (r *contestService) DeleteContest(c model.Contest) error {
	if err := r.mysqlRepo.DeleteContest(c); err != nil {
		return errors.Wrap(err, "can't delete contest")
	}
	return nil
}
func (r *contestService) CreateTeamContest(c model.ContestsTeam) error {
	if err := r.mysqlRepo.InsertTeamContests(c); err != nil {
		return errors.Wrap(err, "can't create team_contests")
	}
	return nil
}
func (r *contestService) DeleteTeamContest(c model.ContestsTeam) error {
	if err := r.mysqlRepo.DeleteTeamContests(c); err != nil {
		return errors.Wrap(err, "can't delete team_contests")
	}
	return nil
}

func (r *contestService) ListContest() ([]model.Contest, error) {
	contests, err := r.mysqlRepo.SelectContest()
	if err != nil {
		return nil, errors.Wrap(err, "can't delete team_contests")
	}
	return contests, nil
}
func (r *contestService) ListContestByTeams(tid int) ([]model.Contest, error) {
	contests, err := r.mysqlRepo.SelectContestsByTeamID(tid)
	if err != nil {
		return nil, errors.Wrap(err, "can't delete team_contests")
	}
	return contests, nil
}

func (r *contestService) JoinListContestQuesionts(cqs []model.ContestQuestions) error {
	// 送信された question_id のリストを取得
	incomingQuestionIDs := make([]int, len(cqs))
	incomingQuestionMap := make(map[int]model.ContestQuestions)
	for i, q := range cqs {
		incomingQuestionIDs[i] = q.QuestionID
		incomingQuestionMap[q.QuestionID] = q
	}
	con, err := r.mysqlRepo.SelectContestQuestionsByContestID(cqs[0].ContestID)
	if err != nil {
		return errors.Wrap(err, "can't select ContestQuestions by ContestID")
	}
	existingQuestionMap := make(map[int]int) // question_id -> point
	for _, q := range con.Questions {
		existingQuestionMap[q.ID] = q.Point
	}
	for _, cq := range cqs {
		if _, exists := existingQuestionMap[cq.QuestionID]; !exists {
			// 新規追加
			err := r.mysqlRepo.InsertContestsQuestions(&cq)
			if err != nil {
				return errors.Wrap(err, "can't Join Contest Questions")
			}
		}
		// 既存のものは放置
	}
	for existingQID := range existingQuestionMap {
		if _, exists := incomingQuestionMap[existingQID]; !exists {
			// 新規追加
			err := r.mysqlRepo.DeleteContestsQuestions(existingQID, cqs[0].ContestID)
			if err != nil {
				return errors.Wrap(err, "can't Delete Contest Questions")
			}
		}
	}

	return nil
}
func (r *contestService) UpdateContestQuesionts(cq *model.ContestQuestions) error {
	err := r.mysqlRepo.UpdateContestsQuestions(cq)
	if err != nil {
		return errors.Wrap(err, "can't Upate Contest Questions")
	}
	return nil
}
func (r *contestService) StartContest(cid int) error {
	teams, err := r.teamRepo.ListTeamUsersByContest(cid, nil)
	if err != nil {
		return errors.Wrap(err, "can't get ListTeamUsers")
	}
	fmt.Printf("%+v", teams)
	// questions, err := r.quesRepo.GetListQuestionsByContest(cid)
	questions, err := r.mysqlRepo.SelectContestQuestionsByContestID(cid)
	if err != nil {
		return errors.Wrap(err, "can't get ListQuestions")
	}
	cluster, err := r.pveRepo.GetClusterResource()
	if err != nil {
		return errors.Wrap(err, "can't get cluster resouece")
	}
	mapcluster := map[string]model.ClusterResources{}
	for _, c := range cluster {
		if c.Type == "qemu" {
			mapcluster[c.Name] = c
		}
	}
	for _, team := range teams {
		for _, ques := range questions.Questions {
			name := fmt.Sprintf("%d-%d-%d", cid, team.ID, ques.ID)
			fmt.Println("name: ", name)

			if _, ok := mapcluster[name]; ok {
				fmt.Println("skip clone vm")
				fmt.Println(mapcluster[name])
				continue
			}
			password, err := generatePassword(16)
			if err != nil {
				return errors.Wrap(err, "can't generate password")
			}
			m := model.QuesionRequest{
				ID:       ques.VMID,
				Name:     name,
				Password: password,
			}
			vmid, err := r.quesRepo.CloneQuestion(m)
			if err != nil {
				return errors.Wrap(err, "can't get ListQuestions")
			}
			cloudinit := model.Cloudinit{
				QuestionID: ques.ID,
				ContestID:  cid,
				Filename:   "",
				TeamID:     team.ID,
				VMID:       vmid,
				Access:     password,
			}
			err = r.mysqlRepo.InsertCloudinit(cloudinit)
			if err != nil {
				errors.Wrap(err, "can't InsertCloudinit")
			}
		}
	}
	return nil
}

func (r *contestService) StopContest(cid int) error {
	cloudinit, err := r.mysqlRepo.SelectCloudinitByContestID(cid)
	if err != nil {
		return errors.Wrap(err, "can't get Cloudinit")
	}

	for _, c := range cloudinit {
		if err = r.quesRepo.DeleteVM(c.VMID); err != nil {
			return errors.Wrap(err, "can't get ListQuestions")
		}
		cloudinit := model.Cloudinit{
			QuestionID: c.QuestionID,
			ContestID:  cid,
			TeamID:     c.TeamID,
		}
		err = r.mysqlRepo.DeleteCloudinit(cloudinit)
		if err != nil {
			errors.Wrap(err, "can't InsertCloudinit")
		}
	}

	return nil
}
func (r *contestService) AllDeleteVM() error {
	// cloudinit, err := r.mysqlRepo.SelectCloudinitByContestID(cid)
	// if err != nil {
	// 	return errors.Wrap(err, "can't get Cloudinit")
	// }

	cluster, err := r.pveRepo.GetClusterResource()
	if err != nil {
		return errors.Wrap(err, "can't get cluster resouece")
	}

	var pattern = regexp.MustCompile(`^\d+-\d+-\d+$`)
	filterdcluster := []model.ClusterResources{}
	for _, c := range cluster {
		if pattern.MatchString(c.Name) {
			filterdcluster = append(filterdcluster, c)
		}
	}
	fmt.Println("フィルタリングされたアイテム:")
	for _, item := range filterdcluster {
		fmt.Println(item.Name, ":", item.Vmid)
		if err = r.quesRepo.DeleteVM(item.Vmid); err != nil {
			return errors.Wrap(err, "can't get ListQuestions")
		}
	}

	// for _, c := range cloudinit {
	// 	if err = r.quesRepo.DeleteVM(c.VMID); err != nil {
	// 		// return errors.Wrap(err, "can't get ListQuestions")
	// 	}
	// 	cloudinit := model.Cloudinit{
	// 		QuestionID: c.QuestionID,
	// 		ContestID:  cid,
	// 		TeamID:     c.TeamID,
	// 	}
	// 	err = r.mysqlRepo.DeleteCloudinit(cloudinit)
	// 	if err != nil {
	// 		errors.Wrap(err, "can't InsertCloudinit")
	// 	}
	// }

	return nil
}

func generatePassword(length int) (string, error) {
	var (
		lowerLetters = "abcdefghijklmnopqrstuvwxyz"
		upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits       = "0123456789"
		// symbols      = "!@#$%^&*"
		allChars = lowerLetters + upperLetters + digits
	)

	password := make([]byte, length)
	for i := range password {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			return "", err
		}
		password[i] = allChars[index.Int64()]
	}

	return string(password), nil
}
func (r *contestService) GetPoints(cid int) ([]model.ResponsePoints, error) {
	var res []model.ResponsePoints
	teams, err := r.teamRepo.ListTeamUsersByContest(cid, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't get point")
	}

	points, err := r.mysqlRepo.SelectPoint(cid)
	if err != nil {
		return nil, errors.Wrap(err, "can't get point")
	}
	for _, team := range teams {
		var tpoints []model.Point
		for _, point := range points {
			if point.TeamID == team.ID {
				tpoint := model.Point{
					Point:      point.Point,
					InsertDate: point.InsertDate,
				}
				tpoints = append(tpoints, tpoint)
			}
		}
		t := model.ResponsePoints{
			TeamID: team.ID,
			Name:   team.Name,
			Points: tpoints,
		}
		res = append(res, t)
	}
	return res, nil
}

func (r *contestService) CheckQuestion(cid int, qid int, tid int, ans string) (bool, error) {
	contest, err := r.mysqlRepo.SelectContestQuestionsByContestID(cid)
	if err != nil {
		return false, errors.Wrap(err, "can't get Questions")
	}
	question := FilterQuestionsByID(contest.Questions, qid)
	if question == nil {
		return false, errors.Wrap(err, "can't filter quesion")
	}
	if question.Answer == ans {
		if err := r.mysqlRepo.InsertPoint(tid, qid, cid, question.Point); err != nil {
			return false, errors.Wrap(err, "can't get Questions")
		}
		return true, nil
	}

	return false, nil
}

// FilterQuestionsByID 指定されたIDでフィルタリングする関数
func FilterQuestionsByID(questions []model.Question, id int) *model.Question {
	for _, q := range questions {
		fmt.Printf("questions: %+v\n", q)
		if q.ID == id {
			return &q
		}
	}
	return nil
}

func (s *contestService) ListQuestionsByContestID(cid int, tid int) (*model.Contest, error) {
	contests, err := s.mysqlRepo.SelectContestQuestionsByContestID(cid)
	if err != nil {
		return nil, errors.Wrap(err, "get questions")
	}
	points, err := s.mysqlRepo.SelectPointByTeamidAndContestid(cid, tid)
	fmt.Printf("%+v", points)
	if err != nil {
		return nil, errors.Wrap(err, "get questions")
	}
	// pointsをマップに変換
	pointMap := make(map[int]int)
	for _, point := range points {
		if _, exists := pointMap[point.QuestionID]; !exists {
			pointMap[point.QuestionID] = point.Point
		}
	}
	// cloudinit, err := s.mysqlRepo.SelectCloudinitByContestIDAndTeamID(cid, tid)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "can't error")
	// }

	// contests.Questionsを更新
	// cmap := map[int]int{}
	// for _, c := range cloudinit {
	// 	cmap[c.QuestionID] = c.VMID
	// }
	for i := range contests.Questions {
		if point, exists := pointMap[contests.Questions[i].ID]; exists {
			contests.Questions[i].CurrentPoint = point
		}
		// ips, err := s.pveRepo.GetIPByVMID(cmap[contests.Questions[i].ID])
		// if err != nil {
		// 	break
		// }
		// contests.Questions[i].IPs = *ips
		// fmt.Println(ips)

	}
	fmt.Printf("conntest:%+v", contests)
	return &contests, nil
}

func (s *contestService) ListQuestionsByContestIDAdmin(cid int) (*model.Contest, error) {
	contests, err := s.mysqlRepo.SelectContestQuestionsByContestID(cid)
	if err != nil {
		return nil, errors.Wrap(err, "get questions")
	}
	fmt.Printf("%+v", contests)
	return &contests, nil
}

func (s *contestService) GetTeamByUserID(cid int, uid int) ([]model.Team, error) {
	fmt.Printf("%d\n", uid)
	teams, err := s.teamRepo.ListTeamUsersByContest(cid, &uid)
	if err != nil {
		return nil, nil
	}
	return teams, nil
}

func (s *contestService) GetCloudinit(cid, tid, qid int) (*model.Cloudinit, error) {
	cloudinit, err := s.mysqlRepo.SelectCloudinitByContestIDAndTeamIDAndQuestionID(cid, tid, qid)
	if err != nil {
		return nil, errors.Wrap(err, "errors")
	}
	ips, err := s.pveRepo.GetIPByVMID(cloudinit.VMID)
	if err != nil {
		return nil, errors.Wrap(err, "errors")
	}
	cloudinit.IPs = *ips
	return cloudinit, nil
}

func (s *contestService) GetClusterResource() ([]model.ClusterResources, error) {
	cluster, err := s.pveRepo.GetClusterResource()
	if err != nil {
		return nil, errors.Wrap(err, "can't get error")
	}
	return cluster, nil
}
