package hooks

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"net/url"
	"service/service/_base"
	"service/service/entities"
	"service/service/info"
	"strings"
	"time"
)

func init() {
	_base.Subscribe(_base.TargetEvent("auth", "jwt"), eventAuthJwt)
}

func eventAuthJwt(r uniform.IRequest, p diary.IPage) {
	var request uniform.AuthJwtRequest
	var response uniform.AuthJwtResponse
	r.Read(&request)

	inverted := false
	tags := make([]string, 0)
	links := make(map[string][]string)
	meta := uniform.M{}

	db := r.Conn().Mongo(p, "")

	switch strings.ToLower(request.Type) {
	default:
		p.Warning("check", "an attempt to auth an unknown type", diary.M{
			"id":   request.Id,
			"type": request.Type,
		})
		uniform.Alert(401, "Incorrect login details")
	case entities.CollectionAdministrators:
		var administrator entities.Administrator
		db.Read(r.Remainder(), info.Database, entities.CollectionAdministrators, request.Id, &administrator, nil)

		links["administrators"] = []string{administrator.Id.Hex()}
		inverted = administrator.Inverted

		allowTags := make([]string, 0)
		denyTags := make([]string, 0)
		if administrator.Role != nil {
			var role entities.AdministratorRole
			db.Read(r.Remainder(), info.Database, entities.CollectionAdministratorRoles, administrator.Role.Id.Hex(), &role, nil)
			if role.AllowTags != nil {
				allowTags = append(allowTags, role.AllowTags...)
			}
			if role.DenyTags != nil {
				denyTags = append(denyTags, role.DenyTags...)
			}
		}

		if inverted {
			tags = uniform.Filter(denyTags, allowTags)
			tags = uniform.Filter(tags, administrator.AllowTags)
			if administrator.DenyTags != nil {
				tags = append(tags, administrator.DenyTags...)
			}
		} else {
			tags = uniform.Filter(allowTags, denyTags)
			tags = uniform.Filter(tags, administrator.DenyTags)
			if administrator.AllowTags != nil {
				tags = append(tags, administrator.AllowTags...)
			}
		}

		response.TwoFactor = !uniform.Contains([]string{"staging", "qa", "development", "dev", "localhost", "local"}, info.Env, false)
		db.Update(r.Remainder(), info.Database, entities.CollectionAdministrators, administrator.Id.Hex(), uniform.M{
			"lastLoginAt": time.Now(),
			"counter": 0,
		}, nil, nil)

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