package service

import (
	"fmt"
	"strconv"

	"github.com/LainInTheWired/ctf_backend/question/model"
	"github.com/LainInTheWired/ctf_backend/question/repository"
	"github.com/cockroachdb/errors"
)

type quesionService struct {
	myrepo     repository.MysqlRepository
	pveapirepo repository.PVEAPIRepository
	teamrepo   repository.TeamRepository
}

type QuesionService interface {
	CreateQuestion(q model.CreateQuestion) error
	DeleteQuestion(qid int) error
	CloneQuestion(q model.CreateQuestion) (int, error)
	GetQuestionsInContest(contestID int) ([]model.Question, error)
	GetQuestions() ([]model.Question, error)
	GetQuesionByID(qid int) (model.Question, error)
	DeleteVM(vmid int) error
	GetQuesionIp(vmid int) (*model.ResponseIPs, error)
	UpdateQuestion(q model.Question) error
	Template(id, vmid int) error
}

func NewQuestionService(r repository.MysqlRepository, p repository.PVEAPIRepository, t repository.TeamRepository) QuesionService {
	return &quesionService{
		myrepo:     r,
		pveapirepo: p,
		teamrepo:   t,
	}
}

func (s *quesionService) CreateQuestion(q model.CreateQuestion) error {
	// モデルの構造体に移し替えてから、repositoryに渡す
	clconf := &model.CloudinitResponse{
		Filename:  q.Name + ".yaml",
		Hostname:  q.Name,
		Sshkeys:   q.Sshkeys,
		SshPwauth: "1",
		Username:  q.Username,
		Password:  q.Password,
	}

	vmconfig := &model.CreateVM{
		Cloneid:  9000,
		Name:     q.Name,
		Memory:   q.Memory,
		IP:       q.IP,
		Gateway:  q.Gateway,
		Disk:     q.Disk,
		Cicustom: q.Name + ".yaml",
		CPU:      q.CPUs,
	}
	fmt.Printf("%+v", vmconfig)

	if err := s.pveapirepo.Cloudinit(clconf); err != nil {
		return errors.Wrap(err, "can't create contest")
	}
	svmid, err := s.pveapirepo.CreateVM(vmconfig)

	if err != nil {
		return errors.Wrap(err, "can't create contest")
	}

	vmid, err := strconv.Atoi(svmid)
	if err != nil {
		return errors.Wrap(err, "can't Atoi vmid")
	}

	ques := &model.Question{
		Name:        q.Name,
		CategoryId:  q.CategoryID,
		Description: q.Description,
		Env:         q.Env,
		VMID:        vmid,
	}

	if err := s.myrepo.InsertQuestion(*ques); err != nil {
		return errors.Wrap(err, "can't create contest")
	}

	return nil
}

func (s *quesionService) Template(id, vmid int) error {
	if err := s.pveapirepo.Template(vmid); err != nil {
		return errors.Wrap(err, "can't to template")
	}
	if err := s.myrepo.UpdateQuestionEnv("true", id); err != nil {
		return errors.Wrap(err, "can't to template")
	}
	return nil
}

func (s *quesionService) DeleteQuestion(qid int) error {
	// モデルの構造体に移し替えてから、repositoryに渡す
	ques, err := s.myrepo.SelectQuesionByQuestionID(qid)
	if err != nil {
		return errors.Wrap(err, "can't Select Questinos")
	}
	if err := s.pveapirepo.DeleteVM(ques.VMID); err != nil {
		return errors.Wrap(err, "can't Delete vm")
	}
	if err := s.myrepo.DeleteQuestion(qid); err != nil {
		return errors.Wrap(err, "can't Delete Question")
	}
	return nil
}

func (s *quesionService) DeleteVM(vmid int) error {
	fmt.Println("vmid", vmid)
	if err := s.pveapirepo.DeleteVM(vmid); err != nil {
		return errors.Wrap(err, "can't Delete vm")
	}
	return nil
}

func (s *quesionService) GetQuestions() ([]model.Question, error) {
	// モデルの構造体に移し替えてから、repositoryに渡す
	q, err := s.myrepo.SelectContestQuestions()
	if err != nil {
		return nil, errors.Wrap(err, "can't create contest")
	}
	return q, nil
}

func (s *quesionService) GetQuestionsInContest(contestID int) ([]model.Question, error) {
	// モデルの構造体に移し替えてから、repositoryに渡す
	q, err := s.myrepo.SelectContestQuestionsByContestID(contestID)
	if err != nil {
		return nil, errors.Wrap(err, "can't create contest")
	}
	return q, nil
}

func (s *quesionService) CloneQuestion(q model.CreateQuestion) (int, error) {
	// モデルの構造体に移し替えてから、repositoryに渡す
	clconf := &model.CloudinitResponse{
		Filename:  q.Name + ".yaml",
		Hostname:  q.Name,
		Sshkeys:   q.Sshkeys,
		SshPwauth: "1",
		Username:  "user",
		Password:  q.Password,
	}

	vmconfig := &model.CreateVM{
		Cloneid:  q.ID,
		Name:     q.Name,
		Memory:   q.Memory,
		IP:       q.IP,
		Gateway:  q.Gateway,
		Disk:     q.Disk,
		Cicustom: q.Name + ".yaml",
		CPU:      q.CPUs,
	}

	if err := s.pveapirepo.Cloudinit(clconf); err != nil {
		return 0, errors.Wrap(err, "can't create contest")
	}
	svmid, err := s.pveapirepo.CreateVM(vmconfig)
	if err != nil {
		return 0, errors.Wrap(err, "can't create contest")
	}

	vmid, err := strconv.Atoi(svmid)
	if err != nil {
		return 0, errors.Wrap(err, "can't Atoi vmid")
	}

	// ques := &model.Question{
	// 	Name:        q.Name,
	// 	CategoryId:  q.CategoryID,
	// 	Description: q.Description,
	// 	Env:         q.Env,
	// 	VMID:        vmid,
	// }

	// if err := s.myrepo.InsertQuestion(*ques); err != nil {
	// 	return errors.Wrap(err, "can't create contest")
	// }

	return vmid, nil
}

func (s *quesionService) GetQuesionByID(qid int) (model.Question, error) {
	question, err := s.myrepo.SelectQuesionByQuestionID(qid)
	if err != nil {
		return model.Question{}, errors.Wrap(err, "can't get question by id")
	}
	return question, nil
}

func (s *quesionService) GetQuesionIp(vmid int) (*model.ResponseIPs, error) {
	ips, err := s.pveapirepo.GetIPByVMID(vmid)
	if err != nil {
		return nil, errors.Wrap(err, "errors")
	}
	return ips, nil
}

func (s *quesionService) UpdateQuestion(q model.Question) error {
	if err := s.myrepo.UpdateQuestion(q); err != nil {
		return errors.Wrap(err, "errors")

	}
	return nil
}
