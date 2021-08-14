package service

import (
	"fmt"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"strings"
)

func init() {
	subscribe(event("auth", "otp"), eventAuthOtp)
}

func eventAuthOtp(r uniform.IRequest, p diary.IPage) {
	var request uniform.AuthOtpRequest
	var response uniform.AuthOtpResponse
	r.Read(&request)

	db := r.Conn().Mongo(p, "")

	id := request.Id
	if id == "" {
		r.Conn().Request(p, "", r.Remainder(), uniform.Request{
			Model: uniform.AuthCheckRequest{
				Type: request.Type,
				Identifier: request.Identifier,
				Reset: true,
			},
		}, func(r uniform.IRequest, p diary.IPage) {
			var entity uniform.AuthCheckResponse
			r.Read(&entity)
			id = entity.Id
		})
	}

	var contact struct {
		Name       string
		Email      string
		Mobile     string
	}

	// populate contact based on auth type and id
	switch strings.ToLower(request.Type) {
	default:
		p.Warning("check", "an attempt to auth an unknown type", diary.M{
			"id":   request.Id,
			"type": request.Type,
		})
		uniform.Alert(401, "Incorrect login details")
	case "administrator":
		db.Read(r.Remainder(), Database, CollectionAdministrators, request.Id, &contact, TagsAdministrator)
		break
	}

	// handle the login otp scenario
	if !request.Reset {
		r.Conn().SendSmsTemplate(p, r.Remainder(), nil, "", "auth.otp.login.code", uniform.M{
			"Name": contact.Name,
			"Project": strings.ToTitle(AppProject),
			"Code": *request.Code,
			"Env": EnvPrefix(),
		}, contact.Mobile)

		if err := r.Reply(uniform.Request{
			Model: response,
		}); err != nil {
			p.Error("reply", err.Error(), diary.M{
				"error": err,
			})
		}
		return
	}

	// handle the password reset otp code/link scenario
	switch strings.ToLower(request.Method) {
	default:
		panic(fmt.Sprintf("unsupported method `%s` given", request.Method))
	case "code":
		if contact.Email != "" {
			r.Conn().SendEmailTemplate(p, r.Remainder(), nil, "", "from", "fromName", "auth.otp.reset.code", uniform.M{
				"Name": contact.Name,
				"Project": strings.ToTitle(AppProject),
				"Code": *request.Code,
				"Env": EnvPrefix(),
			}, contact.Email)
		}
		r.Conn().SendSmsTemplate(p, r.Remainder(), nil, "", "auth.otp.reset.code", uniform.M{
			"Name": contact.Name,
			"Project": strings.ToTitle(AppProject),
			"Code": *request.Code,
			"Env": EnvPrefix(),
		}, contact.Mobile)
		break
	case "link":
		if contact.Email != "" {
			r.Conn().SendEmailTemplate(p, r.Remainder(), nil, "", "from", "fromName", "auth.otp.reset.link", uniform.M{
				"Name": contact.Name,
				"Project": strings.ToTitle(AppProject),
				"Link": fmt.Sprintf("%s/#/password-reset-complete?token=%s", BaseAdministratorPortalUrl, *request.Token),
				"Env": EnvPrefix(),
			}, contact.Email)
		}
		r.Conn().SendSmsTemplate(p, r.Remainder(), nil, "", "auth.otp.reset.link", uniform.M{
			"Name": contact.Name,
			"Project": strings.ToTitle(AppProject),
			"Link": fmt.Sprintf("%s/#/password-reset-complete?token=%s", BaseAdministratorPortalUrl, *request.Token),
			"Env": EnvPrefix(),
		}, contact.Mobile)
		break
	}

	if err := r.Reply(uniform.Request{
		Model: response,
	}); err != nil {
		p.Error("reply", err.Error(), diary.M{
			"error": err,
		})
	}
}