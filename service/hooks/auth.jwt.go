package hooks

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"net/url"
	"service/service"
	"service/service/_base"
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
	meta := service.M{}

	db := r.Conn().Mongo(p, "")

	switch strings.ToLower(request.Type) {
	default:
		p.Warning("check", "an attempt to auth an unknown type", diary.M{
			"id":   request.Id,
			"type": request.Type,
		})
		uniform.Alert(401, "Incorrect login details")
	case "administrator":
		var administrator service.Administrator
		db.Read(r.Remainder(), _base.Database, "administrators", request.Id, &administrator, service.TagsAdministrator)
		links["administrators"] = []string{administrator.Id.Hex()}
		inverted = administrator.Inverted

		allowTags := make([]string, 0)
		denyTags := make([]string, 0)
		if administrator.Role != nil {
			var role service.AdministratorRole
			db.Read(r.Remainder(), _base.Database, service.CollectionAdministratorRoles, administrator.Role.Id.Hex(), &role, service.TagsAdministratorRole)
			if role.AllowTags != nil {
				allowTags = append(allowTags, role.AllowTags...)
			}
			if role.DenyTags != nil {
				denyTags = append(denyTags, role.DenyTags...)
			}
		}

		if inverted {
			tags = service.filter(denyTags, allowTags)
			tags = service.filter(tags, administrator.AllowTags)
			if administrator.DenyTags != nil {
				tags = append(tags, administrator.DenyTags...)
			}
		} else {
			tags = service.filter(allowTags, denyTags)
			tags = service.filter(tags, administrator.DenyTags)
			if administrator.AllowTags != nil {
				tags = append(tags, administrator.AllowTags...)
			}
		}

		response.TwoFactor = !contains([]string{"staging", "qa", "development", "dev", "localhost", "local"}, service.Env, false)
		db.Update(r.Remainder(), _base.Database, "administrators", administrator.Id.Hex(), uniform.M{
			"lastLoginAt": time.Now(),
			"counter": 0,
		}, nil, nil)

		break
	}

	domain, err := url.Parse(service.BaseApiUrl)
	if err != nil {
		panic(err)
	}
	response.Audience = domain.Host
	response.Issuer = service.AppProject
	response.Inverted = inverted
	response.Tags = tags
	response.ExpiresAt = time.Now().Add(service.JwtExpiryTime)
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