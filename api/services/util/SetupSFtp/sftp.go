package SetupSFtp

import (
	"api/services/util/log"
	m_sort "api/services/util/m-sort"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type SFTPClient struct {
	host     string
	username string
	paswword string
	conn     *sftp.Client
}

// 建立 ssh 連線
func setupSSHTunnel(host, username, password string) (*ssh.Client, error) {
	var auths []ssh.AuthMethod

	if aConn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aConn).Signers))
	}
	auths = append(auths, ssh.Password(password))

	return ssh.Dial("tcp", host, &ssh.ClientConfig{
		User:            username,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
}

// 建立 sftp
func SFTPConnect(host, username, password string) (sftpConn *sftp.Client, error error) {
	sshConnect, err := setupSSHTunnel(host, username, password)
	if err != nil {
		log.Error("SSHTunnel 建立失敗 [%s]", host)
		log.Error("SSHTunnel 建立失敗  Error [%s]", err.Error())
		return sftpConn, fmt.Errorf("SSH 建立失敗")
	}

	sftpConn, err = sftp.NewClient(sshConnect, sftp.MaxPacket(1<<15))
	if err != nil {
		log.Error("SFTP 連線失敗 [%s]", host)
		log.Error("SFTP 連線失敗  Error [%s]", err.Error())
		return sftpConn, err
	}

	return sftpConn, nil
}

func NewSFTPClientAndLogin(host, username, password string) (SFTPClient, error) {
	cli := SFTPClient{}
	var auths []ssh.AuthMethod
	if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))
	}
	auths = append(auths, ssh.Password(password))

	config := ssh.ClientConfig{
		User:            username,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", host, &config)
	if err != nil {
		log.Error("unable to connect to [", host, "]:", err)
		return cli, err
	}

	c, err := sftp.NewClient(conn, sftp.MaxPacket(1<<15))
	if err != nil {
		log.Error("unable to start sftp subsytem: %v", err)
		return cli, err
	}

	cli.host = host
	cli.username = username
	cli.paswword = password
	cli.conn = c

	return cli, nil
}

func (receiver *SFTPClient) LogoutAndClose() error {
	return receiver.conn.Close()
}

func (receiver *SFTPClient) GetAEarlierFileInDir(dir, filePrefix, fileExt string) ([]byte, string, error) {
	files, err := receiver.conn.ReadDir(dir)
	if err != nil {
		return nil, "", err
	}
	names := []string{}
	for _, v := range files {
		name := v.Name()
		if strings.HasPrefix(name, filePrefix) && strings.HasSuffix(name, fileExt) {
			names = append(names, name)
		}
	}

	name, err := m_sort.SortStringsAscAndGetFirst(names)
	if err != nil {
		return nil, "", err
	}

	file, err := receiver.conn.Open(dir + "./" + name)
	if err != nil {
		return nil, "", errors.New("[GetAEarlierFileInDir]Open file fail")
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, "", errors.New("[GetAEarlierFileInDir]Read file fail")
	}

	return data, name, nil
}
