package repository

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/LainInTheWired/ctf-backend/pveapi/model"
	"github.com/amoghe/go-crypt"
	"github.com/cockroachdb/errors"
	"gopkg.in/yaml.v3"
)

// GenerateSHA512CryptHash は、指定されたパスワード、ソルト、反復回数を使用してSHA-512 Cryptハッシュを生成します。
func GenerateSHA512CryptHash(password, salt string, rounds int) (string, error) {
	// パスワードとソルトを指定してハッシュを生成
	// "$6$" はSHA-512 Cryptを示すプレフィックス
	// "rounds=4096" は反復回数を指定
	hashed, err := crypt.Crypt(password, fmt.Sprintf("$6$rounds=%d$%s", rounds, salt))
	if err != nil {
		return "", err
	}
	return hashed, nil
}

func (r *pveRepository) CloudinitGenerator(fname string, host string, fqdn string, sshPwauth int, users []model.User) error {
	// salt := "randomsalt" // 任意のソルトを設定

	for i, u := range users {
		users[i].PlainTextPasswd = u.PlainTextPasswd
		users[i].LockPasswd = false
	}
	fmt.Println(users)

	config := model.CloudinitConfig{
		Hostname:  host,
		Users:     users,
		FQDN:      fqdn,
		SshPwauth: sshPwauth,
		Packages: []string{
			"git",
			"curl",
			"qemu-guest-agent",
		},
	}
	fmt.Printf("%+v", users)
	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		fmt.Printf("Error marshalling YAML: %v\n", err)
		return nil
	}
	yamlData = append([]byte("#cloud-config\n"), yamlData...)

	// ファイルに書き出し
	err = os.WriteFile(fname, yamlData, 0644)
	if err != nil {
		fmt.Printf("Error writing YAML to file: %v\n", err)
		return nil
	}
	return nil
}

func (r *pveRepository) TransferFileViaSCP(fname string) error {
	localPath := fname
	remoteUser := "root"
	// remoteHost := "10.0.10.30"
	remoteHost := os.Getenv("CLOUDINIT_FILESERVER_IP")
	remotePath := "/mnt/pve/cephfs/snippets/"
	cmd := exec.Command("scp",
		"-i", "/ssh/id_ed25519",
		"-o", "StrictHostKeyChecking=no",
		localPath,
		fmt.Sprintf("%s@%s:%s", remoteUser, remoteHost, remotePath),
	)

	// コマンドの標準出力と標準エラーを取得
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// コマンドの実行
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed scp")
	}

	if err := DeleteLocalfile(fname); err != nil {
		return errors.Wrap(err, "failed delete localfile")

	}
	return nil
}

func (r *pveRepository) DeleteFile(fname string) error {
	remoteUser := "root"
	remoteHost := os.Getenv("CLOUDINIT_FILESERVER_IP")
	remotePath := "/mnt/pve/cephfs/snippets/"

	// リモート上のファイルの完全なパスを作成
	remoteFilePath := filepath.Join(remotePath, fname)

	// 実行するリモートコマンドを定義
	remoteCmd := fmt.Sprintf("rm -f %s", remoteFilePath)

	// ssh コマンドの構築
	cmd := exec.Command("ssh",
		"-i", "/ssh/id_ed25519",
		"-o", "StrictHostKeyChecking=no",
		fmt.Sprintf("%s@%s", remoteUser, remoteHost),
		remoteCmd,
	)

	// コマンドの標準出力と標準エラーを取得
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// コマンドの実行
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to delete remote file via ssh")
	}
	return nil
}

func DeleteLocalfile(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return errors.Wrap(err, "failed to remove local file")
	}
	return nil
}
