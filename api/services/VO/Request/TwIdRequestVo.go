package Request

type TWIDParams struct {
	IdentityName  string `form:"IdentityName" validate:"required"`
	IdentityId    string `form:"IdentityId" validate:"required"`
	IssueType     string `form:"IssueType" validate:"required"`
	IssueDate     Date   `form:"IssueDate"`
	IssueCounties string `form:"IssueCounties" validate:"required"`
}

type Date struct {
	IssueYear  string `form:"IssueYear" validate:"required"`
	IssueMonth string `form:"IssueMonth" validate:"required"`
	IssueDay   string `form:"IssueDay" validate:"required"`
}

type VerifyEmailParams struct {
	Email  string `json:"Email" form:"Email" validate:"required"`
}