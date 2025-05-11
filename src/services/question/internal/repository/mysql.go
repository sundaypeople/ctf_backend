package repository

import (
	"database/sql"
	"fmt"

	"github.com/cockroachdb/errors"

	"github.com/LainInTheWired/ctf_backend/question/model"
)

type mysqlRepository struct {
	DB *sql.DB
}
type MysqlRepository interface {
	InsertQuestion(q model.Question) error
	DeleteQuestion(qid int) error
	SelectContestQuestionsByContestID(contestID int) ([]model.Question, error)
	SelectContestQuestions() ([]model.Question, error)
	SelectQuesionByQuestionID(qid int) (model.Question, error)
	UpdateQuestion(q model.Question) error
	UpdateQuestionEnv(e string, vmid int) error
}

func NewMysqlRepository(db *sql.DB) MysqlRepository {
	return &mysqlRepository{
		DB: db,
	}
}

func (m *mysqlRepository) InsertQuestion(q model.Question) error {
	// emailが登録されているかチェック
	ins, err := m.DB.Prepare("INSERT INTO questions (name,env,category_id,description,vmid) VALUES(?,?,?,?,?)")
	if err != nil {
		return errors.Wrap(err, "question insert error")
	}
	defer ins.Close()

	fmt.Printf("INSERT INTO questions (name,env,category_id,describe,vmid) VALUES(%s,%s,%d,%s,%d)", q.Name, q.Env, q.CategoryId, q.Description, q.VMID)

	_, err = ins.Exec(q.Name, q.Env, q.CategoryId, q.Description, q.VMID)
	if err != nil {
		return errors.Wrap(err, "can't insert question")
	}
	return nil
}

func (m *mysqlRepository) DeleteQuestion(qid int) error {
	// emailが登録されているかチェック
	ins, err := m.DB.Prepare("DELETE FROM questions WHERE id = ?")
	if err != nil {
		return errors.Wrap(err, "question insert error")
	}
	defer ins.Close()

	_, err = ins.Exec(qid)
	if err != nil {
		return errors.Wrap(err, "can't insert question")
	}
	return nil
}
func (m *mysqlRepository) UpdateQuestion(q model.Question) error {
	// emailが登録されているかチェック
	ins, err := m.DB.Prepare("UPDATE questions SET name = ?, answer = ? , description = ? WHERE id =?")
	if err != nil {
		return errors.Wrap(err, "contest_teams insert error")
	}
	defer ins.Close()

	_, err = ins.Exec(q.Name, q.Answer, q.Description, q.ID)
	if err != nil {
		return errors.Wrap(err, "can't insert contest_questions")
	}
	return nil
}
func (m *mysqlRepository) UpdateQuestionEnv(e string, vmid int) error {
	// emailが登録されているかチェック
	// この関数は暫定的な処置　後で消す
	ins, err := m.DB.Prepare("UPDATE questions SET env = ? WHERE id =?")
	if err != nil {
		return errors.Wrap(err, "contest_teams insert error")
	}
	defer ins.Close()

	_, err = ins.Exec(e, vmid)
	if err != nil {
		return errors.Wrap(err, "can't insert contest_questions")
	}
	return nil
}

func (m *mysqlRepository) SelectContestQuestionsByContestID(contestID int) ([]model.Question, error) {
	var questions []model.Question
	//  emailよりユーザ情報を取得
	// rows, err := m.DB.Query("SELECT id,name,category_id,description,vmid FROM questions WEHERE id = ?", contestID)
	rows, err := m.DB.Query("SELECT q.id,q.name,c.id,c.name,cq.point,q.description,q.vmid FROM contest_questions as cq JOIN questions as q ON q.id = cq.question_id JOIN contests  AS c ON  c.id = cq.contest_id JOIN category AS cg ON cg.id = q.category_id WHERE cg.id = ?;", contestID)
	if err != nil {
		return nil, errors.Wrap(err, "error select contest")
	}

	for rows.Next() {
		q := model.Question{}
		if err := rows.Scan(&q.ID, &q.Name, &q.CategoryId, &q.CategoryName, &q.Point, &q.Description, &q.VMID); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		questions = append(questions, q)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "select errors")
	}

	return questions, nil
}
func (m *mysqlRepository) SelectContestQuestionsByCategoryID(categoryID int) ([]model.Question, error) {
	var questions []model.Question
	//  emailよりユーザ情報を取得
	// rows, err := m.DB.Query("SELECT id,name,category_id,description,vmid FROM questions WEHERE id = ?", contestID)
	rows, err := m.DB.Query("SELECT q.id,q.name,c.name,q.description,q.vmid FROM questions as q JOIN category AS c ON c.id = q.category_id WHERE c.id = ?", categoryID)
	if err != nil {
		return nil, errors.Wrap(err, "error select contest")
	}

	for rows.Next() {
		q := model.Question{}
		if err := rows.Scan(&q.ID, &q.Name, &q.CategoryName, &q.Description, &q.VMID); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		questions = append(questions, q)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "select errors")
	}
	return questions, nil
}

func (m *mysqlRepository) SelectPointsByConteastID(categoryID int) ([]model.Question, error) {
	var questions []model.Question
	//  emailよりユーザ情報を取得
	// rows, err := m.DB.Query("SELECT id,name,category_id,description,vmid FROM questions WEHERE id = ?", contestID)
	rows, err := m.DB.Query("SELECT team_id,question_id,contest_id,point FROM points as  WHERE contest_id = ?", categoryID)
	if err != nil {
		return nil, errors.Wrap(err, "error select contest")
	}

	for rows.Next() {
		q := model.Question{}
		if err := rows.Scan(&q.ID, &q.Name, &q.CategoryName, &q.Description, &q.VMID); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		questions = append(questions, q)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "select errors")
	}

	return questions, nil
}

func (m *mysqlRepository) SelectQuesionByQuestionID(qid int) (model.Question, error) {
	var quesion model.Question
	var Ans sql.NullString
	if err := m.DB.QueryRow("SELECT q.id,q.name,c.name,q.description,q.vmid,q.answer FROM questions as q JOIN category AS c ON c.id = q.category_id WHERE q.id = ?", qid).Scan(&quesion.ID, &quesion.Name, &quesion.CategoryName, &quesion.Description, &quesion.VMID, &Ans); err != nil {
		if err == sql.ErrNoRows {
			return model.Question{}, errors.Wrap(err, "not exist this id")
		}
		return model.Question{}, errors.Wrap(err, "can't select question by id")
	}
	if Ans.Valid {
		quesion.Answer = Ans.String
	}
	return quesion, nil
}

func (m *mysqlRepository) SelectContestQuestions() ([]model.Question, error) {
	var questions []model.Question
	//  emailよりユーザ情報を取得
	// rows, err := m.DB.Query("SELECT id,name,category_id,description,vmid FROM questions WEHERE id = ?", contestID)
	rows, err := m.DB.Query("SELECT q.id,q.name,c.id,c.name,q.description,q.vmid,q.answer,q.env FROM questions as q JOIN category AS c ON c.id = q.category_id")
	if err != nil {
		return nil, errors.Wrap(err, "error select contest")
	}

	for rows.Next() {
		q := model.Question{}
		var Ans sql.NullString
		if err := rows.Scan(&q.ID, &q.Name, &q.CategoryId, &q.CategoryName, &q.Description, &q.VMID, &Ans, &q.Env); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		if Ans.Valid {
			q.Answer = Ans.String
		}
		questions = append(questions, q)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "select errors")
	}

	return questions, nil
}
