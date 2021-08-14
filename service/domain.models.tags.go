package service

var TagsAdministrators = map[string][]string{
	"created-at": { "type:string","format:date-time","not-editable" },
	"id": { "type:string","format:uuid","not-editable","not-header" },
	"allow-tags": { "type:string","array","not-header","not-sortable","not-filterable" },
	"deny-tags": { "type:string","array","not-header","not-sortable","not-filterable" },
	"first-name": { "type:string" },
	"last-name": { "type:string" },
	"email": { "type:string","encrypted","identifier","not-sortable" },
	"mobile": { "type:string","encrypted","identifier","not-sortable" },
	"password": { "type:string","format:password","password","not-header","not-sortable","not-editable","not-filterable" },
	"role": { "type:link","entity:administratorRoles","field:id","optional" },
	"inverted": { "type:boolean" },
	"modified-at": { "type:string","format:date-time","not-editable" },
	"deleted-at": { "type:string","format:date-time","not-editable","optional","not-header" },
}