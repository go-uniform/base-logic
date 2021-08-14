package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"net/url"
	"strings"
	"time"
)

func init() {
	subscribe(event("auth", "jwt"), eventAuthJwt)
}

func eventAuthJwt(r uniform.IRequest, p diary.IPage) {
	var request uniform.AuthJwtRequest
	var response uniform.AuthJwtResponse
	r.Read(&request)

	inverted := false
	tags := make([]string, 0)
	links := make(map[string][]string)
	meta := M{}

	db := r.Conn().Mongo(p, "")

	switch strings.ToLower(request.Type) {
	default:
		p.Warning("check", "an attempt to auth an unknown type", diary.M{
			"id":   request.Id,
			"type": request.Type,
		})
		uniform.Alert(401, "Incorrect login details")
	case "administrator":
		var administrator Administrator
		db.Read(r.Remainder(), Database, CollectionAdministrators, request.Id, &administrator, TagsAdministrator)
		links["administrators"] = []string{administrator.Id.Hex()}
		inverted = administrator.Inverted

		allowTags := make([]string, 0)
		denyTags := make([]string, 0)
		if administrator.Role != nil {
			var role AdministratorRole
			db.Read(r.Remainder(), Database, CollectionAdministratorRoles, administrator.Role.Id.Hex(), &role, TagsAdministratorRole)
			if role.AllowTags != nil {
				allowTags = append(allowTags, role.AllowTags...)
			}
			if role.DenyTags != nil {
				denyTags = append(denyTags, role.DenyTags...)
			}
		}

		if inverted {
			tags = filter(denyTags, allowTags)
			tags = filter(tags, administrator.AllowTags)
			if administrator.DenyTags != nil {
				tags = append(tags, administrator.DenyTags...)
			}
		} else {
			tags = filter(allowTags, denyTags)
			tags = filter(tags, administrator.DenyTags)
			if administrator.AllowTags != nil {
				tags = append(tags, administrator.AllowTags...)
			}
		}

		response.TwoFactor = !contains([]string{"staging", "qa", "development", "dev", "localhost", "local"}, Env, false)
		db.Update(r.Remainder(), Database, CollectionAdministrators, administrator.Id.Hex(), uniform.M{
			"lastLoginAt": time.Now(),
			"counter": 0,
		}, nil, nil)

		break
	}

	domain, err := url.Parse(BaseApiUrl)
	if err != nil {
		panic(err)
	}
	response.Audience = domain.Host
	response.Issuer = AppProject
	response.Inverted = inverted
	response.Tags = tags
	response.ExpiresAt = time.Now().Add(JwtExpiryTime)
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