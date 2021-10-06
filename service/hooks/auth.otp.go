package hooks

import (
	"fmt"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"github.com/go-uniform/uniform/nosql"
	"go.mongodb.org/mongo-driver/bson"
	"service/service/_base"
	"service/service/entities"
	"service/service/info"
	"strings"
)

func init() {
	_base.Subscribe(_base.TargetEvent("auth", "otp"), eventAuthOtp)
}

func eventAuthOtp(r uniform.IRequest, p diary.IPage) {
	var request uniform.AuthOtpRequest
	var response uniform.AuthOtpResponse
	r.Read(&request)

	db := nosql.Request(r.Conn(), p, "")

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
		db.FindOne(r.Remainder(), info.Database, entities.CollectionAdministrators, "", 0, bson.D{
			{"_id", request.Identifier },
		}, &contact)
		break
	}

	// handle the login otp scenario
	if !request.Reset {
		r.Conn().SendSmsTemplate(p, r.Remainder(), info.MustAsset, "", "auth.otp.login.code", uniform.M{
			"Name":    contact.Name,
			"Project": strings.ToTitle(info.AppProject),
			"Code":    *request.Code,
			"Env":     info.EnvPrefix(),
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
			r.Conn().SendEmailTemplate(p, r.Remainder(), info.MustAsset, "", "from", "fromName", "auth.otp.reset.code", uniform.M{
				"Name":    contact.Name,
				"Project": strings.ToTitle(info.AppProject),
				"Code":    *request.Code,
				"Env":     info.EnvPrefix(),
			}, contact.Email)
		}
		r.Conn().SendSmsTemplate(p, r.Remainder(), info.MustAsset, "", "auth.otp.reset.code", uniform.M{
			"Name":    contact.Name,
			"Project": strings.ToTitle(info.AppProject),
			"Code":    *request.Code,
			"Env":     info.EnvPrefix(),
		}, contact.Mobile)
		break
	case "link":
		if contact.Email != "" {
			r.Conn().SendEmailTemplate(p, r.Remainder(), info.MustAsset, "", "from", "fromName", "auth.otp.reset.link", uniform.M{
				"Name":    contact.Name,
				"Project": strings.ToTitle(info.AppProject),
				"Link":    fmt.Sprintf("%s/#/password-reset-complete?token=%s", info.BaseAdministratorPortalUrl, *request.Token),
				"Env":     info.EnvPrefix(),
			}, contact.Email)
		}
		r.Conn().SendSmsTemplate(p, r.Remainder(), info.MustAsset, "", "auth.otp.reset.link", uniform.M{
			"Name":    contact.Name,
			"Project": strings.ToTitle(info.AppProject),
			"Link":    fmt.Sprintf("%s/#/password-reset-complete?token=%s", info.BaseAdministratorPortalUrl, *request.Token),
			"Env":     info.EnvPrefix(),
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