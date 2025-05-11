-- 1. teams テーブルへの挿入
INSERT INTO teams (name) VALUES
('Alpha Team'),
('Beta Team'),
('Gamma Team');

-- 2. users テーブルへの挿入
-- 注: 実際の運用では、パスワードはハッシュ化して保存してください。
INSERT INTO users (email, name, password) VALUES
('a@example.com', 'Alice', '$2a$10$4IcZRIMPxNCkhUGOrW/pW.2T.zW7blr.DVVsPWi0itCH7W.pHqVWe'),
('b@example.com', 'Bob', '$2a$10$4IcZRIMPxNCkhUGOrW/pW.2T.zW7blr.DVVsPWi0itCH7W.pHqVWe'),
('c@example.com', 'Carol', '$2a$10$4IcZRIMPxNCkhUGOrW/pW.2T.zW7blr.DVVsPWi0itCH7W.pHqVWe'),
('d@example.com', 'Dave', '$2a$10$4IcZRIMPxNCkhUGOrW/pW.2T.zW7blr.DVVsPWi0itCH7W.pHqVWe');
-- 2.5 team_users テーブルへの挿入
INSERT INTO team_users (team_id, user_id) VALUES
(1, 1), -- Alpha Team に Alice を追加
(1, 2), -- Alpha Team に Bob を追加
(2, 3), -- Beta Team に Carol を追加
(2, 4); -- Beta Team に Dave を追加

-- 3. roles テーブルへの挿入
INSERT INTO roles (name, namespace,contest) VALUES
('ctf_admin',"ctf" ,'[ "all" ]'),
('contest_admin',"ctf" ,'[  "1" ]'),
('user',"ctf" ,'[ "all" ]');

-- 4. permissions テーブルへの挿入
INSERT INTO permissions (name, description) VALUES
('view', 'view list'),
('edit', 'view and edit'),
('create', 'create anything but not edit anything');

-- 5. user_roles テーブルへの挿入
-- 例: Alice と Bob を admin、Carol と Dave を participant に割り当て
INSERT INTO user_roles (role_id, user_id) VALUES
-- admin ロール (id = 1) を Alice (id = 1) と Bob (id = 2) に割り当て
(1, 1),
(1, 2),
-- participant ロール (id = 2) を Carol (id = 3) と Dave (id = 4) に割り当て
(2, 3),
(3, 4);

-- 6. role_permissions テーブルへの挿入
-- admin ロールには全ての権限を付与
INSERT INTO role_permissions (role_id, permission_id) VALUES
-- manage_contest (id = 1)
(1, 1),
-- submit_flag (id = 2)
(1, 2),
-- view_scores (id = 3)
(1, 3),
-- participant ロールには submit_flag と view_scores の権限を付与
(2, 1),
(2, 2),

(3, 1);

-- 7. contests テーブルへの挿入
INSERT INTO contests (name, start, end) VALUES
('test_db', '2024-11-01 10:00:00', '2025-12-07 18:00:00'),
('Spring CTF 2024', '2024-04-15 09:00:00', '2024-04-21 17:00:00');

-- 8. contest_teams テーブルへの挿入
-- Alpha Team (id = 1) と Beta Team (id = 2) を Winter CTF 2024 (id = 1) に参加
INSERT INTO contest_teams (contest_id, team_id) VALUES
(1, 1),
(1, 2),
-- Gamma Team (id = 3) を Spring CTF 2024 (id = 2) に参加
(2, 1),
(2, 3);

-- 9. questions テーブルへの挿入
INSERT INTO questions (name, category_id, env, description, vmid,answer) VALUES
('Crypto Challenge 1', 1, 'env1', 'Solve the crypto puzzle.', 101,"cc1"),
('Reverse Engineering 1', 2, 'env2', 'Reverse engineer the binary.', 102,"re1"),
('Web Challenge 1', 3, 'env3', 'Find the vulnerability in the web app.', 103,"wc1"),
('Forensics 1', 4, 'env4', 'Analyze the forensic data.', 104,"f1"),
('Crypto Challenge 2', 1, 'env1', 'Advanced crypto puzzle.', 105,"cc2"),
('Web Challenge 2', 3, 'env3', 'Advanced web vulnerability.', 106,"wc2");

-- 10. contest_questions テーブルへの挿入
-- Winter CTF 2024 (id = 1) に Crypto Challenge 1 (id = 1) と Web Challenge 1 (id = 3) を追加
INSERT INTO contest_questions ( contest_id, question_id, point) VALUES
( 1, 1, 100),
( 1, 3, 150),
-- Spring CTF 2024 (id = 2) に Reverse Engineering 1 (id = 2) と Forensics 1 (id = 4) を追加
( 2, 2, 120),
( 2, 4, 130);

-- 11. cloudinit テーブルへの挿入
-- 例: contest_questions_id = 1 (Crypto Challenge 1) に Alpha Team (id = 1) が対応
INSERT INTO cloudinit (contest_id,question_id, team_id, filename) VALUES
(1,1, 1, 'alpha_crypto1_init.sh'),
(1,3, 2, 'beta_web1_init.sh'),
(2,2, 3, 'gamma_re_1_init.sh'),
(2,4, 1, 'alpha_forensics1_init.sh');

-- 12. points テーブルへの挿入
-- 例: Alpha Team が Crypto Challenge 1 でポイントを獲得
INSERT INTO points (team_id, question_id, contest_id,point) VALUES
(1, 1, 1,100), -- Alpha Team が Crypto Challenge 1 を解決
(2, 3, 1,150), -- Beta Team が Web Challenge 1 を解決
(3, 2, 2,120), -- Gamma Team が Reverse Engineering 1 を解決
(1, 4, 2,130); -- Alpha Team が Forensics 1 を解決


INSERT INTO category (name) VALUES
('Crypto'),
('Reverse Engineering'),
('Web'),
('Forensics');