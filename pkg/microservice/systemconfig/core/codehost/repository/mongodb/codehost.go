package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/koderover/zadig/pkg/microservice/systemconfig/config"
	"github.com/koderover/zadig/pkg/microservice/systemconfig/core/codehost/repository/models"
	"github.com/koderover/zadig/pkg/tool/log"
	mongotool "github.com/koderover/zadig/pkg/tool/mongo"
)

type CodehostColl struct {
	*mongo.Collection

	coll string
}

type ListArgs struct {
	Owner   string
	Address string
	Source  string
}

func NewCodehostColl() *CodehostColl {
	name := models.CodeHost{}.TableName()
	coll := &CodehostColl{Collection: mongotool.Database(config.MongoDatabase()).Collection(name), coll: name}

	return coll
}

func (c *CodehostColl) GetCollectionName() string {
	return c.coll
}
func (c *CodehostColl) EnsureIndex(ctx context.Context) error {
	return nil
}

func (c *CodehostColl) AddCodeHost(iCodeHost *models.CodeHost) (*models.CodeHost, error) {

	_, err := c.Collection.InsertOne(context.TODO(), iCodeHost)
	if err != nil {
		log.Error("repository AddCodeHost err : %v", err)
		return nil, err
	}
	return iCodeHost, nil
}

func (c *CodehostColl) DeleteCodeHost() error {
	query := bson.M{"deleted_at": 0}
	change := bson.M{"$set": bson.M{
		"deleted_at": time.Now().Unix(),
	}}

	_, err := c.Collection.UpdateOne(context.TODO(), query, change)
	if err != nil {
		log.Error("repository DeleteCodeHostByID err : %v", err)
		return err
	}
	return nil
}

func (c *CodehostColl) GetCodeHostByID(ID int) (*models.CodeHost, error) {

	codehost := new(models.CodeHost)
	query := bson.M{"id": ID, "deleted_at": 0}
	err := c.Collection.FindOne(context.TODO(), query).Decode(codehost)
	if err != nil {
		return nil, nil
	}
	if v, ok := config.CodeHostMap[codehost.Type]; ok {
		codehost.Type = v
	}
	return codehost, nil
}

func (c *CodehostColl) List(args *ListArgs) ([]*models.CodeHost, error) {
	codeHosts := make([]*models.CodeHost, 0)
	query := bson.M{"deleted_at": 0}
	if args.Address != "" {
		query["address"] = args.Address
	}
	if args.Owner != "" {
		query["namespace"] = args.Owner
	}
	if args.Source != "" {
		query["type"] = args.Source
	}

	cursor, err := c.Collection.Find(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &codeHosts)
	if err != nil {
		return nil, err
	}
	// NOTE: to adapt old data
	for i, v := range codeHosts {
		if v, ok := config.CodeHostMap[v.Type]; ok {
			codeHosts[i].Type = v
		}
	}
	return codeHosts, nil
}

func (c *CodehostColl) CodeHostList() ([]*models.CodeHost, error) {
	codeHosts := make([]*models.CodeHost, 0)

	cursor, err := c.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &codeHosts)
	if err != nil {
		return nil, err
	}
	return codeHosts, nil
}

func (c *CodehostColl) DeleteCodeHostByID(ID int) error {
	query := bson.M{"id": ID, "deleted_at": 0}
	change := bson.M{"$set": bson.M{
		"deleted_at": time.Now().Unix(),
	}}
	_, err := c.Collection.UpdateOne(context.TODO(), query, change)
	if err != nil {
		log.Error("repository update fail")
		return err
	}
	return nil
}

func (c *CodehostColl) UpdateCodeHost(host *models.CodeHost) (*models.CodeHost, error) {
	query := bson.M{"id": host.ID, "deleted_at": 0}
	change := bson.M{"$set": bson.M{
		"type":           host.Type,
		"address":        host.Address,
		"namespace":      host.Namespace,
		"application_id": host.ApplicationId,
		"client_secret":  host.ClientSecret,
		"region":         host.Region,
		"username":       host.Username,
		"password":       host.Password,
		"updated_at":     time.Now().Unix(),
	}}
	_, err := c.Collection.UpdateOne(context.TODO(), query, change)
	return host, err
}

func (c *CodehostColl) UpdateCodeHostByToken(host *models.CodeHost) (*models.CodeHost, error) {
	query := bson.M{"id": host.ID, "deleted_at": 0}
	change := bson.M{"$set": bson.M{
		"is_ready":      "2",
		"access_token":  host.AccessToken,
		"updated_at":    time.Now().Unix(),
		"refresh_token": host.RefreshToken,
	}}
	_, err := c.Collection.UpdateOne(context.TODO(), query, change)
	return host, err
}
