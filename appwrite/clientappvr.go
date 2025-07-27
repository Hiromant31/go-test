package appwrite

import (
	"fmt"
	"go-test/config"

	"github.com/appwrite/sdk-for-go/account"
	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/databases"
	"github.com/appwrite/sdk-for-go/models"
)

type AppwriteClient struct {
	Client     *client.Client
	Account    *account.Account
	Databases  *databases.Databases
	ProjectID  string
	DatabaseID string
}

// Клиент для операций от имени пользователя по JWT
func NewClientWithJWT(jwt string) *AppwriteClient {
	cli := appwrite.NewClient(
		appwrite.WithEndpoint(config.Cfg.AppwriteHost),
		appwrite.WithProject(config.Cfg.AppwriteProject),
		appwrite.WithJWT(jwt),
	)

	return &AppwriteClient{
		Client:     &cli,
		Account:    account.New(cli),
		Databases:  databases.New(cli),
		ProjectID:  config.Cfg.AppwriteProject,
		DatabaseID: config.Cfg.AppwriteDatabaseID,
	}
}

// Клиент с админским API-ключом
func NewAppwriteClient() *AppwriteClient {
	cli := appwrite.NewClient(
		appwrite.WithEndpoint(config.Cfg.AppwriteHost),
		appwrite.WithProject(config.Cfg.AppwriteProject),
		appwrite.WithKey(config.Cfg.AppwriteKey),
	)

	return &AppwriteClient{
		Client:     &cli,
		Account:    account.New(cli),
		Databases:  databases.New(cli),
		ProjectID:  config.Cfg.AppwriteProject,
		DatabaseID: config.Cfg.AppwriteDatabaseID,
	}
}

// Создаём пользователя через Account API (максимум 36 символов; можно ID.unique())
func (c *AppwriteClient) CreateUser(userId, email, password string) (*models.User, error) {
	fmt.Println("UserID:", userId, "Email:", email)
	user, err := c.Account.Create(userId, email, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Логин через email+пароль → возвращает сессию
func (a *AppwriteClient) LoginUser(email, password string) (*models.Session, error) {
	session, err := a.Account.CreateEmailPasswordSession(email, password)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// Генерация JWT (только после успешного входа)
func (a *AppwriteClient) CreateJWT() (*models.Jwt, error) {
	jwtObj, err := a.Account.CreateJWT()
	if err != nil {
		return nil, err
	}
	return jwtObj, nil
}

func (a *AppwriteClient) GetCollections(queries []string, search string) (*models.CollectionList, error) {
	var opts []databases.ListCollectionsOption
	if len(queries) > 0 {
		opts = append(opts, a.Databases.WithListCollectionsQueries(queries))
	}
	if search != "" {
		opts = append(opts, a.Databases.WithListCollectionsSearch(search))
	}
	return a.Databases.ListCollections(a.DatabaseID, opts...)
}

func (a *AppwriteClient) CreateCollection(
	name string,
	permissions []string,
	documentSecurity, enabled bool,
) (*models.Collection, error) {
	return a.Databases.CreateCollection(
		a.DatabaseID, name, name,
		a.Databases.WithCreateCollectionPermissions(permissions),
		a.Databases.WithCreateCollectionDocumentSecurity(documentSecurity),
		a.Databases.WithCreateCollectionEnabled(enabled),
	)
}
