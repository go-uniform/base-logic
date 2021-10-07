package auth

import (
	"encoding/base64"
	"fmt"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"service/service/_base"
	"service/service/info"
	"time"
)

type LoginRequest struct {
	Type       string `bson:"type"`
	Identifier string `bson:"identifier"`
	Password   string `bson:"password"`
}

type LoginResponse struct {
	TwoFactor bool   `bson:"twoFactor"`
	Token     string `bson:"token"`
}

type Entity struct {
	Id       string `bson:"id"`
	Password string `bson:"password"`
}

func init() {
	_base.Subscribe(_base.TargetLocal("auth.login"), func(r uniform.IRequest, p diary.IPage) {
		defer func() {
			if rec := recover(); rec != nil {
				if err, ok := rec.(error); ok {
					p.Notice("error", diary.M{
						"errorMsg": err.Error(),
						"error":    err,
					})
				} else {
					p.Notice("error", diary.M{
						"errorMsg": fmt.Sprint(rec),
						"error":    rec,
					})
				}
				uniform.Alert(500, "Incorrect login details")
				if r.CanReply() {
					if err := r.Reply(uniform.Request{}); err != nil {
						p.Error("reply", err.Error(), diary.M{
							"err": err,
						})
					}
				}
			}
		}()

		var request LoginRequest
		r.Read(&request)

		p.Notice("auth.login", diary.M{
			"type":       request.Type,
			"identifier": request.Identifier,
		})

		var entity Entity
		if err := r.Conn().Request(p, _base.TargetLocal(_base.TargetEvent("auth", "check")), r.Remainder(), uniform.Request{
			Model: uniform.M{
				"type":       request.Type,
				"identifier": request.Identifier,
			},
		}, func(r uniform.IRequest, p diary.IPage) {
			if r.HasError() {
				panic(r.Error())
			}
			r.Read(&entity)
		}); err != nil {
			p.Notice("failed.auth.check", diary.M{
				"type":       request.Type,
				"identifier": request.Identifier,
			})
			panic(err)
		}

		password, err := base64.StdEncoding.DecodeString(entity.Password)
		if err != nil {
			p.Notice("failed.base64", diary.M{
				"type":       request.Type,
				"identifier": request.Identifier,
			})
			panic(err)
		}

		if err := bcrypt.CompareHashAndPassword(password, []byte(request.Password)); err != nil {
			p.Notice("failed.compare", diary.M{
				"type":       request.Type,
				"identifier": request.Identifier,
			})

			if err := r.Conn().Request(p, _base.TargetLocal(_base.TargetEvent("auth", "failed")), r.Remainder(), uniform.Request{
				Model: uniform.M{
					"id":   entity.Id,
					"type": request.Type,
				},
			}, func(r uniform.IRequest, p diary.IPage) {
				if r.HasError() {
					panic(r.Error())
				}
			}); err != nil {
				p.Notice("failed.auth.failed", diary.M{
					"type":       request.Type,
					"identifier": request.Identifier,
				})
				panic(err)
			}
			panic("auth failed")
		}

		rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(info.MustAsset("jwt.key"))
		if err != nil {
			p.Notice("failed.jwt", diary.M{
				"type":       request.Type,
				"identifier": request.Identifier,
			})
			panic(err)
		}

		var out struct {
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
		if err := r.Conn().Request(p, _base.TargetLocal(_base.TargetEvent("auth", "jwt")), r.Remainder(), uniform.Request{
			Model: uniform.M{
				"id":   entity.Id,
				"type": request.Type,
			},
		}, func(r uniform.IRequest, p diary.IPage) {
			if r.HasError() {
				panic(r.Error())
			}
			r.Read(&out)
		}); err != nil {
			p.Notice("failed.auth.jwt", diary.M{
				"type":       request.Type,
				"identifier": request.Identifier,
			})
			panic(err)
		}

		claims := jwt.MapClaims{
			"id":                   entity.Id,
			"type":                 request.Type,
			"permissions.inverted": out.Inverted,
			"permissions.tags":     out.Tags,
			"links":                out.Links,
			"meta":                 out.Meta,
			"aud":                  out.Audience,
			"exp":                  out.ExpiresAt.Unix(),
			"iat":                  time.Now().Unix(),
		}

		if out.ActivateAt != nil {
			claims["nbf"] = out.ActivateAt.Unix()
		}

		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

		signed, err := token.SignedString(rsaPrivateKey)
		if err != nil {
			p.Notice("failed.sign", diary.M{
				"type":       request.Type,
				"identifier": request.Identifier,
			})
			panic(err)
		}

		loginResponse := LoginResponse{
			Token: signed,
		}
		if out.TwoFactor {
			var codeEntity struct {
				Id   string `bson:"id"`
				Code string `bson:"code"`
			}
			if err := r.Conn().Request(p, _base.TargetLocal(_base.TargetEvent("code", "issue")), r.Remainder(), uniform.Request{
				Model: uniform.M{
					"meta": uniform.M{
						"token":      signed,
						"type":       request.Type,
						"identifier": request.Identifier,
					},
				},
			}, func(r uniform.IRequest, p diary.IPage) {
				if r.HasError() {
					panic(r.Error())
				}
				r.Read(&codeEntity)
			}); err != nil {
				p.Notice("failed.code", diary.M{
					"type":       request.Type,
					"identifier": request.Identifier,
				})
				panic(err)
			}

			if err := r.Conn().Request(p, _base.TargetLocal(_base.TargetEvent("auth", "login.send")), r.Remainder(), uniform.Request{
				Model: uniform.M{
					"type":       request.Type,
					"identifier": request.Identifier,
					"method":     "code",
					"code":       codeEntity.Code,
				},
			}, func(r uniform.IRequest, p diary.IPage) {
				if r.HasError() {
					panic(r.Error())
				}
				r.Read(&codeEntity)
			}); err != nil {
				p.Notice("failed.send", diary.M{
					"type":       request.Type,
					"identifier": request.Identifier,
				})
				panic(err)
			}

			loginResponse.TwoFactor = true
			loginResponse.Token = codeEntity.Id
		}

		p.Notice("success", diary.M{
			"type":       request.Type,
			"identifier": request.Identifier,
		})

		if r.CanReply() {
			if err := r.Reply(uniform.Request{
				Model: loginResponse,
			}); err != nil {
				p.Error("reply", err.Error(), diary.M{
					"err": err,
				})
			}
		}
	})
}
