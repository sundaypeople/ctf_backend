-- 既存の 'ctf' データベースを削除（存在する場合）
DROP DATABASE IF EXISTS ctf;

-- 新しく 'ctf' データベースを作成
CREATE DATABASE ctf;

-- 'ctf' データベースを使用
USE ctf;

-- 'teams'
CREATE TABLE teams (
    id           INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    create_date  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'contests'
CREATE TABLE contests (
    id           INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    start        DATETIME NOT NULL,
    end          DATETIME NOT NULL, 
    create_date  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'users'
CREATE TABLE users (
    id           INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    email        VARCHAR(255) NOT NULL UNIQUE,
    name         VARCHAR(100) NOT NULL,
    password     VARCHAR(255) NOT NULL,
    create_date  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'team_users'
CREATE TABLE team_users (
    team_id        INT UNSIGNED NOT NULL,
    user_id        INT UNSIGNED,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (team_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'contest_teams'
CREATE TABLE contest_teams (
    contest_id     INT UNSIGNED NOT NULL,
    team_id        INT UNSIGNED NOT NULL,
    FOREIGN KEY (contest_id) REFERENCES contests(id) ON DELETE CASCADE,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    PRIMARY KEY (contest_id, team_id) 
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'category'
CREATE TABLE category (
    id             INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name           VARCHAR(255) NOT NULL UNIQUE,
    create_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'questions'
CREATE TABLE questions (
    id             INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    category_id    INT UNSIGNED  NOT NULL,
    env            VARCHAR(255),
    description    VARCHAR(255),
    vmid           INT NOT NULL,
    answer         VARCHAR(255),
    create_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'points'
CREATE TABLE points (
    id              INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    team_id         INT UNSIGNED NOT NULL,
    question_id     INT UNSIGNED NOT NULL,
    contest_id      INT UNSIGNED NOT NULL,
    point           INT UNSIGNED NOT NULL,
    insert_date     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    create_date     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    FOREIGN KEY (contest_id) REFERENCES contests(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'contest_questions'
CREATE TABLE contest_questions (
    contest_id INT UNSIGNED NOT NULL,
    question_id INT UNSIGNED NOT NULL,
    point INT NOT NULL,
    create_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (contest_id) REFERENCES contests(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    PRIMARY KEY (contest_id,question_id) 
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- 'cloudinit'
CREATE TABLE cloudinit (
    contest_id INT UNSIGNED NOT NULL,
    question_id INT UNSIGNED NOT NULL,
    team_id                 INT UNSIGNED NOT NULL,
    filename                VARCHAR(255) NOT NULL,
    access                  VARCHAR(255),
    vmid                    INT,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    FOREIGN KEY (contest_id) REFERENCES contests(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    PRIMARY KEY (contest_id,question_id,team_id), 
    create_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- 'roles'
CREATE TABLE roles (
    id              INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    namespace       VARCHAR(255) NOT NULL,
    contest         VARCHAR(255)  NOT NULL,
    create_date     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT unique_name_namespace UNIQUE (name, namespace)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'permissions' (
CREATE TABLE permissions (
    id              INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    description     VARCHAR(255) NOT NULL,
    create_date     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'permission_roles'
CREATE TABLE role_permissions (
    role_id INT UNSIGNED NOT NULL,
    permission_id INT UNSIGNED NOT NULL,
    create_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'user_roles'
CREATE TABLE user_roles (
    role_id INT UNSIGNED NOT NULL,
    user_id INT UNSIGNED NOT NULL,
    create_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



-- -- 7. contests テーブルへの挿入
-- INSERT INTO contests (name, start, end) VALUES
-- ('create question', '2002-11-01 10:00:00', '2002-12-07 18:00:00');

-- -- 2. users テーブルへの挿入
-- -- 注: 実際の運用では、パスワードはハッシュ化して保存してください。
-- INSERT INTO users (email, name, password) VALUES
-- ('admin@admin.com', 'Admin', '$2a$10$4IcZRIMPxNCkhUGOrW/pW.2T.zW7blr.DVVsPWi0itCH7W.pHqVWe');

-- -- 3. roles テーブルへの挿入
-- INSERT INTO roles (name, namespace,contest) VALUES
-- ('ctf_admin',"ctf" ,'[ "all" ]'),
-- ('contest_admin',"ctf" ,'[  "1" ]'),
-- ('user',"ctf" ,'[ "all" ]');

-- -- 4. permissions テーブルへの挿入
-- INSERT INTO permissions (name, description) VALUES
-- ('view', 'view list'),
-- ('edit', 'view and edit'),
-- ('create', 'create anything but not edit anything');

-- INSERT INTO role_permissions (role_id, permission_id) VALUES
-- -- manage_contest (id = 1)
-- (1, 1),
-- -- submit_flag (id = 2)
-- (1, 2),
-- -- view_scores (id = 3)
-- (1, 3),
-- -- participant ロールには submit_flag と view_scores の権限を付与
-- (2, 1),
-- (2, 2),
-- (3, 1);
