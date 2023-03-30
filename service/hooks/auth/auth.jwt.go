package auth

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"github.com/go-uniform/uniform/common/nosql"
	"go.mongodb.org/mongo-driver/bson"
	"net/url"
	"service/service/_base"
	"service/service/entities"
	"service/service/info"
	"strings"
	"time"
)

type JwtRequest struct {
	Type string `bson:"type"`
	Id   string `bson:"id"`
}

type JwtResponse struct {
	TwoFactor  bool                   `bson:"twoFactor"`
	Issuer     string                 `bson:"issuer"`
	Audience   string                 `bson:"audience"`
	ExpiresAt  time.Time              `bson:"expiresAt"`
	ActivateAt *time.Time             `bson:"activateAt"`
	Inverted   bool                   `bson:"inverted"`
	Tags       []string               `bson:"tags"`
	Links      map[string][]string    `bson:"links"`
	Meta       map[string]interface{} `bson:"meta"`
}

func init() {
	_base.Subscribe(_base.TargetLocal(_base.TargetEvent("auth", "jwt")), eventAuthJwt)
}

func eventAuthJwt(r uniform.IRequest, p diary.IPage) {
	var request JwtRequest
	var response JwtResponse
	r.Read(&request)

	inverted := false
	tags := make([]string, 0)
	links := make(map[string][]string)
	meta := uniform.M{}

	db := nosql.Connector(r.Conn(), p, info.AppService)

	switch strings.ToLower(request.Type) {
	default:
		p.Warning("check", "an attempt to auth an unknown type", diary.M{
			"id":   request.Id,
			"type": request.Type,
		})
		uniform.Alert(401, "Incorrect login details")
	case "administrator":
		var administrator entities.Administrator
		db.FindOne(r.Remainder(), info.Database, entities.CollectionAdministrators, "", 0, bson.D{
			{"_id", request.Id},
		}, &administrator)

		links["administrators"] = []string{administrator.Id.Hex()}
		inverted = administrator.Inverted

		allowTags := make([]string, 0)
		denyTags := make([]string, 0)
		if administrator.Role != nil {
			var role entities.AdministratorRole
			db.FindOne(r.Remainder(), info.Database, entities.CollectionAdministratorRoles, "", 0, bson.D{
				{"_id", administrator.Role.Id},
			}, &role)
			if role.AllowTags != nil {
				allowTags = append(allowTags, role.AllowTags...)
			}
			if role.DenyTags != nil {
				denyTags = append(denyTags, role.DenyTags...)
			}
		}

		if inverted {
			tags = uniform.Filter(denyTags, allowTags, false)
			tags = uniform.Filter(tags, administrator.AllowTags, false)
			if administrator.DenyTags != nil {
				tags = append(tags, administrator.DenyTags...)
			}
		} else {
			tags = uniform.Filter(allowTags, denyTags, false)
			tags = uniform.Filter(tags, administrator.DenyTags, false)
			if administrator.AllowTags != nil {
				tags = append(tags, administrator.AllowTags...)
			}
		}

		response.TwoFactor = !uniform.Contains([]string{"staging", "qa", "development", "dev", "localhost", "local"}, info.Env, false)
		now := time.Now()
		administrator.LastLoginAt = &now
		administrator.LoginAttemptCounter = 0
		db.UpdateOne(r.Remainder(), info.Database, entities.CollectionAdministrators, bson.D{{"_id", administrator.Id}}, administrator, nil)

		break
	}

	domain, err := url.Parse(info.BaseApiUrl)
	if err != nil {
		panic(err)
	}
	response.Audience = domain.Host
	response.Issuer = info.AppProject
	response.Inverted = inverted
	response.Tags = tags
	response.ExpiresAt = time.Now().Add(info.JwtExpiryTime)
	response.Links = links
	response.Meta = meta

	if err := r.Reply(uniform.Request{
		Model: response,
	}); err != nil {
		p.Error("reply", err.Error(), diary.M{
			"error": err,
		})
	}
}
